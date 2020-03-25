package http

import "github.com/bobwong89757/cellnet"

type StatusRespond struct {
	StatusCode int
}

func (self *StatusRespond) WriteRespond(ses *httpSession) error {

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	log.Debugf("#http.recv(%s) '%s' %s | [%d] Status",
		peerInfo.Name(),
		ses.req.Method,
		ses.req.URL.Path,
		self.StatusCode)

	ses.resp.WriteHeader(int(self.StatusCode))
	return nil
}
