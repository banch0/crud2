package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// ErrInvalidContentType ...
var ErrInvalidContentType = errors.New("Content-Type must be aplication/json")

// ErrReadRequestBody ...
var ErrReadRequestBody = errors.New("read request body error")

// ErrUnmarshalBody ...
var ErrUnmarshalBody = errors.New("unmarshal body to interface error")

// ErrWriteResponse ...
var ErrWriteResponse = errors.New("write reponse error")

// ErrMarshal ...
var ErrMarshal = errors.New("marshal struct error")

//ReadJSONBody ...
func ReadJSONBody(request *http.Request, dto interface{}) (err error) {
	if request.Header.Get("Content-Type") != "application/json" {
		return ErrInvalidContentType
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("requeset body ReadAll error: ", err)
		return ErrReadRequestBody
	}
	defer request.Body.Close()

	err = json.Unmarshal(body, &dto)
	if err != nil {
		log.Println("Unamrshal error: ", err)
		return ErrUnmarshalBody
	}
	return nil
}

// WriteJSONBody ...
func WriteJSONBody(response http.ResponseWriter, dto interface{}) (err error) {
	response.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(dto)
	if err != nil {
		return ErrMarshal
	}

	_, err = response.Write(body)
	if err != nil {
		return ErrWriteResponse
	}

	return nil
}
