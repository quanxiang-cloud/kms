package sort

import (
	"fmt"
	"net/url"
	"testing"
)

func TestGet(t *testing.T) {
	type Interest struct {
		Z    int
		Time int64
		A    string
		Name string
	}
	mockStruct := struct {
		Name     string `security:"name"`
		age      int
		Gender   string                 `security:"gender"`
		Email    string                 `security:"email"`
		Address  map[string]interface{} `security:"address"`
		Interest []Interest             `security:"interest" xml:"interest"`
	}{
		Name:   "张三",
		age:    18,
		Gender: "男",
		Email:  "11@gmail.com",
		Address: map[string]interface{}{
			"Contry":  "china",
			"Address": "gaoxin",
			"City":    "chengdu",
		},
		Interest: []Interest{{
			Z:    1,
			Time: 10,
			Name: "alex",
			A:    "a",
		}, {
			Time: 12,
			Name: "alex",
			A:    "a",
			Z:    2,
		},
		}}

	dst, err := WordSortDESC(mockStruct, XMLFormat, Gonic)
	if err != nil {
		t.Fatal(err)
	}

	mockMap := map[string]interface{}{
		"name":   "张三",
		"age":    18,
		"gender": "男",
		"email":  "11@gmail.com",
		"address": map[string]interface{}{
			"contry": "china",
		},
	}

	dst, err = WordSortASC(mockMap, XMLFormat, Gonic)
	if err != nil {
		t.Fatal(err)
	}

	dst, err = WordSortASC(`{"email":"11@gmail.com","address":{"address":"gaoxin","city":"chengdu","contry":"china"},"gender":"男","interest":[{"Time":10,"A":"a","Name":"alex"},{"Time":10,"A":"a","Name":"alex"}],"name":"张三"}`,
		XMLFormat, Gonic)
	if err != nil {
		t.Fatal(err)
	}

	_ = dst

}

func TestSA(t *testing.T) {
	unstr := "=====_:"
	fmt.Printf("url.QueryEscape:%s\n", url.QueryEscape(unstr))

	fmt.Println(url.PathEscape(unstr))
	fmt.Println(url.QueryUnescape("%2F"))
}

func TestRecursion(t *testing.T) {
	data := `{"name":"alex","address":{"contry":"china","city":"chengdu"},"inter":["pingpong","basketball","football"]}`

	dst, err := WordSortASC(data,
		JSONFormat, Gonic)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(dst))
}
