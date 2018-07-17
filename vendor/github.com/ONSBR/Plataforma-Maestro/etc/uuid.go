package etc

import (
	"github.com/google/uuid"
)

//GetUUID returns a new UUID
func GetUUID() string {
	uuid, _ := uuid.NewUUID()
	return uuid.String()
}
