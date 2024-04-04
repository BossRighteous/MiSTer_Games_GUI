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
)

func main() {
	rand.Seed(time.Now().UnixNano())

	modeline := groovymister.ModelineFromString("6.700 320 336 367 426 240 244 247 262")
	fmt.Println(modeline)

	udpClient := groovymister.NewUdpClient("192.168.0.168")

	udpClient.CmdInit()
	udpClient.CmdSwitchres(modeline)

	last := time.Now()
	var tickDuration int64 = int64(1000000000 / modeline.FrameRate)
	ticker := time.NewTicker(time.Duration(tickDuration))

	quitChan := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	interlaceBool := false
	if modeline.Interlace != 0 {
		interlaceBool = true
	}
	surface := mistergui.NewSurface(modeline.HActive, modeline.VActive, interlaceBool)
	isBlitting := false
	frameCount := uint32(0)
	gui := mistergui.NewGUI()
	gui.Setup(surface, modeline)

	isRunning := true

	for {
		if !isRunning {
			break
		}
		select {
		case <-sigs:
			quitChan <- true
		case <-quitChan:
			ticker.Stop()
			isRunning = false
		case tick := <-ticker.C:
			frameCount++
			if isBlitting {
				continue
			}
			isBlitting = true
			elapsed := tick.Sub(last)
			last = tick
			udpClient.CmdBlit(surface.BGRbytes(true))

			//udpClient.PollInput()

			drawStart := time.Now()

			gui.OnTick(frameCount, elapsed.Seconds())
			fmt.Println("OnTick time", time.Since(drawStart))
			isBlitting = false
		}

	}
	udpClient.CmdClose()
	println("closed successfully")
}
