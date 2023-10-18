package capture

import (
	"image"
	"testing"
)

func BenchmarkCapture(b *testing.B) {
	area := Area{
		Name: "bla",
		Borders: image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{800, 600},
		},
	}
	s := NewScreen(
		[]Area{area},
		CaptureConfig{
			Monitor: 1,
		},
	)

	spacings := map[string]int{
		"dense":   1,
		"space 2": 2,
		"space 4": 4,
	}

	for name, space := range spacings {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				monitorImage, _ := s.Capture()
				areaImages := s.CutAreaImages(monitorImage)
				for _, areaImage := range areaImages {
					areaImage.GetColor(space, 0)
				}
			}
		})
	}
}
