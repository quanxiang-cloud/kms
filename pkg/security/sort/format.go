package sort

import (
	"bytes"
	"encoding/json"
	x "encoding/xml"
	"fmt"
	"net/url"
)

// FormatType format type
type FormatType string

const (
	// DefaultFormat json
	DefaultFormat FormatType = JSONFormat
	// JSONFormat json
	JSONFormat FormatType = "json"
	// XMLFormat xml
	XMLFormat FormatType = "xml"
	// QueryFormat query
	QueryFormat FormatType = "query"
)

// OneofFormat one of format
func OneofFormat(t string) (FormatType, bool) {
	switch FormatType(t) {
	case JSONFormat:
		return JSONFormat, true
	case XMLFormat:
		return XMLFormat, true
	case QueryFormat:
		return QueryFormat, true
	}
	return "", false
}

// Format format
func Format(entity []Intermediate, _t T, format FormatType) ([]byte, error) {
	return getFormat(format)(entity, _t)
}

type format = func(entity []Intermediate, _t T) ([]byte, error)

func getFormat(format FormatType) format {
	switch format {
	case DefaultFormat:
		return jsonFormat
	case XMLFormat:
		return xmlFormat
	case QueryFormat:
		return queryFormat
	default:
		return jsonFormat
	}
}

func gen(entity []Intermediate, _t T) interface{} {
	size := len(entity)
	if size == 0 {
		return nil
	}
	if _t == _slice {
		value := make([]interface{}, 0, len(entity))
		for _, val := range entity {
			v, ok := val.Value.([]Intermediate)
			if ok {
				value = append(value, gen(v, val.Type))
				continue
			}
			value = append(value, val.Value)
		}
		return value
	}

	value := make(map[string]interface{}, len(entity))
	for _, val := range entity {
		v, ok := val.Value.([]Intermediate)
		if ok {
			value[val.Key] = gen(v, val.Type)
			continue
		}
		value[val.Key] = val.Value
	}
	return value
}

func jsonFormat(entity []Intermediate, _t T) ([]byte, error) {
	return json.Marshal(gen(entity, _t))
}

func xmlFormat(entity []Intermediate, _t T) ([]byte, error) {
	return x.Marshal(form(entity))
}

type form []Intermediate

type xmlSub struct {
	name  string
	value []Intermediate
}

// xmlMap marshals into XML
func (f form) MarshalXML(e *x.Encoder, start x.StartElement) error {

	tokens := []x.Token{start}

	for _, value := range f {
		if value.Key == "" {
			value.Key = "elem"
		}
		t := x.StartElement{Name: x.Name{Space: "", Local: value.Key}}
		val, ok := value.Value.([]Intermediate)
		if ok {
			if value.Type == _slice {
				for index, v := range val {
					name := fmt.Sprintf("%s%d", t.Name.Local, index+1)
					sub, ok := v.Value.([]Intermediate)
					if ok {
						tokens = append(tokens, xmlSub{name: name, value: sub})
						continue
					}
					tokens = append(tokens, x.StartElement{Name: x.Name{Space: "", Local: name}},
						x.CharData(fmt.Sprint(v.Value)),
						x.EndElement{Name: x.Name{Space: "", Local: name}})
				}
				continue
			}
			tokens = append(tokens, xmlSub{name: value.Key, value: val})
			continue
		}

		tokens = append(tokens, t, x.CharData(fmt.Sprint(value.Value)), x.EndElement{Name: t.Name})
	}

	tokens = append(tokens, x.EndElement{Name: start.Name})

	for _, t := range tokens {
		var err error
		switch t.(type) {
		case xmlSub:
			val := t.(xmlSub)

			err = e.EncodeElement(form(val.value), x.StartElement{Name: x.Name{Space: "", Local: val.name}})
		default:
			err = e.EncodeToken(t)
		}

		if err != nil {
			return err
		}
	}

	return e.Flush()
}

func queryFormat(entity []Intermediate, _t T) ([]byte, error) {
	var buf bytes.Buffer
	for _, val := range entity {
		if val.Type == _slice {
			for index, elem := range val.Value.([]Intermediate) {
				buf.WriteString(fmt.Sprintf("%s.%d", val.Key, index+1))
				buf.WriteString("=")
				// TODO 待优化
				buf.WriteString(fmt.Sprintf("%v", elem.Value))
				buf.WriteString("&")
			}
			continue
		}
		buf.WriteString(val.Key)
		buf.WriteString("=")
		// TODO 待优化
		// encode to url 
		// BUG: ? _ encode fail, need rewrite Escape func
		value := val.Value
		if v, ok := value.(string); ok {
			value = url.QueryEscape(v)
		}
		buf.WriteString(fmt.Sprintf("%v", value))
		buf.WriteString("&")
	}

	query := buf.Bytes()
	if size := len(query); size > 1 {
		query = query[:size-1]
	}

	return query, nil
}
