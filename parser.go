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

type Parser struct {
	logger         *logger.Logger
	maxConcurrency int64
}

func (p *Parser) Run(url string, format ContentFormat) (*array.Array[Result], error) {
	results := array.New[Result]()
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

func New() *Parser {
	parser := &Parser{
		logger: logger.New("/tmp/", "mediastreamparser.log"),
	}
	return parser
}
