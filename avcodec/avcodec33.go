// +build ffmpeg33

package avcodec

//#include <libavutil/avutil.h>
//#include <libavcodec/avcodec.h>
//
// #cgo pkg-config: libavcodec libavutil
import "C"

import (
	"unsafe"

	"github.com/imkira/go-libav/avutil"
)

type CodecParameters struct {
	CAVCodecParameters *C.AVCodecParameters
}

func NewCodecParameters() (*CodecParameters, error) {
	cPkt := (*C.AVCodecParameters)(C.avcodec_parameters_alloc())
	if cPkt == nil {
		return nil, ErrAllocationError
	}
	return NewCodecParametersFromC(unsafe.Pointer(cPkt)), nil
}

func NewCodecParametersFromC(cPSD unsafe.Pointer) *CodecParameters {
	return &CodecParameters{CAVCodecParameters: (*C.AVCodecParameters)(cPSD)}
}

func (cParams *CodecParameters) Free() {
	C.avcodec_parameters_free(&cParams.CAVCodecParameters)
}

func (ctx *Context) CopyTo(dst *Context) error {
	// added in lavc 57.33.100
	parameters, err := NewCodecParameters()
	if err != nil {
		return err
	}
	defer parameters.Free()
	cParams := (*C.AVCodecParameters)(unsafe.Pointer(parameters.CAVCodecParameters))
	code := C.avcodec_parameters_from_context(cParams, ctx.CAVCodecContext)
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	code = C.avcodec_parameters_to_context(dst.CAVCodecContext, cParams)
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

func (ctx *Context) DecodeVideo(pkt *Packet, onFrame func(*avutil.Frame)) (int, error) {
	code, err := ctx.SendPacket(pkt)
	if err != nil {
		return code, err
	}

	for {
		frame, err := avutil.NewFrame()
		if err != nil {
			panic(err)
		}

		code, err := ctx.ReceiveFrame(frame)
		if code == avutil.AVERROR_EOF || code == avutil.AVERROR_EAGAIN {
			frame.Free()
			break
		} else if code < 0 {
			frame.Free()
			return code, err
		}

		onFrame(frame)
		frame.Free()
	}

	return 0, nil
}

func (ctx *Context) SendPacket(pkt *Packet) (int, error) {
	cPkt := (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket))
	code := C.avcodec_send_packet(ctx.CAVCodecContext, cPkt)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return int(code), err
}

func (ctx *Context) ReceiveFrame(frame *avutil.Frame) (int, error) {
	cFrame := (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	var err error
	code := C.avcodec_receive_frame(ctx.CAVCodecContext, cFrame)
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return int(code), err
}

func (ctx *Context) DecodeAudio(pkt *Packet, frame *avutil.Frame) (int, error) {
	cFrame := (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	cPkt := (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket))
	code := C.avcodec_send_packet(ctx.CAVCodecContext, cPkt)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		return int(code), err
	}
	code = C.avcodec_receive_frame(ctx.CAVCodecContext, cFrame)
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return int(code), err
}

func (ctx *Context) EncodeVideo(pkt *Packet, frame *avutil.Frame, onData func([]byte)) (int, error) {
	code, err := ctx.SendFrame(frame)
	if err != nil {
		return code, err
	}

	for {
		code, err := ctx.ReceivePacket(pkt)
		if code == avutil.AVERROR_EOF || code == avutil.AVERROR_EAGAIN {
			pkt.Unref()
			break
		} else if code < 0 {
			pkt.Unref()
			return code, err
		}

		data := C.GoBytes(pkt.Data(), C.int(pkt.Size()))
		onData(data)

		pkt.Unref()
	}

	return 0, nil
}

func (ctx *Context) SendFrame(frame *avutil.Frame) (int, error) {
	var err error
	code := C.avcodec_send_frame(ctx.CAVCodecContext, (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame)))
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return int(code), err
}

func (ctx *Context) ReceivePacket(pkt *Packet) (int, error) {
	var err error
	code := C.avcodec_receive_packet(ctx.CAVCodecContext, (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket)))
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return int(code), err
}

func (ctx *Context) EncodeAudio(pkt *Packet, frame *avutil.Frame) (int, error) {
	var err error
	cPkt := (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket))
	code := C.avcodec_send_frame(ctx.CAVCodecContext, (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame)))
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		return int(code), err
	}

	code = C.avcodec_receive_packet(ctx.CAVCodecContext, cPkt)
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return int(code), err
}
