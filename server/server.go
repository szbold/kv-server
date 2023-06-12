package server

import (
	"key-value-server/datastore"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type server struct {
	addr string
	ds   *datastore.DataStore
}

func NewKeyValueServer(addr string) server {
	ds := datastore.NewDataStore()
	return server{addr, &ds}
}

func (s *server) Run() {
	var err error

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go s.dumpInterval(time.Minute * 2)

	go func() {
		<-c
		s.dump()
		os.Exit(0)
	}()

	s.load()

	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		log.Fatal("Error starting server ", err)
		os.Exit(1)
	}

	log.Println("Listening on ", s.addr)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("Error accepting conn ", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *server) load() {
	err := s.ds.Load()

	if err != nil {
		log.Fatal("Error loading data ", err)
		os.Exit(1)
	}
}

func (s *server) dump() {
	err := s.ds.Dump()

	if err != nil {
		log.Fatal("Error dumping data", err)
		os.Exit(1)
	}
}

func (s *server) dumpInterval(interval time.Duration) {
	for {
		time.Sleep(interval)
		err := s.ds.Dump()

		if err != nil {
			log.Println("Error dumping data", err) // dont want to exit after unsuccessfull save
		}
	}
}

func (s *server) handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)

		if err != nil {
			conn.Write([]byte("MISSING COMMAND"))
		}

		query := strings.Trim(string(buf[0:reqLen]), " \n")

		if reqLen == 0 || query == "exit" {
			log.Println(conn.RemoteAddr(), "EXIT")
			return
		}

		res := s.ds.HandleQuery(query)

		if strings.HasPrefix(res, "-") {
			log.Println(conn.RemoteAddr(), res)
		}

		conn.Write([]byte(res))
	}
}
