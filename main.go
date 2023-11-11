package main

import (
	"flag"
	"github.com/fatih/color"
	"github.com/jkittell/mediastreamparser/parser"
	"log"
)

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

		scans := parser.ScanSegments(segments)
		for i := 0; i < scans.Length(); i++ {
			s := scans.Lookup(i)
			if s.StatusCode == 200 {
				color.Green("%d,%s\n", s.StatusCode, s.SegmentURL)
			} else if s.StatusCode == 500 {
				color.Red("%d,%s\n", s.StatusCode, s.SegmentURL)
			} else {
				color.Yellow("%d,%s\n", s.StatusCode, s.SegmentURL)
			}
		}
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
