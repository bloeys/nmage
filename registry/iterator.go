package registry

// Iterator goes through the entire registry it was created from and
// returns all alive items, and nil after its done.
//
// The iterator will still work if items are added/removed to the registry
// after it was created, but the following conditions apply:
//   - If items are removed, iterator will not show the removed items (assuming it didn't return them before their removal)
//   - If items are added, the iterator will either only return older items (i.e. is not affected), or only return newer items (i.e. items that were going to be returned before will now not get returned in favor of newly inserted items), or a mix of old and new items.
//     However, in all cases the iterator will *never* returns more items than were alive at the time of the iterator's creation.
//   - If items were both added and removed, the iterator might follow either of the previous 2 cases or a combination of them
//
// To summarize: The iterator will *never* return more items than were alive at the time of its creation, and will *never* return freed items
//
// Example usage:
//
//	for item, handle := it.Next(); !it.IsDone(); item, handle = it.Next() {
//		// Do stuff
//	}
type Iterator[T any] struct {
	registry       *Registry[T]
	remainingItems uint
	currIndex      int
}

func (it *Iterator[T]) Next() (*T, Handle) {

	if it.IsDone() {
		return nil, 0
	}

	// If IsDone() only checked 'remainingItems', then when Next() returns the last item IsDone() will immediately be true which will cause loops to exit before processing the last item!
	// However, with this check IsDone will remain false until Next() is called at least one more time after returning the last item which ensures the last item is processed in the loop.
	//
	// In cases where iterator is created on an empty registry, IsDone() will report true and the above check will return early
	if it.remainingItems == 0 {
		it.currIndex = -1
		return nil, 0
	}

	for ; it.currIndex < len(it.registry.Handles); it.currIndex++ {

		handle := it.registry.Handles[it.currIndex]
		if !handle.HasFlag(HandleFlag_Alive) {
			continue
		}

		it.remainingItems--
		it.currIndex++
		return &it.registry.Items[it.currIndex], handle
	}

	// If we reached here means we iterated to the end and didn't find anything, which probably
	// means that the registry changed since we were created, and that remainingItems is not accurate.
	//
	// As such, we zero remaining items so that this iterator is considered done
	it.currIndex = -1
	it.remainingItems = 0
	return nil, 0
}

func (it *Iterator[T]) IsDone() bool {

	if it.remainingItems != 0 {
		return false
	}

	// We have two cases here:
	// 1. Index of zero means Next() never returned an item. Remaining items of zero without returning anything means we have an empty registry and so its safe to report done
	// 2. Negative index means Next() has detected we reached the end and that its safe to report being done
	return it.currIndex <= 0
}
