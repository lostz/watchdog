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
		logrus.Printf("server handshake %s ", err.Error())
		return err
	}
	err = s.readHandshakeResponse(salt, packet)
	if err != nil {
		logrus.Printf("server readHandshakeResponse %s", err.Error())
	}
	return nil

}

func (s *Server) writeInitialHandshake(salt string, packet *protocol.Packet) error {
	handshake := protocol.NewPacketHandShake(atomic.AddUint32(&baseConnID, 1), salt, packet)
	err := handshake.ToPacket()
	return err
}

func (s *Server) readHandshakeResponse(salt []byte, packet *protocol.Packet) error {
	response := &protocol.PacketHandshakeResponse{Packet: packet}
	err := response.FromPacket()
	if err != nil {
		logrus.Printf("readHandshakeResponse %s", err.Error())
		packetErr := protocol.NewPacketErr(packet, err.Error())
		pErr := packetErr.ToPacket()
		if pErr != nil {
			logrus.Printf("write packetErr %s ", err.Error())
			return pErr
		}
		return err
	}
	username := response.Username()
	if passwd, find := s.userList[username]; find {
		checkAuth := protocol.CalcPassword(salt, []byte(passwd))
		logrus.Printf(response.AuthResponse())
		if bytes.Equal([]byte(response.AuthResponse()), checkAuth) {
			return nil
		}
		logrus.Printf("readHandshakeResponseau, passwd:%s,client_user:%s", passwd, username)
		packetErr := protocol.NewDefaultPacketErr(packet, protocol.ER_ACCESS_DENIED_ERROR, username, packet.Conn().RemoteAddr().String(), "Yes")
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
