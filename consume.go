package csslexer

import (
	"strings"

	"go.baoshuo.dev/cssutil"
)

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
func (l *Lexer) consumeEscape() rune {
	var res rune = 0

	next := l.r.Peek(0)

	if cssutil.IsHexDigit(next) {
		l.r.Move(1)
		res = hexDigitToValue(next)

		for i := 1; i < 6; i++ {
			c := l.r.Peek(0)
			if cssutil.IsHexDigit(c) {
				l.r.Move(1)
				res = res*16 + hexDigitToValue(c)
			} else {
				break
			}
		}

		if !isValidCodePoint(res) {
			res = '\uFFFD' // U+FFFD REPLACEMENT CHARACTER
		}

		// If the next input code point is whitespace, consume it as well.
		l.consumeSingleWhitespace()
	} else if next != EOF {
		l.r.Move(1) // consume the escape character
		res = next
	} else {
		res = '\uFFFD' // U+FFFD REPLACEMENT CHARACTER for EOF
	}

	return res
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-name
func (l *Lexer) consumeName() string {
	var result strings.Builder

	for {
		next := l.r.Peek(0)

		if cssutil.IsIdentCodePoint(next) {
			l.r.Move(1)
			result.WriteRune(next)
		} else if twoCharsAreValidEscape(next, l.r.Peek(1)) {
			l.r.Move(1) // consume the backslash
			escaped := l.consumeEscape()
			result.WriteRune(escaped)
		} else {
			break
		}
	}

	return result.String()
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-number
func (l *Lexer) consumeNumber() string {
	offset := l.r.CurrentOffset()

	next := l.r.Peek(0)

	// If the next rune is '+' or '-', consume it as part of the number.
	if next == '+' || next == '-' {
		l.r.Move(1)
	}

	// consume the integer part of the number
	l.r.MoveWhilePredicate(cssutil.IsDigit)

	// float
	next = l.r.Peek(0)
	if next == '.' && cssutil.IsDigit(l.r.Peek(1)) {
		l.r.Move(1) // consume the '.'
		l.r.MoveWhilePredicate(cssutil.IsDigit)
	}

	// scientific notation
	next = l.r.Peek(0)
	if next == 'e' || next == 'E' {
		next_next := l.r.Peek(1)

		if cssutil.IsDigit(next_next) {
			l.r.Move(1) // consume 'e' or 'E'
			l.r.MoveWhilePredicate(cssutil.IsDigit)
		} else if (next_next == '+' || next_next == '-') && cssutil.IsDigit(l.r.Peek(2)) {
			l.r.Move(2) // consume 'e' or 'E' and the sign
			l.r.MoveWhilePredicate(cssutil.IsDigit)
		}
	}

	return l.r.CurrentAfterOffsetString(offset)
}

func (l *Lexer) consumeSingleWhitespace() string {
	offset := l.r.CurrentOffset()

	next := l.r.Peek(0)
	if next == '\r' && l.r.Peek(1) == '\n' {
		l.r.Move(2) // consume CRLF
	} else if cssutil.IsWhitespace(next) {
		l.r.Move(1) // consume the whitespace character
	}

	return l.r.CurrentAfterOffsetString(offset)
}

func (l *Lexer) consumeWhitespace() string {
	offset := l.r.CurrentOffset()

	for {
		next := l.r.Peek(0)

		if cssutil.IsWhitespace(next) {
			l.consumeSingleWhitespace()
		} else if next == EOF {
			break
		} else {
			break
		}
	}

	return l.r.CurrentAfterOffsetString(offset)
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-the-remnants-of-a-bad-url
func (l *Lexer) consumeBadUrlRemnants() {
	for {
		next := l.r.Peek(0)

		if next == ')' {
			l.r.Move(1)
			break
		}
		if next == EOF {
			break
		}

		if twoCharsAreValidEscape(next, l.r.Peek(1)) {
			l.r.Move(1) // consume the backslash
			l.consumeEscape()
			continue
		}

		l.r.Move(1)
	}
}
