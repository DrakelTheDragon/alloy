package response

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSONify writes a JSON response to the io.Writer or http.ResponseWriter.
func JSONify(w io.Writer, r Response) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	if rw, ok := w.(http.ResponseWriter); ok {
		writeHeaderJSON(rw, r)
	}

	_, err = w.Write(data)
	return err
}

func writeHeaderJSON(w http.ResponseWriter, r Response) {
	for k, v := range r.Header() {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code())
}
