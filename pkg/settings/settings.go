package settings

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type Settings struct {
	MiSTerHost          string
	UdpMtuSize          uint16
	Modeline            string
	FrameRate           float64
	Interlace           bool
	CollectionsPath     string
	GroovyRBFPath       string
	GroovyCoreDelayMS   int
	LoadGroovyCore      bool
	GroovyClientDelayMS int
	IsDev               bool
}

func ParseIniSettings(iniPath string) *Settings {
	settingsDefault := Settings{
		MiSTerHost:          "127.0.0.1",
		UdpMtuSize:          65496, // 2^16 - 40 header bytes on loopback
		Modeline:            "6.700 320 336 367 426 240 244 247 262",
		FrameRate:           60,
		Interlace:           false,
		CollectionsPath:     "/media/fat/Scripts/mistergamesgui/collections",
		GroovyRBFPath:       "/media/fat/_Utility/Groovy_20240912.rbf",
		GroovyCoreDelayMS:   2000,
		LoadGroovyCore:      true,
		GroovyClientDelayMS: 10000,
		IsDev:               false,
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

		iniCollectionsPath := section.Key("collections_path").String()
		if iniCollectionsPath != "" {
			settingsDefault.CollectionsPath = iniCollectionsPath
		}

		iniGroovyRbfPath := section.Key("groovy_rbf_path").String()
		if iniGroovyRbfPath != "" {
			settingsDefault.GroovyRBFPath = iniGroovyRbfPath
		}

		iniCoreDelayMs, err := section.Key("groovy_core_delay_ms").Int()
		if err == nil && iniCoreDelayMs > 0 {
			settingsDefault.GroovyCoreDelayMS = iniCoreDelayMs
		}

		iniLoadCore, err := section.Key("load_groovy_core").Bool()
		if err == nil {
			settingsDefault.LoadGroovyCore = iniLoadCore
		}

		iniClientDelayMs, err := section.Key("groovy_client_delay_ms").Int()
		if err == nil && iniClientDelayMs > 0 {
			settingsDefault.GroovyClientDelayMS = iniClientDelayMs
		}

		iniIsDev, err := section.Key("is_dev").Bool()
		if err == nil {
			settingsDefault.IsDev = iniIsDev
		}

	}
	return &settingsDefault
}
