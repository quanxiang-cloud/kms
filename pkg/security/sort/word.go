package sort

import (
	"sort"
)

type word []Intermediate

func (s word) Len() int           { return len(s) }
func (s word) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s word) Less(i, j int) bool { return s[i].Key < s[j].Key }

// WordSortASC word sort with asc
func WordSortASC(entity interface{}, format FormatType, mapper MapperType) ([]byte, error) {
	value, _type, err := Get(entity, mapper)
	if err != nil {
		return nil, err
	}

	wordSort(value)
	b, err := Format(value, _type, format)

	return b, err
}

// WordSortDESC word sort with desc
func WordSortDESC(entity interface{}, format FormatType, mapper MapperType) ([]byte, error) {
	value, _type, err := Get(entity, mapper)
	if err != nil {
		return nil, err
	}

	wordSort(value)

	// 数据倒置
	length := len(value)
	half := length / 2
	for i := 0; i < half; i++ {
		value[i], value[length-1-i] = value[length-1-i], value[i]
	}

	return Format(value, _type, format)
}

func wordSort(entity word) {
	sort.Sort(entity)
	for _, k := range entity {
		val, ok := k.Value.(word)
		if !ok {
			continue
		}
		wordSort(val)
		k.Value = val
	}
}
