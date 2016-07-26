package mysql

import (
	"log"
	"testing"
)

func TestConn(t *testing.T) {
	conn := &Conn{}
	err := conn.Connect("10.88.147.1:6001", "portal", "portal", "")
	if err != nil {
		log.Println(err.Error())
	}

}
