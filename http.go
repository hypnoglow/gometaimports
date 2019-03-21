package main

import (
	"net/http"
)

// ImportsHandler handles HTTP request that generates go-import template.
type ImportsHandler struct {
	Generator *Generator
	Logger    Logger
}

// ServeHTTP implements http.Handler.
func (h ImportsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.Logger.Infof("Serve HTTP request: %s", req.RequestURI)

	if err := h.Generator.Generate(w, req.URL.Path); err != nil {
		if h.Logger != nil {
			h.Logger.Errorf("%v", err)
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
