package tcp

import (
	"github.com/BobWong/cellnet"
	"github.com/BobWong/cellnet/proc"
)

func init() {

	transmitter := new(TCPMessageTransmitter)
	hooker := new(MsgHooker)

	proc.RegisterProcessor("tcp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(transmitter)
		bundle.SetHooker(hooker)
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))

	})
}
