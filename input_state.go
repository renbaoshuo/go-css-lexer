package csslexer

type InputState struct {
	pos   int // The current position in the input stream.
	start int // The start position of the current token being read.
}

// NewInputState creates a new InputState instance with the given position and start.
func NewInputState(pos, start int) *InputState {
	return &InputState{
		pos:   pos,
		start: start,
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
