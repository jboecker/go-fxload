package fxload

/*
Load firmware to an FX2LP device using the 0xA0 "Firmware Load" vendor request

The 0xA0 vendor request is documented on page 57 of the
EZ-USB Technical Reference Manual, Document # 001-13670 Rev. *F
http://www.cypress.com/file/126446/download

Example:

	ctx := gousb.NewContext()
	defer ctx.Close()

	// find an EZ-USB FX2LP Development Board
	dev, err := ctx.OpenDeviceWithVIDPID(0x04b4, 0x8613)
	if err != nil {
		return err
	}

	fxload.DownloadFirmwareFile(dev, "firmware.ihx")
*/

import (
	"os"

	"github.com/kierdavis/ihex-go"

	"github.com/google/gousb"
)

func Poke(dev *gousb.Device, addr uint16, data []byte) {
	dev.Control(gousb.RequestTypeVendor, 0xA0, addr, 0, data)
}

func StopCPU(dev *gousb.Device) {
	Poke(dev, 0xe600, []byte{0x01})
}

func StartCPU(dev *gousb.Device) {
	Poke(dev, 0xe600, []byte{0x00})
}

func DownloadFirmwareFile(dev *gousb.Device, filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//fmt.Println("Stop CPU")
	StopCPU(dev)

	d := ihex.NewDecoder(file)
	//tlen := 0
	for d.Scan() {
		rec := d.Record()
		//fmt.Printf("%v %v %v\n", rec.Type, rec.Address, rec.Data)
		//tlen += len(rec.Data)
		Poke(dev, rec.Address, rec.Data)
	}

	//fmt.Printf("total length %v\n", tlen)

	//fmt.Println("Start CPU")
	StartCPU(dev)
	return nil
}
