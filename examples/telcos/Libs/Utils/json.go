package Utils

import (
	"bytes"
	"encoding/json"
)

// PrettyJSON parse object to json string and make it pretty
func PrettyJSON(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(true)

	// try to encode json
	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
