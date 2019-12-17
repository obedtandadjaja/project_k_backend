package controller

import (
	"io"
	"net/http"
)

func Health(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	io.WriteString(w, "OK")
	return nil
}
