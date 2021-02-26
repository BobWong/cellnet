package msglog

import (
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/log"
)

// 萃取消息中的消息
type PacketMessagePeeker interface {
	Message() interface{}
}

func WriteRecvLogger(protocol string, ses cellnet.Session, msg interface{}) {

	if peeker, ok := msg.(PacketMessagePeeker); ok {
		msg = peeker.Message()
	}

	if IsMsgLogValid(cellnet.MessageToID(msg)) {
		peerInfo := ses.Peer().(cellnet.PeerProperty)

		log.GetLog().Debug("#%s.recv(%s)@%d len: %d %s | %s",
			protocol,
			peerInfo.Name(),
			ses.ID(),
			cellnet.MessageSize(msg),
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}

}

func WriteSendLogger(protocol string, ses cellnet.Session, msg interface{}) {

	if peeker, ok := msg.(PacketMessagePeeker); ok {
		msg = peeker.Message()
	}

	if IsMsgLogValid(cellnet.MessageToID(msg)) {
		peerInfo := ses.Peer().(cellnet.PeerProperty)

		log.GetLog().Debug("#%s.send(%s)@%d len: %d %s | %s",
			protocol,
			peerInfo.Name(),
			ses.ID(),
			cellnet.MessageSize(msg),
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}

}
