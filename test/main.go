package main

import (
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
	Map       map[string]string `json:"map"`
	StructMap struct {
		Map map[string]string `json:"map"`
	} `json:"struct_map"`
	Num   int         `json:"num"`
	Num32 int32       `json:"num32"`
	Num64 int64       `json:"num64"`
	F32   float32     `json:"f32"`
	F64   float64     `json:"f64"`
	Bool  bool        `json:"bool"`
	Inter interface{} `json:"inter"`
}

var jsonOps = `[
    {"op":"add", "path":"/Foo", "value":"xxx"},
    {"op":"add", "path":"/Inter", "value":{"K":123, "k2":"value..."}},
    {"op":"replace", "path":"/User", "value":{"name":"elsonwu", "fullname":"elson wu"}},
    {"op":"replace", "path":"/Courses", "value":[{"cid":"0000"}]},
    {"op":"add", "path":"/Courses", "value":{"cid":"1111"}},
    {"op":"add", "path":"/Courses", "value":{"cid":"2222"}},
    {"op":"add", "path":"/Courses", "value":{"cid":"3333"}},
    {"op":"add", "path":"/Courses", "value":{"cid":"4444"}},
    {"op":"add", "path":"/Courses", "value":{"cid":"5555"}},
    {"op":"remove", "path":"/Courses/4"},
    {"op":"remove", "path":"/Courses/4"},
    {"op":"add", "path":"/People/Work/Place", "value":"Guangzhou"},
    {"op":"add", "path":"/Num", "value":123},
    {"op":"add", "path":"/Num32", "value":123123},
    {"op":"add", "path":"/Num64", "value":123321},
    {"op":"remove", "path":"/Num64", "value":123321},
    {"op":"add", "path":"/F32", "value":123},
    {"op":"add", "path":"/F64", "value":123},
    {"op":"add", "path":"/Bool", "value":true},
    {"op":"add", "path":"/map", "value":{"xxx":"yyy"}},
    {"op":"remove", "path":"/map"},
    {"op":"add", "path":"/map", "value":{"kkk":"jjj"}},
    {"op":"replace", "path":"/map", "value":{"mmm":"nnn"}},
    {"op":"add", "path":"/map", "value":{"aaa":"123"}},
    {"op":"add", "path":"/StructMap/map", "value":{"age":"27"}},
    {"op":"add", "path":"/StructMap/map", "value":{"cn":"ew10"}}
]`

func main() {
	dd := new(d)
	err := jsonpatch.Run(jsonOps, dd)
	fmt.Printf("%#v \n\n[error] %#v \n", dd, err)
	return
}
