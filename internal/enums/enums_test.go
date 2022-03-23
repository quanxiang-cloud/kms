package enums

import (
	"fmt"
	"testing"
)

func TestEnums(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	type testingCase struct {
		sample interface{}
		expect interface{}
		fun    string
	}

	testingCases := []testingCase{
		{
			sample: "testing",
			expect: Enum("testing"),
			fun:    "reg",
		},
		{
			sample: nil,
			expect: []string{"testing"},
			fun:    "array",
		},
		{
			sample: "testing",
			expect: nil,
			fun:    "reg",
		},
	}

	es := EnumSet{}
	for _, tc := range testingCases {
		switch tc.fun {
		case "reg":
			ret := es.Reg(tc.sample.(string))
			if ret != tc.expect.(Enum) {
				fail(tc, ret, tc.expect)
			}
		case "array":
			ret := es.GetAll()
			if len(ret) != len(tc.expect.([]string)) {
				fail(tc, ret, tc.expect)
				continue
			}
			for i := 0; i < len(ret); i++ {
				if ret[i] != (tc.expect.([]string)[i]) {
					fail(tc, ret, tc.expect)
					break
				}
			}
		}
	}
}

func fail(testingCase, ret, expect interface{}) {
	fmt.Printf("testing fail, case: %+v\nresult: %+v\nexpect: %+V\n", testingCase, ret, expect)
}
