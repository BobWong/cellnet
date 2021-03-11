package kcp

import (
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/log"
	"github.com/bobwong89757/cellnet/peer"
	"github.com/bobwong89757/kcp-go/v6"
	"net"
	"sync"
	"time"
)

type DataReader interface {
	ReadData() []byte
}

type DataWriter interface {
	WriteData(data []byte)
}

// Socket会话
type kcpSession struct {
	*peer.CoreProcBundle
	peer.CoreContextSet
	peer.CoreSessionIdentify

	pInterface cellnet.Peer

	pkt []byte

	// Socket原始连接
	//remote      *net.UDPAddr
	//conn        *net.UDPConn
	connGuard   sync.RWMutex
	timeOutTick time.Time
	kcpSession *kcp.UDPSession
	key         *connTrackKey
}

func (self *kcpSession) SetKcpSession(udpSes *kcp.UDPSession) {
	self.connGuard.Lock()
	self.kcpSession = udpSes
	self.connGuard.Unlock()
}

func (self *kcpSession) GetKcpSession() *kcp.UDPSession {
	self.connGuard.RLock()
	defer self.connGuard.RUnlock()
	return self.kcpSession
}

func (self *kcpSession) IsAlive() bool {
	return time.Now().Before(self.timeOutTick)
}

func (self *kcpSession) ID() int64 {
	return 0
}

func (self *kcpSession) LocalAddress() net.Addr {
	return self.GetKcpSession().LocalAddr()
}

func (self *kcpSession) Peer() cellnet.Peer {
	return self.pInterface
}

// 取原始连接
func (self *kcpSession) Raw() interface{} {
	return self
}

func (self *kcpSession) Recv(data []byte) {
	n,err := self.kcpSession.Read(data)
	if err != nil {
		log.GetLog().Error("kcp读取错误 %v",err)
	}
	self.pkt = data[:n]
	msg, err := self.ReadMessage(self)

	if msg != nil && err == nil {
		self.ProcEvent(&cellnet.RecvMsgEvent{self, msg})
	}
}

func (self *kcpSession) ReadData() []byte {
	return self.pkt
}

func (self *kcpSession) WriteData(data []byte) {

	c := self.GetKcpSession()
	if c == nil {
		return
	}

	// Connector中的Session
	if self.kcpSession.RemoteAddr() == nil {
		c.Write(data)

		// Acceptor中的Session
	} else {
		self.kcpSession.Write(data)
		//c.WriteToUDP(data, self.remote)
	}
}

// 发送封包
func (self *kcpSession) Send(msg interface{}) {

	self.SendMessage(&cellnet.SendMsgEvent{self, msg})
}

func (self *kcpSession) Close() {
	self.kcpSession.Close()
}
