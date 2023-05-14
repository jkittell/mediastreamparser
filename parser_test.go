package mediastreamparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func testURLs() []string {
	var urls []string
	file, err := os.Open("urls.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return urls
}

func TestParser_Run(t *testing.T) {
	parser := New()
	for n, url := range testURLs() {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			var format ContentFormat
			if strings.Contains(url, ".m3u8") {
				format = HLS
			}
			if strings.Contains(url, ".mpd") {
				format = DASH
			}
			results, err := parser.Run(url, format)
			if err != nil {
				fmt.Println(err)
				t.FailNow()
			}
			fmt.Println(results.Length(), "number of results")
			if results.Length() == 0 {
				t.Fail()
			}
		})
	}
}

func BenchmarkParser_Run(b *testing.B) {
	parser := New()
	for n, url := range testURLs() {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var format ContentFormat
				if strings.Contains(url, ".m3u8") {
					format = HLS
				}
				if strings.Contains(url, ".mpd") {
					format = DASH
				}
				parser.Run(url, format)
			}
		})
	}
}
