package credentials

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/obedtandadjaja/project_k_backend/services/auth/controller"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models/credential"
)

type UpdateRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func Update(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseUpdateRequest(r)
	if err != nil {
		return controller.HandlerError{400, err.Error(), err}
	}

	err = processUpdateRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func parseUpdateRequest(r *http.Request) (*UpdateRequest, error) {
	var request UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processUpdateRequest(sr *controller.SharedResources, request *UpdateRequest, r *http.Request) error {
	vars := mux.Vars(r)

	cred, err := credential.FindBy(sr.DB, map[string]interface{}{
		"uuid": vars["uuid"],
	})
	if err != nil {
		return controller.HandlerError{404, "Credential not found", err}
	}

	err = cred.Update(sr.DB, map[string]interface{}{
		"email": request.Email,
		"phone": request.Phone,
	})
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return controller.HandlerError{
					400,
					"Email has been taken",
					err,
				}
			}
		}

		return controller.HandlerError{
			500,
			"Failed to update credential",
			err,
		}
	}

	return nil
}
