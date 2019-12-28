package avcodec

//#include <libavcodec/avcodec.h>
//
// #cgo pkg-config: libavcodec
import "C"
import (
	"errors"
)

type ParserContext struct {
	parserContext *C.AVCodecParserContext
	CodecContext  *Context
}

func NewParserContext(codecContext *Context) (*ParserContext, error) {
	parser := C.av_parser_init(C.int(codecContext.Codec().CAVCodec.id))
	if parser == nil {
		return nil, errors.New("could not open codec")
	}

	//parser.flags = parser.flags | C.PARSER_FLAG_COMPLETE_FRAMES

	return &ParserContext{
		parserContext: parser,
		CodecContext:  codecContext,
	}, nil
}

func (p *ParserContext) Parse(data []byte, size int, packet *Packet) (int, error) {
	ret := int(C.av_parser_parse2(p.parserContext, p.CodecContext.CAVCodecContext, &packet.CAVPacket.data, &packet.CAVPacket.size,
		(*C.uchar)(C.CBytes(data)), C.int(size), C.AV_NOPTS_VALUE, C.AV_NOPTS_VALUE, C.AV_NOPTS_VALUE))

	if ret < 0 {
		return ret, errors.New("Error while parsing")
	}

	return ret, nil
}

func (p *ParserContext) Free() {
	if p.parserContext != nil {
		C.av_parser_close(p.parserContext)
	}
}
