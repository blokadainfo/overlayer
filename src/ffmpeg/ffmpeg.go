package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/blokadainfo/overlayer/src/config"
)

func CreateFFmpegCommand(oc config.OverlayConfig, sc config.StreamConfig) *exec.Cmd {
	referenceHeight := 1080 // Video height at which the logo will remain unscaled (exactly width x height pixels as defined in the config)

	cmdArgs := []string{
		"-i", sc.Src,
		"-i", oc.LogoPath,
		"-filter_complex", fmt.Sprintf("[1:v]format=rgba,colorchannelmixer=aa=0.75[logo];[logo][0:v]scale=w=%d*(rh/%d):h=%d*(rh/%d)[logo_scaled];[0:v][logo_scaled]overlay=24*(main_h/%d):(main_h-overlay_h)/2", oc.LogoWidth, referenceHeight, oc.LogoHeight, referenceHeight, referenceHeight),
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
