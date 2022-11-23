package handlers

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/google/logger"
)

var view = jet.NewSet(
	jet.NewOSFileSystemLoader("./src/html"),
	jet.InDevelopmentMode(),
)

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) (err error) {
	t, err := view.GetTemplate(tmpl)
	if err != nil {
		logger.Error("---Error-> GetTemplate--->", err.Error())
	}

	err = t.Execute(w, data, nil)
	if err != nil {
		logger.Error("---Error-> Execute--->", err.Error())
	}

	return nil
}
