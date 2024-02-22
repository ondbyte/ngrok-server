## a simple wrapper for ngrok

this has a single responsibility
take a http handler and a addr and serve the handler locally at the addr and expose the locally served handler at a url to the internet

### PS: you can use `"golang.ngrok.com/ngrok"` directly to have more control

#### example:
```go

func TestNgrok(t *testing.T) {
	authToken := "<auth-token>"
	domain := "your-ngrok-domain.ngrok-free.app"
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*15)
	s, err := ngrok_server.NewWithDomain(ctx, authToken, domain)
	// you can use
	// s, err := ngrok_server.NewWithRandomDomain(ctx, authToken)
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
```