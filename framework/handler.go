package framework

import (
	"net/http"
)

func OptionsHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// later add domain check options
		WriteResponse(w, http.StatusOK, nil)
		return
	})
}
