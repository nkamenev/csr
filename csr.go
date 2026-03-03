package csr

import (
	"cmp"
	"slices"
)

// Relation is a compact adjacency structure representing a binary relation
// over a set of entities.
//
// Each entity may have two neighbor lists:
//   - left  neighbors: entities that point to this entity
//   - right neighbors: entities this entity points to
//
// Entities are stored in sorted order.
//
// Internally, all neighbor lists are packed into a single flat slice.
// For each entity, two consecutive boundary offsets are stored in `bounds`:
//   - bounds[2*i]   — end of the left-neighbor segment
//   - bounds[2*i+1] — end of the right-neighbor segment
//
// Lookup complexity:
//   - entity lookup: O(log |V|)
//   - neighbor slice access: O(1)
//
// The returned neighbor slices reference internal storage and must be treated
// as read-only.
type Relation[T cmp.Ordered] struct {
	entities  []T
	bounds    []int
	neighbors []T
}

// NewRelation constructs a Relation from two equally sized slices representing
// a set of binary atoms (left[i] -> right[i]).
//
// Duplicate atoms are ignored.
// Only entities that participate in at least one atom are stored.
// Neighbor lists preserve atom insertion order.
//
// If the input slices are empty, have different lengths, or are nil,
// NewRelation returns nil.
func NewRelation[T cmp.Ordered](left, right []T) *Relation[T] {
	if len(left) == 0 || len(right) == 0 || len(left) != len(right) {
		return nil
	}

	type neighbors struct {
		left, right []T
	}

	type atom struct {
		left, right T
	}

	uniqAtoms := make(map[atom]struct{}, len(left))
	uniqEnts := map[T]*neighbors{}
	var neighCount int
	for i, la := range left {
		curr := atom{la, right[i]}
		if _, ok := uniqAtoms[curr]; ok {
			continue
		}
		uniqAtoms[curr] = struct{}{}
		// initialize neighbors struct on first encounter of entity
		if _, ok := uniqEnts[curr.left]; !ok {
			uniqEnts[curr.left] = &neighbors{}
		}
		if _, ok := uniqEnts[curr.right]; !ok {
			uniqEnts[curr.right] = &neighbors{}
		}

		// add neighbors to left/right for each entity
		uniqEnts[curr.left].right = append(uniqEnts[curr.left].right, curr.right)
		uniqEnts[curr.right].left = append(uniqEnts[curr.right].left, curr.left)
	}

	// count total neighbors to allocate slice exactly once
	for _, ent := range uniqEnts {
		neighCount += len(ent.left) + len(ent.right)
	}

	// collect all uniqEntsue entities and sort them
	entities := make([]T, 0, len(uniqEnts))
	for id := range uniqEnts {
		entities = append(entities, id)
	}
	slices.Sort(entities)

	// prepare bounds and neighbors slices
	bounds := make([]int, 0, len(entities)<<1)
	neigh := make([]T, 0, neighCount)
	for _, e := range entities {
		neigh = append(neigh, uniqEnts[e].left...)
		bounds = append(bounds, len(neigh))
		neigh = append(neigh, uniqEnts[e].right...)
		bounds = append(bounds, len(neigh))
	}

	return &Relation[T]{
		entities:  entities,
		bounds:    bounds,
		neighbors: neigh,
	}
}

// GetLeftNeighbors returns all left neighbors of ent.
//
// A left neighbor is an entity L such that (L -> ent) exists in the relation.
//
// If ent is not present or has no left neighbors, nil is returned.
//
// The returned slice aliases internal storage and must not be modified.
func (r *Relation[T]) GetLeftNeighbors(ent T) []T {
	return r.getNeighbors(ent, false)
}

// GetRightNeighbors returns all right neighbors of ent.
//
// A right neighbor is an entity R such that (ent -> R) exists in the relation.
//
// If ent is not present or has no right neighbors, nil is returned.
//
// The returned slice aliases internal storage and must not be modified.
func (r *Relation[T]) GetRightNeighbors(ent T) []T {
	return r.getNeighbors(ent, true)
}

// getNeighbors performs a binary search over sorted entities and
// returns the corresponding neighbor segment.
//
// isRight selects between left (false) and right (true) neighbor lists.
func (r *Relation[T]) getNeighbors(ent T, isRight bool) []T {
	idx, ok := slices.BinarySearch(r.entities, ent)
	if !ok {
		return nil
	}

	// each entity has two bounds: left at 2*i, right at 2*i+1
	bnd := idx << 1
	if isRight {
		bnd++
	}

	// start is previous bound (or 0 for first block)
	start := 0
	if bnd > 0 {
		start = r.bounds[bnd-1]
	}
	end := r.bounds[bnd]

	if start == end {
		return nil
	}
	return r.neighbors[start:end]
}
