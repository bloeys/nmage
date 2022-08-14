package entity

import "github.com/bloeys/nmage/assert"

type Registry struct {
	EntityCount uint64
	Entities    []Entity
}

func (r *Registry) NewEntity() *Entity {

	assert.T(r.EntityCount < uint64(len(r.Entities)), "Can not add more entities to registry because it is full")

	for i := 0; i < len(r.Entities); i++ {

		// @TODO: Implement generational indices
		e := &r.Entities[i]
		if e.ID == 0 {
			r.EntityCount++
			e.ID = uint64(i) + 1

			assert.T(e.ID != 0, "Entity ID must not be zero")
			return e
		}
	}

	panic("failed to create new entity because we did not find a free spot in the registry. Why did the assert not go off?")
}

func NewRegistry(size uint32) *Registry {
	assert.T(size > 0, "Registry size must be more than zero")
	return &Registry{
		Entities: make([]Entity, size),
	}
}
