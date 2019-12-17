// Consider removing delete credential, since there is no use case for it. If we do want accounts
// to be deactivated, it should be a soft delete instead

package credentials

import (
	"encoding/json"
	"net/http"

	"github.com/obedtandadjaja/project_k_backend/services/auth/controller"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models/credential"
)

type DeleteRequest struct {
	CredentialUuid string `json:"credential_uuid"`
}

type DeleteResponse struct {
	CredentialUuid string `json:"credential_uuid"`
}

func Delete(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseDeleteRequest(r)
	if err != nil {
		return controller.HandlerError{400, "", err}
	}

	response, err := processDeleteRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseDeleteRequest(r *http.Request) (*DeleteRequest, error) {
	var request DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processDeleteRequest(sr *controller.SharedResources, request *DeleteRequest, r *http.Request) (*DeleteResponse, error) {
	var response DeleteResponse

	cred, err := credential.FindBy(sr.DB, map[string]interface{}{
		"uuid": request.CredentialUuid,
	})

	if err != nil {
		return &response, controller.HandlerError{404, "Credential not found", err}
	}

	err = cred.Delete(sr.DB)
	if err != nil {
		return &response, controller.HandlerError{400, "Failed to delete credential", err}
	}

	response.CredentialUuid = cred.Uuid
	return &response, nil
}
