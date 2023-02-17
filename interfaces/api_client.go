package interfaces

import (
	"github.com/GiladLeef/SharkCrack/models"
)

type ApiClient interface {
	GetHashName() (int, string)
	GetPasswords(int) (int, []string)
	SubmitHashes(models.HashSubmission) int
}
