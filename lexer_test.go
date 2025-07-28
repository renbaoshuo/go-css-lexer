package csslexer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const (
	testDataDir    = "tests"
	sourceCssFile  = "source.css"
	tokensJsonFile = "tokens.json"
)

type testToken struct {
	Type string `json:"type"`
	Raw  string `json:"raw"`
}

func TestLexer(t *testing.T) {
	// ensure the test data directory exists
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		t.Fatalf("test data directory does not exist: %s", testDataDir)
	}

	categoryDirs, err := os.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("failed to read test data directory: %v", err)
	}

	for _, categoryDir := range categoryDirs {
		if !categoryDir.IsDir() {
			continue
		}

		testCategory := categoryDir.Name()
		t.Run(testCategory, func(t *testing.T) {
			categoryIdDirs, err := os.ReadDir(filepath.Join(testDataDir, testCategory))
			if err != nil {
				t.Fatalf("failed to read test category directory: %v", err)
			}

			for _, categoryIdDir := range categoryIdDirs {
				if !categoryIdDir.IsDir() {
					continue
				}

				testId := categoryIdDir.Name()
				t.Run(testId, func(t *testing.T) {
					testPath := filepath.Join(testDataDir, testCategory, testId)
					sourceFile := filepath.Join(testPath, sourceCssFile)
					tokensFile := filepath.Join(testPath, tokensJsonFile)

					sources, err := os.ReadFile(sourceFile)
					if err != nil {
						t.Fatalf("failed to read test source file: %v", err)
					}

					tokensRaw, err := os.ReadFile(tokensFile)
					if err != nil {
						t.Fatalf("failed to read test tokens file: %v", err)
					}
					var tokens []testToken
					if err := json.Unmarshal(tokensRaw, &tokens); err != nil {
						t.Fatalf("failed to unmarshal tokens: %v", err)
					}

					lexer := NewLexer(NewInputBytes(sources))

					for i := 0; i < len(tokens); i++ {
						expectedToken := tokens[i]
						tokenType, tokenRaw := lexer.Next()

						// t.Logf("Expect token %d: Type=%s, Raw=%s", i, expectedToken.Type, expectedToken.Raw)
						// t.Logf("Lexer returned: Type=%s, Raw=%s", tokenType, string(tokenRaw))

						if tokenType != convertTestTokenName(expectedToken.Type) {
							t.Errorf("expected token type '%s' (raw: '%s'), got '%s' (raw: '%s') at index %d in test %s/%s",
								expectedToken.Type, expectedToken.Raw, tokenType, string(tokenRaw), i, testCategory, testId)
						}

						if string(tokenRaw) != expectedToken.Raw {
							t.Errorf("expected '%s' token raw '%s', got '%s' at index %d in test %s/%s",
								expectedToken.Type, expectedToken.Raw, string(tokenRaw), i, testCategory, testId)
						}
					}

					if tokenType, _ := lexer.Next(); tokenType != EOFToken {
						t.Errorf("expected EOF token, got '%s' at the end of test %s/%s", tokenType, testCategory, testId)
					}
				})
			}
		})
	}
}

func convertTestTokenName(name string) TokenType {
	switch name {
	case "ident-token":
		return IdentToken
	case "function-token":
		return FunctionToken
	case "at-keyword-token":
		return AtKeywordToken
	case "hash-token":
		return HashToken
	case "string-token":
		return StringToken
	case "bad-string-token":
		return BadStringToken
	case "url-token":
		return UrlToken
	case "bad-url-token":
		return BadUrlToken
	case "delim-token":
		return DelimiterToken
	case "number-token":
		return NumberToken
	case "percentage-token":
		return PercentageToken
	case "dimension-token":
		return DimensionToken
	case "whitespace-token":
		return WhitespaceToken
	case "CDO-token":
		return CDOToken
	case "CDC-token":
		return CDCToken
	case "colon-token":
		return ColonToken
	case "semicolon-token":
		return SemicolonToken
	case "comma-token":
		return CommaToken
	case "(-token":
		return LeftParenthesisToken
	case ")-token":
		return RightParenthesisToken
	case "[-token":
		return LeftBracketToken
	case "]-token":
		return RightBracketToken
	case "{-token":
		return LeftBraceToken
	case "}-token":
		return RightBraceToken
	case "":
		return EOFToken

	// Additional CSS token types
	case "comment":
		return CommentToken
	// case "":
	// 	return IncludeMatchToken
	// case "":
	// 	return DashMatchToken
	// case "":
	// 	return PrefixMatchToken
	// case "":
	// 	return SuffixMatchToken
	// case "":
	// 	return SubstringMatchToken
	// case "":
	// 	return ColumnToken
	// case "":
	// 	return UnicodeRangeToken

	default:
		return DefaultToken
	}
}
