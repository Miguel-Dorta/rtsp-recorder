# rtsp-recorder
A very basic/simple RTSP recorder for Linux

### Requirements
- ffmpeg
- User and group called `camrec` (only for standard installation, but it's recommended)

### Build & Installation
##### Build
- Run `build.sh`

##### Standard installation
- Put the built binary in `/opt/bin`
- Create each RTSP Stream config in `/opt/etc/rtsp-recorder/<alias>.conf`
- Put the systemd config in `/etc/systemd/system`
- (Optional) Create a directory for the recordings in `/var/rec`
- Start/Enable the services for each alias

##### "I want the commands, mate"
```shell script
chmod +x build.sh
./build.sh

mkdir -p /opt/bin /opt/etc/rtsp-recorder /var/rec
cp dist/rtsp-recorder /opt/bin
cp configs/rtsp-recorder@.service /etc/systemd/system
nano /opt/etc/rtsp-recorder/alias.conf # And write the config there
systemctl enable rtsp-recorder@alias --now
```

### Config example
```shell script
#CAMERA_ALIAS=alias     # This is not necessary with the standard config because it's defined in the systemd service
RECORDING_TIME=10       # Duration of the recordings, in minutes
RECORDING_TIMEOUT=60    # Timeout for gratefull stop, in seconds
SAVING_PATH=/var/rec    # Path to save the recordings
RTSP_URL='rtsp://user:pass@127.0.0.1:554/stream'    # RTSP stream URL
VERBOSE=false           # Verbose output 
```
