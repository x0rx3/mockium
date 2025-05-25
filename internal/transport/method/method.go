package method

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Method string

func (inst *Method) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	method := Method(strings.ToUpper(s))
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
		*inst = method
		return nil
	case "":
		return fmt.Errorf("unxpected method")
	default:
		return fmt.Errorf("unxpected method")
	}
}
