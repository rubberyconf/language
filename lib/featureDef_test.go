package lib

import (
	"context"
	"fmt"
	"testing"
)

func TestInMemoryOptions(t *testing.T) {

	ctx := context.Background()
	for _, tt := range tests {
		testname := fmt.Sprintf("key: %s, value: %s, duration: %d", tt.key, tt.value, tt.duration)
		t.Run(testname, func(t *testing.T) {
			if !completed || err != nil {
				t.Errorf(" error storing value key: %s", tt.key)
			}

		})
	}
}