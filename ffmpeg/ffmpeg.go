package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/blokadainfo/overlayer/config"
)

func CreateFFmpegCommand(oc config.OverlayConfig, sc config.StreamConfig) *exec.Cmd {
	cmdArgs := []string{
		"-i", sc.Src,
		"-i", oc.LogoPath,
		"-filter_complex", fmt.Sprintf("[1:v]scale=%d:%d,format=rgba,colorchannelmixer=aa=0.75[logo];[0:v][logo]overlay=24:(main_h-overlay_h)/2", oc.LogoWidth, oc.LogoHeight),
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-c:a", "aac",
	}

	switch inferProtocol(sc.Dst) {
	case "rtmp":
		cmdArgs = append(cmdArgs, "-f", "flv")
	case "srt":
		cmdArgs = append(cmdArgs, "-f", "mpegts")
	}

	cmdArgs = append(cmdArgs, sc.Dst)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func inferProtocol(url string) string {
	if strings.HasPrefix(url, "rtmp://") || strings.HasPrefix(url, "rtmps://") {
		return "rtmp"
	} else if strings.HasPrefix(url, "srt://") {
		return "srt"
	}

	panic(fmt.Sprintf("Unknown protocol in URL: %s", url))
}
