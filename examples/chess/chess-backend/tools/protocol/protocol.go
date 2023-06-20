package protocol

import (
	packtool "chess-backend/tools/packet"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/ringbuffer"
)

// 拆包/封包协议, 4字节先导
type Protocol struct{}

func (p *Protocol) UnPacket(c *gev.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
	if buffer.Length() < 4 {
		return nil, nil
	}

	packetLength := buffer.PeekUint32()
	if packetLength == 0 {
		return nil, []byte{}
	}

	if buffer.Length() < 4+int(packetLength) {
		return nil, []byte{}
	}

	buffer.Retrieve(4)
	packetBytes := make([]byte, packetLength)
	copy(packetBytes, buffer.Bytes())
	buffer.Retrieve(int(packetLength))

	return nil, packetBytes
}

func (p *Protocol) Packet(c *gev.Connection, data interface{}) []byte {
	bs := data.([]byte)
	return packtool.DoPackWith4BytesHeader(bs)
}
