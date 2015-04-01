package framework

import (
	"github.com/rainingclouds/lemonade/logger"
	"html/template"
	"os"
	"path/filepath"
)

func loadTemplates(basePath string) (*template.Template, error) {
	var templates *template.Template
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// don't process folders themselves
		if info.IsDir() {
			return nil
		}
		templateName := path[len(basePath):]
		if templates == nil {
			templates = template.New(templateName)
			templates.Delims("[[[", "]]]")
			_, err = templates.ParseFiles(path)
		} else {
			_, err = templates.New(templateName).ParseFiles(path)
		}
		logger.Debug("Processed template", templateName)
		return err
	})
	return templates, err
}
