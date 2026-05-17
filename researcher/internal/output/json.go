package output

import (
	"encoding/json"
	"io"
)

func WriteJSON(w io.Writer, v any, pretty bool) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(v)
}
