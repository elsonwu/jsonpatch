package main

import (
	"encoding/json"
	"fmt"

	"github.com/elsonwu/jsonpatch"
)

var js = `{ "foo": "bar" }`

type d struct {
	Foo  string `json:"foo"`
	User struct {
		Name string `json:"name"`
	} `json:"user"`
	Courses []struct {
		Cid string `json:"cid"`
	} `json:"courses"`
	People struct {
		Work struct {
			Place string `json:"place"`
		} `json:"work"`
	} `json:"people"`
	Map   map[string]string `json:"map"`
	Num   int               `json:"num"`
	Num32 int32             `json:"num32"`
	Num64 int64             `json:"num64"`
	F32   float32           `json:"f32"`
	F64   float64           `json:"f64"`
	Bool  bool              `json:"bool"`
	Inter interface{}       `json:"inter"`
}

// var jsonOps = `[{"op":"add","path":"/foo","value":"xxx"}]`
var jsonOps = `[
    {"op":"add", "path":"/foo", "value":"xxx"},
    {"op":"add", "path":"/inter", "value":{"k":123, "k2":"value..."}},
    {"op":"replace", "path":"/user", "value":{"name":"elsonwu", "fullname":"elson wu"}},
    {"op":"replace", "path":"/courses/-", "value":[{"cid":"001234"}]},
    {"op":"add", "path":"/people/work/place", "value":"Guangzhou"},
    {"op":"add", "path":"/num", "value":123},
    {"op":"add", "path":"/num32", "value":123123},
    {"op":"add", "path":"/num64", "value":123321},
    {"op":"add", "path":"/f32", "value":123},
    {"op":"add", "path":"/f64", "value":123},
    {"op":"add", "path":"/bool", "value":true}
]`

func main() {
	dd := new(d)
	ops := []jsonpatch.Patch{}
	if err := json.Unmarshal([]byte(jsonOps), &ops); err != nil {
		fmt.Println(err)
		return
	}

	for _, opt := range ops {
		if f, err := jsonpatch.FindField(dd, opt); err == nil {
			if e := jsonpatch.Do(f, opt); e != nil {
				fmt.Printf("[ERROR Do] -----> %#v\n", e)
			}
		} else {
			fmt.Printf("[ERROR FindField] -----> %#v\n", err)
		}
	}

	fmt.Printf("%#v \n", dd)
	return
}
