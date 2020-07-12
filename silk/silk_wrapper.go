package silk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/orcaman/writerseeker"
	"io"
	"io/ioutil"
)

/*
#include "SKP_Silk_SDK_API.h"
#include <stdlib.h>
*/
import "C"

var (
	ErrInvalid    = errors.New("not a silk stream")
	ErrCodecError = errors.New("codec error")
)

func DecodeSilkBuffToWave(src []byte, sampleRate int) (dst []byte, err error) {
	reader := bytes.NewBuffer(src)
	f, err := reader.ReadByte()
	if err != nil {
		return
	}
	header := make([]byte, 9)
	var n int
	if f == 2 {
		n, err = reader.Read(header)
		if err != nil {
			return
		}
		if n != 9 {
			err = ErrInvalid
			return
		}
		if string(header) != "#!SILK_V3" {
			err = ErrInvalid
			return
		}
	} else if f == '#' {
		n, err = reader.Read(header)
		if err != nil {
			return
		}
		if n != 8 {
			err = ErrInvalid
			return
		}
		if string(header) != "!SILK_V3" {
			err = ErrInvalid
			return
		}
	} else {
		err = ErrInvalid
		return
	}
	var decControl C.SKP_SILK_SDK_DecControlStruct
	decControl.API_sampleRate = C.int32_t(sampleRate)
	decControl.framesPerPacket = 1
	var decSize C.int32_t
	C.SKP_Silk_SDK_Get_Decoder_Size(&decSize)
	dec := C.malloc(C.size_t(decSize))
	defer C.free(dec)
	if C.SKP_Silk_SDK_InitDecoder(dec) != 0 {
		err = ErrCodecError
		return
	}
	// 40ms
	frameSize := sampleRate / 1000 * 40
	in := make([]byte, frameSize)
	buf := make([]int16, frameSize)
	out := &writerseeker.WriterSeeker{}
	enc := wav.NewEncoder(out, sampleRate, 16, 1, 1)
	audioBuf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  sampleRate,
		},
	}
	for {
		var nByte C.int16_t
		err = binary.Read(reader, binary.LittleEndian, &nByte)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		if int(nByte) > frameSize {
			err = ErrInvalid
			return
		}
		n, err = reader.Read(in[:nByte])
		if err != nil {
			return
		}
		if n != int(nByte) {
			err = ErrInvalid
			return
		}
		C.SKP_Silk_SDK_Decode(dec, &decControl, 0,
			(*C.SKP_uint8)(&in[0]), C.int(n),
			(*C.SKP_int16)(&buf[0]), &nByte,
		)
		for _, w := range buf[:int(nByte)] {
			audioBuf.Data = append(audioBuf.Data, int(w))
		}
	}
	if err = enc.Write(audioBuf); err != nil {
		return
	}
	if err = enc.Close(); err != nil {
		return
	}
	dst, err = ioutil.ReadAll(out.Reader())
	return
}
