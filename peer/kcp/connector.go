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

type udpConnector struct {
	peer.SessionManager

	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle

	remoteAddr *net.UDPAddr

	defaultSes *kcpSession

	tryConnTimes int // 尝试连接次数

	sesEndSignal sync.WaitGroup

	reconDur time.Duration
}

func (self *udpConnector) Start() cellnet.Peer {

	self.WaitStopFinished()
	if self.IsRunning() {
		return self
	}

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

func (self *udpConnector) SetSessionManager(raw interface{}) {
	self.SessionManager = raw.(peer.SessionManager)
}

func (self *udpConnector) IsReady() bool {
	return self.SessionCount() != 0
}

func (self *udpConnector) connect() {
	self.SetRunning(true)
	for {
		self.tryConnTimes++

		// 尝试用Socket连接地址
		sess, err := kcp.DialWithOptions(self.remoteAddr.String(), nil, 0, 0)
		if err != nil {
			log.GetLog().Error("#udp.connect failed(%s) %v", self.Name(), err.Error())
			return
		}

		self.defaultSes.SetKcpSession(sess)

		// 发生错误时退出
		if err != nil {

			if self.tryConnTimes <= reportConnectFailedLimitTimes {
				log.GetLog().Error("#kcp.connect failed(%s) %v", self.Name(), err.Error())

				if self.tryConnTimes == reportConnectFailedLimitTimes {
					log.GetLog().Error("(%s) continue reconnecting, but mute log", self.Name())
				}
			}

			// 没重连就退出
			if self.ReconnectDuration() == 0 || self.IsStopping() {

				self.ProcEvent(&cellnet.RecvMsgEvent{
					Ses: self.defaultSes,
					Msg: &cellnet.SessionConnectError{},
				})
				break
			}

			// 有重连就等待
			time.Sleep(self.ReconnectDuration())

			// 继续连接
			continue
		}

		self.sesEndSignal.Add(1)

		self.defaultSes.Start()

		self.tryConnTimes = 0

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnet.SessionConnected{}})

		self.sesEndSignal.Wait()

		self.defaultSes.SetKcpSession(nil)

		// 没重连就退出/主动退出
		if self.IsStopping() || self.ReconnectDuration() == 0 {
			break
		}

		// 有重连就等待
		time.Sleep(self.ReconnectDuration())

		// 继续连接
		continue

	}

	self.SetRunning(false)

	self.EndStopping()
}

func (self *udpConnector) Stop() {

	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()

	if c := self.defaultSes.GetKcpSession(); c != nil {
		c.Close()
	}

	// 等待线程结束
	self.WaitStopFinished()
}

func (self *udpConnector) ReconnectDuration() time.Duration {

	return self.reconDur
}

func (self *udpConnector) SetReconnectDuration(v time.Duration) {
	self.reconDur = v
}

func (self *udpConnector) Port() int {

	conn := self.defaultSes.GetKcpSession()

	if conn == nil {
		return 0
	}

	return conn.LocalAddr().(*net.TCPAddr).Port
}

const reportConnectFailedLimitTimes = 3

func (self *udpConnector) TypeName() string {
	return "kcp.Connector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &udpConnector{
			SessionManager: new(peer.CoreSessionManager),
		}
		//p.defaultSes = &kcpSession{
		//	pInterface:     p,
		//	CoreProcBundle: &p.CoreProcBundle,
		//}
		p.defaultSes = newSession(nil,p,func() {
			p.sesEndSignal.Done()
		})

		return p
	})
}
