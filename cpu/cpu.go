package cpu

import (
	"time"

	"barista.run/bar"
	"barista.run/modules/cputemp"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/KarolosLykos/archista/utils"
	"github.com/martinlindhe/unit"
)

func GetCPUTemp() *cputemp.Module {
	return cputemp.New().
		RefreshInterval(2 * time.Second).
		Output(func(temp unit.Temperature) bar.Output {
			out := outputs.Pango(
				pango.Icon("mdi-fan"), utils.Spacer,
				pango.Textf("%2d℃", int(temp.Celsius())),
			)
			utils.Threshold(out,
				temp.Celsius() > 90,
				temp.Celsius() > 70,
				temp.Celsius() > 60,
			)
			return out
		})
}
