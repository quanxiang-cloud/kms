package enums

import (
	"bytes"
	"fmt"
	"sort"
)

var allEnumSet []*EnumSet

func init() {
	AfterReg()
}

// AfterReg sort the enums
func AfterReg() {
	for _, v := range allEnumSet {
		v.Sort()
	}
}

// Enum represents a enum
type Enum string

// Val return enum value
func (e Enum) Val() string {
	return string(e)
}

// Equal equal
func (e Enum) Equal(s string) bool {
	return string(e) == s
}

//------------------------------------------------------------------------------

// NewEnumSet create a new EnumSet
func NewEnumSet(exists *EnumSet) *EnumSet {
	ret := &EnumSet{}
	if exists != nil {
		ret.list = append(ret.list, exists.list...)
	}
	allEnumSet = append(allEnumSet, ret)
	return ret
}

// EnumSet represents a set of enum
type EnumSet struct {
	list   []string // enum set
	sorted []string // ordered list
}

// Reg regist a new enum to the set
func (es *EnumSet) Reg(val string) Enum {
	e := Enum(val)
	for _, v := range es.list {
		if e.Equal(v) {
			panic(fmt.Sprintf("dupicate enum value %s", val))
		}
	}

	es.list = append(es.list, e.Val())
	return e
}

// ShowAll show enum list of the set
func (es *EnumSet) ShowAll() string {
	var b bytes.Buffer
	for i := 0; i < len(es.list); i++ {
		v := es.list[i]
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(v)
	}
	return b.String()
}

// ShowSorted show enum list of the set with sort
func (es *EnumSet) ShowSorted() string {
	var b bytes.Buffer
	for i := 0; i < len(es.sorted); i++ {
		v := es.sorted[i]
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(v)
	}
	return b.String()
}

// GetAll return enum list of the set
func (es *EnumSet) GetAll() []string {
	return es.list
}

// Sort sort the enums
func (es *EnumSet) Sort() {
	if len(es.sorted) > 0 {
		return
	}
	es.sorted = append([]string{}, es.list...)
	sort.Strings(es.sorted)
}

// Verify check if a enum is valid, binary search
func (es *EnumSet) Verify(e string) bool {
	s := es.sorted
	low, high := 0, len(s)-1
	for low <= high {
		mid := (low + high) / 2
		switch {
		case e == s[mid]:
			return true
		case e > s[mid]:
			low = mid + 1
		case e < s[mid]:
			high = mid - 1
		}
	}

	return false
}
