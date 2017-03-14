package main

import (
	"net"
	"net/http"

	"log"

	"sync"

	"os"

	"os/signal"

	"github.com/gorilla/mux"
	"github.com/wang502/chord"
)

// closing all servers, notify wait group after an OS interrup signal comes in
func cleanup(listeners []net.Listener, wg *sync.WaitGroup) error {
	log.Println("Cleaning up.....")
	for _, listener := range listeners {
		log.Printf("closing server at address %s ....", listener.Addr())
		if err := listener.Close(); err != nil {
			return err
		}
		log.Printf("closed server at address %s", listener.Addr())
		wg.Done()
	}
	return nil
}

func main() {
	host := "http://localhost"
	ports := []string{":2000", ":2001", ":3000", ":4000", ":5000"}

	listenersSlice := make([]net.Listener, len(ports))
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	var wg sync.WaitGroup

	for i := 0; i < len(ports); i++ {
		transporter := chord.NewTransporter()
		chordServer := chord.NewServer("TestNode"+ports[i], chord.DefaultConfig(host+ports[i]), transporter)
		router := mux.NewRouter()
		transporter.Install(chordServer, router)

		server := &http.Server{Addr: ports[i], Handler: router}
		listener, err := net.Listen("tcp", ports[i])
		listenersSlice[i] = listener

		log.Printf("Listening at: %s", host+ports[i])

		if err != nil {
			log.Println(err)
			return
		}
		wg.Add(1)
		go func(port string, server *http.Server) {
			err := server.Serve(listener.(*net.TCPListener))
			if err != nil {
				log.Printf("http server error, %s", err)
			}
		}(ports[i], server)
	}

	wg.Add(1)
	go func() {
		// after receive an OS interrup signal, start the cleanup
		<-signalChan
		if err := cleanup(listenersSlice, &wg); err != nil {
			log.Println("failed to close every server")
		}
		wg.Done()
	}()

	wg.Wait()
}
