package server

import (
	"net"
	"sync/atomic"

	"github.com/Sirupsen/logrus"
	"github.com/lostz/watchdog/protocol"
)

var baseConnID uint32 = 1212

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
	packet := protocol.NewPacket(c)
	err := s.writeInitialHandshake(string(salt), packet)
	if err != nil {
		logrus.Printf("Server handshake %s ", err.Error())
		return err
	}
	return nil

}

func (s *Server) writeInitialHandshake(salt string, packet *protocol.Packet) error {
	handshake := protocol.NewPacketHandShake(atomic.AddUint32(&baseConnID, 1), salt, packet)
	err := handshake.ToPacket()
	return err
}

func (s *Server) readHandshakeResponse(c net.Conn) error {
	response := &protocol.PacketHandshakeResponse{}
	err := response.FromPacket()
	if err != nil {
		return err
	}

	return nil
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
