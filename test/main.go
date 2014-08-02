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
	Num   int     `json:"num"`
	Num32 int32   `json:"num32"`
	Num64 int64   `json:"num64"`
	F32   float32 `json:"f32"`
	F64   float64 `json:"f64"`
	Bool  bool    `json:"bool"`
}

func main() {
	dd := new(d)
	ops := []jsonpatch.Patch{jsonpatch.Patch{
		Op:    "replace",
		Path:  "/foo",
		Value: `xxx`,
	}, jsonpatch.Patch{
		Op:    "replace",
		Path:  "/user",
		Value: `{"name":"elsonwu"}`,
	}, jsonpatch.Patch{
		Op:    "replace",
		Path:  "/courses/-",
		Value: `[{"cid":"123321"}]`,
	}, jsonpatch.Patch{
		Op:    "add",
		Path:  "/people/work/place",
		Value: `China/guangzhou`,
	}, jsonpatch.Patch{
		Op:    "add",
		Path:  "/people/work/place",
		Value: `-----`,
	}, jsonpatch.Patch{
		Op:    "add",
		Path:  "num",
		Value: `123`,
	}, jsonpatch.Patch{
		Op:    "add",
		Path:  "num32",
		Value: `321`,
	}, jsonpatch.Patch{
		Op:    "add",
		Path:  "num64",
		Value: `321123`,
	}, jsonpatch.Patch{
		Op:    "add",
		Path:  "f32",
		Value: `321123`,
	}, jsonpatch.Patch{
		Op:    "add",
		Path:  "f64",
		Value: `321123`,
	}, jsonpatch.Patch{
		Op:    "add",
		Path:  "bool",
		Value: `1`,
	}}

	for _, opt := range ops {
		if f, err := jsonpatch.FindField(dd, opt); err == nil {
			if e := jsonpatch.Do(f, opt); e != nil {
				fmt.Printf("[ERROR do] -----> %#v\n", e)
			}
		} else {
			fmt.Printf("[ERROR field] -----> %#v\n", err)
		}
	}

	fmt.Printf("%#v \n", dd)
	return
}
