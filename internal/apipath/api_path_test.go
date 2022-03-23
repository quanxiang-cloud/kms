package apipath

import (
	"fmt"
	"testing"
)

func TestSplitAPIPath(t *testing.T) {
	type testingCase struct {
		sample interface{}
		expect interface{}
	}

	splitCases := []testingCase{
		{
			sample: "testing/",
			expect: []string{
				"/testing",
				"",
			},
		},
		{
			sample: "testing",
			expect: []string{
				"",
				"testing",
			},
		},
		{
			sample: "/testing",
			expect: []string{
				"",
				"tesing",
			},
		},
		{
			sample: "/testing/case/a/",
			expect: []string{
				"/testing/case/a",
				"",
			},
		},
	}

	for _, sc := range splitCases {
		path, name := Split(sc.sample.(string))
		if path != sc.expect.([]string)[0] && name != sc.expect.([]string)[1] {
			fmt.Printf("tesing fail, case : %+v\nresult: {path: %s, name: %s}\nexpect: {path: %s, name: %s}\n", sc, path, name, sc.expect.([]string)[0], sc.expect.([]string)[1])
		}
	}
}

func TestJoinAPIPath(t *testing.T) {
	type testingCase struct {
		sample interface{}
		expect interface{}
	}

	joinCases := []testingCase{
		{
			sample: []string{
				"testing",
				"",
			},
			expect: "/testing",
		},
		{
			sample: []string{
				"testing",
				"case",
			},
			expect: "/testing/case",
		},
		{
			sample: []string{
				"testing",
				"case/",
			},
			expect: "/testing/case",
		},
		{
			sample: []string{
				"testing/",
				"/case",
			},
			expect: "/testing/case",
		},
	}

	for _, jc := range joinCases {
		c := jc.sample.([]string)
		ret := Join(c[0], c[1])
		if ret != jc.expect.(string) {
			fmt.Printf("testing fail, case: %+v\nresult: %s\nexpect: %s\n", jc, ret, jc.expect)
		}
	}
}
