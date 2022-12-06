package entity

import (
	"github.com/bloeys/nmage/assert"
)

var (
	// The number of slots required to be in the free list before the free list
	// is used for creating new entries
	FreeListUsageThreshold uint32 = 20
)

type freeListitem struct {
	EntityIndex uint64
	nextFree    *freeListitem
}

type Registry struct {
	EntityCount uint64
	Entities    []Entity

	FreeList     *freeListitem
	FreeListSize uint32
}

func (r *Registry) NewEntity() *Entity {

	assert.T(r.EntityCount < uint64(len(r.Entities)), "Can not add more entities to registry because it is full")

	entityToUseIndex := uint64(0)
	var entityToUse *Entity = nil

	if r.FreeList != nil && r.FreeListSize > FreeListUsageThreshold {

		entityToUseIndex = r.FreeList.EntityIndex
		entityToUse = &r.Entities[entityToUseIndex]
		r.FreeList = r.FreeList.nextFree
		r.FreeListSize--
	} else {

		for i := 0; i < len(r.Entities); i++ {

			e := &r.Entities[i]
			if e.HasFlag(EntityFlag_Alive) {
				continue
			}

			entityToUse = e
			entityToUseIndex = uint64(i)
			break
		}
	}

	if entityToUse == nil {
		panic("failed to create new entity because we did not find a free spot in the registry. Why did the assert not go off?")
	}

	r.EntityCount++
	entityToUse.ID = NewEntityId(GetGeneration(entityToUse.ID)+1, EntityFlag_Alive, entityToUseIndex)
	assert.T(entityToUse.ID != 0, "Entity ID must not be zero")
	return entityToUse
}

func (r *Registry) GetEntity(id EntityHandle) *Entity {

	index := GetIndex(id)
	gen := GetGeneration(id)

	e := &r.Entities[index]
	eGen := GetGeneration(e.ID)

	if gen != eGen {
		return nil
	}

	return e
}

// FreeEntity calls Destroy on all the entities components, resets the component list, resets the entity flags, then ads this entity to the free list
func (r *Registry) FreeEntity(id EntityHandle) {

	e := r.GetEntity(id)
	if e == nil {
		return
	}

	for i := 0; i < len(e.Comps); i++ {
		e.Comps[i].Destroy()
	}

	r.EntityCount--
	eIndex := GetIndex(e.ID)

	e.Comps = []Comp{}
	e.ID = NewEntityId(GetGeneration(e.ID), EntityFlag_None, eIndex)

	r.FreeList = &freeListitem{
		EntityIndex: eIndex,
		nextFree:    r.FreeList,
	}
	r.FreeListSize++
}

func NewRegistry(size uint32) *Registry {
	assert.T(size > 0, "Registry size must be more than zero")
	return &Registry{
		Entities: make([]Entity, size),
	}
}
