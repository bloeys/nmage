package entity

import "github.com/bloeys/nmage/assert"

type Comp interface {
	// This ensures that implementors of the Comp interface
	// always embed BaseComp
	base()

	Name() string
	Init(parent *Entity)
	Update()
	Destroy()
}

func AddComp[T Comp](e *Entity, c T) {

	assert.T(!HasComp[T](e), "Entity with id %v already has component with name %s", e.ID, c.Name())

	e.Comps = append(e.Comps, c)
	c.Init(e)
}

func HasComp[T Comp](e *Entity) bool {

	for i := 0; i < len(e.Comps); i++ {

		_, ok := e.Comps[i].(T)
		if ok {
			return true
		}
	}

	return false
}

func GetComp[T Comp](e *Entity) (out T) {

	for i := 0; i < len(e.Comps); i++ {

		comp, ok := e.Comps[i].(T)
		if ok {
			return comp
		}
	}

	return out
}

// DestroyComp calls Destroy on the component and then removes it from the entities component list
func DestroyComp[T Comp](e *Entity) {

	for i := 0; i < len(e.Comps); i++ {

		comp, ok := e.Comps[i].(T)
		if ok {
			comp.Destroy()
			e.Comps = append(e.Comps[:i], e.Comps[i+1:]...)
			return
		}
	}
}
