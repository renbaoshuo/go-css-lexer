package csslexer

import (
	"io"
)

type InputState struct {
	inputStream *Input // The input stream of runes.
	pos         int    // The current position in the input stream.
	start       int    // The start position of the current token being read.
}

// State returns the current input state.
// It captures the current position and start position in the input stream.
// This is useful for saving the state of the lexer and restoring it later.
func (z *Input) State() InputState {
	return InputState{
		inputStream: z,
		pos:         z.pos,
		start:       z.start,
	}
}

// Pos returns the current position in the input stream.
func (s *InputState) Pos() int {
	return s.pos
}

// Start returns the start position of the current token being read.
func (s *InputState) Start() int {
	return s.start
}

// Restore restores the input state from the given InputState.
//
// This method is used to restore the input state after parsing a token.
// It allows the lexer to backtrack to a previous state if needed.
func (s *InputState) Restore() {
	s.inputStream.pos = s.pos
	s.inputStream.start = s.start

	if s.pos >= len(s.inputStream.runes) {
		s.inputStream.err = io.EOF
	} else {
		s.inputStream.err = nil
	}
}
