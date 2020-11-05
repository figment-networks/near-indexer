package near

import (
	"encoding/base64"
	"encoding/json"
)

func argsToBase64(input interface{}) (string, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
