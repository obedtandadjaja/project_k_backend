package controller

import (
	"database/sql"
	"log"
	"net/http"
)

type HandlerError struct {
	Code    int
	Message string
	Err     error
}

func (error HandlerError) Error() string {
	return error.Message
}

type SharedResources struct {
	DB  *sql.DB
	Env string
}

type Handler struct {
	SharedResources *SharedResources
	Handler         func(sr *SharedResources, w http.ResponseWriter, r *http.Request) error
}

// to satisfy http.Handler
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.Handler(h.SharedResources, w, r)

	if err != nil {
		switch e := err.(type) {
		case HandlerError:
			if error := e.Err; error != nil {
				log.Printf("ERROR: %s\n", error)
			}

			log.Printf("HTTP %d - %s\n", e.Code, e)

			// on prod error codes >= 500 should not be returned
			if h.SharedResources.Env == "production" && e.Code >= 500 {
				http.Error(w, "Internal Server Error", e.Code)
			} else if h.SharedResources.Env != "production" && e.Code >= 500 {
				http.Error(w, e.Err.Error(), e.Code)
			} else {
				http.Error(w, e.Message, e.Code)
			}
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
