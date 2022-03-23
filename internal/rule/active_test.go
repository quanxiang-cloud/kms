package rule

import (
	"fmt"
	"testing"
)

func TestRule(t *testing.T) {
	type testingCase struct {
		sample interface{}
		expect interface{}
		fun    string
	}

	testingCases := []*testingCase{
		{
			sample: 1,
			expect: true,
			fun:    "IsActive",
		},
		{
			sample: []int{1, int(OpCreate)},
			expect: nil,
			fun:    "ValidateActive",
		},
		{
			sample: []int{0, int(OpDelete)},
			expect: nil,
			fun:    "ValidateActive",
		},
		{
			sample: []int{1, int(OpQuery)},
			expect: nil,
			fun:    "ValidateActive",
		},
		{
			sample: []int{1, int(OpSignature)},
			expect: nil,
			fun:    "ValidateActive",
		},
		{
			sample: []int{0, int(OpUpdate)},
			expect: nil,
			fun:    "ValidateActive",
		},
		{
			sample: []int{1, int(OpDelete)},
			expect: nil,
			fun:    "ValidateActive",
		},
		{
			sample: []int{0, int(OpSignature)},
			expect: nil,
			fun:    "ValidateActive",
		},
		{
			sample: []int{1, int(OpUpdate)},
			expect: nil,
			fun:    "ValidateActive",
		},
	}

	for _, tc := range testingCases {
		switch tc.fun {
		case "IsActive":
			ret := IsActive(tc.sample.(int))
			if ret != tc.expect.(bool) {
				fail(tc, ret, tc.expect)
			}
		case "ValidateActive":
			cs := tc.sample.([]int)
			err := ValidateActive(cs[0], Operation(cs[1]))
			if err != nil {
				fail(tc, err, tc.expect)
			}
		}
	}
}

func fail(testingCase, ret, expect interface{}) {
	fmt.Printf("testing fail, case: %+v\nresult: %+v\nexpect: %+V\n", testingCase, ret, expect)
}
