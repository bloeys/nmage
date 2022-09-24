package entity

type EntityFlag byte

const (
	EntityFlag_Unknown EntityFlag = 0
	EntityFlag_Dead    EntityFlag = 1 << (iota - 1)
	EntityFlag_Alive
)

const (
	GenerationShiftBits = 64 - 8
	FlagsShiftBits      = 64 - 16
	IndexBitMask        = 0x00_00_FFFF_FFFF_FFFF
)

type Entity struct {

	// Byte 1: Generation; Byte 2: Flags; Bytes 3-8: Index
	ID    uint64
	Comps []Comp
}

func GetGeneration(id uint64) byte {
	return byte(id >> GenerationShiftBits)
}

func GetFlags(id uint64) byte {
	return byte(id >> FlagsShiftBits)
}

func GetIndex(id uint64) uint64 {
	return id & IndexBitMask
}

func (e *Entity) HasFlag(ef EntityFlag) bool {
	return GetFlags(e.ID)&byte(ef) > 0
}

func NewEntityId(generation, flags byte, index uint64) uint64 {
	return index | (uint64(generation) << GenerationShiftBits) | (uint64(flags) << FlagsShiftBits)
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
