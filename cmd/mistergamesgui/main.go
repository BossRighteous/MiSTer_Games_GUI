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
	gInputPacket := groovymister.GroovyInputPacket{}

	rand.Seed(time.Now().UnixNano())

	settings := settings.ParseIniSettings("mistergamesgui.ini")
	fmt.Println(settings)

	modeline := groovymister.ModelineFromSettings(settings.Modeline, settings.FrameRate, settings.Interlace)
	fmt.Println(modeline)

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

	frameCount := uint32(0)
	gui := mistergui.GUI{}
	gui.Setup(modeline, settings, &udpClient, quitChan)

	isRunning := true

	fmt.Println("Starting main loop")
	for {
		if !isRunning {
			break
		}
		select {
		case <-signalChan:
			quitChan <- true
		case <-quitChan:
			inputQuitChan <- true
			isRunning = false
		case gInputPacket = <-inputChan:
			frameCount++
			gui.OnTick(mistergui.TickData{
				FrameCount:  frameCount,
				InputPacket: gInputPacket,
			})
		}
	}
	udpClient.CmdClose()
	println("closed successfully")
}
