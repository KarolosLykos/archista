package sound

import (
	"os/exec"
	"strings"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/modules/volume"
	"barista.run/modules/volume/alsa"
	"barista.run/outputs"
	"barista.run/pango"

	"github.com/KarolosLykos/archista/utils"
)

func GetVolume() *volume.Module {
	return volume.New(alsa.DefaultMixer()).Output(func(v volume.Volume) bar.Output {
		if v.Mute {
			return outputs.
				Pango(pango.Icon("mdi-volume-off"), pango.Textf("%s", "MUT")).
				Color(colors.Hex("#FF0000")).
				OnClick(click.Left(func() {
					exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "toggle").Run()
				}))
		}
		iconName := "mute"
		pct := v.Pct()
		if pct > 66 {
			iconName = "high"
		} else if pct > 33 {
			iconName = "low"
		}
		return outputs.Pango(
			pango.Icon("mdi-volume-"+iconName),
			utils.Spacer,
			pango.Textf("%2d%%", pct),
		)
	})
}

func GetSource() *volume.Module {
	return volume.New(alsa.DefaultMixer()).Output(func(v volume.Volume) bar.Output {
		b, _ := exec.Command("pacmd", "list-sinks").Output()
		if strings.Contains(string(b), "active port: <analog-output-headphones>") {
			return outputs.
				Pango(pango.Icon("mdi-headphones")).
				OnClick(click.Left(func() {
					exec.Command("pacmd", "set-sink-port", "1", "analog-output-lineout").Run()
				}))
		} else {
			return outputs.
				Pango(pango.Icon("mdi-speaker")).
				OnClick(click.Left(func() {
					exec.Command("pacmd", "set-sink-port", "1", "analog-output-headphones").Run()
				}))
		}
	})
}
