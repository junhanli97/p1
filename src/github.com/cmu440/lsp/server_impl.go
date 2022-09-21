// Contains the implementation of a LSP server.

package lsp

import (
	"errors"
	"strconv"

	"github.com/cmu440/lspnet"
)

type cli struct {
	clinetclose chan int
	addr        *lspnet.UDPAddr
}
type ser struct {
	msg     *Message
	addr    *lspnet.UDPAddr
	msgType MsgType
	acked   bool
}
type server struct {
	// TODO: Implement this!
	conn             lspnet.UDPConn
	curid            int
	closeconnrequest chan int //closeConnSignal for client
	clinetclosedone  chan int
	// cliExitDone for client
	closesignal chan int // for server
	//clinetconnclosed chan int
	closecompelet chan int // for server
	sops          chan int
	clients       map[int]cli
	ack           chan int
	MsgConnect    chan ser
	MsgData       chan ser
	MsgAck        chan ser
	readchan      chan ser
	writerequest  chan ser
	writechan     chan ser
}

// NewServer creates, initiates, and returns a new server. This function should
// NOT block. Instead, it should spawn one or more goroutines (to handle things
// like accepting incoming client connections, triggering epoch events at
// fixed intervals, synchronizing events using a for-select loop like you saw in
// project 0, etc.) and immediately return. It should return a non-nil error if
// there was an error resolving or listening on the specified port number.
func NewServer(port int, params *Params) (Server, error) {
	s := new(server)
	s.curid = 0
	s.closeconnrequest = make(chan int)
	s.clinetclosedone = make(chan int)
	s.closesignal = make(chan int)
	//s.clinetconnclosed = make(chan int)
	s.closecompelet = make(chan int)
	s.sops = make(chan int)
	s.clients = make(map[int]cli)
	s.ack = make(chan int)
	s.MsgConnect = make(chan ser)
	s.MsgData = make(chan ser)
	s.MsgAck = make(chan ser)
	s.readchan = make(chan ser)
	s.writechan = make(chan ser)
	s.writerequest = make(chan ser)
	// start listening
	addr, err := lspnet.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, nil
	}
	conn, err := lspnet.ListenUDP("udp", addr)
	if err != nil {
		return nil, nil
	}
	s.conn = *conn
	go serverops(s)
	go serverchan(s)
	return s, nil

}

func serverops(s *server) {
	for {
		select {
		case <-s.sops: //eventHandlersExit
			return

		// case p := <-s.writeSignal:
		case id := <-s.closeconnrequest:
			c, ok := s.clients[id]
			if !ok {
				return
			}
			c.clinetclose <- 1
		case ack := <-s.ack:
			c := s.clients[ack.connId]
		}
	}
}
func serverchan(s *server) {
	m := Message{}
	for {
		switch m.Type {
		case MsgAck:
			s.ack <- 1
		case MsgConnect:

		}

	}
}

func (s *server) Read() (int, []byte, error) {
	// TODO: remove this line when you are ready to begin implementing this method.
	select {} // Blocks indefinitely.
	return -1, nil, errors.New("not yet implemented")
}

func (s *server) Write(connId int, payload []byte) error {
	return errors.New("not yet implemented")
}

func (s *server) CloseConn(connId int) error {
	s.closeconnrequest <- connId //closeConnSignal
	<-s.clinetclosedone          //cliExitDone
	return nil
}

func (s *server) Close() error {
	s.closesignal <- 1
	err := <-s.closecompelet //closeDone
	if err == 1 {
		return errors.New("clients are lost during this time")
	}
	s.sops <- 1
	s.conn.Close()
	return nil
}
