// Package  mysql provides ...
package mysql

import (
	"log"
	"testing"
)

func TestPoolConn(t *testing.T) {
	_, err := NewConnPool(10, 20, "10.88.147.1:6004", "testdb", "testdb", "")
	if err != nil {
		log.Println(err.Error())
		return

	}

}
