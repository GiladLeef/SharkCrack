package interfaces

import (
	"github.com/GiladLeef/SharkCrack/models"
)

type ClientStopQueue interface {
	Get() (models.ClientStopReason, error)
	Put(models.ClientStopReason) error
}
