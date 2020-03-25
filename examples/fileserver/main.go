package main

import (
	"flag"
	"fmt"
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/peer"
	_ "github.com/bobwong89757/cellnet/peer/http"
	"github.com/bobwong89757/cellnet/proc"
	_ "github.com/bobwong89757/cellnet/proc/http"
)

var shareDir = flag.String("share", ".", "folder to share")
var port = flag.Int("port", 9091, "listen port")

func main() {

	flag.Parse()

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("http.Acceptor", "httpfile", fmt.Sprintf(":%d", *port), nil).(cellnet.HTTPAcceptor)
	p.SetFileServe(".", *shareDir)

	proc.BindProcessorHandler(p, "http", nil)

	p.Start()
	queue.StartLoop()

	queue.Wait()
}
