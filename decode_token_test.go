package csslexer

import (
	"testing"
)

func TestDecodeToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenType TokenType
		input     string
		expect    string
	}{
		{
			name:      "ident with escape",
			tokenType: IdentToken,
			input:     `foo\20 bar`,
			expect:    `foo bar`,
		},
		{
			name:      "string with escape",
			tokenType: StringToken,
			input:     `"foo\22 bar"`,
			expect:    `foo"bar`,
		},
		{
			name:      "url with escape in prefix and content",
			tokenType: UrlToken,
			input:     `u\72 l(foo\20 bar)`,
			expect:    `url(foo bar)`,
		},
		{
			name:      "hash with escape",
			tokenType: HashToken,
			input:     `#foo\20 bar`,
			expect:    `#foo bar`,
		},
		{
			name:      "dimension with escape",
			tokenType: DimensionToken,
			input:     `10p\78`,
			expect:    `10px`,
		},
		{
			name:      "at-keyword with escape",
			tokenType: AtKeywordToken,
			input:     `@f\6fobar`,
			expect:    `@foobar`,
		},
		{
			name:      "function with escape",
			tokenType: FunctionToken,
			input:     `\63 alc(`,
			expect:    `calc(`,
		},
		{
			name:      "default type",
			tokenType: DefaultToken,
			input:     "foo",
			expect:    "foo",
		},
	}

	for _, tc := range tests {
		raw := []rune(tc.input)
		got := DecodeToken(tc.tokenType, raw)
		if got != tc.expect {
			t.Errorf("%s: got '%s', want '%s'", tc.name, got, tc.expect)
		}
	}
}
