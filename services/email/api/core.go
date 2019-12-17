package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/sirupsen/logrus"
)

type Server struct {
	mailgun *mailgun.MailgunImpl
	router  *httprouter.Router
}

func Start(AppUrl string, mailgun *mailgun.MailgunImpl) {
	server := Server{
		mailgun: mailgun,
		router:  httprouter.New(),
	}
	server.router.POST("/api/v1/send", server.Send)
	server.router.GET("/api/health", server.Health)

	logrus.Info("App running on " + AppUrl)
	logrus.Fatal(http.ListenAndServe(AppUrl, server.router))
}
