package jsonpatch

import "encoding/json"

func ToPatchs(jsonPos []byte) ([]Patch, error) {
	ops := []Patch{}
	if err := json.Unmarshal([]byte(jsonPos), &ops); err != nil {
		return nil, err
	}

	return ops, nil
}

func Run(jsonPos []byte, model interface{}) error {
	ops, err := ToPatchs(jsonPos)
	if err != nil {
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
