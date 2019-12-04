package hsutil

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(i interface{}) (err error) {
	b, err := json.MarshalIndent(i, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
