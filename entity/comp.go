package entity

type Comp interface {
	// This ensures that implementors of the Comp interface
	// always embed BaseComp
	base()

	Name() string
}

var _ Comp = &BaseComp{}

type BaseComp struct {
}

func (b *BaseComp) base() {
}

func (b *BaseComp) Name() string {
	return "Base Component"
}

func AddComp(e *Entity, c Comp) {
	e.Comps = append(e.Comps, c)
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

func GetAllCompOfType[T Comp](e *Entity) (out []T) {

	out = []T{}
	for i := 0; i < len(e.Comps); i++ {

		comp, ok := e.Comps[i].(T)
		if ok {
			out = append(out, comp)
		}
	}

	return out
}
