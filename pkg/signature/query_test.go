package signature

import (
	"encoding/json"
	"fmt"
	"kms/internal/polysign"
	"testing"
)

func TestToQuery(t *testing.T) {
	var testCase = map[string]interface{}{
		"a": "f",
		"b": []string{"foo", "bar"},
		"c": map[string]interface{}{
			"c1": 1,
			"c2": map[string]interface{}{
				"c21": "foo",
				"c22": []interface{}{
					[]interface{}{1, "foo"},
					map[string]interface{}{
						"c221": "foo",
						"c222": true,
					},
				},
			},
		},
		"d": []interface{}{
			map[string]interface{}{
				"c21": "foo",
				"c22": []interface{}{
					[]interface{}{1, "foo"},
					map[string]interface{}{
						"c221": "foo",
						"c222": true,
					},
				},
			},
			[]interface{}{
				[]interface{}{1, "foo"},
				map[string]interface{}{
					"c221": []string{"foo", "bar"},
					"c222": true,
				},
			},
		},
	}

	b, err := json.MarshalIndent(testCase, "", "  ")
	if err != nil {
		panic(err)
	}

	var d interface{}
	if err := json.Unmarshal(b, &d); err != nil {
		panic(err)
	}

	query, err := ToQuery(d)
	expect := `a=f&b.1=foo&b.2=bar&c.c1=1&c.c2.c21=foo&c.c2.c22.1.1=1&c.c2.c22.1.2=foo&c.c2.c22.2.c221=foo&c.c2.c22.2.c222=true&d.1.c21=foo&d.1.c22.1.1=1&d.1.c22.1.2=foo&d.1.c22.2.c221=foo&d.1.c22.2.c222=true&d.2.1.1=1&d.2.1.2=foo&d.2.2.c221.1=foo&d.2.2.c221.2=bar&d.2.2.c222=true`
	if expect != query || err != nil {
		fmt.Println(string(b))
		t.Errorf("TestToQuery:\nexpect %s\ngot    %s err=%v\n", expect, query, err)
	}
}

func TestToQueryExt(t *testing.T) {
	type testCase struct {
		body      interface{}
		expect    string
		expectErr bool
	}
	testCases := []*testCase{
		&testCase{
			body: map[string]interface{}{
				polysign.XHeaderPolySignKeyID:  "foo",
				polysign.XPolyRaiseUpFieldName: "stringBody",
			},
			expect: "$body$=stringBody&X-Polysign-Access-Key-Id=foo",
		},
		&testCase{
			body: map[string]interface{}{
				polysign.XHeaderPolySignKeyID:  "foo",
				polysign.XPolyRaiseUpFieldName: []string{"foo", "bar"},
			},
			expect: "$body$.1=foo&$body$.2=bar&X-Polysign-Access-Key-Id=foo",
		},
		&testCase{
			body: map[string]interface{}{
				polysign.XHeaderPolySignKeyID: "foo",
				polysign.XPolyRaiseUpFieldName: map[string]interface{}{
					"a": "foo",
					polysign.XPolyBodyHideArgs: map[string]interface{}{
						"app": "foo",
					},
					polysign.XPolyCustomerBodyRoot: "bar",
				},
			},
			expect: "$body$=bar&$polyapi_hide$.app=foo&X-Polysign-Access-Key-Id=foo&a=foo",
		},
		&testCase{
			body: map[string]interface{}{
				polysign.XHeaderPolySignKeyID: "foo",
				polysign.XPolyRaiseUpFieldName: map[string]interface{}{
					"a":                           "foo",
					polysign.XHeaderPolySignKeyID: "bar",
				},
			},
			expect:    "",
			expectErr: true,
		},
	}
	for i, v := range testCases {
		got, err := ToQuery(v.body)
		if got != v.expect || (err != nil) != v.expectErr {
			t.Errorf("case %d, expect %q\ngot %q err=%v", i+1, v.expect, got, err)
		}
	}
}
