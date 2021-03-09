package kcp

import (
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/log"
	"github.com/bobwong89757/cellnet/peer"
	"github.com/bobwong89757/kcp-go/v6"
	"net"
)

type udpConnector struct {
	peer.CoreSessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle

	remoteAddr *net.UDPAddr

	defaultSes *kcpSession
}

func (self *udpConnector) Start() cellnet.Peer {

	var err error
	self.remoteAddr, err = net.ResolveUDPAddr("udp", self.Address())

	if err != nil {

		log.GetLog().Error("#resolve udp address failed(%s) %v", self.Name(), err.Error())
		return self
	}

	go self.connect()

	return self
}

func (self *udpConnector) Session() cellnet.Session {
	return self.defaultSes
}

func (self *udpConnector) IsReady() bool {

	return self.defaultSes.GetKcpSession() != nil
}

func (self *udpConnector) connect() {
	sess, err := kcp.DialWithOptions(self.remoteAddr.String(), nil, 10, 3)
	if err != nil {
		log.GetLog().Error("#udp.connect failed(%s) %v", self.Name(), err.Error())
		return
	}

	self.defaultSes.SetKcpSession(sess)

	ses := self.defaultSes

	self.ProcEvent(&cellnet.RecvMsgEvent{ses, &cellnet.SessionConnected{}})

	recvBuff := make([]byte, MaxUDPRecvBuffer)

	self.SetRunning(true)

	for self.IsRunning() {
		ses.Recv(recvBuff)
		//if n, err := io.ReadFull(sess, recvBuff); err == nil {
		//	ses.Recv(recvBuff[:n])
		//} else {
		//	break
		//}

	}
}

func (self *udpConnector) Stop() {

	self.SetRunning(false)

	if c := self.defaultSes.GetKcpSession(); c != nil {
		c.Close()
	}
}

func (self *udpConnector) TypeName() string {
	return "kcp.Connector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &udpConnector{}
		p.defaultSes = &kcpSession{
			pInterface:     p,
			CoreProcBundle: &p.CoreProcBundle,
		}

		return p
	})
}
