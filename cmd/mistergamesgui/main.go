package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mistergui"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/settings"
)

func main() {
	quitChan := make(chan bool, 1)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	var frameBuffer []uint8
	gInput := groovymister.GroovyInput{}

	rand.Seed(time.Now().UnixNano())

	settings := settings.ParseIniSettings("mistergamesgui.ini")
	fmt.Println(settings)

	modeline := groovymister.ModelineFromSettings(settings.Modeline, settings.FrameRate, settings.Interlace)
	fmt.Println(modeline)
	frameBuffer = make([]uint8, int(modeline.HActive)*int(modeline.VActive)*mistergui.BGR8BytesPerPixel)

	udpClient := groovymister.NewUdpClient(settings.MiSTerHost, int32(settings.UdpMtuSize))

	udpClient.CmdInit()
	udpClient.CmdSwitchres(modeline)
	inputChan, inputQuitChan := udpClient.PollInput()

	last := time.Now()
	var tickDuration int64 = int64(1000000000 / modeline.FrameRate)
	ticker := time.NewTicker(time.Duration(tickDuration))

	frameCount := uint32(0)
	gui := mistergui.NewGUI()
	gui.Setup(modeline)

	isRunning := true

	for {
		if !isRunning {
			break
		}
		select {
		case <-signalChan:
			quitChan <- true
		case <-quitChan:
			inputQuitChan <- true
			ticker.Stop()
			isRunning = false
		case gInput := <-inputChan:
			fmt.Println("Inputs", gInput)
		case frameBuffer = <-gui.FrameBufferChan:
			fmt.Println("buffer event recv")
			//update frame buffer from gui event
		case tick := <-ticker.C:
			frameCount++
			elapsed := tick.Sub(last)
			last = tick
			udpClient.CmdBlit(frameBuffer)

			gui.TickChan <- mistergui.TickData{
				FrameCount: frameCount,
				Delta:      elapsed.Seconds(),
				Input:      gInput,
			}
			fmt.Println(elapsed.Seconds())
		}

	}
	udpClient.CmdClose()
	println("closed successfully")
}
