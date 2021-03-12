package kcp

import (
	"encoding/binary"
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/codec"
	"github.com/bobwong89757/cellnet/peer/kcp"
)

func SendPacket(writer kcp.DataWriter, ctx cellnet.ContextSet, msg interface{}) error {

	var (
		msgData []byte
		msgID   int
		meta    *cellnet.MessageMeta
	)

	switch m := msg.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	default: // 发普通编码包
		var err error

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(msg, ctx)

		if err != nil {
			return err
		}

		msgID = meta.ID
	}

	pktData := make([]byte, HeaderSize+len(msgData))

	// 写入消息长度做验证
	binary.LittleEndian.PutUint16(pktData, uint16(HeaderSize+len(msgData)))

	// Type
	binary.LittleEndian.PutUint16(pktData[2:], uint16(msgID))

	// Value
	copy(pktData[HeaderSize:], msgData)

	writer.WriteData(pktData)

	codec.FreeCodecResource(meta.Codec, msgData, ctx)

	return nil
}
