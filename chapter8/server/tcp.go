package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// construiremos un servidor TCP simple que maneje
// solicitudes simultáneamente y se cierre correctamente

type TCPServer struct {
	Listener    net.Listener
	HandlerFunc func(context.Context, net.Conn)

	wg sync.WaitGroup
}

func (srv *TCPServer) Listen() error {
	baseContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		conn, err := srv.Listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			fmt.Println(err)
		}

		srv.wg.Add(1)
		go func() {
			defer srv.wg.Done()
			srv.HandlerFunc(baseContext, conn)
		}()
	}
}

func (srv *TCPServer) StopListener() error {
	return srv.Listener.Close()
}

func (srv *TCPServer) WaitForConnections(timeout time.Duration) {
	toCh := time.After(timeout)
	doneCh := make(chan struct{})

	go func() {
		srv.wg.Wait()
		close(doneCh)
	}()

	select {
	case <-toCh:
	case <-doneCh:
	}
}

func LocalMainTCP() {
	var srv TCPServer
	var err error

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		fmt.Println("Terminating on signal")
		err := srv.StopListener()
		if err != nil {
			log.Fatal(err)
		}
		srv.WaitForConnections(5 * time.Second)
	}()

	// -----

	srv.Listener, err = net.Listen("tcp", "")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening at %s\n", srv.Listener.Addr())
	srv.HandlerFunc = func(ctx context.Context, conn net.Conn) {
		defer conn.Close()
		defer fmt.Println("Connection closed")
		fmt.Println("Handling connection")
		// Echo server
		_, _ = io.Copy(conn, conn)
	}
	err = srv.Listen()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Waiting for connections to terminate")
	srv.WaitForConnections(5 * time.Second)
	fmt.Println("Done")
}
