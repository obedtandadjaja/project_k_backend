package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type SendRequest struct {
	Sender        string   `json:"sender"`
	Recipients    []string `json:"recipients"`
	Subject       string   `json:"subject"`
	Body          string   `json:"body"`
	BodyHtml      string   `json:"bodyHtml"`
	BccRecipients []string `json:"bccRecipients"`
}

type SendResponse struct {
	Message string `json:"message"`
	EmailId string `json:"emailId"`
}

func parseSendRequest(r *http.Request) (*SendRequest, error) {
	var request SendRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func (server Server) Send(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := parseSendRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	logrus.Warn(request)

	message := server.mailgun.NewMessage(request.Sender, request.Subject, request.Body, request.Recipients...)
	if request.BodyHtml != "" {
		message.SetHtml(request.BodyHtml)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	msg, id, err := server.mailgun.Send(ctx, message)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	response := SendResponse{
		Message: msg,
		EmailId: id,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Warn(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
