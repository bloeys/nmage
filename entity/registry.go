package entity

import (
	"github.com/bloeys/nmage/assert"
)

type freeListitem struct {
	EntityIndex uint64
	nextFree    *freeListitem
}

type Registry struct {
	EntityCount uint64
	Entities    []Entity
	FreeList    *freeListitem
}

func (r *Registry) NewEntity() *Entity {

	assert.T(r.EntityCount < uint64(len(r.Entities)), "Can not add more entities to registry because it is full")

	entityToUseIndex := uint64(0)
	var entityToUse *Entity = nil

	if r.FreeList != nil {

		entityToUseIndex = r.FreeList.EntityIndex
		entityToUse = &r.Entities[entityToUseIndex]
		r.FreeList = r.FreeList.nextFree
	} else {

		for i := 0; i < len(r.Entities); i++ {

			e := &r.Entities[i]
			if GetFlags(e.ID) != byte(EntityFlag_Unknown) && !e.HasFlag(EntityFlag_Dead) {
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
	entityToUse.ID = NewEntityId(GetGeneration(entityToUse.ID)+1, byte(EntityFlag_Alive), entityToUseIndex)
	assert.T(entityToUse.ID != 0, "Entity ID must not be zero")
	return entityToUse
}

func (r *Registry) GetEntity(id uint64) *Entity {

	index := GetIndex(id)
	gen := GetGeneration(id)

	e := &r.Entities[index]
	eGen := GetGeneration(e.ID)

	if gen != eGen {
		return nil
	}

	return e
}

func (r *Registry) FreeEntity(id uint64) {

	e := r.GetEntity(id)
	if e == nil {
		return
	}

	r.EntityCount--
	eIndex := GetIndex(e.ID)

	e.Comps = []Comp{}
	e.ID = NewEntityId(GetGeneration(e.ID), byte(EntityFlag_Dead), eIndex)

	r.FreeList = &freeListitem{
		EntityIndex: eIndex,
		nextFree:    r.FreeList,
	}
}

func NewRegistry(size uint32) *Registry {
	assert.T(size > 0, "Registry size must be more than zero")
	return &Registry{
		Entities: make([]Entity, size),
	}
}
