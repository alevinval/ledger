package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartitionSelector(t *testing.T) {
	partitionSelector := PartitionSelector{Size: 3}

	selected := partitionSelector.Next()
	assert.Equal(t, 0, selected)

	selected = partitionSelector.Next()
	assert.Equal(t, 1, selected)

	selected = partitionSelector.Next()
	assert.Equal(t, 2, selected)

	selected = partitionSelector.Next()
	assert.Equal(t, 0, selected)
}
