package jsonpatch

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	OP_ADD     = "add"
	OP_REMOVE  = "remove"
	OP_REPLACE = "replace"
)

// it support add / replace / remove only
// don't support test / move / copy
type Patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func FindField(m interface{}, opt Patch) (f reflect.Value, err error) {
	refV := reflect.ValueOf(m)
	refT := reflect.TypeOf(m)
	if refT.Kind() == reflect.Ptr {
		refV = refV.Elem()
		refT = refT.Elem()
	}

	fields := strings.Split(opt.Path, "/")
	if 1 < len(fields) {
		fields = fields[1:]
	}

	find := false
	for _, fname := range fields {
		n := refT.NumField()
		for i := 0; i < n; i++ {
			if fname == refT.Field(i).Tag.Get("json") {
				f = refV.Field(i)
				refV = refV.Field(i)
				refT = refT.Field(i).Type
				find = true
				if refT.Kind() != reflect.Struct {
					goto end
				}

				break
			}
		}
	}

end:

	if find {
		return f, nil
	}

	return f, errors.New("field " + opt.Path + " not found")
}

func Do(f reflect.Value, opt Patch) (err error) {
	if f.Kind() == reflect.Ptr {
		f = f.Elem()
	}

	switch f.Kind() {
	case reflect.Struct:
		v := reflect.New(f.Type()).Interface()
		err = json.Unmarshal([]byte(opt.Value), v)
		if err == nil {
			f.Set(reflect.ValueOf(v).Elem())
		}
	case reflect.Slice, reflect.Array:
		v := reflect.New(f.Type()).Interface()
		err = json.Unmarshal([]byte(opt.Value), v)
		if err == nil {
			f.Set(reflect.ValueOf(v).Elem())
		}
	case reflect.String:
		f.SetString(opt.Value)
	case reflect.Bool:
		if boolean, er := strconv.ParseBool(opt.Value); er == nil {
			f.SetBool(boolean)
		} else {
			err = er
		}
	case reflect.Float32, reflect.Float64:
		if float, er := strconv.ParseFloat(opt.Value, 64); er == nil {
			f.SetFloat(float)
		} else {
			err = er
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, er := strconv.ParseInt(opt.Value, 10, 64); er == nil {
			f.SetInt(i)
		} else {
			err = er
		}
	}

	return err
}
