package util

import (
	"encoding/json"
	"fmt"
	"io"
)

func DecodeToStruct(r io.Reader, result interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(result)
	if err != nil {
		return fmt.Errorf("action=util.decodeToStruct err=%v", err)
	}
	return nil
}
