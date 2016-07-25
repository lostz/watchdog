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
	packet, err := s.handShake(c)
	if err != nil {
		packet.Close()
	}
	err = s.writeProxyInfo(packet)
	if err != nil {
		packet.Close()
	}
	s.run(packet)

}

func (s *Server) handShake(c net.Conn) (*protocol.Packet, error) {
	salt := protocol.RandomBuf(20)
	packet := protocol.NewPacket(c)
	err := s.writeInitialHandshake(string(salt), packet)
	if err != nil {
		logrus.Printf("server handshake %s ", err.Error())
		return packet, err
	}
	err = s.readHandshakeResponse(salt, packet)
	if err != nil {
		logrus.Printf("server readHandshakeResponse %s", err.Error())
		return packet, err
	}
	err = s.writeOk(packet)
	if err != nil {
		logrus.Printf("server handShake weireOK %s", err.Error())
		return packet, err
	}
	return packet, nil

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
		if bytes.Equal([]byte(response.AuthResponse()), checkAuth) {
			logrus.Printf("readHandshakeResponse client_user:%s authOk", username)
			return nil
		}
		packetErr := protocol.NewDefaultPacketErr(packet, protocol.ER_ACCESS_DENIED_ERROR, username, packet.Conn().RemoteAddr().String(), "Yes")
		pErr := packetErr.ToPacket()
		if pErr != nil {
			logrus.Printf("write packetErr %s", err.Error())
			return pErr
		}
	}
	return errors.New("auth error")
}

func (s *Server) writeOk(packet *protocol.Packet) error {
	ok := protocol.NewPacketOk(packet)
	return ok.ToPacket()
}

//Client first send select @@version_comment limit 1
func (s *Server) writeProxyInfo(packet *protocol.Packet) error {
	packet.CleanSequenceId()
	packetQuery := protocol.NewPacketQuery(packet)
	err := packetQuery.FromPacket()
	if err != nil {
		logrus.Errorf("read packet query %s", err.Error())
		return err
	}
	packetColumnCount := protocol.NewPacketColumnCount(packet)
	packetColumnCount.SetColumnCount(1)
	err = packetColumnCount.ToPacket()
	if err != nil {
		logrus.Errorf("send  packet column count %s", err.Error())
		return err
	}
	packetColumnDefinition := protocol.DefaultProxyColumnDefinition(packet)
	err = packetColumnDefinition.ToPacket()
	if err != nil {
		logrus.Errorf("send packet column definition %s", err.Error())
		return err
	}
	packetEOF := protocol.NewPacketEOF(packet)
	packetEOF.SetStatus(protocol.SERVER_STATUS_AUTOCOMMIT)
	err = packetEOF.ToPacket()
	if err != nil {
		logrus.Errorf("send  packet eof %s", err.Error())
		return err
	}
	packetResultsetRow := protocol.NewPacketResultsetRow(packet)
	packetResultsetRow.SetData([]byte(protocol.VersionComment))
	err = packetResultsetRow.ToPacket()
	if err != nil {
		logrus.Errorf("send packet resultset row %s", err.Error())
		return err
	}
	packetEOF = protocol.NewPacketEOF(packet)
	packetEOF.SetStatus(protocol.SERVER_STATUS_AUTOCOMMIT)
	err = packetEOF.ToPacket()
	if err != nil {
		logrus.Errorf("send  packet eof %s", err.Error())
		return err
	}

	return nil
}

func (s *Server) run(packet *protocol.Packet) {
	for {
		data, err := packet.ReadPacket()
		if err != nil {
			return
		}
	}

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
