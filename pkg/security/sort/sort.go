package sort

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"reflect"
)

// OrderType sort type
type OrderType string

const (
	// DefaultSort default sort asc
	DefaultSort = SortASC
	// SortDESC sort desc
	SortDESC OrderType = "desc"
	// SortASC sort asc
	SortASC OrderType = "asc"
)

var (
	// ErrInvalidType invalid type
	ErrInvalidType = errors.New("invalid type")
)

// OneofOrder one of order
func OneofOrder(t string) (OrderType, bool) {
	switch OrderType(t) {
	case SortDESC:
		return SortDESC, true
	case SortASC:
		return SortASC, true
	}

	return "", false
}

// Structured 处理函数
type Structured func(interface{}, MapperType) ([]Intermediate, T, error)

func none(interface{}, MapperType) ([]Intermediate, T, error) {
	return nil, _default, ErrInvalidType
}

func getInterface(value interface{}, mapper MapperType) ([]Intermediate, T, error) {
	return get(value, getMapper(mapper))
}

// Intermediate 中间态数据
// 分离数据 并没有排序
type Intermediate struct {
	Key   string
	Value interface{}
	Type  T
}

// T  type
type T string

const (
	_default T = "default"
	_map     T = "map"
	_slice   T = "slice"
)

func get(value interface{}, mapper mapperGet) ([]Intermediate, T, error) {
	if value == nil {
		return nil, _default, nil
	}
	typeOfValue := reflect.TypeOf(value)
	valueOfValue := reflect.ValueOf(value)

	var keys []Intermediate

	switch typeOfValue.Kind() {
	case reflect.Ptr:
		return get(valueOfValue.Elem(), mapper)
	case reflect.Map:
		numField := valueOfValue.Len()
		keys = make([]Intermediate, 0, numField)

		iter := valueOfValue.MapRange()
		for iter.Next() {
			if !iter.Value().CanInterface() {
				continue
			}
			name := mapper(iter.Key().String())
			value := iter.Value().Interface()

			sub, subKind, err := get(value, mapper)
			if err != nil {
				return nil, subKind, err
			}
			if sub != nil {
				value = sub
			}
			keys = append(keys, Intermediate{
				Key:   name,
				Type:  subKind,
				Value: value,
			})
		}
		return keys, _map, nil
	case reflect.Struct:
		numField := typeOfValue.NumField()
		keys = make([]Intermediate, 0, numField)

		for index := 0; index < numField; index++ {
			if !valueOfValue.Field(index).CanInterface() {
				continue
			}
			var name string
			if name = typeOfValue.Field(index).Tag.Get("security"); name == "" {
				name = mapper(typeOfValue.Field(index).Name)
			}

			value := valueOfValue.Field(index).Interface()

			sub, subKind, err := get(value, mapper)

			if err != nil {
				return nil, subKind, err
			}
			if sub != nil {
				value = sub
			}
			keys = append(keys, Intermediate{
				Key:   name,
				Type:  subKind,
				Value: value,
			})
		}
		return keys, _default, nil
	case reflect.Array, reflect.Slice:
		for i := 0; i < valueOfValue.Len(); i++ {
			item := valueOfValue.Index(i)
			if !item.CanInterface() {
				continue
			}

			value := item.Interface()

			sub, subKind, err := get(value, mapper)

			if err != nil {
				return nil, subKind, err
			}
			if sub != nil {
				value = sub
			}
			keys = append(keys, Intermediate{
				Key:   "",
				Type:  subKind,
				Value: value,
			})
		}
		return keys, _slice, nil

	default:
		return nil, _default, nil
	}

}

type recur interface{}

// Get get
func Get(entity interface{}, mapper MapperType) ([]Intermediate, T, error) {
	var err error
	entity, err = mapping(entity)
	if err != nil {
		return nil, _default, err
	}

	typeOfValue := reflect.TypeOf(entity)
	valueOfValue := reflect.ValueOf(entity)
	if typeOfValue.Kind() == reflect.Ptr {
		return Get(valueOfValue.Elem().Interface(), mapper)
	}

	switch typeOfValue.Kind() {
	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		return getInterface(entity, mapper)
	default:
		return none(entity, mapper)
	}

}

var (
	// ErrType 类型不是json or xml
	ErrType = errors.New("type is not json or xml")
)

func mapping(entity interface{}) (recur, error) {
	if temp, ok := entity.(string); ok {
		entity = []byte(temp)
	}

	var r recur
	if temp, ok := entity.([]byte); ok {
		if bytes.HasPrefix(temp, []byte("{")) || bytes.HasPrefix(temp, []byte("[{")) {
			err := json.Unmarshal(temp, &r)
			return r, err
		} else if bytes.HasPrefix(temp, []byte("<")) {
			err := xml.Unmarshal(temp, &r)
			return r, err
		} else {
			return r, ErrType
		}
	}

	return entity, nil
}
