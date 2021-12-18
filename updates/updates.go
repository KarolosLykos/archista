package updates

import (
	"barista.run/bar"
	"barista.run/colors"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/KarolosLykos/archista/utils"
	"github.com/KarolosLykos/archista/yay"
)

func GetUpdates() *yay.Module {
	m := yay.New()
	return m.Output(func(y yay.Yay) bar.Output {
		if y.Updates == 0 {
			return outputs.Pango(pango.Icon("mdi-package-variant").Color(colors.Hex("#13ca91")), utils.Spacer)
		}

		return outputs.
			Pango(
				utils.Spacer,
				pango.Icon("mdi-package-variant-closed"),
				pango.Textf("%d", y.Updates),
				utils.Spacer,
			).
			Color(colors.Hex("#a04f4f"))
	})
}
