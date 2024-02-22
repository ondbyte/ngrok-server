package ngrok_server

import (
	"context"
	"fmt"
	"net/http"

	_ngrok "golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

// helps to access the tunnel url before actually serving the your app locally and becoming blocking go routine
type Server interface {
	// serves the given handler on the local addr and reroutes any incoming traffic on the Serve.Url() to this addr
	Serve(addr string, h http.Handler) error
	// url of this local server which is exposed to the internet
	Url() (remote string)
	//stops the ngrok session
	Stop() error
}

type server struct {
	listener _ngrok.Tunnel
}

// starts a server which has exposed itself to the internet at Server.Url(),
// pass the ngrok token and ngrok domain
func New(ctx context.Context, authToken, domain string) (Server, error) {
	listener, err := _ngrok.Listen(ctx,
		config.HTTPEndpoint(config.WithDomain(domain)),
		_ngrok.WithAuthtoken(authToken),
	)
	if err != nil {
		return nil, fmt.Errorf("err while _grok.Listen : %v", err)
	}

	s := &server{
		listener: listener,
	}
	return s, nil
}

func (s *server) Url() (remote string) {
	return s.listener.URL()
}

func (s *server) Serve(addr string, h http.Handler) error {
	server := &http.Server{Handler: h, Addr: addr}
	err := server.Serve(s.listener)
	if err != nil {
		return fmt.Errorf("err while server.Serve: %v", err)
	}
	return nil
}

func (s *server) Stop() error {
	return s.listener.Session().Close()
}
