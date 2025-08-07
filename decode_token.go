package csslexer

import (
	"slices"
	"strings"
)

// DecodeToken decodes a CSS token based on its type and raw rune slice.
// It handles escape sequences for identifier-like tokens, strings, and URLs.
//
// It returns the decoded string representation of the token.
//
// NOTE: 1) For string-token, it removes surrounding quotes. 2) For url-token,
// it removes whitespace around the content inside parentheses.
func DecodeToken(tokenType TokenType, raw []rune) string {
	switch tokenType {
	case IdentToken, FunctionToken, AtKeywordToken, HashToken, DimensionToken:
		return decodeIdentLikeToken(raw)
	case StringToken:
		return decodeStringToken(raw)
	case UrlToken:
		return decodeUrlToken(raw)
	default:
		return string(raw)
	}
}

// decodeIdentLikeToken decodes escape sequences in identifier-like tokens.
//
// This includes ident, function, at-keyword, hash, and dimension tokens
func decodeIdentLikeToken(raw []rune) string {
	if len(raw) == 0 {
		return ""
	}
	return decodeEscapeSequences(raw)
}

// decodeStringToken decodes escape sequences in string tokens.
//
// NOTE: It removes surrounding quotes of the string.
func decodeStringToken(raw []rune) string {
	if len(raw) < 2 {
		return string(raw)
	}

	// Remove quotes and decode escape sequences
	content := raw[1 : len(raw)-1] // Remove surrounding quotes
	return decodeEscapeSequences(content)
}

// decodeUrlToken decodes escape sequences in URL tokens.
//
// NOTE: It removes the whitespace around the content inside parentheses.
func decodeUrlToken(raw []rune) string {
	if len(raw) < 2 {
		return string(raw)
	}

	// Find the first parenthesis
	parenIdx := slices.Index(raw, '(')
	if parenIdx == -1 {
		// No parenthesis found, decode the entire raw
		return decodeEscapeSequences(raw)
	}

	prefix := decodeEscapeSequences(raw[:parenIdx])

	// the content inside the parentheses
	start := parenIdx + 1
	end := len(raw)
	if end > 0 && raw[end-1] == ')' {
		end--
	}

	// Skip leading and trailing whitespace
	for start < end && isHTMLWhitespace(raw[start]) {
		start++
	}
	for end > start && isHTMLWhitespace(raw[end-1]) {
		end--
	}

	content := raw[start:end]
	contentStr := decodeEscapeSequences(content)

	return prefix + "(" + contentStr + ")"
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
