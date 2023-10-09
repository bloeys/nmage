package registry

type HandleFlag byte

const (
	HandleFlag_None  HandleFlag = 0
	HandleFlag_Alive HandleFlag = 1 << (iota - 1)
)

const (
	GenerationShiftBits = 64 - 8
	FlagsShiftBits      = 64 - 16
	IndexBitMask        = 0x00_00_FFFF_FFFF_FFFF
)

// Byte 1: Generation; Byte 2: Flags; Bytes 3-8: Index
type Handle uint64

// IsZero reports whether the handle is in its default 'zero' state.
// A zero handle is an invalid handle that does NOT point to any entity
func (h Handle) IsZero() bool {
	return h == 0
}

func (h Handle) HasFlag(ef HandleFlag) bool {
	return h.Flags()&ef > 0
}

func (h Handle) Generation() byte {
	return byte(h >> GenerationShiftBits)
}

func (h Handle) Flags() HandleFlag {
	return HandleFlag(h >> FlagsShiftBits)
}

func (h Handle) Index() uint64 {
	return uint64(h & IndexBitMask)
}

func NewHandle(generation byte, flags HandleFlag, index uint64) Handle {
	return Handle(index | (uint64(generation) << GenerationShiftBits) | (uint64(flags) << FlagsShiftBits))
}
