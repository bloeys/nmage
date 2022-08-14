package entity

type Entity struct {

	// Byte 1: Generation; Byte 2: Flags; Bytes 3-8: Index
	ID    uint64
	Comps []Comp
}

type Comp interface {
	Name() string
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
