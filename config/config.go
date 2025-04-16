package config

type Config struct {
	SEGMENT_DURATION      int
	PRE_TRANSCODE_CONTENT bool
	PRE_TRANSCODE_AD      bool
	SEGMENTATION_TYPE     string
	INPUT_VIDEO_PATH      string
	AD_URL                string
}

func LoadConfig() Config {
	config := Config{
		SEGMENT_DURATION:      6,
		PRE_TRANSCODE_CONTENT: false,
		PRE_TRANSCODE_AD:      false,
		SEGMENTATION_TYPE:     "ts",
		INPUT_VIDEO_PATH:      "inputs/video.mp4",
		AD_URL:                "inputs/ad.mp4",
	}
	return config
}
