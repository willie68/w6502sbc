package emulator

import (
	"fmt"
	"testing"
)

func Test2Kom(t *testing.T) {
	count := 0
	for _, j := range functions {
		if j != nil {
			count++
		}
	}
	fmt.Printf("funccount: %d", count)
}
