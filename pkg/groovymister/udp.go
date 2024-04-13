package groovymister

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"time"
)

const (
	cmdHeaderClose     byte = 1
	cmdHeaderInit      byte = 2
	cmdHeaderSwitchRes byte = 3
	cmdHeaderBlit      byte = 6
)

type UdpClient struct {
	host         string
	conn         net.PacketConn
	addr         *net.UDPAddr
	frame        uint32
	mtuBlockSize int32
}

func (client *UdpClient) SendPacket(buffer []byte) {
	//fmt.Println("Sending Packet length", len(buffer))
	_, err := client.conn.WriteTo(buffer, client.addr)
	if err != nil {
		panic(err)
	}
}

func (client *UdpClient) SendMTU(buffer []byte) {
	//fmt.Println(buffer[3000:3100])
	bytesToSend := int32(len(buffer))
	//fmt.Println(bytesToSend, client.mtuBlockSize)
	chunkMaxSize := int32(client.mtuBlockSize)
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
	buffer := make([]byte, 5)
	buffer[0] = cmdHeaderInit
	buffer[1] = 0 // lz4 compression flag
	buffer[2] = 0 // sound rate
	buffer[3] = 0 // sound channel
	buffer[4] = 0 // rgb mode
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
	buffer[25] = 0
	if modeline.Interlace {
		buffer[25] = 1
	}
	client.SendPacket(buffer)
}

func (client *UdpClient) CmdBlit(frameBuffer []byte) {
	client.frame++
	buffer := make([]byte, 7)
	buffer[0] = cmdHeaderBlit
	binary.LittleEndian.PutUint32(buffer[1:5], client.frame)
	binary.LittleEndian.PutUint16(buffer[5:7], 0) // vsyncAuto
	//buffer[7] = 0                                 // lz4 blockSize & 0xff
	//buffer[8] = 0                                 // lz4 blockSize >> 8
	client.SendPacket(buffer)
	start := time.Now()
	client.SendMTU(frameBuffer)
	fmt.Println("blit took", time.Since(start))
}

func (client *UdpClient) PollInput() (chan GroovyInput, chan bool) {
	inputChan := make(chan GroovyInput, 20)
	inputQuitChan := make(chan bool, 1)
	go func() {
		gInput := GroovyInput{}
		for {
			select {
			case <-inputQuitChan:
				return
			default:
				buf := make([]byte, 9)
				rlen, _, err := client.conn.ReadFrom(buf)
				if err != nil {
					fmt.Println(err)
				}
				if rlen == 9 {
					print(buf)
					nInput := InputFromBuffer(buf)
					if nInput.JoyFrame > gInput.JoyFrame || (nInput.JoyFrame == gInput.JoyFrame && nInput.JoyOrder > gInput.JoyOrder) {
						gInput = nInput
						inputChan <- nInput
					}
				}
			}
		}
	}()
	return inputChan, inputQuitChan
}

func NewUdpClient(host string, mtuBlockSize int32) UdpClient {
	var client UdpClient
	client.host = host
	conn, err := net.ListenPacket("udp4", ":32101")
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:32100", host))
	if err != nil {
		panic(err)
	}
	client.conn = conn
	client.addr = addr
	client.mtuBlockSize = mtuBlockSize
	return client
}
