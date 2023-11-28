package cors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCORS(t *testing.T) {
	fmt.Println("TESTING...")
	a := NewCORS()
	assert.NotNil(t, a, "they should be equal")
}
