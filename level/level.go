package level

import (
	"github.com/bloeys/nmage/assert"
	"github.com/bloeys/nmage/entity"
)

type Level struct {
	*entity.Registry
	Name string
}

func NewLevel(name string, maxEntities uint32) *Level {

	assert.T(name != "", "Level name can not be empty")
	return &Level{
		Name:     name,
		Registry: entity.NewRegistry(maxEntities),
	}
}
