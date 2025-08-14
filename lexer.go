package csslexer

import (
	"go.baoshuo.dev/cssutil"
)

// Lexer is the state for the CSS lexer.
type Lexer struct {
	r *Input // The input stream of runes to be lexed.
	p *Token // The peeked token, nil if no token is peeked.
}

// NewLexer creates a new Lexer instance with the given Input.
func NewLexer(r *Input) *Lexer {
	return &Lexer{
		r: r,
		p: nil,
	}
}

// Peek returns the next token without advancing the position.
// It returns a copy of the token.
func (l *Lexer) Peek() Token {
	if l.p == nil {
		tokenType, data := l.readNextToken()
		// Get a token from the pool
		p := tokenPool.Get().(*Token)
		p.Type = tokenType
		p.Data = data
		l.p = p
	}
	return Token{Type: l.p.Type, Data: l.p.Data}
}

// Next reads the next token from the input stream.
func (l *Lexer) Next() Token {
	if l.p != nil {
		token := Token{Type: l.p.Type, Data: l.p.Data}
		tokenPool.Put(l.p) // Return the token to the pool
		l.p = nil
		return token
	}
	tokenType, data := l.readNextToken()
	return Token{Type: tokenType, Data: data}
}

// readNextToken reads the next token from the input stream.
// This is the internal method that actually parses tokens.
func (l *Lexer) readNextToken() (TokenType, string) {
	switch l.r.Peek(0) {
	case EOF:
		return EOFToken, ""

	case '\t', '\n', '\r', '\f', ' ':
		l.consumeWhitespace()
		return WhitespaceToken, l.r.ShiftString()

	case '\'', '"':
		return l.consumeStringToken()

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// We consider numbers starts with '+' or '-' in other cases,
		// so we don't handle them here.
		return l.consumeNumericToken()

	case '(':
		l.r.Move(1)
		return LeftParenthesisToken, l.r.ShiftString()

	case ')':
		l.r.Move(1)
		return RightParenthesisToken, l.r.ShiftString()

	case '[':
		l.r.Move(1)
		return LeftBracketToken, l.r.ShiftString()

	case ']':
		l.r.Move(1)
		return RightBracketToken, l.r.ShiftString()

	case '{':
		l.r.Move(1)
		return LeftBraceToken, l.r.ShiftString()

	case '}':
		l.r.Move(1)
		return RightBraceToken, l.r.ShiftString()

	case '+', '.':
		if l.nextCharsAreNumber() {
			return l.consumeNumericToken()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case '-':
		if l.nextCharsAreNumber() {
			return l.consumeNumericToken()
		}
		if l.r.Peek(1) == '-' && l.r.Peek(2) == '>' {
			l.r.Move(3) // consume "-->"
			return CDCToken, l.r.ShiftString()
		}
		if l.nextCharsAreIdentifier() {
			return l.consumeIdentLikeToken()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case '*':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "*="
			return SubstringMatchToken, l.r.ShiftString()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case '<':
		if l.r.Peek(1) == '!' && l.r.Peek(2) == '-' && l.r.Peek(3) == '-' {
			l.r.Move(4) // consume "<!--"
			return CDOToken, l.r.ShiftString()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case ',':
		l.r.Move(1)
		return CommaToken, l.r.ShiftString()

	case '/':
		if l.r.Peek(1) == '*' {
			l.consumeUntilCommentEnd()
			return CommentToken, l.r.ShiftString()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case '\\':
		if twoCharsAreValidEscape(l.r.Peek(0), l.r.Peek(1)) {
			return l.consumeIdentLikeToken()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case ':':
		l.r.Move(1)
		return ColonToken, l.r.ShiftString()

	case ';':
		l.r.Move(1)
		return SemicolonToken, l.r.ShiftString()

	case '#':
		l.r.Move(1)
		if cssutil.IsIdentCodePoint(l.r.Peek(0)) || twoCharsAreValidEscape(l.r.Peek(0), l.r.Peek(1)) {
			name := l.consumeName()
			l.r.Shift() // Shift the input after consuming the name
			return HashToken, name
		}
		return DelimiterToken, l.r.ShiftString()

	case '^':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "^="
			return PrefixMatchToken, l.r.ShiftString()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case '$':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "$="
			return SuffixMatchToken, l.r.ShiftString()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case '|':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "|="
			return DashMatchToken, l.r.ShiftString()
		}
		// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#the-column-combinator
		if l.r.Peek(1) == '|' {
			l.r.Move(2) // consume "||"
			return ColumnToken, l.r.ShiftString()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case '~':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "~="
			return IncludeMatchToken, l.r.ShiftString()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	case '@':
		l.r.Move(1)
		if l.nextCharsAreIdentifier() {
			name := l.consumeName()
			l.r.Shift() // Shift the input after consuming the name
			return AtKeywordToken, name
		}
		return DelimiterToken, l.r.ShiftString()

	case 'u', 'U':
		if l.r.Peek(1) == '+' &&
			(cssutil.IsHexDigit(l.r.Peek(2)) || l.r.Peek(2) == '?') {
			l.r.Move(2) // consume "u+"
			return l.consumeUnicodeRangeToken()
		}
		return l.consumeIdentLikeToken()

	case 1, 2, 3, 4, 5, 6, 7, 8, 11, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27, 28, 29, 30, 31, '!', '%', '&', '=', '>', '?', '`', 127:
		l.r.Move(1)
		return DelimiterToken, l.r.ShiftString()

	default:
		return l.consumeIdentLikeToken()
	}
}
