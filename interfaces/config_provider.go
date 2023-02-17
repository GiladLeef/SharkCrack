package interfaces

import (
	"github.com/GiladLeef/SharkCrack/models"
)

type ConfigProvider interface {
	GetConfig() *models.Config
}
