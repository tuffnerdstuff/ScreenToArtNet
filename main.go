package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var args = struct {
	Mode      *string
	Src       *string
	Dst       *string
	Pause     *int
	Screen    *int
	Spacing   *int
	Threshold *int
	Config    *string
}{
	flag.String("mode", "run", "tool mode {run|preview}"),
	flag.String("src", "", "artnet source"),
	flag.String("dst", "", "artnet destination"),
	flag.Int("pause", 0, "pause time in ms"),
	flag.Int("screen", 0, "screen identifier"),
	flag.Int("spacing", 1, "spacing of pixels for averaging"),
	flag.Int("threshold", 0, "threshold of color (0<255)"),
	flag.String("config", "config.json", "config file"),
}

func parseArgs() {
	if len(os.Args) == 1 {
		flag.PrintDefaults()
		return
	}
	flag.Parse()
}

func handleInterrupts() {
	// Make sure we clean everything up.
	abort := make(chan os.Signal)
	signal.Notify(abort, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-abort
		fmt.Printf("\r%v received, stopping...\n", s)
		os.Exit(0)
	}()
}

func executeMode() error {

	var err error = nil
	switch *args.Mode {
	case "run":
		err = run()
	case "preview":
		err = preview()
	default:
		fmt.Printf("unknown mode: %s", *args.Mode)
	}

	return err
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("encountered error:\n%s\n", err.Error())
		os.Exit(-1)
	}
}

func main() {
	parseArgs()

	handleInterrupts()

	err := executeMode()

	handleError(err)
}
