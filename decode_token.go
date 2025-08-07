package csslexer

import (
	"strings"
)

// DecodeToken decodes a CSS token based on its type and raw rune slice.
//
// It handles escape sequences for identifier-like tokens, strings, and URLs.
//
// Return: the decoded string representation of the token.
func DecodeToken(tokenType TokenType, raw []rune) string {
	switch tokenType {
	case IdentToken, FunctionToken, AtKeywordToken, HashToken, DimensionToken, StringToken, UrlToken:
		return decodeEscapeSequences(raw)
	default:
		return string(raw)
	}
}

// decodeEscapeSequences is a common function to decode escape sequences
// in a rune slice.
func decodeEscapeSequences(runes []rune) string {
	if len(runes) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(runes)) // Pre-allocate capacity

	i := 0
	for i < len(runes) {
		if runes[i] == '\\' && i+1 < len(runes) && twoCharsAreValidEscape(runes[i], runes[i+1]) {
			// Handle escape sequence
			decoded := consumeEscapeSequence(runes, &i)
			builder.WriteRune(decoded)
		} else {
			builder.WriteRune(runes[i])
			i++
		}
	}

	return builder.String()
}

// consumeEscapeSequence consumes an escape sequence from the rune slice
// starting at index i and returns the decoded rune. The index i is updated
// to point to the next character after the escape sequence.
func consumeEscapeSequence(runes []rune, i *int) rune {
	*i++ // Skip the backslash

	if *i >= len(runes) {
		return '\uFFFD' // Replacement character for EOF
	}

	c := runes[*i]

	if isASCIIHexDigit(c) {
		// Hex escape sequence: \123456
		var res rune = 0
		hexCount := 0

		for hexCount < 6 && *i < len(runes) && isASCIIHexDigit(runes[*i]) {
			res = res*16 + hexDigitToValue(runes[*i])
			*i++
			hexCount++
		}

		// Skip trailing whitespace after hex escape
		if *i < len(runes) && isHTMLWhitespace(runes[*i]) {
			*i++
		}

		if !isValidCodePoint(res) {
			return '\uFFFD' // Replacement character
		}

		return res
	} else {
		// Single character escape: \n, \", \\, etc.
		*i++
		return c
	}
}
