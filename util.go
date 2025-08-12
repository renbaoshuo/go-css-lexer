package csslexer

import (
	"strings"

	"go.baoshuo.dev/cssutil"
)

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#starts-with-a-valid-escape
func twoCharsAreValidEscape(first, second rune) bool {
	return first == '\\' && !cssutil.IsNewline(second)
}

// https://infra.spec.whatwg.org/#surrogate
func isSurrogate(c rune) bool {
	// Surrogate pairs are in the range U+D800 to U+DFFF
	return c >= 0xD800 && c <= 0xDFFF
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#maximum-allowed-code-point
func isLowerThanMaxCodePoint(c rune) bool {
	return c <= 0x10FFFF
}

func isValidCodePoint(c rune) bool {
	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-escaped-code-point
	// If this number is zero, or is for a surrogate, or is greater than the maximum allowed code point, return U+FFFD REPLACEMENT CHARACTER
	return c != 0 && !isSurrogate(c) && isLowerThanMaxCodePoint(c)
}

// hexDigitToValue converts a hex digit character to its numeric value
func hexDigitToValue(c rune) rune {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	if c >= 'a' && c <= 'f' {
		return c - 'a' + 10
	}
	if c >= 'A' && c <= 'F' {
		return c - 'A' + 10
	}
	return 0
}

// equalIgnoringASCIICase compares two rune slices for equality, ignoring ASCII case.
func equalIgnoringASCIICase(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			continue
		} else if cssutil.IsLetter(a[i]) && cssutil.IsLetter(b[i]) {
			if cssutil.IsUpperCaseLetter(a[i]) {
				if a[i]+32 != b[i] {
					return false
				}
			} else if cssutil.IsUpperCaseLetter(b[i]) {
				if b[i]+32 != a[i] {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return true
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

	if cssutil.IsHexDigit(c) {
		// Hex escape sequence: \123456
		var res rune = 0
		hexCount := 0

		for hexCount < 6 && *i < len(runes) && cssutil.IsHexDigit(runes[*i]) {
			res = res*16 + hexDigitToValue(runes[*i])
			*i++
			hexCount++
		}

		// Skip trailing whitespace after hex escape
		if *i < len(runes) && cssutil.IsWhitespace(runes[*i]) {
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
