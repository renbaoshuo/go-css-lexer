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
	// Preprocess.
	for i, r := range runes {
		if r == 0x00 { // U+0000 NULL CHARACTER
			runes[i] = '\uFFFD' // Replace with U+FFFD REPLACEMENT CHARACTER
		}
	}

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

// CurrentOffset returns the current offset in the input stream.
//
// It calculates the offset as the difference between the current position
// and the start position.
func (z *Input) CurrentOffset() int {
	return z.pos - z.start
}

// Current returns the current token as a slice of runes.
func (z *Input) Current() []rune {
	if z.start >= z.pos {
		return nullRune
	}
	return z.runes[z.start:z.pos:z.pos]
}

// CurrentString returns the current token as a string.
func (z *Input) CurrentString() string {
	return string(z.Current())
}

// CurrentSuffix returns the current token after applying the
// offset.
//
// If the offset is greater than the current position, it returns an
// empty slice.
func (z *Input) CurrentSuffix(offset int) []rune {
	if z.start+offset >= z.pos {
		return nullRune
	}
	return z.runes[z.start+offset : z.pos : z.pos]
}

// CurrentSuffixString returns the current token as a string after
// applying the offset.
func (z *Input) CurrentSuffixString(offset int) string {
	return string(z.CurrentSuffix(offset))
}

// Shift resets the start position to the current position.
func (z *Input) Shift() {
	z.start = z.pos
}

// MoveWhilePredicate advances the position while the predicate function returns true for the current rune.
func (z *Input) MoveWhilePredicate(pred func(rune) bool) {
	for pred(z.Peek(0)) {
		z.Move(1)
	}
}
