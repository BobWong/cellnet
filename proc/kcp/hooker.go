package kcp

import (
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/log"
	"github.com/bobwong89757/cellnet/msglog"
	"github.com/bobwong89757/cellnet/relay"
	"github.com/bobwong89757/cellnet/rpc"
)

// 带有RPC和relay功能
type MsgHooker struct {
}

func (self MsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	var handled bool
	var err error

	inputEvent, handled, err = rpc.ResolveInboundEvent(inputEvent)

	if err != nil {
		log.GetLog().Error("rpc.ResolveInboundEvent:", err)
		return
	}

	if !handled {

		inputEvent, handled, err = relay.ResoleveInboundEvent(inputEvent)

		if err != nil {
			log.GetLog().Error("relay.ResoleveInboundEvent:", err)
			return
		}

		if !handled {
			msglog.WriteRecvLogger("kcp", inputEvent.Session(), inputEvent.Message())
		}
	}

	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	handled, err := rpc.ResolveOutboundEvent(inputEvent)

	if err != nil {
		log.GetLog().Error("rpc.ResolveOutboundEvent:", err)
		return nil
	}

	if !handled {

		handled, err = relay.ResolveOutboundEvent(inputEvent)

		if err != nil {
			log.GetLog().Error("relay.ResolveOutboundEvent:", err)
			return nil
		}

		if !handled {
			msglog.WriteSendLogger("kcp", inputEvent.Session(), inputEvent.Message())
		}
	}

	return inputEvent
}
