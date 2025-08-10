package csslexer

import (
	"sync"
)

// runeSlicePool is a sync.Pool for reusing slices of runes.
var runeSlicePool = sync.Pool{
	New: func() interface{} {
		s := make([]rune, 0, 16)
		return &s
	},
}

// tokenPool is a sync.Pool for reusing token instances.
var tokenPool = sync.Pool{
	New: func() interface{} {
		return &Token{
			Type: DefaultToken,
			Data: nil,
		}
	},
}
