package main

import (
	"net/http"
	"path"
	"strings"
)

// ImportsHandler handles HTTP request that generates go-import template.
type ImportsHandler struct {
	Generator *Generator
	Logger    Logger
}

// ServeHTTP implements http.Handler.
func (h ImportsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(strings.TrimPrefix(req.URL.Path, "/"), "_status") {
		// Internal endpoints.
		base := path.Base(strings.TrimSuffix(req.URL.Path, "/"))
		if base == "healthz" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Logger.Infof("Serve HTTP request: %s", req.RequestURI)

	if err := h.Generator.Generate(w, h.getHost(req), req.URL.Path); err != nil {
		if h.Logger != nil {
			h.Logger.Errorf("%v", err)
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h ImportsHandler) getHost(req *http.Request) string {
	host := req.Header.Get("X-Forwarded-Host")
	if index := strings.IndexByte(host, ','); index >= 0 {
		host = host[0:index]
	}
	host = strings.TrimSpace(host)
	if len(host) > 0 {
		return host
	}

	return req.Host
}
