package model

import (
	"encoding/json"
)

const COOKIE = "Cookie"

type MustHeader struct {
	Cookie []Cookie
	Other  map[string]any
}

type MustHeaderTemplate struct {
	Cookie []CookieTemplate `json:"Cookie"`
	Other  map[string]any   `json:"-"`
}

func (inst *MustHeaderTemplate) UnmarshalJSON(data []byte) error {
	type Alias MustHeaderTemplate
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(inst),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	delete(rawMap, COOKIE)

	inst.Other = make(map[string]any)
	for key, val := range rawMap {
		var v any
		if err := json.Unmarshal(val, &v); err != nil {
			return err
		}
		inst.Other[key] = v
	}

	return nil
}
