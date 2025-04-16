package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/md-mudassir7/hls-go/internal/redis"
)

func Start() {
	http.Handle("/playlist.m3u8", withCORS(http.HandlerFunc(masterHandler)))
	for i := 0; i < 3; i++ {
		idx := i
		http.Handle(fmt.Sprintf("/variant_%d.m3u8", i), withCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			variantHandler(w, r, idx)
		})))
	}
	http.Handle("/segments/", withCORS(http.StripPrefix("/segments/", http.FileServer(http.Dir("segments")))))
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func masterHandler(w http.ResponseWriter, r *http.Request) {
	playlist := "#EXTM3U\n"
	playlist += "#EXT-X-VERSION:3\n"
	playlist += "#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360\n/variant_0.m3u8\n"
	// playlist += "#EXT-X-STREAM-INF:BANDWIDTH=2800000,RESOLUTION=1280x720\n/variant_1.m3u8\n"
	// playlist += "#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080\n/variant_2.m3u8\n"
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	log.Println("Serving Master playlist", playlist)
	_, _ = w.Write([]byte(playlist))
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func variantHandler(w http.ResponseWriter, r *http.Request, variant int) {
	playlist := "#EXTM3U\n"
	playlist += "#EXT-X-VERSION:3\n"
	playlist += "#EXT-X-TARGETDURATION:6\n"

	mode := redis.GetCurrentMode()
	index := redis.GetVariantIndex(variant)

	playlist += fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", index)

	if mode == "ad" {
		seg, err := redis.GetAdSegment()
		if err == nil && seg != "" {
			playlist += "#EXT-X-DISCONTINUITY\n"
			for i := 0; i < 3; i++ {
				playlist += "#EXTINF:6.0,\n"
				playlist += fmt.Sprintf("/segments/ad/%s\n", seg)
			}
		}
	} else {
		segs, err := redis.GetSegmentsForVariant(variant, index, 3)
		if err == nil && len(segs) > 0 {
			for _, s := range segs {
				playlist += "#EXTINF:6.0,\n"
				playlist += fmt.Sprintf("/segments/%d/%s\n", variant, s)
			}
			redis.IncrementVariantIndex(variant, 3)
		}
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	log.Println("Serving variant playlist", playlist)
	_, _ = w.Write([]byte(strings.TrimSpace(playlist)))
}
