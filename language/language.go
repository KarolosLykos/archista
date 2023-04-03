package language

import (
	"barista.run/bar"
	"barista.run/colors"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/martinohmann/barista-contrib/modules/keyboard"
	"github.com/martinohmann/barista-contrib/modules/keyboard/xkbmap"

	"github.com/KarolosLykos/archista/utils"
)

func New() bar.Module {
	return xkbmap.New("uk", "gr").Output(func(l keyboard.Layout) bar.Output {
		return outputs.
			Pango(
				pango.Icon("mdi-keyboard-outline").Color(colors.Hex("#13ca91")),
				utils.Spacer,
				pango.Textf("%s", l.Name).Small(),
			).
			OnClick(func(e bar.Event) {
				switch e.Button {
				case bar.ButtonLeft:
					l.Next()
				case bar.ButtonRight:
					l.Previous()
				}
			})
	})
}
