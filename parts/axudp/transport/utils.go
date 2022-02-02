package transport

import (
	"encoding/binary"
)

func dataTail(b []byte) [][]byte {
	var res [][]byte
	size := getIntFromBytes(b[0:2])
	l := len(b)
	//log.Debug().Int("len", l).Int("size", size+2).Hex("size-hex", b[0:2]).Msg("len")
	if l >= size+2 {
		res = append(res, b[2:size+2])
		if l > size+2 {
			res = append(res, dataTail(b[:size+2])...)
		}
	}
	return res
}

func addSize(data []byte) []byte {
	ld := len(data)
	//log.Debug().Int("datalen", len(data)).Hex("size-hex", getBytesFromInt(ld)).Hex("data", data).Msg("addSize")
	return append(getBytesFromInt(ld), data...)
	//return res
}

func getBytesFromInt(len int) []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len))
	return bs
}

func getIntFromBytes(lens []byte) int {
	return int(binary.LittleEndian.Uint16(lens))
}

func bytesJoin(s ...[]byte) []byte {
	n := 0
	for _, v := range s {
		n += len(v)
	}

	b, i := make([]byte, n), 0
	for _, v := range s {
		i += copy(b[i:], v)
	}
	return b
}

func bytesToBooleans(b []byte) []bool {
	t := make([]bool, 8*len(b))
	for i, x := range b {
		for j := 0; j < 8; j++ {
			if (x<<uint(j))&0x80 == 0x80 {
				t[8*i+j] = true
			}
		}
	}
	return t
}
