package csslexer

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#would-start-an-identifier
func (l *Lexer) nextCharsAreIdentifier() bool {
	first := l.r.Peek(0)
	second := l.r.Peek(1)

	if isNameStartCodePoint(first) {
		return true
	}

	if twoCharsAreValidEscape(first, second) {
		return true
	}

	if first == '-' {
		return isNameStartCodePoint(second) || second == '-' ||
			twoCharsAreValidEscape(second, l.r.Peek(2))
	}

	return false
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#starts-with-a-number
func (l *Lexer) nextCharsAreNumber() bool {
	first := l.r.Peek(0)

	if isASCIIDigit(first) {
		return true
	}

	second := l.r.Peek(1)

	if first == '+' || first == '-' {
		if isASCIIDigit(second) {
			return true
		}

		if second == '.' {
			third := l.r.Peek(2)

			if isASCIIDigit(third) {
				return true
			}
		}
	}

	if first == '.' {
		return isASCIIDigit(second)
	}

	return false
}
