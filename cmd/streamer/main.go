package main

import (
	"log"

	"github.com/md-mudassir7/hls-go/config"
	"github.com/md-mudassir7/hls-go/internal/ffmpeg"
	"github.com/md-mudassir7/hls-go/internal/redis"
	"github.com/md-mudassir7/hls-go/internal/server"
)

func main() {
	log.Println("[INIT] Segmenting video and loading Redis")
	config := config.LoadConfig()
	redis.FlushAll()

	if err := ffmpeg.GenerateAndLoadSegments("inputs/video.mp4", "inputs/ad.mp4", config.SEGMENT_DURATION, config.PRE_TRANSCODE); err != nil {
		log.Fatalf("Failed to segment and load: %v", err)
	}

	log.Println("[SERVER] Serving HLS at :8081")
	server.Start()
}
