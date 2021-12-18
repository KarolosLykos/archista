package update

import (
	"time"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/outputs"
	"barista.run/pango"

	"github.com/martinohmann/barista-contrib/modules/updates"
	"github.com/martinohmann/barista-contrib/modules/updates/yay"

	"github.com/KarolosLykos/archista/utils"
)

func GetUpdates() *updates.Module {
	m := updates.
		New(yay.New()).
		Every(30 * time.Minute)

	var refresh = false
	return m.Output(func(info updates.Info) bar.Output {
		if info.Updates == 0 {
			out := outputs.Pango(pango.Icon("mdi-package-variant").Color(colors.Hex("#13ca91")), utils.Spacer).
				OnClick(click.Left(func() {
					refresh = true
					m.Refresh()
				}))

			if refresh {
				refresh = false
				m.Refresh()

				return outputs.Pango(pango.Icon("mdi-reload").Color(colors.Hex("#13ca91")), utils.Spacer)
			}

			return out

		}

		return outputs.
			Pango(
				utils.Spacer,
				pango.Icon("mdi-package-variant-closed"),
				pango.Textf("%d", info.Updates),
				utils.Spacer,
			).
			Color(colors.Hex("#a04f4f"))
	})
}
