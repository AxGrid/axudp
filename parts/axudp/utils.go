package axudp

import (
	"axudp/shared/axtools"
	"encoding/binary"
)

func booleansToBytes(t []bool) []byte {
	b := make([]byte, (len(t)+7)/8)
	for i, x := range t {
		if x {
			b[i/8] |= 0x80 >> uint(i%8)
		}
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

func getBytesFromInt(len int) []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len))
	return bs
}

func getIntFromBytes(lens []byte) int {
	return int(binary.LittleEndian.Uint16(lens))
}

func getChunk(b []byte, l int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(b)/l+1)
	for len(b) >= l {
		chunk, b = b[:l], b[l:]
		chunks = append(chunks, chunk)
	}
	if len(b) > 0 {
		chunks = append(chunks, b[:len(b)])
	}
	return chunks
}

func dataTail(b []byte) [][]byte {
	var res [][]byte
	size := getIntFromBytes(b[0:2])
	l := len(b)
	if l >= size+2 {
		res = append(res, b[2:size+2])
		if l > size+2 {
			res = append(res, dataTail(b[:size+2])...)
		}
	}
	return res
}

func mergeBoolLists(into *[]bool, b []bool) {
	for i := 0; i < axtools.MinInt(len(*into), len(b)); i++ {
		if b[i] {
			(*into)[i] = true
		}
	}
}

func booleanIndex(a []bool) int {
	for i := 0; i < len(a); i++ {
		if a[i] {
			return i
		}
	}
	return -1
}

func allBooleans(a []bool) bool {
	for _, x := range a {
		if !x {
			return false
		}
	}
	return true
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
