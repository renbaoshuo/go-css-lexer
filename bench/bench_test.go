package bench

import (
	"bytes"
	"compress/gzip"
	"embed"
	"io"
	"testing"

	"go.baoshuo.dev/csslexer"
)

// the testdata is copied from
// https://github.com/GoalSmashers/css-minification-benchmark/tree/master/data
//
//go:embed testdata
var fs embed.FS

func BenchmarkLexer(b *testing.B) {
	files, err := fs.ReadDir("testdata")
	if err != nil {
		b.Fatalf("failed to read testdata directory: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		b.Run(file.Name(), func(b *testing.B) {
			dataGz, err := fs.ReadFile("testdata/" + file.Name())
			if err != nil {
				b.Fatalf("failed to read file %s: %v", file.Name(), err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				input := csslexer.NewInputBytes(ungzip(dataGz))
				lexer := csslexer.NewLexer(input)

				b.StartTimer()
				for {
					token := lexer.Next()
					if token.Type == csslexer.EOFToken {
						break
					}
				}
			}
			b.StopTimer()

			totalBytes := len(ungzip(dataGz)) * b.N
			totalMiB := totalBytes / 1024 / 1024
			b.ReportMetric(float64(totalMiB)/b.Elapsed().Seconds(), "MiB/s")
		})
	}
}

func BenchmarkDecodeToken(b *testing.B) {
	files, err := fs.ReadDir("testdata")
	if err != nil {
		b.Fatalf("failed to read testdata directory: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		b.Run(file.Name(), func(b *testing.B) {
			dataGz, err := fs.ReadFile("testdata/" + file.Name())
			if err != nil {
				b.Fatalf("failed to read file %s: %v", file.Name(), err)
			}
			data := ungzip(dataGz)

			// 预解析所有 token，不计时
			input := csslexer.NewInputBytes(data)
			lexer := csslexer.NewLexer(input)
			var tokens []csslexer.Token
			for {
				t := lexer.Next()
				if t.Type == csslexer.EOFToken {
					break
				}
				tokens = append(tokens, t)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, t := range tokens {
					_ = t.DecodeData()
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkLexerWithPeek(b *testing.B) {
	files, err := fs.ReadDir("testdata")
	if err != nil {
		b.Fatalf("failed to read testdata directory: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		b.Run(file.Name(), func(b *testing.B) {
			dataGz, err := fs.ReadFile("testdata/" + file.Name())
			if err != nil {
				b.Fatalf("failed to read file %s: %v", file.Name(), err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				input := csslexer.NewInputBytes(ungzip(dataGz))
				lexer := csslexer.NewLexer(input)

				b.StartTimer()
				for {
					_ = lexer.Peek()
					token := lexer.Next()
					if token.Type == csslexer.EOFToken {
						break
					}
				}
			}
			b.StopTimer()

			totalBytes := len(ungzip(dataGz)) * b.N
			totalMiB := totalBytes / 1024 / 1024
			b.ReportMetric(float64(totalMiB)/b.Elapsed().Seconds(), "MiB/s")
		})
	}
}

func ungzip(gz []byte) []byte {
	reader, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		return nil
	}
	defer reader.Close()

	var out bytes.Buffer
	if _, err := io.Copy(&out, reader); err != nil {
		return nil
	}
	return out.Bytes()
}
