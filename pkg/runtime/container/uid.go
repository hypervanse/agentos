package container

import (
	"github.com/google/uuid"
)

func GenerateContainerUID() string {
	return uuid.NewString()
}
