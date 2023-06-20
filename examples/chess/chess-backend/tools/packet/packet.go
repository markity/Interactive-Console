package packet

import (
	"encoding/binary"
)

// 做封包操作
func DoPackWith4BytesHeader(bs []byte) []byte {
	buf := make([]byte, 4, 4+len(bs))
	binary.BigEndian.PutUint32(buf, uint32(len(bs)))
	buf = append(buf, bs...)
	return buf
}

// 检查字节是否能够组成一个包, 如果能够组成, 拿出包的本体
// 如果bs的前四个字节为0, 此时返回包为空的包->([], true)
func IsBytesCompleteWith4BytesHeader(bs []byte) ([]byte, bool) {
	// 不够4个字节, 不能组成包
	lenBS := len(bs)
	if lenBS < 4 {
		return nil, false
	}

	bytesLength := binary.BigEndian.Uint32(bs[:4])
	if uint32(lenBS) < 4+bytesLength {
		return nil, false
	}

	newBS := make([]byte, bytesLength)
	copy(newBS, bs[4:4+bytesLength])

	return newBS, true
}
