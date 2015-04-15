package interceptors

import (
	"errors"
	"github.com/rainingclouds/lemonades/framework"
	"github.com/rainingclouds/lemonades/logger"
	"github.com/rainingclouds/lemonades/models"
	"net/http"
)

func UserAuthenticate(handler func(http.ResponseWriter, *framework.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionKey := r.Header.Get("Session-Key")
		if sessionKey == "" {
			framework.WriteError(w, r, http.StatusUnauthorized, errors.New("Illegal request"))
			return
		}
		user, err := models.GetUserBySessionKey(sessionKey)
		if err != nil {
			logger.Debug("While finding user")
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
		sessionKey := r.Header.Get("Auth-Key")
		if sessionKey == "" {
			framework.WriteError(w, r, http.StatusUnauthorized, errors.New("Auth key is not present"))
			return
		}
		admin, err := models.GetAdminByAuthKey(sessionKey)
		if err != nil {
			framework.WriteError(w, r, http.StatusUnauthorized, err)
			return
		}
		req := framework.Request{Request: r}
		req.Push("admin", admin)
		handler(w, &req)
	})
}
