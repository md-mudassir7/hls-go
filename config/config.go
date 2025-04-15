package config

type Config struct {
	SEGMENT_DURATION int
	PRE_TRANSCODE    bool
}

func LoadConfig() Config {
	config := Config{
		SEGMENT_DURATION: 6,
		PRE_TRANSCODE:    false,
	}
	return config
}
