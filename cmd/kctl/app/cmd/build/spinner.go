package build

import (
	"context"
	"fmt"
	"os"
	"time"

	spinner "github.com/schollz/progressbar/v3"
)

type spin struct {
	*spinner.ProgressBar
}

func newSpinner() *spin {
	s := spinner.NewOptions64(
		-1,
		spinner.OptionClearOnFinish(),
		spinner.OptionSetWriter(os.Stderr),
		spinner.OptionSetWidth(10),
		spinner.OptionThrottle(60*time.Millisecond),
		spinner.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		spinner.OptionSpinnerType(14),
		spinner.OptionFullWidth(),
	)
	s.RenderBlank()

	return &spin{s}
}

// Run is intended to run inside a goroutine
func (s *spin) Run(ctx context.Context, done chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		default:
			s.Add(1)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
