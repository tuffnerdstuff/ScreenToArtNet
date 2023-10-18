package capture

import (
	"fmt"
	"image"
	"image/color"

	"github.com/kbinani/screenshot"
)

// Screen represents a tiled screen.
type Screen struct {
	// Areas holds the screen areas.
	Areas []Area
	// Borders holds the capturing borders.
	Borders image.Rectangle

	// Config holds the configuration for the screen capture.
}

// CaptureConfig holds the configuration for the screen capture.
type CaptureConfig struct {
	// Spacing holds the averaging spacing.
	Spacing int
	// Threshold holds the averaged color threshold.
	Threshold int

	// Monitor holds the monitor used for capture.
	Monitor int
}

type Area struct {
	Name    string
	Borders image.Rectangle
}

type AreaImage struct {
	Area  Area
	Image *image.RGBA
}

// NewScreen returns a new screen, tiled with the given configuration.
func NewScreen(areas []Area, config CaptureConfig) *Screen {
	return &Screen{
		Areas:   areas,
		Borders: screenshot.GetDisplayBounds(config.Monitor),
	}
}

func (s *Screen) Capture() (monitor *image.RGBA, err error) {
	monitor, err = screenshot.CaptureRect(s.Borders)
	if err != nil {
		return nil, err
	}

	return monitor, nil
}

func (s *Screen) CutAreaImages(monitor *image.RGBA) []AreaImage {
	areas := make([]AreaImage, len(s.Areas))
	for i, b := range s.Areas {
		areas[i].Area = b
		areas[i].Image = monitor.SubImage(b.Borders).(*image.RGBA)
	}
	return areas
}

func (a *AreaImage) GetColor(space int, threshold int) (color.RGBA, error) {
	image := a.Image
	var r uint64
	var g uint64
	var b uint64

	var count uint64 = 1

	if space < 1 {
		return color.RGBA{}, fmt.Errorf("invalid spacing for averaging (%v)", space)
	}
	if threshold < 0 || threshold > 255 {
		return color.RGBA{}, fmt.Errorf("invalid threshold for averaging (%v)", threshold)
	}

	for x := image.Rect.Min.X; x < image.Rect.Max.X; x = x + space {
		for y := image.Rect.Min.Y; y < image.Rect.Max.Y; y = y + space {
			pixel := color.RGBAModel.Convert(image.At(x, y)).(color.RGBA)
			lr := pixel.R
			lg := pixel.G
			lb := pixel.B

			average := (uint32(lr) + uint32(lg) + uint32(lb)) / 3
			if average < uint32(threshold) {
				continue
			}

			r = r + uint64(lr)
			g = g + uint64(lg)
			b = b + uint64(lb)

			count++
		}
	}

	return color.RGBA{
		R: uint8(r / count),
		G: uint8(g / count),
		B: uint8(b / count),
		A: 255,
	}, nil
}
