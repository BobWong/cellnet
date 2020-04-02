package gorillaws

import (
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/proc"
)

func init() {

	proc.RegisterProcessor("gorillaws.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(WSMessageTransmitter))
		bundle.SetHooker(new(MsgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))

	})
}
