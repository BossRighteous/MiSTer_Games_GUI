package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mistergui"
)

func main() {

	modeline := groovymister.ModelineFromString("6.700 320 336 367 426 240 244 247 262")
	fmt.Println(modeline)

	udpClient := groovymister.NewUdpClient("192.168.0.168", "32100")

	udpClient.CmdInit()
	udpClient.CmdSwitchres(modeline)

	last := time.Now()
	var tickDuration int64 = int64(1000000000 / modeline.FrameRate)
	ticker := time.NewTicker(time.Duration(tickDuration))

	quitChan := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	surface := mistergui.NewSurface(modeline.HActive, modeline.VActive, modeline.Interlace)
	surface.Fill(uint8(rand.IntN(255)), uint8(rand.IntN(255)), uint8(rand.IntN(255)))
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
			fmt.Println("elapsed", elapsed)
			udpClient.CmdBlit(surface.BGRbytes())

			gui.OnTick(frameCount, elapsed.Seconds())
			isBlitting = false
		}

	}
	udpClient.CmdClose()
	println("closed successfully")
}
