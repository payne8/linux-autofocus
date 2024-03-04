package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/blackjack/webcam"
)

var device = flag.String("input", "/dev/video4", "Input video device")

func main() {
	getControls()
}

func getControls() {
	flag.Parse()
	cam, err := webcam.Open(*device)
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

	fmap := cam.GetSupportedFormats()
	fmt.Println("Available Formats: ")
	for p, s := range fmap {
		var pix []byte
		for i := 0; i < 4; i++ {
			pix = append(pix, byte(p>>uint(i*8)))
		}
		fmt.Printf("ID:%08x ('%s') %s\n   ", p, pix, s)
		for _, fs := range cam.GetSupportedFrameSizes(p) {
			fmt.Printf(" %s", fs.GetString())
		}
		fmt.Printf("\n")
	}

	var autofocusID webcam.ControlID
	cmap := cam.GetControls()
	fmt.Println("Available controls: ")
	for id, c := range cmap {
		fmt.Printf("ID:%08x %-32s Type: %1d Min: %6d Max: %6d Step: %6d\n", id, c.Name, c.Type, c.Min, c.Max, c.Step)
		if c.Name == "Focus, Auto" {
			autofocusID = id
		}
	}

	autofocusControl, ok := cmap[webcam.ControlID(0x009a090c)]
	if !ok {
		fmt.Println("Autofocus control not found")
	} else {
		fmt.Println("Autofocus control found")
		fmt.Println(autofocusControl)
	}

	// Set autofocus to auto
	err = cam.SetControl(autofocusID, 1)
	if err != nil {
		fmt.Println(fmt.Errorf("error setting autofocus to auto | %w", err))
	} else {
		fmt.Println("Autofocus set to auto")
	}

}

func basicTurnOn() {
	timeout := uint32(5) //5 seconds

	// ...
	cam, err := webcam.Open(*device) // Open webcam
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()
	// ...
	// Setup webcam image format and frame size here (see examples or documentation)
	// ...
	err = cam.StartStreaming()
	if err != nil {
		panic(err.Error())
	}
	for {
		err = cam.WaitForFrame(timeout)

		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Fprint(os.Stderr, err.Error())
			continue
		default:
			panic(err.Error())
		}

		frame, err := cam.ReadFrame()
		if len(frame) != 0 {
			// Process frame
		} else if err != nil {
			panic(err.Error())
		}
	}
}
