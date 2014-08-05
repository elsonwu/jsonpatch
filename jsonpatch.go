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
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func findStructField(refV reflect.Value, fieldName string) (f reflect.Value, err error) {
	nf := refV.NumField()
	t := refV.Type()
	for i := 0; i < nf; i++ {
		f = refV.Field(i)
		ft := t.Field(i)
		if ft.Name == fieldName || ft.Tag.Get("json") == fieldName {
			return f, nil
		}
	}

	return f, errors.New("field " + fieldName + "does not exist")
}

func findMapField(refV reflect.Value, fieldName string) (f reflect.Value, err error) {
	if !refV.IsValid() {
		return refV, errors.New("value is invalid")
	}

	keys := refV.MapKeys()
	for i := 0; i < len(keys); i++ {
		if keys[i].Kind() == reflect.String {
			if str, ok := keys[i].Interface().(string); ok && str == fieldName {
				return refV.MapIndex(keys[i]), nil
			}
		}
	}

	return refV, errors.New("key " + fieldName + " does not exist")
}

func findField(refV reflect.Value, fields []string, offset int) (f reflect.Value, err error) {
	if len(fields) <= offset {
		return refV, nil
	}

	if refV.Kind() == reflect.Ptr {
		refV = refV.Elem()
	}

	if refV.Kind() == reflect.Struct {
		f, err = findStructField(refV, fields[offset])
	} else if refV.Kind() == reflect.Map {
		f, err = findMapField(refV, fields[offset])
	} else {
		return refV, errors.New("field " + fields[offset] + " must be map or struct")
	}

	if !f.IsValid() {
		return f, err
	}

	if err != nil {
		return refV, err
	}

	return findField(f, fields, offset+1)
}

func FindField(m interface{}, opt Patch) (f reflect.Value, err error) {
	refV := reflect.ValueOf(m)
	fields := strings.Split(strings.Trim(opt.Path, "/"), "/")
	if 0 == len(fields) {
		return f, errors.New("opt.Path is invalid")
	}

	return findField(refV, fields, 0)
}

func Do(f reflect.Value, opt Patch) (err error) {
	if f.Kind() == reflect.Ptr {
		f = f.Elem()
	}

	fields := strings.Split(strings.Trim(opt.Path, "/"), "/")
	fieldName := fields[len(fields)-1:][0]

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
		if f.IsNil() {
			f.Set(reflect.MakeMap(f.Type()))
		}

		if opt.Op == OP_ADD {
			// the value is a map / array
			if jsv, er := json.Marshal(opt.Value); er == nil {
				v := reflect.New(f.Type()).Interface()
				if err := json.Unmarshal(jsv, v); err == nil {
					subMap := reflect.ValueOf(v).Elem()
					for _, k := range subMap.MapKeys() {
						f.SetMapIndex(k, subMap.MapIndex(k))
					}

					return nil
				}
			} else {
				err = er
			}
		} else if opt.Op == OP_REPLACE {
			// the value is a map / array
			if jsv, er := json.Marshal(opt.Value); er == nil {
				v := reflect.New(f.Type()).Interface()
				if err := json.Unmarshal(jsv, v); err == nil {
					f.Set(reflect.ValueOf(v).Elem())
				} else {
					return err
				}
			} else {
				err = er
			}
		} else if opt.Op == OP_REMOVE {
			f.Set(reflect.New(f.Type()).Elem())
			return nil
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

	case reflect.Slice, reflect.Array:
		if opt.Op == OP_REPLACE {
			//if the fieldName is digit, then we think he/she want to replace the value by index
			if i, err := strconv.Atoi(fieldName); err == nil {
				if i < f.Len() {
					if jsv, er := json.Marshal(opt.Value); er == nil {
						v := reflect.New(f.Type().Elem()).Interface()
						if json.Unmarshal(jsv, v) == nil {
							f.Index(i).Set(reflect.ValueOf(v).Elem())
							return nil
						}
					} else {
						return er
					}
				} else {
					return errors.New("index " + fieldName + " is out of range")
				}
			} else {
				//if the fieldName is string only, then we think he/she want to replace the whole array/slice
				if jsv, er := json.Marshal(opt.Value); er == nil {
					v := reflect.New(f.Type()).Interface()
					if json.Unmarshal(jsv, v) == nil {
						f.Set(reflect.ValueOf(v).Elem())
						return nil
					}
				} else {
					return er
				}
			}
		} else if opt.Op == OP_ADD {
			if jsv, er := json.Marshal(opt.Value); er == nil {
				v := reflect.New(f.Type().Elem()).Interface()
				if json.Unmarshal(jsv, v) == nil {
					f.Set(reflect.Append(f, reflect.ValueOf(v).Elem()))
					return nil
				}
			} else {
				return er
			}
		} else if opt.Op == OP_REMOVE {
			// remove from the index, such as /courses/1
			if i, err := strconv.Atoi(fieldName); err == nil {
				if f.Len() > i {
					newF := reflect.AppendSlice(f.Slice(0, i), f.Slice(i+1, f.Len()))
					f.Set(newF)
					return nil
				}
			} else {
				// remove the whole array, such as /courses
				v := reflect.New(f.Type()).Interface()
				f.Set(reflect.ValueOf(v).Elem())
				return nil
			}
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

	case reflect.String:
		if opt.Op == OP_REPLACE || opt.Op == OP_ADD {
			if str, ok := opt.Value.(string); ok {
				f.SetString(str)
				return nil
			} else {
				return errors.New(opt.Path + " must be string ")
			}
		} else if opt.Op == OP_REMOVE {
			f.SetString("")
			return nil
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

	case reflect.Bool:
		if opt.Op != OP_REPLACE && opt.Op != OP_ADD {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

		if boolean, ok := opt.Value.(bool); ok {
			f.SetBool(boolean)
			return nil
		} else {
			return errors.New(opt.Path + " must be boolean")
		}

	case reflect.Float32, reflect.Float64:
		if opt.Op == OP_REPLACE || opt.Op == OP_ADD {
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

			return nil
		} else if opt.Op == OP_REMOVE {
			f.Set(reflect.New(f.Type()).Elem())
			return nil
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if opt.Op == OP_REPLACE || opt.Op == OP_ADD {
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

			return nil
		} else if opt.Op == OP_REMOVE {
			f.Set(reflect.New(f.Type()).Elem())
			return nil
		} else {
			return errors.New(opt.Path + " unsupport op " + opt.Op)
		}
	}

	return
}
