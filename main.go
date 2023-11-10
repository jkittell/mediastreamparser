package main

import (
	"flag"
	"github.com/fatih/color"
	"github.com/jkittell/array"
	"github.com/jkittell/mediastreamparser/parser"
	"github.com/jkittell/toolbox"
	"log"
)

func scanSegments(segments *array.Array[parser.Segment]) {
	for i := 0; i < segments.Length(); i++ {
		s := segments.Lookup(i)
		statusCode, _, err := toolbox.SendRequest(toolbox.HEAD, s.SegmentURL, "", nil)
		if err != nil {
			color.Red("%s,%s", err.Error(), s.SegmentURL)
		}
		if statusCode == 200 {
			color.Green("%d,%s\n", statusCode, s.SegmentURL)
		} else if statusCode == 500 {
			color.Red("%d,%s\n", statusCode, s.SegmentURL)
		} else {
			color.Yellow("%d,%s\n", statusCode, s.SegmentURL)
		}
	}
}

func main() {
	url := flag.String("url", "", "url of hls playlist or dash manifest")
	scan := flag.Bool("scan", false, "scan for segments")
	info := flag.Bool("info", false, "segment info")
	flag.Parse()

	if *scan {
		segments, err := parser.GetSegments(*url)
		if err != nil {
			log.Println(err)
		}

		scanSegments(segments)
	}

	if *info {
		segments, err := parser.GetSegments(*url)
		if err != nil {
			log.Println(err)
		}

		for i := 0; i < segments.Length(); i++ {
			s := segments.Lookup(i)
			b, _ := s.JSON()
			color.Magenta(string(b))
		}
	}
}
