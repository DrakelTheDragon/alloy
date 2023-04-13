package response

import "net/http"

// JSONify writes a JSON response to the http.ResponseWriter.
func JSONify(w http.ResponseWriter, r *Response) error {
	data, err := r.MarshalJSON()
	if err != nil {
		return err
	}

	data = append(data, '\n')

	if r.header != nil {
		for k, v := range r.header {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code())

	_, err = w.Write(data)
	return err
}
