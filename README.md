# csr

`csr` provides a compact adjacency structure for representing a binary relation
between ordered entities.

It stores all neighbors in a single flat slice and allows efficient lookup
via binary search over sorted entities.

## Features

- Generic over `cmp.Ordered`
- Duplicate atoms are ignored
- Entities stored in sorted order
- Neighbor lists preserve insertion order
- Lookup complexity: `O(log |V|)` for entity search, `O(1)` for slice access
- Zero-copy neighbor slices (read-only)

---

## Example

```go
package main

import (
	"fmt"

	"github.com/your/module/csr"
)

func main() {
	left  := []int{1, 1, 2, 3, 3}
	right := []int{2, 3, 3, 4, 4} // duplicate (3 -> 4) ignored

	r := csr.NewRelation(left, right)

	fmt.Println("Right neighbors of 1:", r.GetRightNeighbors(1))
	fmt.Println("Left neighbors of 3:", r.GetLeftNeighbors(3))
}
