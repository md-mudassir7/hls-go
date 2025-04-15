# VOD to HLS Linear along with Ad insertion


### Command to start the server
```
go run cmd/streamer/main.go
```

### Then Access the playlist at, it should switch ads every 30 seconds
```
http://localhost:8081/playlist.m3u8
```


## Note
1. Make sure that PRE_TRANSCODE configuration is set to True in config/config.go if running for the first time