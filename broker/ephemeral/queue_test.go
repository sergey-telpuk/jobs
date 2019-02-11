package ephemeral

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueue_ConfigureNil(t *testing.T) {
	q := newQueue(0)
	assert.NoError(t, q.configure(nil, nil))
}
