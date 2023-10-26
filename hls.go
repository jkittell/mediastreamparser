package mediastreamparser

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jkittell/array"
	"github.com/jkittell/toolbox"
	"regexp"
	"strconv"
	"strings"
)

func decodeVariant(masterPlaylistURL, variantName, variantURL string, results *array.Array[Result]) error {
	_, playlist, err := toolbox.SendRequest(toolbox.GET, variantURL, "", nil)
	if err != nil {
		return err
	}

	// store byte range then continue to next line for Segment
	var ByteRangeStart int
	var ByteRangeSize int

	segmentFormats := []string{".ts", ".fmp4", ".cmfv", ".cmfa", ".aac", ".ac3", ".ec3", ".webvtt"}
	scanner := bufio.NewScanner(bytes.NewReader(playlist))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "#EXT-X-BYTERANGE") {
			// #EXT-X-BYTERANGE:44744@2304880
			// -H "Range: bytes=0-1023"
			// parse byte range here
			byteRangeValues := strings.Split(line, ":")
			if len(byteRangeValues) != 2 {
				return err
			}
			byteRange := strings.Split(byteRangeValues[1], "@")
			startNumber, err := strconv.Atoi(byteRange[1])
			if err != nil {
				return err
			}
			sizeNumber, err := strconv.Atoi(byteRange[0])
			if err != nil {
				return err
			}

			ByteRangeStart = startNumber
			ByteRangeSize = sizeNumber
			continue
		} else {
			ByteRangeStart = -1
			ByteRangeSize = -1
		}

		for _, format := range segmentFormats {
			var match bool
			if strings.Contains(line, format) {
				match = true
				if match {
					if strings.Contains(line, "#EXT-X-MAP:URI=") {
						re := regexp.MustCompile(`"[^"]+"`)
						initSegment := re.FindString(line)
						if initSegment != "" {
							SegmentName := strings.Trim(initSegment, "\"")
							var SegmentURL string
							if !strings.Contains(SegmentName, "http") {
								baseURL := toolbox.BaseURL(variantURL)
								SegmentURL = fmt.Sprintf("%s/%s", baseURL, SegmentName)
							} else {
								SegmentURL = SegmentName
							}

							result := Result{
								playlistURL:    masterPlaylistURL,
								streamName:     variantName,
								streamURL:      variantURL,
								segmentName:    SegmentName,
								segmentURL:     SegmentURL,
								byteRangeStart: ByteRangeStart,
								byteRangeSize:  ByteRangeSize,
							}
							results.Push(result)
						} else {
							return err
						}
					} else {
						SegmentName := line
						var SegmentURL string
						if !strings.Contains(SegmentName, "http") {
							baseURL := toolbox.BaseURL(variantURL)
							SegmentURL = fmt.Sprintf("%s/%s", baseURL, SegmentName)
						} else {
							SegmentURL = SegmentName
						}
						result := Result{
							playlistURL:    masterPlaylistURL,
							streamName:     variantName,
							streamURL:      variantURL,
							segmentName:    SegmentName,
							segmentURL:     SegmentURL,
							byteRangeStart: ByteRangeStart,
							byteRangeSize:  ByteRangeSize,
						}
						results.Push(result)
					}
				}
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

func decodeMaster(url string) (map[string]string, error) {
	streams := make(map[string]string)
	_, playlist, err := toolbox.SendRequest(toolbox.GET, url, "", nil)
	if err != nil {
		return streams, err
	}

	baseURL := toolbox.BaseURL(url)
	scanner := bufio.NewScanner(bytes.NewReader(playlist))
	for scanner.Scan() {
		var streamURL string
		line := scanner.Text()
		if !strings.Contains(line, "#EXT") && strings.Contains(line, "m3u8") {
			if !strings.Contains(line, "http") {
				streamURL = fmt.Sprintf("%s/%s", baseURL, line)
			} else {
				streamURL = line
			}
			streams[line] = streamURL
		} else if strings.Contains(line, "#EXT-X-I-FRAME-STREAM-INF") || strings.Contains(line, "#EXT-X-MEDIA") {
			regEx := regexp.MustCompile("URI=\"(.*?)\"")
			match := regEx.MatchString(line)
			if match {
				s1 := regEx.FindString(line)
				_, s2, _ := strings.Cut(s1, "=")
				s3 := strings.Trim(s2, "\"")
				URI := s3
				if !strings.Contains(line, "http") {
					streamURL = fmt.Sprintf("%s/%s", baseURL, URI)
				} else {
					streamURL = line
				}
				streams[URI] = streamURL
			}
		}
	}
	return streams, err
}

func parseHLS(url string) (*array.Array[Result], error) {
	results := array.New[Result]()
	variants, err := decodeMaster(url)
	if err != nil {
		return results, err
	}

	if len(variants) > 0 {
		for variantName, variantURL := range variants {
			err = decodeVariant(url, variantName, variantURL, results)
			if err != nil {
				return results, err
			}
		}
	} else {
		err = decodeVariant(url, "", url, results)
	}

	return results, err
}
