package parser

import (
	"github.com/jkittell/array"
	"github.com/jkittell/ccur"
	"github.com/jkittell/toolbox"
)

type Scan struct {
	Index      int
	StatusCode int
	SegmentURL string
	Err        error
}

func scanSegmentStatusCode(s Scan) Scan {
	statusCode, _, err := toolbox.SendRequest(toolbox.HEAD, s.SegmentURL, "", nil)
	if err != nil {
		s.Err = err
	}
	s.StatusCode = statusCode
	return s
}

func ScanSegments(segments *array.Array[Segment]) *array.Array[Scan] {
	sortedResults := array.New[Scan]()
	results := array.New[Scan]()
	scans := array.New[Scan]()
	for i := 0; i < segments.Length(); i++ {
		seg := segments.Lookup(i)
		scan := Scan{
			Index:      i,
			StatusCode: 0,
			SegmentURL: seg.SegmentURL,
			Err:        nil,
		}
		scans.Push(scan)
		// TODO until I get a sort func on the array
		sortedResults.Push(scan)
	}

	done := make(chan bool)
	defer close(done)
	in := ccur.Source[Scan](done, scans)
	out := ccur.FanOut[Scan](done, in, scanSegmentStatusCode, 100)
	for o := range out {
		results.Push(o)
	}

	for i := 0; i < results.Length(); i++ {
		res := results.Lookup(i)
		sortedResults.Set(res.Index, res)
	}

	return sortedResults
}
