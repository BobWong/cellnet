package msglog

import (
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/golog/logs"
)

// 萃取消息中的消息
type PacketMessagePeeker interface {
	Message() interface{}
}

func WriteRecvLogger(log *logs.BeeLogger, protocol string, ses cellnet.Session, msg interface{}) {

	if peeker, ok := msg.(PacketMessagePeeker); ok {
		msg = peeker.Message()
	}

	if IsMsgLogValid(cellnet.MessageToID(msg)) {
		peerInfo := ses.Peer().(cellnet.PeerProperty)

		log.Debug("#%s.recv(%s)@%d len: %d %s | %s",
			protocol,
			peerInfo.Name(),
			ses.ID(),
			cellnet.MessageSize(msg),
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}

}

func WriteSendLogger(log *logs.BeeLogger, protocol string, ses cellnet.Session, msg interface{}) {

	if peeker, ok := msg.(PacketMessagePeeker); ok {
		msg = peeker.Message()
	}

	if IsMsgLogValid(cellnet.MessageToID(msg)) {
		peerInfo := ses.Peer().(cellnet.PeerProperty)

		log.Debug("#%s.send(%s)@%d len: %d %s | %s",
			protocol,
			peerInfo.Name(),
			ses.ID(),
			cellnet.MessageSize(msg),
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}

}
