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

// token represents a parsed token with its type and data.
//
// It is used internally by the lexer to represent the tokens
// it generates.
type token struct {
	Type TokenType
	Data []rune
}

// tokenPool is a sync.Pool for reusing token instances.
var tokenPool = sync.Pool{
	New: func() interface{} {
		return &token{
			Type: DefaultToken,
			Data: nil,
		}
	},
}
