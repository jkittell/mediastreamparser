package mediastreamparser

import (
	"fmt"
	"testing"
)

func TestParser_Example(t *testing.T) {
	url := "http://devimages.apple.com/iphone/samples/bipbop/bipbopall.m3u8"
	scanner := New()
	results, err := scanner.Scan(url, HLS)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for i := 0; i < results.Length(); i++ {
		data, _ := results.Lookup(i).MarshallJSON()
		fmt.Println(string(data))
	}

	fmt.Println(results.Length(), "number of results")
}
