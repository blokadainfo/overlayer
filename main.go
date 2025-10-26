package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config defines the structure of the YAML config file.
type Config struct {
	Input  string     `yaml:"input"`
	Output string     `yaml:"output"`
	Logo   LogoConfig `yaml:"logo"`
}

type LogoConfig struct {
	Path   string `yaml:"path"`
	Width  int    `yaml:"width"`
	Height int    `yaml:"height"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// inferProtocol infers the protocol (RTMP or SRT) from the stream URL.
func inferProtocol(url string) string {
	if strings.HasPrefix(url, "rtmp://") {
		return "rtmp"
	} else if strings.HasPrefix(url, "srt://") {
		return "srt"
	}
	return "" // Unknown protocol
}

func createFFmpegCommand(config *Config) *exec.Cmd {
	// Determine the protocols based on the URL
	outputProtocol := inferProtocol(config.Output)

	// FFmpeg command to read the input stream and overlay the logo
	cmdArgs := []string{
		"-re",
		"-analyzeduration", "5000000",
		"-probesize", "10000000",
		"-i", config.Input,
		"-i", config.Logo.Path,
		"-filter_complex", fmt.Sprintf("[1:v]scale=%d:%d,format=rgba,colorchannelmixer=aa=0.75[logo];[0:v][logo]overlay=24:(main_h-overlay_h)/2", config.Logo.Width, config.Logo.Height),
		"-c:v", "libx264",
		"-preset", "utrafast",
		"-tune", "zerolatency",
		"-c:a", "copy",
		"-flush_packets", "0",
		"-max_muxing_queue_size", "1024",
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "5",
		"-err_detect", "ignore_err",
		"-loglevel", "warning",
	}

	// Set output format based on protocol
	switch outputProtocol {
	case "rtmp":
		cmdArgs = append(cmdArgs, "-f", "flv") // Use FLV format for RTMP
	case "srt":
		cmdArgs = append(cmdArgs, "-f", "mpegts") // Use MPEG-TS format for SRT
	}

	// Output stream URL
	cmdArgs = append(cmdArgs, config.Output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func startStream(config *Config) {
	for {
		cmd := createFFmpegCommand(config)

		log.Println("Starting FFmpeg process...")
		if err := cmd.Start(); err != nil {
			log.Printf("Error starting FFmpeg: %v", err)
			time.Sleep(5 * time.Second) // Retry after a delay
			continue
		}

		err := cmd.Wait()
		if err != nil {
			log.Printf("FFmpeg process exited with error: %v", err)
		} else {
			log.Println("FFmpeg process ended successfully.")
		}

		log.Println("Retrying in 5 seconds...")
		time.Sleep(5 * time.Second) // Retry after a delay
	}
}

func main() {
	// Load config
	configFile := "config.yaml"
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	// Start the stream handling
	startStream(config)
}
