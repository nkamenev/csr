package csr

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewRelationFromAtoms(t *testing.T) {

	tests := map[string]struct {
		l, r      []int
		wantEnts  []int
		wantBnds  []int
		wantNeigh []int
	}{
		"empty input": {
			wantEnts:  nil,
			wantBnds:  nil,
			wantNeigh: nil,
		},
		"single atom": {
			l:         []int{1},
			r:         []int{2},
			wantEnts:  []int{1, 2},
			wantBnds:  []int{0, 1, 2, 2},
			wantNeigh: []int{2, 1},
		},
		"two atoms with shared entity": {
			l:         []int{1, 1},
			r:         []int{2, 3},
			wantEnts:  []int{1, 2, 3},
			wantBnds:  []int{0, 2, 3, 3, 4, 4},
			wantNeigh: []int{2, 3, 1, 1},
		},
		"self loop": {
			l:         []int{1},
			r:         []int{1},
			wantEnts:  []int{1},
			wantBnds:  []int{1, 2},
			wantNeigh: []int{1, 1},
		},
		"duplicate atoms": {
			l:         []int{1, 1},
			r:         []int{2, 2},
			wantEnts:  []int{1, 2},
			wantBnds:  []int{0, 1, 2, 2},
			wantNeigh: []int{2, 1},
		},
		"multiple entities": {
			l:        []int{1, 1, 2, 3},
			r:        []int{2, 3, 3, 1},
			wantEnts: []int{1, 2, 3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewRelation(tt.l, tt.r)

			if tt.l == nil {
				if r != nil {
					t.Fatalf("expected nil relation, got %+v", r)
				}
				return
			}

			if !reflect.DeepEqual(r.entities, tt.wantEnts) {
				t.Fatalf("entities = %v, want %v", r.entities, tt.wantEnts)
			}

			if tt.wantBnds != nil && !reflect.DeepEqual(r.bounds, tt.wantBnds) {
				fmt.Printf("neighbors = %v\n", r.neighbors)
				t.Fatalf("bounds = %v, want %v", r.bounds, tt.wantBnds)
			}

			if tt.wantNeigh != nil && !reflect.DeepEqual(r.neighbors, tt.wantNeigh) {
				t.Fatalf("neighbors = %v, want %v", r.neighbors, tt.wantNeigh)
			}
		})
	}
}

func TestGetLeftNeighbors(t *testing.T) {
	tests := map[string]struct {
		l, r []int
		ent  int
		want []int
	}{
		"nil relation": {
			ent:  1,
			want: nil,
		},
		"entity not present": {
			l:    []int{1},
			r:    []int{2},
			ent:  3,
			want: nil,
		},
		"no left neighbors": {
			l:    []int{1},
			r:    []int{2},
			ent:  1,
			want: nil,
		},
		"single left neighbor": {
			l:    []int{1},
			r:    []int{2},
			ent:  2,
			want: []int{1},
		},
		"multiple left neighbors": {
			l:    []int{1, 2, 4},
			r:    []int{3, 3, 3},
			ent:  3,
			want: []int{1, 2, 4},
		},
		"self loop": {
			l:    []int{1},
			r:    []int{1},
			ent:  1,
			want: []int{1},
		},
		"mixed graph": {
			l:    []int{1, 3, 2, 5},
			r:    []int{2, 2, 4, 2},
			ent:  2,
			want: []int{1, 3, 5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewRelation(tt.l, tt.r)

			if r == nil {
				if tt.want != nil {
					t.Fatalf("expected %v, got nil relation", tt.want)
				}
				return
			}

			got := r.GetLeftNeighbors(tt.ent)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetLeftNeighbors(%v) = %v, want %v",
					tt.ent, got, tt.want)
			}
		})
	}
}

func TestGetRightNeighbors(t *testing.T) {
	tests := map[string]struct {
		l, r []int
		ent  int
		want []int
	}{
		"nil relation": {
			ent:  1,
			want: nil,
		},
		"entity not present": {
			l:    []int{1},
			r:    []int{2},
			ent:  3,
			want: nil,
		},
		"no right neighbors": {
			l:    []int{1},
			r:    []int{2},
			ent:  2,
			want: nil,
		},
		"single right neighbor": {
			l:    []int{1},
			r:    []int{2},
			ent:  1,
			want: []int{2},
		},
		"multiple right neighbors": {
			l:    []int{3, 3, 3},
			r:    []int{1, 2, 4},
			ent:  3,
			want: []int{1, 2, 4},
		},
		"self loop": {
			l:    []int{1},
			r:    []int{1},
			ent:  1,
			want: []int{1},
		},
		"mixed graph": {
			l:    []int{2, 2, 4, 2},
			r:    []int{1, 3, 2, 5},
			ent:  2,
			want: []int{1, 3, 5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewRelation(tt.l, tt.r)
			if r == nil {
				if tt.want != nil {
					t.Fatalf("expected %v, got nil relation", tt.want)
				}
				return
			}

			got := r.GetRightNeighbors(tt.ent)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetRightNeighbors(%v) = %v, want %v",
					tt.ent, got, tt.want)
			}
		})
	}
}
