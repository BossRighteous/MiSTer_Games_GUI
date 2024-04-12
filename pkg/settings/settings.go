package settings

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type Settings struct {
	MiSTerHost        string
	UdpMtuSize        uint16
	Modeline          string
	FrameRate         float64
	Interlace         bool
	MetaPath          string
	InputSource       byte
	GamesPathOverride string
}

func ParseIniSettings(iniPath string) *Settings {
	settingsDefault := Settings{
		MiSTerHost:        "127.0.0.1",
		UdpMtuSize:        65496, // 2^16 - 40 header bytes on loopback
		Modeline:          "6.700 320 336 367 426 240 244 247 262",
		FrameRate:         60,
		Interlace:         false,
		MetaPath:          "/media/fat/.mistergamesgui",
		InputSource:       0,
		GamesPathOverride: "",
	}

	if iniPath != "" {
		cfg, err := ini.Load(iniPath)
		if err != nil {
			fmt.Println(err)
			return &settingsDefault
		}

		section := cfg.Section("mistergamesgui")

		iniMiSTerHost := section.Key("mister_host").String()
		if iniMiSTerHost != "" {
			settingsDefault.MiSTerHost = iniMiSTerHost
		}

		iniUdpMtuSize, err := section.Key("udp_mtu_size").Uint()
		if err != nil {
			fmt.Println(err)
		}
		if err == nil && iniUdpMtuSize >= 1470 && iniUdpMtuSize <= 65496 {
			settingsDefault.UdpMtuSize = uint16(iniUdpMtuSize)
		}

		iniModeline := section.Key("modeline").String()
		if iniModeline != "" {
			settingsDefault.Modeline = iniModeline
		}

		iniFrameRate, err := section.Key("frame_rate").Float64()
		if err == nil && iniFrameRate >= 0 {
			settingsDefault.FrameRate = iniFrameRate
		}

		iniInterlace, err := section.Key("interlace").Bool()
		if err == nil {
			settingsDefault.Interlace = iniInterlace
		}

		iniInputSource, err := section.Key("input_source").Uint()
		if err == nil && iniInputSource <= 255 {
			settingsDefault.InputSource = byte(iniInputSource)
		}

		iniMetaPath := section.Key("meta_path").String()
		if iniMetaPath != "" {
			settingsDefault.MetaPath = iniMetaPath
		}

		iniGamesPathOverride := section.Key("games_path_override").String()
		if iniGamesPathOverride != "" {
			settingsDefault.GamesPathOverride = iniGamesPathOverride
		}

	}
	return &settingsDefault
}
