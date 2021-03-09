package kcp

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/msglog"
	"github.com/bobwong89757/cellnet/peer/udp"
	"github.com/bobwong89757/cellnet/proc"
)

type KCPMessageTransmitter struct {
}

func (KCPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	data := ses.Raw().(udp.DataReader).ReadData()

	msg, err = RecvPacket(data)

	msglog.WriteRecvLogger("kcp", ses, msg)

	return
}

func (KCPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	writer := ses.(udp.DataWriter)

	msglog.WriteSendLogger("kcp", ses, msg)

	// ses不再被复用, 所以使用session自己的contextset做内存池, 避免串台
	return sendPacket(writer, ses.(cellnet.ContextSet), msg)
}

func init() {
	var convid uint32
	binary.Read(rand.Reader, binary.LittleEndian, &convid)
	proc.RegisterProcessor("kcp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(KCPMessageTransmitter))
		bundle.SetCallback(userCallback)

	})
}
