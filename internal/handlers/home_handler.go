package handlers

import (
	"net/http"

	"github.com/google/logger"
)

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "index.jet", nil)
	if err != nil {
		logger.Error("---Error-> renderPage--->", err.Error())
		return
	}
}
