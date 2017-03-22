# lebulb-remote
[![GoDoc](https://godoc.org/github.com/Xenograph/lebulb-remote?status.svg)](https://godoc.org/github.com/Xenograph/lebulb-remote)

## Overview
Small library I wrote to interface with the **Lighting EVER Dimmable A19 E26 RGB LED Bulb**. It expects that there is an lirc device
at /dev/lirc0. I wrote this to avoid using lircd, which seemed to have issues when sending a sequence of commands in fast
succession. It may work with other Lighting EVER products, I haven't tested any though.

## Usage
Usage is very straightforward:

```go
lebulb.SendCommand(lebulb.CmdPowerOn)
```

