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
				Pango(pango.Icon("mdi-volume-off")).
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
			iconName = "medium"
		} else if pct > 1 {
			iconName = "low"
		}
		return outputs.Pango(
			pango.Icon("mdi-volume-"+iconName).Color(colors.Hex("#13ca91")),
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
				Pango(pango.Icon("mdi-headphones").Color(colors.Hex("#13ca91"))).
				OnClick(click.Left(func() {
					exec.Command("pacmd", "set-sink-port", "@DEFAULT_SINK@", "analog-output-lineout").Run()
				}))
		} else {
			return outputs.
				Pango(pango.Icon("mdi-speaker").Color(colors.Hex("#13ca91"))).
				OnClick(click.Left(func() {
					exec.Command("pacmd", "set-sink-port", "@DEFAULT_SINK@", "analog-output-headphones").Run()
				}))
		}
	})
}
