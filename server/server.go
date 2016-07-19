package server

import (
	"net"
	"sync/atomic"

	"github.com/Sirupsen/logrus"
	"github.com/lostz/watchdog/protocol"
)

var baseConnId uint32 = 1212

// Server ...
type Server struct {
	listener net.Listener
	running  bool
}

//NewServer ...
func NewServer() (*Server, error) {
	s := &Server{}
	return s, nil

}

func (s *Server) onConn(c net.Conn) {
	logrus.Println("Accept")
	err := s.handShake(c)
	if err != nil {
	}

}

func (s *Server) handShake(c net.Conn) error {
	salt := protocol.RandomBuf(20)
	err := s.writeInitialHandshake(c, string(salt))
	if err != nil {
		logrus.Printf("Server handshake %s ", err.Error())
		return err

	}
	return nil

}

func (s *Server) writeInitialHandshake(c net.Conn, salt string) error {
	handshake := protocol.NewPacketHandShake(atomic.AddUint32(&baseConnId, 1), salt)
	err := protocol.WritePacket(c, handshake)
	return err
}

func (s *Server) readHandshakeResponse(c net.Conn) error {

}

//Run ...
func (s *Server) Run() error {
	s.running = true
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			logrus.Printf("Can not Accept ")
			continue
		}
		go s.onConn(conn)

	}

	return nil
}
