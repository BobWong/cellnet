package tests

import (
	"fmt"
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/peer"
	httppeer "github.com/bobwong89757/cellnet/peer/http"
	"github.com/bobwong89757/cellnet/proc"
	"net/http"
	"testing"
)

const pageAddress = "127.0.0.1:10087"

func TestPrintPage(t *testing.T) {

	p := peer.NewGenericPeer("http.Acceptor", "httpserver", pageAddress, nil)

	proc.BindProcessorHandler(p, "http", func(raw cellnet.Event) {

		switch {
		case raw.Session().(httppeer.RequestMatcher).Match("GET", "/"):

			raw.Session().Send(&httppeer.HTMLRespond{
				StatusCode:    http.StatusOK,
				PageTemplate:  "index",
				TemplateModel: "world",
			})
		}

	})

	p.Start()

	validPage(t, fmt.Sprintf("http://%s", pageAddress), "<h1>Hello world</h1>")

	p.Stop()

}
