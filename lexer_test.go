package csslexer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
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

// CategoryStats holds test statistics for a category
type CategoryStats struct {
	Total  int
	Passed int
	Tests  map[string]bool // Map of test IDs to pass/fail status
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

	// Statistics for test results
	categoryStats := make(map[string]*CategoryStats)
	totalTests := 0
	totalPassed := 0

	for _, categoryDir := range categoryDirs {
		if !categoryDir.IsDir() {
			continue
		}

		testCategory := categoryDir.Name()
		if testCategory == "fuzz" { // Skip fuzz test directory
			continue
		}

		// Initialize category statistics
		categoryStats[testCategory] = &CategoryStats{
			Tests: make(map[string]bool),
		}

		// Run tests for this category
		t.Run(testCategory, func(t *testing.T) {
			categoryIdDirs, err := os.ReadDir(filepath.Join(testDataDir, testCategory))
			if err != nil {
				t.Fatalf("failed to read test category directory: %v", err)
			}

			stats := categoryStats[testCategory]
			stats.Total = len(categoryIdDirs)
			totalTests += len(categoryIdDirs)

			for _, categoryIdDir := range categoryIdDirs {
				if !categoryIdDir.IsDir() {
					continue
				}

				testId := categoryIdDir.Name()
				passed := t.Run(testId, func(t *testing.T) {
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
						token := lexer.Next()

						// t.Logf("Expect token %d: Type=%s, Raw=%q", i, expectedToken.Type, expectedToken.Raw)
						// t.Logf("Lexer returned: Type=%s, Raw=%q", tokenType, string(tokenRaw))

						if token.Type != convertTestTokenName(expectedToken.Type) {
							t.Errorf("expected token type '%s' (raw: %q), got '%s' (raw: %q) at index %d",
								expectedToken.Type, expectedToken.Raw, token.Type, string(token.Data), i)
						}

						if string(token.Data) != expectedToken.Raw {
							t.Errorf("expected '%s' token raw %q, got %q at index %d",
								expectedToken.Type, expectedToken.Raw, string(token.Data), i)
						}
					}

					if token := lexer.Next(); token.Type != EOFToken {
						t.Errorf("expected EOF token, got '%s' at the end of test", token.Type)
					}
				})

				stats.Tests[testId] = passed
				if passed {
					stats.Passed++
					totalPassed++
				}
			}
		})
	}

	generateTestReport(totalTests, totalPassed, categoryStats)

	fmt.Println("===== Test Results Summary =====")
	fmt.Printf("Total tests: %d, Passed: %d, Pass rate: %.2f%%\n", totalTests, totalPassed, float64(totalPassed)/float64(totalTests)*100)
	fmt.Println("Results have been written to test_result.md")
}

func generateTestReport(totalTests, totalPassed int, categoryStats map[string]*CategoryStats) {
	passRate := float64(totalPassed) / float64(totalTests) * 100

	var content strings.Builder

	content.WriteString("# CSS Lexer Test Results\n\n")
	content.WriteString("## Summary\n\n")
	content.WriteString(fmt.Sprintf("- **Date**: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("- **Total Tests**: %d\n", totalTests))
	content.WriteString(fmt.Sprintf("- **Passed Tests**: %d\n", totalPassed))
	content.WriteString(fmt.Sprintf("- **Pass Rate**: %.2f%%\n\n", passRate))

	content.WriteString("## Test Results by Category\n\n")
	content.WriteString("| Category | Total | Passed | Pass Rate |\n")
	content.WriteString("|----------|-------|--------|----------|\n")

	categories := make([]string, 0, len(categoryStats))
	for category := range categoryStats {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	for _, category := range categories {
		stats := categoryStats[category]
		categoryPassRate := float64(stats.Passed) / float64(stats.Total) * 100
		content.WriteString(fmt.Sprintf("| %s | %d | %d | %.2f%% |\n",
			category, stats.Total, stats.Passed, categoryPassRate))
	}

	content.WriteString("\n## Detailed Test Results\n\n")
	for _, category := range categories {
		stats := categoryStats[category]
		content.WriteString(fmt.Sprintf("### %s\n\n", category))
		content.WriteString("| Test ID | Status |\n")
		content.WriteString("|---------|--------|\n")

		testIDs := make([]string, 0, len(stats.Tests))
		for testID := range stats.Tests {
			testIDs = append(testIDs, testID)
		}
		sort.Strings(testIDs)

		for _, testID := range testIDs {
			status := "❌ Failed"
			if stats.Tests[testID] {
				status = "✅ Passed"
			}
			content.WriteString(fmt.Sprintf("| %s | %s |\n", testID, status))
		}
		content.WriteString("\n")
	}

	err := os.WriteFile("test_result.md", []byte(content.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing test results: %v\n", err)
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
