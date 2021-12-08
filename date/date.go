package date

import (
	"time"

	"barista.run/bar"
	"barista.run/modules/clock"
	"barista.run/outputs"

	"github.com/KarolosLykos/archista/utils"
)

func GetLocalDate() *clock.Module {
	return clock.Local().
		Output(time.Second, func(now time.Time) bar.Output {
			return outputs.Pango(
				utils.Spacer,
				now.Format("Mon Jan 2"),
			)
		})

}

func GetLocalTime() *clock.Module {
	return clock.Local().
		Output(time.Second, func(now time.Time) bar.Output {
			return outputs.Pango(
				utils.Spacer,
				now.Format("15:04:05"),
			)
		})

}
