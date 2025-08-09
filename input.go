package csslexer

import (
	"io"
)

// NOTE: We use rune instead of byte to handle Unicode characters properly.

// EOF is a special rune that represents the end of the input.
const EOF = rune(0)

// nullRune is an empty slice of runes, used to represent no runes.
var nullRune = []rune{}

// Input represents a stream of runes read from a source.
type Input struct {
	runes []rune // The runes in the input stream.
	pos   int    // The current position in the input stream.
	start int    // The start position of the current token being read.
	err   error  // Any error encountered while reading the input.
}

// NewInput creates a new Input instance from the given string.
func NewInput(input string) *Input {
	return NewInputRunes([]rune(input))
}

// NewInputRunes creates a new Input instance from the given slice of runes.
func NewInputRunes(runes []rune) *Input {
	return &Input{
		runes: runes,
		pos:   0,
		start: 0,
		err:   nil,
	}
}

// NewInputBytes creates a new Input instance from the given byte slice.
func NewInputBytes(input []byte) *Input {
	return NewInputRunes([]rune(string(input)))
}

// NewInputReader creates a new Input instance from the given io.Reader.
func NewInputReader(r io.Reader) *Input {
	var b []byte

	if r != nil {
		var err error
		b, err = io.ReadAll(r)
		if err != nil {
			return &Input{
				runes: nullRune,
				pos:   0,
				start: 0,
				err:   err,
			}
		}
	}

	return NewInputBytes(b)
}

// PeekErr checks if there is an error at the current position plus the specified offset.
func (z *Input) PeekErr(pos int) error {
	if z.err != nil {
		return z.err
	}

	if z.pos+pos >= len(z.runes) {
		return io.EOF
	}

	return nil
}

// Err returns the error at the current position.
func (z *Input) Err() error {
	return z.PeekErr(0)
}

// Peek returns the next rune in the input stream without advancing the position.
func (z *Input) Peek(n int) rune {
	if z.pos+n >= len(z.runes) {
		return EOF
	}
	return z.runes[z.pos+n]
}

// Move advances the position by the specified number of runes.
func (z *Input) Move(n int) {
	if z.pos+n >= len(z.runes) {
		z.pos = len(z.runes)
		z.err = io.EOF
		return
	}
	z.pos += n
}

// Current returns the current token as a slice of runes.
func (z *Input) Current() []rune {
	if z.start >= z.pos {
		return nullRune
	}
	return z.runes[z.start:z.pos:z.pos]
}

// Shift returns the current token and resets the start position to the current position.
func (z *Input) Shift() []rune {
	// Shift returns the current token and resets the start position.
	current := z.Current()
	z.start = z.pos
	return current
}

// MoveWhilePredicate advances the position while the predicate function returns true for the current rune.
func (z *Input) MoveWhilePredicate(pred func(rune) bool) {
	for pred(z.Peek(0)) {
		z.Move(1)
	}
}

// SaveState saves the current position and start position.
// It returns the current position and start position as integers.
func (z *Input) State() InputState {
	return *NewInputState(z.pos, z.start)
}

// RestoreState restores the input state from the given InputState.
// It sets the current position and start position to the values in the InputState.
// If the position is beyond the end of the input, it sets the error to io.EOF.
// Otherwise, it clears the error.
//
// This method is used to restore the input state after parsing a token.
// It allows the lexer to backtrack to a previous state if needed.
func (z *Input) RestoreState(s InputState) {
	z.pos = s.Pos()
	z.start = s.Start()

	if z.pos >= len(z.runes) {
		z.err = io.EOF
	} else {
		z.err = nil
	}
}
