package build

import (
	"fmt"
	"os"
	"time"

	spinner "github.com/schollz/progressbar/v3"
)

func newSpinner() *spinner.ProgressBar {
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

	return s
}
