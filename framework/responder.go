package framework

import (
	"github.com/rainingclouds/lemonades/logger"
	"html/template"
	"net/http"
	"net/http/httputil"
	"strings"
)

var templates *template.Template

func SetTemplate(basePath string) error {
	var err error
	templates, err = loadTemplates(basePath)
	return err
}

func WriteError(w http.ResponseWriter, r *http.Request, c int, err error) {
	requstDump, _ := httputil.DumpRequest(r, true)
	logger.Get().Warning(err, "\n", strings.Trim(string(requstDump), "\n\r"))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(c)
	res := JSONResponse{"message": err.Error(), "success": false}
	w.Write(res.ByteArray())
}

func WriteHtml(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteResponse(w http.ResponseWriter, c int, r JSONResponse) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(c)
	w.Write(r.ByteArray())
}
