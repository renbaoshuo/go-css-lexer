package csslexer

import (
	"go.baoshuo.dev/cssutil"
)

func isValidCodePoint(c rune) bool {
	// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-escaped-code-point
	// If this number is zero, or is for a surrogate, or is greater than the maximum allowed code point, return U+FFFD REPLACEMENT CHARACTER
	return c != 0 && !cssutil.IsSurrogate(c) && cssutil.IsLowerThanMaxCodePoint(c)
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
