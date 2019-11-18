package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/obedtandadjaja/project_k_backend/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
