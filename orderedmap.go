// Package orderedmap implements an ordered map, i.e. a map that also keeps track of
// the order in which keys were inserted.
//
// All operations are constant-time.
//
// Github repo: https://github.com/DominicTobias/go-ordered-map
//
package orderedmap

import (
	"github.com/DominicTobias/go-ordered-map/list"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V

	element *list.Element[*Pair[K,V]]
}

type OrderedMap[K comparable, V any] struct {
	pairs map[K]*Pair[K,V]
	list  *list.List[*Pair[K,V]]
}

// New creates a new OrderedMap.
func New[K comparable, V any]() *OrderedMap[K,V] {
	return &OrderedMap[K,V]{
		pairs: make(map[K]*Pair[K,V]),
		list:  list.New[*Pair[K,V]](),
	}
}

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is present in the map.
func (om *OrderedMap[K,V]) Get(key K) (V, bool) {
	if pair, present := om.pairs[key]; present {
		return pair.Value, present
	}
	var empty V
	return empty, false
}

// GetPair looks for the given key, and returns the pair associated with it,
// or nil if not found. The Pair struct can then be used to iterate over the ordered map
// from that point, either forward or backward.
func (om *OrderedMap[K,V]) GetPair(key K) *Pair[K,V] {
	return om.pairs[key]
}

// Set sets the key-value pair, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (om *OrderedMap[K,V]) Set(key K, value V) (V, bool) {
	if pair, present := om.pairs[key]; present {
		oldValue := pair.Value
		pair.Value = value
		return oldValue, true
	}

	pair := &Pair[K,V]{
		Key:   key,
		Value: value,
	}
	// cannot use pair (variable of type *Pair[K, V]) as type Pair[K, V] in argument to om.list.PushBack
	pair.element = om.list.PushBack(pair)
	om.pairs[key] = pair

	var empty V
	return empty, false
}

// Delete removes the key-value pair, and returns what `Get` would have returned
// on that key prior to the call to `Delete`.
func (om *OrderedMap[K,V]) Delete(key K) (V, bool) {
	if pair, present := om.pairs[key]; present {
		om.list.Remove(pair.element)
		delete(om.pairs, key)
		return pair.Value, true
	}

	var empty V
	return empty, false
}

// Len returns the length of the ordered map.
func (om *OrderedMap[K,V]) Len() int {
	return len(om.pairs)
}

// Oldest returns a pointer to the oldest pair. It's meant to be used to iterate on the ordered map's
// pairs from the oldest to the newest, e.g.:
// for pair := orderedMap.Oldest(); pair != nil; pair = pair.Next() { fmt.Printf("%v => %v\n", pair.Key, pair.Value) }
func (om *OrderedMap[K,V]) Oldest() *Pair[K,V] {
	// cannot use om.list.Front() (value of type *list.Element[Pair[K, V]]) as type *list.Element[*Pair[K, V]] in argument to listElementToPair[K, V]
	return listElementToPair[K,V](om.list.Front())
}

// Newest returns a pointer to the newest pair. It's meant to be used to iterate on the ordered map's
// pairs from the newest to the oldest, e.g.:
// for pair := orderedMap.Oldest(); pair != nil; pair = pair.Next() { fmt.Printf("%v => %v\n", pair.Key, pair.Value) }
func (om *OrderedMap[K,V]) Newest() *Pair[K,V] {
	return listElementToPair[K,V](om.list.Back())
}

// Next returns a pointer to the next pair.
func (p *Pair[K,V]) Next() *Pair[K,V] {
	return listElementToPair[K,V](p.element.Next())
}

// Previous returns a pointer to the previous pair.
func (p *Pair[K,V]) Prev() *Pair[K,V] {
	return listElementToPair[K,V](p.element.Prev())
}

func listElementToPair[K comparable, V any](element *list.Element[*Pair[K,V]]) *Pair[K,V] {
	if element == nil {
		return nil
	}
	return element.Value
}
