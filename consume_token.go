package csslexer

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-numeric-token
func (l *Lexer) consumeNumericToken() (TokenType, []rune) {
	l.consumeNumber()

	if l.nextCharsAreIdentifier() {
		l.consumeName()
		return DimensionToken, l.r.Shift()
	} else if l.r.Peek(0) == '%' {
		l.r.Move(1) // consume '%'
		return PercentageToken, l.r.Shift()
	}

	return NumberToken, l.r.Shift()
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#urange
func (l *Lexer) consumeUnicodeRangeToken() (TokenType, []rune) {
	// range start
	start_length_remaining := 6
	for next := l.r.Peek(0); start_length_remaining > 0 && next != EOF && isASCIIHexDigit(next); next = l.r.Peek(0) {
		l.r.Move(1) // consume the hex digit
		start_length_remaining--
	}

	if start_length_remaining > 0 && l.r.Peek(0) == '?' { // wildcard range
		for start_length_remaining > 0 && l.r.Peek(0) == '?' {
			l.r.Move(1) // consume the '?'
			start_length_remaining--
		}
	} else if l.r.Peek(0) == '-' && isASCIIHexDigit(l.r.Peek(1)) { // range end
		l.r.Move(1) // consume the '-'

		end_length_remaining := 6
		for next := l.r.Peek(0); end_length_remaining > 0 && next != EOF && isASCIIHexDigit(next); next = l.r.Peek(0) {
			l.r.Move(1) // consume the hex digit
			end_length_remaining--
		}
	}

	return UnicodeRangeToken, l.r.Shift()
}

var urlRunes = []rune{'u', 'r', 'l'}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-ident-like-token
func (l *Lexer) consumeIdentLikeToken() (TokenType, []rune) {
	name := l.consumeName()

	if l.r.Peek(0) == '(' {
		l.r.Move(1) // consume the opening parenthesis
		if equalIgnoringASCIICase(name, urlRunes) {
			// The spec is slightly different so as to avoid dropping whitespace
			// tokens, but they wouldn't be used and this is easier.
			l.consumeWhitespace()

			next := l.r.Peek(0)
			if next != '"' && next != '\'' {
				return l.consumeURLToken()
			}
		}

		return FunctionToken, l.r.Shift()
	}

	return IdentToken, l.r.Shift()
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-string-token
func (l *Lexer) consumeStringToken() (TokenType, []rune) {
	until := l.r.Peek(0) // the opening quote, already checked valid by the caller
	l.r.Move(1)

	for {
		next := l.r.Peek(0)

		if next == until {
			l.r.Move(1)
			return StringToken, l.r.Shift()
		}

		if next == EOF {
			return StringToken, l.r.Shift()
		}

		if isCSSNewline(next) {
			return BadStringToken, l.r.Shift()
		}

		if next == '\\' {
			next_next := l.r.Peek(1)

			if next_next == EOF {
				l.r.Move(1) // consume the backslash
				continue
			}

			if isCSSNewline(next_next) {
				l.r.Move(1)
				l.consumeSingleWhitespace()
			} else if twoCharsAreValidEscape(next, next_next) {
				l.r.Move(1) // consume the backslash
				l.consumeEscape()
			} else {
				l.r.Move(1)
			}
		} else {
			l.r.Move(1) // consume the current rune
		}
	}
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-url-token
func (l *Lexer) consumeURLToken() (TokenType, []rune) {
	for {
		next := l.r.Peek(0)

		if next == ')' {
			l.r.Move(1)
			return UrlToken, l.r.Shift()
		}

		if next == EOF {
			return UrlToken, l.r.Shift()
		}

		if isHTMLWhitespace(next) {
			l.consumeWhitespace()

			next_next := l.r.Peek(0)
			if next_next == ')' {
				l.r.Move(1) // consume the closing parenthesis
				return UrlToken, l.r.Shift()
			}
			if next_next == EOF {
				return UrlToken, l.r.Shift()
			}

			// If the next character is not a closing parenthesis, there's an error and we should mark it as a bad URL token.
			break
		}

		if next == '"' || next == '\'' || isNonPrintableCodePoint(next) {
			l.r.Move(1) // consume the invalid character
			break
		}

		if next == '\\' {
			if twoCharsAreValidEscape(next, l.r.Peek(1)) {
				l.r.Move(1) // consume the backslash
				l.consumeEscape()
				continue
			} else {
				break
			}
		}

		l.r.Move(1) // consume the current rune
	}

	l.consumeBadUrlRemnants()
	return BadUrlToken, l.r.Shift()
}
