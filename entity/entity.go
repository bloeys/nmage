package entity

type EntityFlag byte

const (
	EntityFlag_None  EntityFlag = 0
	EntityFlag_Alive EntityFlag = 1 << (iota - 1)
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

func GetFlags(id uint64) EntityFlag {
	return EntityFlag(id >> FlagsShiftBits)
}

func GetIndex(id uint64) uint64 {
	return id & IndexBitMask
}

func (e *Entity) HasFlag(ef EntityFlag) bool {
	return GetFlags(e.ID)&ef > 0
}

func NewEntityId(generation byte, flags EntityFlag, index uint64) uint64 {
	return index | (uint64(generation) << GenerationShiftBits) | (uint64(flags) << FlagsShiftBits)
}
