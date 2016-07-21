package server

import (
	"bytes"
	"errors"
	"net"
	"sync/atomic"

	"github.com/Sirupsen/logrus"
	"github.com/lostz/watchdog/protocol"
)

var baseConnID uint32 = 1212

// Server ...
type Server struct {
	listener net.Listener
	userList map[string]string
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

func (s *Server) readHandshakeResponse(c net.Conn, salt string) error {
	response := &protocol.PacketHandshakeResponse{}
	err := response.FromPacket()
	if err != nil {
		packetErr := protocol.NewPacketErr(err.Error())
		pErr := packetErr.ToPacket()
		if pErr != nil {
			logrus.Printf("write packetErr %s", err.Error())
			return pErr
		}
		return err
	}
	username := response.Username()
	if passwd, find := s.userList[username]; find {
		checkAuth := protocol.CalcPassword([]byte(salt), []byte(passwd))
		if bytes.Equal([]byte(response.AuthResponse()), checkAuth) {
			return nil
		}
		logrus.Printf("readHandshakeResponseau, checkAuth:%s,client_user:%s", checkAuth, username)
		packetErr := protocol.NewDefaultPacketErr(protocol.ER_ACCESS_DENIED_ERROR, username, c.RemoteAddr().String(), "Yes")
		pErr := packetErr.ToPacket()
		if pErr != nil {
			logrus.Printf("write packetErr %s", err.Error())
			return pErr
		}
	}
	return errors.New("auth error")
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
