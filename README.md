Pi is the repository for Go code intended to run on a Raspberry Pi

See the cmd/grovepi directory for the binary which runs on the pi and collects sensor information

The binary takes a flag `--config` which is a path to a config file which is valid json like:
```
{
    "light": "A0",
    "sound": "A1"
}
```

The key is the name of the sensor, the only two implemented are light and sound. The value is pin.

The defined pins are like A0-2 and D1-7.
