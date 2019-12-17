package refresh_tokens

import "github.com/obedtandadjaja/project_k_backend/services/auth/models/session"

type GetAllRequest struct {
	Uuid string `json:"credential_uuid"`
}

type GetAllResponse struct {
	RefreshTokens []*session.Session `json:"sessions"`
}
