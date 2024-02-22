package ngrok_server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	ngrok_server "github.com/ondbyte/ngrok-server"
)

func TestNgrok(t *testing.T) {
	authToken := "<auth-token>"
	domain := "square-yakooza-random.ngrok-free.app"
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*15)
	s, err := ngrok_server.NewWithDomain(ctx, authToken, domain)
	// you can use
	// s, err := ngrok_server.NewWithRandomDomain(ctx, authToken, domain)
	// for randomely generated url
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan bool, 1)
	go func() {
		h := http.NewServeMux()
		h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			ch <- true
		})
		s.Serve(":8888", h)
	}()
	url := s.Url()
	_, err = http.Get(url)
	if err != nil {
		t.Fatal("failed to make a get request to url ", s.Url())
	}
	done := <-ch
	if !done {
		t.Fatal("should have been true")
	}

}
