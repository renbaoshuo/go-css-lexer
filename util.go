package csslexer

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#non-ascii-code-point
func isASCII(c rune) bool {
	return c < 128
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#letter
func isASCIIAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#digit
func isASCIIDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#hex-digit
func isASCIIHexDigit(c rune) bool {
	return isASCIIDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#newline
func isCSSNewline(c rune) bool {
	return c == '\n' || c == '\r' || c == '\f'
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#ident-start-code-point
func isNameStartCodePoint(r rune) bool {
	return isASCIIAlpha(r) || r == '_' || !isASCII(r)
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#ident-code-point
func isNameCodePoint(r rune) bool {
	return isNameStartCodePoint(r) || isASCIIDigit(r) || r == '-'
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#non-printable-code-point
func isNonPrintableCodePoint(r rune) bool {
	return (r >= 0x00 && r <= 0x08) || r == 0x0B || (r >= 0x0E && r <= 0x1F) || r == 0x7F
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#starts-with-a-valid-escape
func twoCharsAreValidEscape(first, second rune) bool {
	return first == '\\' && !isCSSNewline(second)
}

func isHTMLSpecialWhitespace(c rune) bool {
	return c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

func isHTMLWhitespace(c rune) bool {
	return c == ' ' || isHTMLSpecialWhitespace(c)
}

func equalIgnoringASCIICase(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			continue
		} else if isASCII(a[i]) && isASCII(b[i]) {
			if a[i] >= 'A' && a[i] <= 'Z' {
				if a[i]+32 != b[i] {
					return false
				}
			} else if b[i] >= 'A' && b[i] <= 'Z' {
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
