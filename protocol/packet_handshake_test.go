package protocol

import (
	"fmt"
	"testing"
)

var data1 = StringToPacket(`
50 00 00 00 0a 35 2e 36    2e 34 2d 6d 37 2d 6c 6f    P....5.6.4-m7-lo
67 00 56 0a 00 00 52 42    33 76 7a 26 47 72 00 ff    g.V...RB3vz&Gr..
ff 08 02 00 0f c0 15 00    00 00 00 00 00 00 00 00    ................
00 2b 79 44 26 2f 5a 5a    33 30 35 5a 47 00 6d 79    .+yD&/ZZ305ZG.my
73 71 6c 5f 6e 61 74 69    76 65 5f 70 61 73 73 77    sql_native_passw
6f 72 64 00                                           ord
`)

func Test_HandShake(t *testing.T) {
	pk := &Packet{}
	p := &PacketHandshake{}
	p.Packet = pk
	p.FromPacket(data1[4:])
	p.ToPacket()
	fmt.Println(data1)

}
