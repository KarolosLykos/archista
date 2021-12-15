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
				now.Format("Mon Jan 2"),
			).Color(colors.Hex("#6b6b6b"))
		})

}

func GetLocalTime() *clock.Module {
	return clock.Local().
		Output(time.Second, func(now time.Time) bar.Output {
			return outputs.Pango(
				utils.Spacer,
				now.Format("15:04:05"),
				utils.Spacer,
				pango.Icon("mdi-arch").Color(colors.Hex("#3c7664")),
				utils.Spacer,
			)
		})

}
