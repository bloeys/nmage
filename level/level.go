package level

import (
	"github.com/bloeys/nmage/assert"
)

type Level struct {
	Name string
}

func NewLevel(name string) *Level {

	assert.T(name != "", "Level name can not be empty")
	return &Level{
		Name: name,
	}
}
