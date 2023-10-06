package entity

import "github.com/bloeys/nmage/registry"

var _ Comp = &BaseComp{}

type BaseComp struct {
	Handle registry.Handle
}

func (b BaseComp) baseComp() {
}

func (b *BaseComp) Init(parentHandle registry.Handle) {
	b.Handle = parentHandle
}

func (b BaseComp) Name() string {
	return "Base Component"
}

func (b BaseComp) Update() {
}

func (b BaseComp) Destroy() {
}
