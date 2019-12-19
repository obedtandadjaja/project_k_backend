// Consider removing delete credential, since there is no use case for it. If we do want accounts
// to be deactivated, it should be a soft delete instead

package credentials

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type DeleteRequest struct {
	CredentialID uuid.UUID `json:"credential_id"`
}

type DeleteResponse struct {
	CredentialID uuid.UUID `json:"credential_id"`
}

func Delete(tx *pop.Connection, request *DeleteRequest) (*DeleteResponse, error) {
	var response DeleteResponse

	cred := &models.Credential{}
	if err := tx.Find(cred, request.CredentialID); err != nil {
		return &response, err
	}

	if err := tx.Destroy(cred); err != nil {
		return &response, err
	}

	response.CredentialID = cred.ID
	return &response, nil
}
