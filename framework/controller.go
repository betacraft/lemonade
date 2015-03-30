package framework

import (
	"encoding/json"
	"errors"
	"github.com/mholt/binding"
	"github.com/rainingclouds/lemonade/logger"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

// a function to read thhe JSON body and return a map of string and interface
func ReadBody(r *http.Request) (map[string]interface{}, error) {
	bodyMap := make(map[string]interface{})
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return bodyMap, err
	}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return bodyMap, err
	}
	return bodyMap, nil
}

// a function to read json body and return an interface
func ReadJSONBody(r *http.Request) (interface{}, error) {
	var response interface{}
	err := json.NewDecoder(r.Body).Decode(response)
	return response, err
}

func Bind(r *http.Request, fm binding.FieldMapper) error {
	err := binding.Bind(r, fm)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func WriteError(w http.ResponseWriter, r *http.Request, c int, err error) {
	requstDump, _ := httputil.DumpRequest(r, true)
	logger.Get().Warning(err, "\n", strings.Trim(string(requstDump), "\n\r"))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(c)
	res := JSONResponse{"message": err.Error(), "success": false}
	w.Write(res.ByteArray())
}

func WriteResponse(w http.ResponseWriter, c int, r JSONResponse) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(c)
	w.Write(r.ByteArray())
}
