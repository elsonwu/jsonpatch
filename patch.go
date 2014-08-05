package jsonpatch

import "encoding/json"

func Run(jsonPos string, model interface{}) error {
	ops := []Patch{}
	if err := json.Unmarshal([]byte(jsonPos), &ops); err != nil {
		return err
	}

	for _, opt := range ops {
		f, err := FindField(model, opt)
		if !f.IsValid() && err != nil {
			return err
		}

		if err := Do(f, opt); err != nil {
			return err
		}
	}

	return nil
}
