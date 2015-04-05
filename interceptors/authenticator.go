package interceptors

import (
	"errors"
	"github.com/rainingclouds/lemonade/framework"
	"github.com/rainingclouds/lemonade/logger"
	"github.com/rainingclouds/lemonade/models"
	"net/http"
)

func UserAuthenticate(handler func(http.ResponseWriter, *framework.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionKey, err := r.Cookie("lemonades_session_key")
		if err != nil {
			framework.WriteError(w, r, http.StatusUnauthorized, errors.New("Illegal request"))
			return
		}
		if sessionKey.Value == "" {
			framework.WriteError(w, r, http.StatusUnauthorized, errors.New("Illegal request"))
			return
		}
		logger.Debug("Session key is", sessionKey.Value)
		user, err := models.GetUserBySessionKey(sessionKey.Value)
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
