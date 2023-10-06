package entity

import "github.com/bloeys/nmage/registry"

type Entity interface {
	GetHandle() registry.Handle
}
