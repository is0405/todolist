package httputil

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	
	response, err := marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			log.Print(writeErr)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, writeErr := w.Write(response)
	if writeErr != nil {
		log.Print(writeErr)
	}
}

func marshal(payload interface{}) ([]byte, error) {
	if isNil(payload) {
		return []byte(`{}`), nil
	}
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(payload)
	if err != nil {
		return []byte(`{}`), nil
	}
	return buffer.Bytes(), err
}

func isNil(p interface{}) bool {
	if p == nil {
		return true
	}
	switch reflect.TypeOf(p).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Array:
		return reflect.ValueOf(p).IsNil()
	}
	return false
}

func RespondErrorJson(w http.ResponseWriter, code int, err error) {
	log.Printf("code=%d, err=%s", code, err)
	
	if e, ok := err.(*HTTPError); ok {
		RespondJSON(w, code, e)
	} else if err != nil {
		he := HTTPError{
			Message: err.Error(),
		}
		RespondJSON(w, code, he)
	}

}
