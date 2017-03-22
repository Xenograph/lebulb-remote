package lebulb

import "os"
import "encoding/binary"
import "syscall"
import "unsafe"

const (
	// ioctl request codes (taken from lirc.h)
	lircSetSendCarrierIoctlReq   uintptr = 0x40046913
	lircSetSendDutyCycleIoctlReq uintptr = 0x40046915

	// Carrier wave parameters
	carrierFrequencyHz uint32 = 38000
	dutyCycle          uint32 = 50 // Percentage [0, 100]

	// Command signal parameters
	headerPulseUsec uint32 = 9059
	headerSpaceUsec uint32 = 4458
	zeroPulseUsec   uint32 = 607
	zeroSpaceUsec   uint32 = 517
	onePulseUsec    uint32 = 607
	oneSpaceUsec    uint32 = 1642
	trailPulseUsec  uint32 = 604

	// Device to write commands to
	devicePath string = "/dev/lirc0"
)

// Valid commands which can be sent to the lightbulb
const (
	CmdPowerOff         uint32 = 0x00F740BF
	CmdPowerOn          uint32 = 0x00F7C03F
	CmdBrightnessUp     uint32 = 0x00F700FF
	CmdBrightnessDown   uint32 = 0x00F7807F
	CmdFlashFast        uint32 = 0x00F7D02F
	CmdFlashSlow        uint32 = 0x00F7F00F
	CmdFadeFast         uint32 = 0x00F7C837
	CmdFadeSlow         uint32 = 0x00F7E817
	CmdColorRed         uint32 = 0x00F720DF
	CmdColorGreen       uint32 = 0x00F7A05F
	CmdColorBlue        uint32 = 0x00F7609F
	CmdColorWhite       uint32 = 0x00F7E01F
	CmdColorOrange      uint32 = 0x00F710EF
	CmdColorLightOrange uint32 = 0x00F730CF
	CmdColorYellow      uint32 = 0x00F708F7
	CmdColorLightYellow uint32 = 0x00F728D7
	CmdColorGreen2      uint32 = 0x00F7906F
	CmdColorGreen3      uint32 = 0x00F7B04F
	CmdColorGreen4      uint32 = 0x00F78877
	CmdColorGreen5      uint32 = 0x00F7A857
	CmdColorBlue2       uint32 = 0x00F750AF
	CmdColorPurple      uint32 = 0x00F7708F
	CmdColorLightPurple uint32 = 0x00F748B7
	CmdColorPink        uint32 = 0x00F76897
)

// SendCommand sends the given command to the lightbulb.
func SendCommand(cmd uint32) error {
	device, err := os.OpenFile(devicePath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer device.Close()
	setDeviceParams(device)
	binary.Write(device, binary.LittleEndian, encodeCommand(cmd))
	return nil
}

func setDeviceParams(device *os.File) {
	// Set carrier frequency
	pval := carrierFrequencyHz
	syscall.Syscall(
		syscall.SYS_IOCTL,
		device.Fd(),
		lircSetSendCarrierIoctlReq,
		uintptr(unsafe.Pointer(&pval)))

	// Set duty cycle
	pval = dutyCycle
	syscall.Syscall(
		syscall.SYS_IOCTL,
		device.Fd(),
		lircSetSendDutyCycleIoctlReq,
		uintptr(unsafe.Pointer(&pval)))
}

func encodeCommand(cmd uint32) []uint32 {
	var data []uint32
	data = append(data, headerPulseUsec)
	data = append(data, headerSpaceUsec)
	for ; cmd != 0; cmd <<= 1 {
		if (cmd&0x80000000)>>31 == 0 {
			data = append(data, zeroPulseUsec)
			data = append(data, zeroSpaceUsec)
		} else {
			data = append(data, onePulseUsec)
			data = append(data, oneSpaceUsec)
		}
	}
	data = append(data, trailPulseUsec)
	return data
}

