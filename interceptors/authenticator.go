package interceptors

import (
	"errors"
	"github.com/rainingclouds/lemonade/framework"
	"github.com/rainingclouds/lemonade/models"
	"net/http"
)

func UserAuthenticate(handler func(http.ResponseWriter, *framework.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authKey := r.Header.Get("Auth-Key")
		if authKey == "" {
			framework.WriteError(w, r, http.StatusUnauthorized, errors.New("Auth key is not present"))
			return
		}
		user, err := models.GetUserByAuthKey(authKey)
		if err != nil {
			framework.WriteError(w, r, http.StatusUnauthorized, err)
			return
		}
		req := framework.Request{Request: r}
		req.Push("user", user)
		handler(w, &req)
	})
}

func AdminAuthenticate(handler func(http.ResponseWriter, *framework.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authKey := r.Header.Get("Auth-Key")
		if authKey == "" {
			framework.WriteError(w, r, http.StatusUnauthorized, errors.New("Auth key is not present"))
			return
		}
		admin, err := models.GetAdminByAuthKey(authKey)
		if err != nil {
			framework.WriteError(w, r, http.StatusUnauthorized, err)
			return
		}
		req := framework.Request{Request: r}
		req.Push("admin", admin)
		handler(w, &req)
	})
}
