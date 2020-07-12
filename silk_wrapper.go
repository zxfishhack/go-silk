package go_silk

import (
	"bytes"
	"errors"
)

/*
##include "silk/SKP_Silk_SDK_API.h"
*/
import "C"

var (
	ErrInvalid = errors.New("not a silk stream")
)

func DecodeSilkBuffToWave(src []byte) (dst []byte, err error) {
	reader := bytes.NewBuffer(src)
	f, err := reader.ReadByte()
	if err != nil {
		return
	}
	header := make([]byte, 9)
	if f == 2 {
		n, err := reader.Read(header)
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
		n, err := reader.Read(header)
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

}
