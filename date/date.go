package date

import (
	"time"

	"barista.run/bar"
	"barista.run/colors"
	"barista.run/modules/clock"
	"barista.run/outputs"
	"barista.run/pango"

	"github.com/KarolosLykos/archista/utils"
)

func GetLocalDate() *clock.Module {
	return clock.Local().
		Output(time.Second, func(now time.Time) bar.Output {
			return outputs.Pango(
				utils.Spacer,
				pango.Icon("mdi-calendar-month-outline").Color(colors.Hex("#13ca91")),
				utils.Spacer,
				pango.Text(now.Format("Mon Jan 2")).Small(),
			)
		})
}

func GetLocalTime() *clock.Module {
	return clock.Local().
		Output(time.Second, func(now time.Time) bar.Output {
			return outputs.Pango(
				utils.Spacer,
				pango.Icon("mdi-clock-outline").Color(colors.Hex("#13ca91")),
				utils.Spacer,
				pango.Text(now.Format("15:04:05")).Small(),
			)
		})
}
