package framework

import (
	"net/http"
)

// this struct basically adds a context to the http.Request so that
// authenticator or any other middleward could push out the data
// to main request handler
type Request struct {
	*http.Request
	context map[string]interface{}
}

func (r *Request) Push(key string, value interface{}) {
	if r.context == nil {
		r.context = map[string]interface{}{}
	}
	r.context[key] = value
}

func (r *Request) MustGet(key string) interface{} {
	return r.context[key]
}
