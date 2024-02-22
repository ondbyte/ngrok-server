package ngrok_server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	ngrok_server "github.com/ondbyte/ngro/ngrok-server"
)

func TestNgrok(t *testing.T) {
	authToken := "1Tk4pgsmPWeH9qPUpbE9gffUTNs_6TMzjgFknNh6KppUxrYet"
	domain := ""
	s, err := ngrok_server.New(context.TODO(), authToken, domain)
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
	time.Sleep(time.Second * 3)
	go func() {
		time.Sleep(time.Second * 10)
		ch <- false
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

	err = s.Stop()
	if err != nil {
		t.Fatal(err)
	}

	s, err = ngrok_server.New(context.TODO(), authToken, "")
	if err != nil {
		t.Fatal(err)
	}
	ch = make(chan bool, 1)
	go func() {
		h := http.NewServeMux()
		h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			ch <- true
		})
		panic(s.Serve(":8888", h))
	}()

	time.Sleep(time.Second * 3)
	_, err = http.Get(url)
	if err != nil {
		t.Fatal("failed to make a get request to url ", s.Url())
	}
	done = <-ch
	if !done {
		t.Fatal("should have been true")
	}

}
