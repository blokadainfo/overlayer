package stream

import (
	"log"
	"time"

	"github.com/blokadainfo/overlayer/config"
	"github.com/blokadainfo/overlayer/ffmpeg"
)

func StartStream(oc config.OverlayConfig, sc config.StreamConfig) {
	for {
		cmd := ffmpeg.CreateFFmpegCommand(oc, sc)

		log.Println("Starting FFmpeg process...")
		if err := cmd.Start(); err != nil {
			log.Printf("Error starting FFmpeg: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err := cmd.Wait()
		if err != nil {
			log.Printf("FFmpeg process exited with error: %v", err)
		} else {
			log.Println("FFmpeg process ended successfully.")
		}

		log.Println("Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
