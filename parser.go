package mediastreamparser

import (
	"github.com/jkittell/array"
	"github.com/jkittell/logger"
)

type ContentFormat byte
type ContainerFormat byte

const (
	HLS ContentFormat = iota
	DASH
)

type Scanner struct {
	logger         *logger.Logger
	maxConcurrency int64
}

func (s *Scanner) Scan(url string, format ContentFormat) (*array.Array[Result], error) {
	var results *array.Array[Result]
	var err error

	switch format {
	case HLS:
		return parseHLS(url)
	case DASH:
		return parseDASH(url)
	default:
		return results, err
	}
}

func New() *Scanner {
	scanner := &Scanner{
		logger: logger.New("/tmp/", "ottscanner.log"),
	}
	return scanner
}
