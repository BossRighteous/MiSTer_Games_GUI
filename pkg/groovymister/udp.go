package groovymister

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"time"
)

const (
	mtuBlockSize       int16 = 1470
	cmdHeaderClose     byte  = 1
	cmdHeaderInit      byte  = 2
	cmdHeaderSwitchRes byte  = 3
	cmdHeaderBlit      byte  = 6
)

type UdpClient struct {
	host  string
	conn  net.PacketConn
	addr  *net.UDPAddr
	frame uint32
}

func (client *UdpClient) SendPacket(buffer []byte) {
	//fmt.Println("Sending Packet length", len(buffer))
	_, err := client.conn.WriteTo(buffer, client.addr)
	if err != nil {
		panic(err)
	}
}

func (client *UdpClient) SendMTU(buffer []byte) {
	bytesToSend := int32(len(buffer))
	chunkMaxSize := int32(mtuBlockSize)
	var chunkSize int32 = 0
	var offset int32 = 0
	for bytesToSend > 0 {
		chunkSize = chunkMaxSize
		if bytesToSend <= chunkMaxSize {
			chunkSize = bytesToSend
		}
		bytesToSend = bytesToSend - chunkSize
		client.SendPacket(buffer[offset : offset+chunkSize])
		offset += chunkSize
	}
}

func (client *UdpClient) CmdClose() {
	buffer := make([]byte, 1)
	buffer[0] = cmdHeaderClose
	client.SendPacket(buffer)
}

func (client *UdpClient) CmdInit() {
	buffer := make([]byte, 4)
	buffer[0] = cmdHeaderInit
	buffer[1] = 0 // lz4 compression flag
	buffer[2] = 0 // sound rate
	buffer[3] = 0 // sound channel
	client.SendPacket(buffer)
}

func (client *UdpClient) CmdSwitchres(modeline *Modeline) {
	buffer := make([]byte, 26)
	buffer[0] = cmdHeaderSwitchRes
	binary.LittleEndian.PutUint64(buffer[1:9], math.Float64bits(modeline.PixelClock))
	binary.LittleEndian.PutUint16(buffer[9:11], modeline.HActive)
	binary.LittleEndian.PutUint16(buffer[11:13], modeline.HBegin)
	binary.LittleEndian.PutUint16(buffer[13:15], modeline.HEnd)
	binary.LittleEndian.PutUint16(buffer[15:17], modeline.HTotal)
	binary.LittleEndian.PutUint16(buffer[17:19], modeline.VActive)
	binary.LittleEndian.PutUint16(buffer[19:21], modeline.VBegin)
	binary.LittleEndian.PutUint16(buffer[21:23], modeline.VEnd)
	binary.LittleEndian.PutUint16(buffer[23:25], modeline.VTotal)
	buffer[25] = modeline.Interlace
	client.SendPacket(buffer)
}

func (client *UdpClient) CmdBlit(frameBuffer []byte) {
	client.frame++
	buffer := make([]byte, 9)
	buffer[0] = cmdHeaderBlit
	binary.LittleEndian.PutUint32(buffer[1:5], client.frame)
	binary.LittleEndian.PutUint16(buffer[5:7], 0) // vsyncAuto
	buffer[7] = 0                                 // lz4 blockSize & 0xff
	buffer[8] = 0                                 // lz4 blockSize >> 8
	client.SendPacket(buffer)
	start := time.Now()
	client.SendMTU(frameBuffer)
	fmt.Println("blit took", time.Now().Sub(start))
}

func NewUdpClient(host string, port string) UdpClient {
	var client UdpClient
	client.host = host
	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		panic(err)
	}
	client.conn = conn
	client.addr = addr
	return client
}
