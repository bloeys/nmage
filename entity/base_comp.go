package entity

import "github.com/bloeys/nmage/assert"

var _ Comp = &BaseComp{}

type BaseComp struct {
	Entity *BaseEntity
}

func (b BaseComp) baseComp() {
}

func (b *BaseComp) Init(parent *BaseEntity) {
	assert.T(parent != nil, "Component was initialized with a nil parent. That is not allowed.")
	b.Entity = parent
}

func (b BaseComp) Name() string {
	return "Base Component"
}

func (b BaseComp) Update() {
}

func (b BaseComp) Destroy() {
}
