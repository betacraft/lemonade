package framework

import (
	"fmt"
	"github.com/rainingclouds/lemonades/logger"
	"html/template"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

var templates *template.Template
var isProd = os.Getenv("ENV") == "prod"

func SetTemplate(basePath string) error {
	var err error
	templates, err = loadTemplates(basePath)
	return err
}

func WriteError(w http.ResponseWriter, r *http.Request, c int, err error) {
	requstDump, _ := httputil.DumpRequest(r, true)
	logger.Get().Warning(err, "\n", strings.Trim(string(requstDump), "\n\r"))
	w.Header().Add("Content-Type", "application/json")
	if isProd {
		w.Header().Add("Access-Control-Allow-Origin", "http://www.lemonades.in")
	} else {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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

func WriteText(w http.ResponseWriter, text string) {
	w.Header().Add("Content-Type", "text/html")
	if isProd {
		w.Header().Add("Access-Control-Allow-Origin", "http://www.lemonades.in")
	} else {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Session-Key")
	w.WriteHeader(200)
	fmt.Fprint(w, text)
}

func WriteResponse(w http.ResponseWriter, c int, r JSONResponse) {
	w.Header().Add("Content-Type", "application/json")
	if isProd {
		w.Header().Add("Access-Control-Allow-Origin", "http://www.lemonades.in")
	} else {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Session-Key")
	w.WriteHeader(c)
	w.Write(r.ByteArray())
}
