package server

import (
	"log"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	s := &Server{}
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:3306")
	if err != nil {
		log.Println(err.Error())
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	userList := map[string]string{"root": "test"}
	s.userList = userList
	s.listener = listener
	log.Println("running")
	s.Run()

}
