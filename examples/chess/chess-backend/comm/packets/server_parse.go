package packets

import "encoding/json"

// 提供给服务端的解析函数
func ServerParse(bs []byte) interface{} {
	header := PacketHeader{}
	err := json.Unmarshal(bs, &header)
	if err != nil {
		return nil
	}

	if header.Type == nil {
		return nil
	}

	switch *header.Type {
	case PacketTypeHeartbeat:
		return &PacketHeartbeat{}
	case PacketTypeClientMove:
		p := PacketClientMove{}
		json.Unmarshal(bs, &p)
		return &p
	case PacketTypeClientStartMatch:
		p := PacketClientStartMatch{}
		json.Unmarshal(bs, &p)
		return &p
	case PacketTypeClientPawnUpgrade:
		p := PacketClientPawnUpgrade{}
		json.Unmarshal(bs, &p)
		return &p
	case PacketTypeClientWhetherAcceptDraw:
		p := PacketClientWhetherAcceptDraw{}
		json.Unmarshal(bs, &p)
		return &p
	case PacketTypeClientDoSurrender:
		p := PacketClientDoSurrender{}
		json.Unmarshal(bs, &p)
		return &p
	default:
		return nil
	}
}
