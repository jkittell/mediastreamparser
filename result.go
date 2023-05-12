package mediastreamparser

import "encoding/json"

type Result struct {
	// master playlist or manifest URL
	playlistURL string

	// Stream Info
	streamName string
	streamURL  string

	// Segment Info
	segmentName    string
	segmentURL     string
	byteRangeStart int
	byteRangeSize  int
}

func (r Result) MarshallJSON() ([]byte, error) {
	return json.Marshal(struct {
		PlaylistURL    string `json:"playlist_url"`
		StreamName     string `json:"stream_name"`
		StreamURL      string `json:"stream_url"`
		SegmentName    string `json:"segment_name"`
		SegmentURL     string `json:"segment_url"`
		ByteRangeStart int    `json:"byte_range_start"`
		ByteRangeSize  int    `json:"byte_range_size"`
	}{
		PlaylistURL:    r.playlistURL,
		StreamName:     r.streamName,
		StreamURL:      r.streamURL,
		SegmentName:    r.segmentName,
		SegmentURL:     r.segmentURL,
		ByteRangeStart: r.byteRangeStart,
		ByteRangeSize:  r.byteRangeSize,
	})
}

func (r Result) PlaylistURL() string {
	return r.playlistURL
}

func (r Result) StreamName() string {
	return r.streamName
}

func (r Result) StreamURL() string {
	return r.streamURL
}

func (r Result) SegmentName() string {
	return r.segmentName
}

func (r Result) SegmentURL() string {
	return r.segmentURL
}

func (r Result) ByteRangeStart() int {
	return r.byteRangeStart
}

func (r Result) ByteRangeSize() int {
	return r.byteRangeSize
}
