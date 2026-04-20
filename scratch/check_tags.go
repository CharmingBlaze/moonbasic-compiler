//go:build ignore

package main
import (
	"fmt"
	"moonbasic/vm/heap"
)
func main() {
	fmt.Printf("TagEntityRef=%d\n", heap.TagEntityRef)
	fmt.Printf("TagProp=%d\n", heap.TagProp)
}
