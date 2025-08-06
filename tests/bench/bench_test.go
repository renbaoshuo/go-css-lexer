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
					tokenType, _ := lexer.Next()
					if tokenType == csslexer.EOFToken {
						break
					}
				}
				if err := lexer.Err(); err != nil && err != io.EOF {
					b.Fatalf("lexer error: %v", err)
				}
			}
			b.StopTimer()

			totalBytes := len(ungzip(dataGz)) * b.N
			totalMiB := totalBytes / 1024 / 1024
			b.ReportMetric(float64(totalMiB)/b.Elapsed().Seconds(), "MiB/s")

		})
	}
}
