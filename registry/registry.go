package registry

import (
	"math"

	"github.com/bloeys/nmage/assert"
)

type freeListitem struct {
	ItemIndex uint64
	nextFree  *freeListitem
}

// Registry is a storage data structure that can efficiently create/get/free items using generational indices.
// Each item stored in the registry is associated with a 'handle' object that is used to get and free objects
//
// The registry 'owns' all items it stores and returns pointers to items in its array. All items are allocated upfront.
//
// It is NOT safe to concurrently create or free items. However, it is SAFE to concurrently get items
type Registry[T any] struct {
	ItemCount uint
	Handles   []Handle
	Items     []T

	FreeList     *freeListitem
	FreeListSize uint32

	// The number of slots required to be in the free list before the free list
	// is used for creating new entries
	FreeListUsageThreshold uint32
}

func (r *Registry[T]) New() (*T, Handle) {

	assert.T(r.ItemCount < uint(len(r.Handles)), "Can not add more entities to registry because it is full")

	var index uint64 = math.MaxUint64

	// Find index to use for the new item
	if r.FreeList != nil && r.FreeListSize > r.FreeListUsageThreshold {

		index = r.FreeList.ItemIndex

		r.FreeList = r.FreeList.nextFree
		r.FreeListSize--
	} else {

		for i := 0; i < len(r.Handles); i++ {

			handle := r.Handles[i]

			if handle.HasFlag(HandleFlag_Alive) {
				continue
			}

			index = uint64(i)
			break
		}
	}

	if index == math.MaxUint64 {
		panic("failed to create new entity because we did not find a free spot in the registry. Why did the item count assert not go off?")
	}

	var newItem T
	newHandle := NewHandle(r.Handles[index].Generation()+1, HandleFlag_Alive, index)
	assert.T(newHandle != 0, "Entity handle must not be zero")

	r.ItemCount++
	r.Handles[index] = newHandle
	r.Items[index] = newItem

	// It is very important we return directly from the items array, because if we return
	// a pointer to newItem, and T is a value not a pointer, then newItem and what's stored in items will be different
	return &r.Items[index], newHandle
}

func (r *Registry[T]) Get(id Handle) *T {

	index := id.Index()
	assert.T(index < uint64(len(r.Handles)), "Failed to get entity because of invalid entity handle. Handle index is %d while registry only has %d slots. Handle: %+v", index, r.ItemCount, id)

	handle := r.Handles[index]
	if handle.Generation() != id.Generation() || !handle.HasFlag(HandleFlag_Alive) {
		return nil
	}

	item := &r.Items[index]
	return item
}

// Free resets the entity flags then adds this entity to the free list
func (r *Registry[T]) Free(id Handle) {

	index := id.Index()
	assert.T(index < uint64(len(r.Handles)), "Failed to free entity because of invalid entity handle. Handle index is %d while registry only has %d slots. Handle: %+v", index, r.ItemCount, id)

	// Nothing to do if already free
	handle := r.Handles[index]
	if handle.Generation() != id.Generation() || !handle.HasFlag(HandleFlag_Alive) {
		return
	}

	// Generation is incremented on aquire, so here we just reset flags
	r.ItemCount--
	r.Handles[index] = NewHandle(id.Generation(), HandleFlag_None, index)

	// Add to free list
	r.FreeList = &freeListitem{
		ItemIndex: index,
		nextFree:  r.FreeList,
	}
	r.FreeListSize++
}

func (r *Registry[T]) NewIterator() Iterator[T] {
	return Iterator[T]{
		registry:       r,
		remainingItems: r.ItemCount,
		currIndex:      0,
	}
}

func NewRegistry[T any](size uint32) *Registry[T] {
	assert.T(size > 0, "Registry size must be more than zero")
	return &Registry[T]{
		Handles:                make([]Handle, size),
		Items:                  make([]T, size),
		FreeListUsageThreshold: 30,
	}
}
