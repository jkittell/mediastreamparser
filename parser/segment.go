package parser

import (
	"bytes"
	"encoding/json"
)

type Segment struct {
	PlaylistURL    string `json:"playlist_url"`
	StreamName     string `json:"stream_name"`
	StreamURL      string `json:"stream_url"`
	SegmentName    string `json:"segment_name"`
	SegmentURL     string `json:"segment_url"`
	ByteRangeStart int    `json:"byte_range_start"`
	ByteRangeSize  int    `json:"byte_range_size"`
}

func (s *Segment) JSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(s)
	return buffer.Bytes(), err
}
