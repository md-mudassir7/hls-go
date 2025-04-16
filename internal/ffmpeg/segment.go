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

func GenerateTranscodedSegments(videoPath, adPath string, segmentDuration int, transcodeContent bool, transcodeAd bool) error {
	if transcodeContent {
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
	}
	if transcodeAd {
		for i := 0; i < 3; i++ {
			adDir := fmt.Sprintf("segments/ad/%d", i)
			_ = os.MkdirAll(adDir, 0755)
			res := variants[i]
			cmd := exec.Command("ffmpeg", "-i", adPath, "-vf", fmt.Sprintf("scale=%s", res),
				"-c:v", "libx264", "-c:a", "aac",
				"-f", "segment", "-segment_time", fmt.Sprint(segmentDuration),
				"-reset_timestamps", "1", fmt.Sprintf("%s/segment_%%03d.ts", adDir))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
			if err := redis.LoadSegments(adDir, i, true); err != nil {
				return err
			}
		}
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
