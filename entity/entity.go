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

type Entity interface {
	baseEntity()
	GetHandle() EntityHandle
}

type EntityHandle uint64

var _ Entity = &BaseEntity{}

type BaseEntity struct {
	// Byte 1: Generation; Byte 2: Flags; Bytes 3-8: Index
	ID EntityHandle
}

func (be BaseEntity) baseEntity() {}

func (be BaseEntity) GetHandle() EntityHandle {
	return be.ID
}

func (e *BaseEntity) HasFlag(ef EntityFlag) bool {
	return GetFlags(e.ID)&ef > 0
}

func GetGeneration(id EntityHandle) byte {
	return byte(id >> GenerationShiftBits)
}

func GetFlags(id EntityHandle) EntityFlag {
	return EntityFlag(id >> FlagsShiftBits)
}

func GetIndex(id EntityHandle) uint64 {
	return uint64(id & IndexBitMask)
}

func NewEntityId(generation byte, flags EntityFlag, index uint64) EntityHandle {
	return EntityHandle(index | (uint64(generation) << GenerationShiftBits) | (uint64(flags) << FlagsShiftBits))
}
