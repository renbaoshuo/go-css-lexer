package csslexer

import (
	"sync"
)

// tokenPool is a sync.Pool for reusing token instances.
var tokenPool = sync.Pool{
	New: func() interface{} {
		return &Token{
			Type: DefaultToken,
			Data: "",
		}
	},
}
