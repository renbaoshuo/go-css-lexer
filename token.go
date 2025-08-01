package csslexer

import (
	"fmt"
)

// TokenType represents the type of a token in the CSS lexer.
type TokenType int

const (
	// DefaultToken is the default token type, used when no specific type is matched.
	// It is not being used in the lexer.
	DefaultToken TokenType = iota

	// Standard CSS token types
	IdentToken
	FunctionToken
	AtKeywordToken
	HashToken
	StringToken
	BadStringToken
	UrlToken
	BadUrlToken
	DelimiterToken
	NumberToken
	PercentageToken
	DimensionToken
	WhitespaceToken
	CDOToken
	CDCToken
	ColonToken
	SemicolonToken
	CommaToken
	LeftParenthesisToken
	RightParenthesisToken
	LeftBracketToken
	RightBracketToken
	LeftBraceToken
	RightBraceToken
	EOFToken

	// Additional CSS token types
	CommentToken
	IncludeMatchToken
	DashMatchToken
	PrefixMatchToken
	SuffixMatchToken
	SubstringMatchToken
	ColumnToken
	UnicodeRangeToken
)

func (tt TokenType) String() string {
	switch tt {
	case DefaultToken:
		return "Default"

	case IdentToken:
		return "Ident"
	case FunctionToken:
		return "Function"
	case AtKeywordToken:
		return "AtKeyword"
	case HashToken:
		return "Hash"
	case StringToken:
		return "String"
	case BadStringToken:
		return "BadString"
	case UrlToken:
		return "Url"
	case BadUrlToken:
		return "BadUrl"
	case DelimiterToken:
		return "Delimiter"
	case NumberToken:
		return "Number"
	case PercentageToken:
		return "Percentage"
	case DimensionToken:
		return "Dimension"
	case WhitespaceToken:
		return "Whitespace"
	case CDOToken:
		return "CDO"
	case CDCToken:
		return "CDC"
	case ColonToken:
		return "Colon"
	case SemicolonToken:
		return "Semicolon"
	case CommaToken:
		return "Comma"
	case LeftParenthesisToken:
		return "LeftParenthesis"
	case RightParenthesisToken:
		return "RightParenthesis"
	case LeftBracketToken:
		return "LeftBracket"
	case RightBracketToken:
		return "RightBracket"
	case LeftBraceToken:
		return "LeftBrace"
	case RightBraceToken:
		return "RightBrace"
	case EOFToken:
		return "EOF"

	case CommentToken:
		return "Comment"
	case IncludeMatchToken:
		return "IncludeMatch"
	case DashMatchToken:
		return "DashMatch"
	case PrefixMatchToken:
		return "PrefixMatch"
	case SuffixMatchToken:
		return "SuffixMatch"
	case SubstringMatchToken:
		return "SubstringMatch"
	case ColumnToken:
		return "Column"
	case UnicodeRangeToken:
		return "UnicodeRange"

	default:
		return fmt.Sprintf("Unknown(%d)", tt)
	}
}
