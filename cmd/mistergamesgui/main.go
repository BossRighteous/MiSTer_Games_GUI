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
	fmt.Println("Starting main process")
	quitChan := make(chan bool, 1)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	var frameBuffer []uint8
	gInputPacket := groovymister.GroovyInputPacket{}

	rand.Seed(time.Now().UnixNano())

	settings := settings.ParseIniSettings("mistergamesgui.ini")
	fmt.Println(settings)

	modeline := groovymister.ModelineFromSettings(settings.Modeline, settings.FrameRate, settings.Interlace)
	fmt.Println(modeline)
	frameBuffer = make([]uint8, int(modeline.HActive)*int(modeline.VActive)*mistergui.BGR8BytesPerPixel)

	// Load Core and wait
	/*
		if !settings.IsDev {
			coreErr := groovymister.LaunchGroovyCore(settings.GroovyRBFPath)
			if coreErr != nil {
				fmt.Println(coreErr)
				return
			}
			time.Sleep(time.Duration(settings.GroovyClientDelayMS) * time.Millisecond)
		}*/

	udpClient := groovymister.NewUdpClient(settings.MiSTerHost, int32(settings.UdpMtuSize))

	udpClient.CmdInit()
	udpClient.CmdSwitchres(modeline)
	inputChan, inputQuitChan := udpClient.PollInput()

	last := time.Now()
	var tickDuration int64 = int64(1000000000 / modeline.FrameRate)
	ticker := time.NewTicker(time.Duration(tickDuration))

	frameCount := uint32(0)
	gui := mistergui.NewGUI()
	gui.Setup(modeline, settings)

	isRunning := true

	fmt.Println("Starting main loop")
	for {
		if !isRunning {
			break
		}
		select {
		case <-signalChan:
			quitChan <- true
		case <-gui.QuitChan:
			quitChan <- true
		case <-quitChan:
			inputQuitChan <- true
			ticker.Stop()
			isRunning = false
		case gInputPacket = <-inputChan:
			//fmt.Println("Inputs", gInputPacket)
		case frameBuffer = <-gui.FrameBufferChan:
			//fmt.Println("buffer event recv")
			//update frame buffer from gui event
		case tick := <-ticker.C:
			frameCount++
			elapsed := tick.Sub(last)
			last = tick
			udpClient.CmdBlit(frameBuffer)
			gui.TickChan <- mistergui.TickData{
				FrameCount:  frameCount,
				Delta:       elapsed.Seconds(),
				InputPacket: gInputPacket,
			}
		}
	}
	udpClient.CmdClose()
	println("closed successfully")
}
