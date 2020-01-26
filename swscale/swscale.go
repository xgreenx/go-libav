package swscale

//#include <string.h>
//#include <libswscale/swscale.h>
//#include <libavutil/avutil.h>
//#include <libavutil/frame.h>
//
// typedef struct SwsContext SwsContext;
//
// #cgo pkg-config: libswscale libavutil
import "C"
import (
	"errors"
	"github.com/imkira/go-libav/avutil"
	"unsafe"
)

var (
	ErrAllocationError     = errors.New("allocation error")
	ErrInvalidArgumentSize = errors.New("invalid argument size")
)

type LogLevel int

const (
	LogLevelQuiet   LogLevel = C.AV_LOG_QUIET
	LogLevelPanic   LogLevel = C.AV_LOG_PANIC
	LogLevelFatal   LogLevel = C.AV_LOG_FATAL
	LogLevelError   LogLevel = C.AV_LOG_ERROR
	LogLevelWarning LogLevel = C.AV_LOG_WARNING
	LogLevelInfo    LogLevel = C.AV_LOG_INFO
	LogLevelVerbose LogLevel = C.AV_LOG_VERBOSE
	LogLevelDebug   LogLevel = C.AV_LOG_DEBUG
)

func init() {
	SetLogLevel(LogLevelQuiet)
}

func Version() (int, int, int) {
	return int(C.LIBAVUTIL_VERSION_MAJOR), int(C.LIBAVUTIL_VERSION_MINOR), int(C.LIBAVUTIL_VERSION_MICRO)
}

func SetLogLevel(level LogLevel) {
	C.av_log_set_level(C.int(level))
}

type DataDescription struct {
	Width  int
	Height int
	PixFmt avutil.PixelFormat
}

type Context struct {
	CSwsContext *C.SwsContext
}

func NewContext(input *DataDescription, output *DataDescription) (*Context, error) {
	cCtx := C.sws_getContext(
		C.int(input.Width),
		C.int(input.Height),
		(C.enum_AVPixelFormat)(input.PixFmt),
		C.int(output.Width),
		C.int(output.Height),
		(C.enum_AVPixelFormat)(output.PixFmt),
		C.SWS_BICUBIC, nil, nil, nil,
	)
	if cCtx == nil {
		return nil, ErrAllocationError
	}

	return NewContextFromC(unsafe.Pointer(cCtx)), nil
}

func NewContextFromC(cCtx unsafe.Pointer) *Context {
	return &Context{
		CSwsContext: (*C.SwsContext)(cCtx),
	}
}

func (ctx *Context) Scale(inFrame *avutil.Frame, offset int, height int, outFrame *avutil.Frame) {
	var inF *C.AVFrame = (*C.AVFrame)(unsafe.Pointer(inFrame.CAVFrame))
	var outF *C.AVFrame = (*C.AVFrame)(unsafe.Pointer(outFrame.CAVFrame))
	C.sws_scale(ctx.CSwsContext, &inF.data[0], &inF.linesize[0], C.int(offset), C.int(height), &outF.data[0], &outF.linesize[0])
}

func (ctx *Context) Free() {
	if ctx.CSwsContext != nil {
		defer C.sws_freeContext(ctx.CSwsContext)
		ctx.CSwsContext = nil
	}
}
