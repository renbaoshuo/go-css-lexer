package csslexer

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-comment
func (l *Lexer) consumeUntilCommentEnd() {
	for {
		next := l.r.Peek(0)

		if next == EOF {
			break
		}

		if next == '*' && l.r.Peek(1) == '/' {
			l.r.Move(2) // consume '*/'
			return
		}

		l.r.Move(1) // consume the current rune
	}
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-escaped-code-point
func (l *Lexer) consumeEscape() {
	next := l.r.Peek(0)

	if next == '\\' {
		l.r.Move(1)
		next = l.r.Peek(0)
	}

	if isASCIIHexDigit(next) {
		for i := 1; i < 6; i++ {
			c := l.r.Peek(i)
			if isASCIIHexDigit(c) {
				l.r.Move(1)
			} else {
				break
			}
		}
	}
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-name
func (l *Lexer) consumeName() {
	for {
		next := l.r.Peek(0)

		if isNameCodePoint(next) {
			l.r.Move(1)
		} else if twoCharsAreValidEscape(next, l.r.Peek(1)) {
			l.consumeEscape()
		} else {
			break
		}
	}
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-number
func (l *Lexer) consumeNumber() {
	next := l.r.Peek(0)

	// If the next rune is '+' or '-', consume it as part of the number.
	if next == '+' || next == '-' {
		l.r.Move(1)
	}

	// consume the integer part of the number
	l.r.MoveWhilePredicate(isASCIIDigit)

	// float
	next = l.r.Peek(0)
	if next == '.' && isASCIIDigit(l.r.Peek(1)) {
		l.r.Move(1) // consume the '.'
		l.r.MoveWhilePredicate(isASCIIDigit)
	}

	// scientific notation
	next = l.r.Peek(0)
	if next == 'e' || next == 'E' {
		next_next := l.r.Peek(1)
		if next_next == '+' || next_next == '-' {
			l.r.Move(2) // consume 'e' or 'E' and the sign
			l.r.MoveWhilePredicate(isASCIIDigit)
		} else {
			l.r.Move(1) // consume 'e' or 'E'
			l.r.MoveWhilePredicate(isASCIIDigit)
		}
	}
}

func (l *Lexer) consumeSingleWhitespace() {
	next := l.r.Peek(0)
	if next == '\r' && l.r.Peek(1) == '\n' {
		l.r.Move(2) // consume CRLF
	} else if isHTMLWhitespace(next) {
		l.r.Move(1) // consume the whitespace character
	}
}

func (l *Lexer) consumeWhitespace() {
	for {
		next := l.r.Peek(0)

		if isHTMLWhitespace(next) {
			l.consumeSingleWhitespace()
		} else if next == EOF {
			return
		} else {
			break
		}
	}
}
