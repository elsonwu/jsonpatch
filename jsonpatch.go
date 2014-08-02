package jsonpatch

import (
	"encoding/json"
	"errors"
	"reflect"
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
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func FindField(m interface{}, opt Patch) (f reflect.Value, err error) {
	refV := reflect.ValueOf(m)
	refT := reflect.TypeOf(m)
	if refT.Kind() == reflect.Ptr {
		refV = refV.Elem()
		refT = refT.Elem()
	}

	fields := strings.Split(strings.Trim(opt.Path, "/"), "/")
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
		if opt.Op == OP_REPLACE || opt.Op == OP_ADD {
			if jsv, er := json.Marshal(opt.Value); er == nil {
				err = json.Unmarshal(jsv, v)
				if err == nil {
					f.Set(reflect.ValueOf(v).Elem())
				}

			} else {
				return er
			}
		} else if opt.Op == OP_REMOVE {
			f.Set(reflect.ValueOf(v).Elem())
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

	case reflect.Interface:
		v := reflect.New(f.Type()).Interface()
		if opt.Op == OP_REPLACE || opt.Op == OP_ADD {
			if jsv, er := json.Marshal(opt.Value); er == nil {
				err = json.Unmarshal(jsv, v)
				if err == nil {
					f.Set(reflect.ValueOf(v).Elem())
				}

			} else {
				return er
			}
		} else if opt.Op == OP_REMOVE {
			f.Set(reflect.ValueOf(v).Elem())
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

	case reflect.Map:
		v := reflect.New(f.Type()).Interface()
		if opt.Op == OP_REPLACE || opt.Op == OP_ADD {
			if jsv, er := json.Marshal(opt.Value); er == nil {
				err = json.Unmarshal(jsv, v)
				if err == nil {
					f.Set(reflect.ValueOf(v).Elem())
				}

			} else {
				return er
			}
		} else if opt.Op == OP_REMOVE {
			f.Set(reflect.ValueOf(v).Elem())
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

	case reflect.Slice, reflect.Array:
		v := reflect.New(f.Type()).Interface()
		if opt.Op == OP_REPLACE || opt.Op == OP_ADD {
			if jsv, er := json.Marshal(opt.Value); er == nil {
				err = json.Unmarshal(jsv, v)
				if err == nil {
					f.Set(reflect.ValueOf(v).Elem())
				}

			} else {
				return er
			}
		} else if opt.Op == OP_REMOVE {
			f.Set(reflect.ValueOf(v).Elem())
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

	case reflect.String:
		if opt.Op == OP_REPLACE || opt.Op == OP_ADD {
			if str, ok := opt.Value.(string); ok {
				f.SetString(str)
			} else {
				return errors.New(opt.Path + " must be string ")
			}
		} else if opt.Op == OP_REMOVE {
			f.SetString("")
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}
	case reflect.Bool:
		if opt.Op != OP_REPLACE && opt.Op != OP_ADD {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

		if boolean, ok := opt.Value.(bool); ok {
			f.SetBool(boolean)
		} else {
			return errors.New(opt.Path + " must be boolean")
		}
	case reflect.Float32, reflect.Float64:
		if opt.Op != OP_REPLACE && opt.Op != OP_ADD {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

		if float1, ok := opt.Value.(float64); ok {
			f.SetFloat(float1)
		} else if float2, ok := opt.Value.(float32); ok {
			f.SetFloat(float64(float2))
		} else if i, ok := opt.Value.(int); ok {
			f.SetFloat(float64(i))
		} else if i8, ok := opt.Value.(int8); ok {
			f.SetFloat(float64(i8))
		} else if i16, ok := opt.Value.(int16); ok {
			f.SetFloat(float64(i16))
		} else if i32, ok := opt.Value.(int32); ok {
			f.SetFloat(float64(i32))
		} else if i64, ok := opt.Value.(int64); ok {
			f.SetFloat(float64(i64))
		} else {
			return errors.New(opt.Path + " must be float")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if opt.Op != OP_REPLACE && opt.Op != OP_ADD {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

		if float, ok := opt.Value.(float64); ok {
			f.SetInt(int64(float))
		} else if float2, ok := opt.Value.(float32); ok {
			f.SetInt(int64(float2))
		} else if i, ok := opt.Value.(int); ok {
			f.SetInt(int64(i))
		} else if i8, ok := opt.Value.(int8); ok {
			f.SetInt(int64(i8))
		} else if i16, ok := opt.Value.(int16); ok {
			f.SetInt(int64(i16))
		} else if i32, ok := opt.Value.(int32); ok {
			f.SetInt(int64(i32))
		} else if i64, ok := opt.Value.(int64); ok {
			f.SetInt(i64)
		} else {
			return errors.New(opt.Path + " must be int")
		}
	}

	return err
}
