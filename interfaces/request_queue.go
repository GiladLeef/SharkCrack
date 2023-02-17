package interfaces

import (
	"github.com/GiladLeef/SharkCrack/models"
)

type RequestQueue interface {
	Size() int
	Get() (models.HashingRequest, error)
	Put(models.HashingRequest) error
}
