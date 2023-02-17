package interfaces

import (
	"github.com/GiladLeef/SharkCrack/models"
)

type SubmissionQueue interface {
	Size() int
	Get() (models.HashSubmission, error)
	Put(models.HashSubmission) error
}
