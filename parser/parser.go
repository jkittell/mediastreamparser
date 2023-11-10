package parser

import (
	"errors"
	"fmt"
	"github.com/jkittell/array"
	"strings"
)

func GetSegments(url string) (*array.Array[Segment], error) {
	if strings.Contains(url, "m3u8") {
		return parseHLS(url)
	} else if strings.Contains(url, "mpd") {
		return parseDASH(url)
	} else {
		return array.New[Segment](), errors.New(fmt.Sprintf("unable to parse %s", url))
	}
}
