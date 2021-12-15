package update

import (
	"os/exec"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/KarolosLykos/archista/utils"
	"github.com/martinohmann/barista-contrib/modules/updates"
	"github.com/martinohmann/barista-contrib/modules/updates/yay"
)

func GetUpdates() *updates.Module {
	return updates.New(yay.New()).Output(func(info updates.Info) bar.Output {
		if info.Updates == 0 {
			return outputs.Pango(pango.Icon("mdi-package-variant"), utils.Spacer)
		}

		return outputs.Pango(
			utils.Spacer,
			pango.Icon("mdi-package-variant-closed"),
			pango.Textf("%d", info.Updates),
			utils.Spacer,
		).Color(colors.Hex("#a04f4f")).OnClick(click.Right(func() { exec.Command("urxvt").Run() }))
	})
}
