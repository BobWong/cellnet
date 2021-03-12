package kcp

import (
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/log"
	"github.com/bobwong89757/cellnet/peer"
	"github.com/bobwong89757/cellnet/util"
	"github.com/bobwong89757/kcp-go/v6"
	"net"
	"sync"
	"sync/atomic"
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
	closing int64
	// 退出同步器
	exitSync sync.WaitGroup

	pInterface cellnet.Peer

	pkt []byte

	// 发送队列
	sendQueue *cellnet.Pipe
	endNotify func()

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

//func (self *kcpSession) Recv(data []byte) {
//	n,err := self.kcpSession.Read(data)
//	if err != nil {
//		log.GetLog().Error("kcp读取错误 %v",err)
//	}
//	self.pkt = data[:n]
//	msg, err := self.ReadMessage(self)
//
//	if msg != nil && err == nil {
//		self.ProcEvent(&cellnet.RecvMsgEvent{self, msg})
//	}
//}

func (self *kcpSession) ReadData() []byte {
	recvBuff := make([]byte, MaxUDPRecvBuffer)
	n, err := self.kcpSession.Read(recvBuff)
	//n, err := self.kcpSession.Read(self.pkt)
	if err != nil {
		log.GetLog().Error("kcp读取错误 %v",err)
	}
	self.pkt = recvBuff[:n]
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

//// 发送封包
//func (self *kcpSession) Send(msg interface{}) {
//
//	self.SendMessage(&cellnet.SendMsgEvent{self, msg})
//}

func (self *kcpSession) Close() {
	closing := atomic.SwapInt64(&self.closing, 1)
	if closing != 0 {
		return
	}
	conn := self.kcpSession.GetConn()
	if conn != nil {
		self.kcpSession.Close()
		self.kcpSession.GetConn().SetReadDeadline(time.Now())
	}
}

func (self *kcpSession) Send(msg interface{}) {
	// 只能通过Close关闭连接
	if msg == nil {
		return
	}

	// 已经关闭，不再发送
	if self.IsManualClosed() {
		return
	}
	self.sendQueue.Add(msg)

}

func (self *kcpSession) protectedReadMessage() (msg interface{}, err error) {

	defer func() {

		if err := recover(); err != nil {
			log.GetLog().Error("IO panic: %s", err)
			self.kcpSession.Close()
		}

	}()

	msg, err = self.ReadMessage(self)

	return
}

func (self *kcpSession) IsManualClosed() bool {
	return atomic.LoadInt64(&self.closing) != 0
}

// 接收循环
func (self *kcpSession) recvLoop() {

	var capturePanic bool

	if i, ok := self.Peer().(cellnet.PeerCaptureIOPanic); ok {
		capturePanic = i.CaptureIOPanic()
	}

	for self.kcpSession.GetConn() != nil {

		var msg interface{}
		var err error

		if capturePanic {
			msg, err = self.protectedReadMessage()
		} else {
			msg, err = self.ReadMessage(self)
		}

		if err != nil {
			if !util.IsEOFOrNetReadError(err) {
				log.GetLog().Error("session closed, sesid: %d, err: %s", self.ID(), err)
			}

			self.sendQueue.Add(nil)

			// 标记为手动关闭原因
			closedMsg := &cellnet.SessionClosed{}
			if self.IsManualClosed() {
				closedMsg.Reason = cellnet.CloseReason_Manual
			}

			self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self, Msg: closedMsg})
			break
		}

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self, Msg: msg})
	}

	// 通知完成
	self.exitSync.Done()
}

// 发送循环
func (self *kcpSession) sendLoop() {

	var writeList []interface{}

	for {
		writeList = writeList[0:0]
		exit := self.sendQueue.Pick(&writeList)

		// 遍历要发送的数据
		for _, msg := range writeList {

			self.SendMessage(&cellnet.SendMsgEvent{Ses: self, Msg: msg})
		}

		if exit {
			break
		}
	}

	// 完整关闭
	conn := self.GetKcpSession().GetConn()
	if conn != nil {
		conn.Close()
	}

	// 通知完成
	self.exitSync.Done()
}



func (self *kcpSession) Start() {
	atomic.StoreInt64(&self.closing, 0)

	// connector复用session时，上一次发送队列未释放可能造成问题
	self.sendQueue.Reset()
	// 需要接收和发送线程同时完成时才算真正的完成
	self.exitSync.Add(2)

	// 将会话添加到管理器, 在线程处理前添加到管理器(分配id), 避免ID还未分配,就开始使用id的竞态问题
	self.Peer().(peer.SessionManager).Add(self)

	go func() {

		// 等待2个任务结束
		self.exitSync.Wait()

		// 将会话从管理器移除
		self.Peer().(peer.SessionManager).Remove(self)

		if self.endNotify != nil {
			self.endNotify()
		}

	}()

	// 启动并发接收goroutine
	go self.recvLoop()

	// 启动并发发送goroutine
	go self.sendLoop()
}

func newSession(session *kcp.UDPSession, p cellnet.Peer, endNotify func()) *kcpSession {
	ses := &kcpSession{}
	ses.pInterface = p
	ses.endNotify = endNotify
	ses.sendQueue = cellnet.NewPipe()
	ses.CoreProcBundle = p.(interface {
		GetBundle() *peer.CoreProcBundle
	}).GetBundle()
	ses.kcpSession = session
	return ses
}
