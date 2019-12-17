package api

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (server Server) Health(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}
