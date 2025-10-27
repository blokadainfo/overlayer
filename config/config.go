package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Overlay OverlayConfig
	Streams []StreamConfig
}

type OverlayConfig struct {
	LogoPath   string
	LogoHeight int
	LogoWidth  int
}

type StreamConfig struct {
	Src string
	Dst string
}

func LoadConfig() (Config, error) {
	overlayConfig := OverlayConfig{
		LogoPath: os.Getenv("OVERLAY_LOGO_PATH"),
	}

	if overlayConfig.LogoPath == "" {
		return Config{}, fmt.Errorf("OVERLAY_LOGO_PATH is required")
	}

	logoHeight, err := strconv.Atoi(os.Getenv("OVERLAY_LOGO_HEIGHT"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid OVERLAY_LOGO_HEIGHT value: %v", err)
	}
	overlayConfig.LogoHeight = logoHeight

	logoWidth, err := strconv.Atoi(os.Getenv("OVERLAY_LOGO_WIDTH"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid OVERLAY_LOGO_WIDTH value: %v", err)
	}
	overlayConfig.LogoWidth = logoWidth

	var streams []StreamConfig
	for i := 0; ; i++ {
		src := os.Getenv(fmt.Sprintf("STREAMS_%d_SRC", i))
		dst := os.Getenv(fmt.Sprintf("STREAMS_%d_DST", i))
		if src == "" || dst == "" {
			break
		}
		streams = append(streams, StreamConfig{Src: src, Dst: dst})
	}

	return Config{
		Overlay: overlayConfig,
		Streams: streams,
	}, nil
}
