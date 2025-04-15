package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/md-mudassir7/hls-go/internal/redis"
)

var variants = map[int]string{
	0: "640x360",
	1: "1280x720",
	2: "1920x1080",
}

func GenerateAndLoadSegments(videoPath, adPath string, segmentDuration int, transcode bool) error {
	if transcode {
		for i := 0; i < 3; i++ {
			res := variants[i]
			dir := fmt.Sprintf("segments/%d", i)
			_ = os.MkdirAll(dir, 0755)
			cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", fmt.Sprintf("scale=%s", res),
				"-c:v", "libx264", "-c:a", "aac",
				"-f", "segment", "-segment_time", fmt.Sprint(segmentDuration),
				"-reset_timestamps", "1", fmt.Sprintf("%s/seg_%d_%%03d.ts", dir, i))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}

		adDir := "segments/ad"
		_ = os.MkdirAll(adDir, 0755)
		cmd := exec.Command("ffmpeg", "-i", adPath, "-c:v", "libx264", "-c:a", "aac", "-f", "mpegts", fmt.Sprintf("%s/segment_000.ts", adDir))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Transcode ad.mp4 to .ts and place in segments/ad/segment_000.ts
	adDir := "segments/ad"
	_ = os.MkdirAll(adDir, 0755)
	cmd := exec.Command("ffmpeg", "-i", adPath, "-y", "-c:v", "libx264", "-c:a", "aac", "-f", "mpegts", fmt.Sprintf("%s/segment_000.ts", adDir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	if err := redis.LoadSegments(adDir, 0, true); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		dir := fmt.Sprintf("segments/%d", i)
		if err := redis.LoadSegments(dir, i, false); err != nil {
			return err
		}
	}
	if err := redis.LoadSegments("segments/ad", 0, true); err != nil {
		return err
	}
	return nil
}
