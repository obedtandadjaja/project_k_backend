package helpers

import "github.com/gofrs/uuid"

func ParseUUID(str string) uuid.UUID {
	uuid, _ := uuid.FromString(str)
	return uuid
}
