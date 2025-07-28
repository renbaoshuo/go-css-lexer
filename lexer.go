package csslexer

// Lexer is the state for the CSS lexer.
type Lexer struct {
	r *Input // The input stream of runes to be lexed.
}

// NewLexer creates a new Lexer instance with the given Input.
func NewLexer(r *Input) *Lexer {
	return &Lexer{
		r: r,
	}
}

// Err returns the error encountered during lexing.
//
// If no error has occurred, it returns nil.
// If the input stream has reached the end, it returns io.EOF.
func (l *Lexer) Err() error {
	return l.r.Err()
}

// Next reads the next token from the input stream.
func (l *Lexer) Next() (TokenType, []rune) {
	switch l.r.Peek(0) {
	case EOF:
		return EOFToken, nil

	case '\t', '\n', '\r', '\f', ' ':
		l.consumeWhitespace()
		return WhitespaceToken, l.r.Shift()

	case '\'', '"':
		return l.consumeStringToken()

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// We consider numbers starts with '+' or '-' in other cases,
		// so we don't handle them here.
		return l.consumeNumericToken()

	case '(':
		l.r.Move(1)
		return LeftParenthesisToken, l.r.Shift()

	case ')':
		l.r.Move(1)
		return RightParenthesisToken, l.r.Shift()

	case '[':
		l.r.Move(1)
		return LeftBracketToken, l.r.Shift()

	case ']':
		l.r.Move(1)
		return RightBracketToken, l.r.Shift()

	case '{':
		l.r.Move(1)
		return LeftBraceToken, l.r.Shift()

	case '}':
		l.r.Move(1)
		return RightBraceToken, l.r.Shift()

	case '+', '.':
		if l.nextCharsAreNumber() {
			return l.consumeNumericToken()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case '-':
		if l.nextCharsAreNumber() {
			return l.consumeNumericToken()
		}
		if l.r.Peek(1) == '-' && l.r.Peek(2) == '>' {
			l.r.Move(3) // consume "-->"
			return CDCToken, l.r.Shift()
		}
		if l.nextCharsAreIdentifier() {
			return l.consumeIdentLikeToken()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case '*':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "*="
			return SubstringMatchToken, l.r.Shift()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case '<':
		if l.r.Peek(1) == '!' && l.r.Peek(2) == '-' && l.r.Peek(3) == '-' {
			l.r.Move(4) // consume "<!--"
			return CDOToken, l.r.Shift()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case ',':
		l.r.Move(1)
		return CommaToken, l.r.Shift()

	case '/':
		if l.r.Peek(1) == '*' {
			l.consumeUntilCommentEnd()
			return CommentToken, l.r.Shift()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case '\\':
		if twoCharsAreValidEscape(l.r.Peek(0), l.r.Peek(1)) {
			return l.consumeIdentLikeToken()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case ':':
		l.r.Move(1)
		return ColonToken, l.r.Shift()

	case ';':
		l.r.Move(1)
		return SemicolonToken, l.r.Shift()

	case '#':
		l.r.Move(1)
		if isNameCodePoint(l.r.Peek(0)) || twoCharsAreValidEscape(l.r.Peek(0), l.r.Peek(1)) {
			l.consumeName()
			return HashToken, l.r.Shift()
		}
		return DelimiterToken, l.r.Shift()

	case '^':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "^="
			return PrefixMatchToken, l.r.Shift()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case '$':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "$="
			return SuffixMatchToken, l.r.Shift()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case '|':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "|="
			return DashMatchToken, l.r.Shift()
		}
		// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#the-column-combinator
		if l.r.Peek(1) == '|' {
			l.r.Move(2) // consume "||"
			return ColumnToken, l.r.Shift()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case '~':
		if l.r.Peek(1) == '=' {
			l.r.Move(2) // consume "~="
			return IncludeMatchToken, l.r.Shift()
		}
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	case '@':
		l.r.Move(1)
		if l.nextCharsAreIdentifier() {
			l.consumeName()
			return AtKeywordToken, l.r.Shift()
		}
		return DelimiterToken, l.r.Shift()

	case 'u', 'U':
		if l.r.Peek(1) == '+' &&
			(isASCIIHexDigit(l.r.Peek(2)) || l.r.Peek(2) == '?') {
			l.r.Move(2) // consume "u+"
			return l.consumeUnicodeRangeToken()
		}
		return l.consumeIdentLikeToken()

	case 1, 2, 3, 4, 5, 6, 7, 8, 11, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27, 28, 29, 30, 31, '!', '%', '&', '=', '>', '?', '`', 127:
		l.r.Move(1)
		return DelimiterToken, l.r.Shift()

	default:
		return l.consumeIdentLikeToken()
	}
}
