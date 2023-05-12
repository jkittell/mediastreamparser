package mediastreamparser

import (
	"fmt"
	"github.com/jkittell/array"
	"github.com/jkittell/toolbox"
	"github.com/unki2aut/go-mpd"
	"regexp"
)

/*
<?xml version="1.0" encoding="UTF-8"?>
<MPD xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="urn:mpeg:dash:schema:mpd:2011" xsi:schemaLocation="urn:mpeg:dash:schema:mpd:2011 http://standards.iso.org/ittf/PubliclyAvailableStandards/MPEG-DASH_schema_files/DASH-MPD.xsd" availabilityStartTime="2023-03-23T17:20:36Z" type="dynamic" mediaPresentationDuration="PT0H15M10S" publishTime="2023-03-23T17:35:17Z" maxSegmentDuration="PT10S" minBufferTime="PT10S" profiles="urn:scte:dash:2015#ts">
    <Period id="1" start="PT0S">
        <AdaptationSet id="1" mimeType="video/mp2t" segmentAlignment="true" bitStreamSwitching="true" lang="und" maxWidth="960" maxHeight="540" maxFrameRate="29">
            <Representation id="Stream1-1" audioSamplingRate="48000" codecs="avc1.64001F,ac-3,mp4a.40.2" width="960" height="540" frameRate="29" sar="16:9" startWithSAP="1" bandwidth="3199008">
                <SegmentTemplate timescale="90000" media="CCURStream_$RepresentationID$_$Number$.ts?ccur_ts_audio_Stream=Stream1&amp;ccur_ts_audio_track=10&amp;ccur_ts_audio_Stream=Stream1&amp;ccur_ts_audio_track=11" startNumber="167791408" presentationTimeOffset="1744652613072" duration="900000"/>
            </Representation>
        </AdaptationSet>
    </Period>
</MPD>

*/

// calculateDashSegmentTimestamp is used to calculate the timestamp values for the segment in the dash segment timeline
func calculateDashSegmentTimestamp(timestampOfFirstSegment *uint64, segmentDuration uint64, segmentRepeat *int64) []uint64 {
	var timestamps []uint64
	var timestamp uint64

	var i int64

	// Loop the number of times indicated by the segment repeat value
	for i = 0; i < *segmentRepeat; i++ {
		// If it's the first loop use the timestamp of the first segment
		// otherwise increment the timestamp by the segment duration
		if i > 0 {
			timestamp = timestamp + segmentDuration
		} else {
			timestamp = *timestampOfFirstSegment
		}
		timestamps = append(timestamps, timestamp)
	}
	return timestamps
}

func getSegmentsFromSegmentTimeline(dashSegmentTimestamps []uint64, baseURL, representationId, media string, results *array.Array[Result]) {
	for _, timestamp := range dashSegmentTimestamps {
		var representationRegex = `\$RepresentationID\$`
		var timeRegex = `\$Time\$`

		var segmentName string
		var r = regexp.MustCompile(representationRegex)
		segmentName = r.ReplaceAllString(media, representationId)

		var n = regexp.MustCompile(timeRegex)
		segmentName = n.ReplaceAllString(segmentName, fmt.Sprint(timestamp))

		result := Result{
			playlistURL:    "",
			streamName:     "",
			streamURL:      "",
			segmentName:    segmentName,
			segmentURL:     fmt.Sprintf("%s/%s", baseURL, segmentName),
			byteRangeStart: -1,
			byteRangeSize:  -1,
		}
		results.Push(result)
	}
}

func getSegmentsFromSegmentTemplate(segmentDuration, timescale, startNumber, manifestDuration uint64, baseURL, representationId, media string, results *array.Array[Result]) {
	// get the segment size
	// duration="900000" / timescale="90000"
	// so 10 second segments
	segmentSize := segmentDuration / timescale

	// the media presentation duration is the size of the sliding window
	// mediaPresentationDuration="PT0H15M10S"

	// divide the media presentation duration / segment size
	// to get the number of segments
	numberOfSegments := manifestDuration / segmentSize

	// start number - N where N is number of segments to get the last
	// segment in the window then increment N times
	N := startNumber + numberOfSegments

	for i := startNumber; i < N; i++ {
		segmentNumber := fmt.Sprintf("%d", i)
		var representationRegex = `\$RepresentationID\$`
		var numberRegex = `\$Number\$`

		var segmentName string
		var r = regexp.MustCompile(representationRegex)
		segmentName = r.ReplaceAllString(media, representationId)

		var n = regexp.MustCompile(numberRegex)
		segmentName = n.ReplaceAllString(segmentName, segmentNumber)

		result := Result{
			playlistURL:    "",
			streamName:     "",
			streamURL:      "",
			segmentName:    segmentName,
			segmentURL:     fmt.Sprintf("%s/%s", baseURL, segmentName),
			byteRangeStart: -1,
			byteRangeSize:  -1,
		}
		results.Push(result)
	}
}

func getManifest(url string) *mpd.MPD {
	manifestFile, err := toolbox.SendRequest(toolbox.GET, url, "", nil)
	if err != nil {
		panic(err)
	}
	dashManifest := new(mpd.MPD)
	err = dashManifest.Decode(manifestFile)
	if err != nil {
		panic(err)
	}
	return dashManifest
}

func parseDASH(url string) (*array.Array[Result], error) {
	var results *array.Array[Result]
	representations := make(map[string]string)

	var segmentDuration uint64
	var timescale uint64
	var manifestDuration uint64
	var startNumber uint64
	var baseURL string
	var representationId string
	var media string

	baseURL = toolbox.BaseURL(url)

	dashManifest := getManifest(url)

	/*
		data, err := json.MarshalIndent(dashManifest, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))

	*/
	manifestDuration = 900
	//fmt.Println(mpdStr)
	//d, err := duration.ParseISO8601(mpdStr)
	//if err != nil {
	//panic(err)
	//}

	//x := d.M * 60
	//y := d.TS
	//xy := x + y

	for _, period := range dashManifest.Period {
		for _, set := range period.AdaptationSets {
			for _, rep := range set.Representations {
				representationId = *rep.ID
				representations[representationId] = url
				timescale = *rep.SegmentTemplate.Timescale
				media = *rep.SegmentTemplate.Media

				if rep.SegmentTemplate.StartNumber != nil {
					startNumber = *rep.SegmentTemplate.StartNumber

					segmentDuration = *rep.SegmentTemplate.Duration

					// get segments for this representation
					getSegmentsFromSegmentTemplate(segmentDuration, timescale, startNumber, manifestDuration, baseURL, representationId, media, results)
					//representation.Segments = segments

					//representations.Push(representation)
				} else {
					var dashSegmentTimestamps []uint64

					for _, timeline := range rep.SegmentTemplate.SegmentTimeline.S {
						timestamps := calculateDashSegmentTimestamp(timeline.T, timeline.D, timeline.R)
						dashSegmentTimestamps = append(dashSegmentTimestamps, timestamps...)
					}
					//segments := getSegmentsFromSegmentTimeline(dashSegmentTimestamps, baseURL, representationId, media)
					//representation.Segments = segments
					//representations.Push(representation)
				}
			}
		}
	}

	return results, nil
}
