package sound

import (
	"fmt"
	"github.com/KarolosLykos/pulse"
	"github.com/KarolosLykos/pulse/proto"
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
					if err := exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "toggle").
						Run(); err != nil {
						return
					}
				}))
		}
		iconName := "mute"
		pct := v.Pct()
		switch {
		case pct > 66:
			iconName = "high"
		case pct > 33:
			iconName = "medium"
		case pct > 1:
			iconName = "low"
		}

		return outputs.Pango(
			pango.Icon("mdi-volume-"+iconName).Color(colors.Hex("#13ca91")),
			utils.Spacer,
			pango.Textf("%2d%%", pct),
		)
	})
}

type re struct {
	SinkIndex uint32
	Port      string
}

func (r *re) IsReplyTo() uint32 {
	return proto.OpSetSinkPort
}
func GetSource() *volume.Module {
	c, err := pulse.NewClient()
	if err != nil {
		panic(err)
	}

	defaultSink, err := c.DefaultSink()
	if err != nil {
		panic(err)
	}

	sinks, err := c.ListSinks()
	if err != nil {
		panic(err)
	}

	var newSink *pulse.Sink
	for _, s := range sinks {
		if s.Name() != defaultSink.Name() {
			newSink = s
		}
	}

	fmt.Println(newSink)

	if err := c.SetDefaultSink(uint32(proto.Undefined), string(proto.Undefined)); err != nil {
		panic(err)
	}

	if err := c.SetDefaultSink(proto.Undefined, defaultSink.ID()); err != nil {
		panic(err)
	}

	return volume.New(alsa.DefaultMixer()).Output(func(v volume.Volume) bar.Output {
		b, _ := exec.Command("pacmd", "list-sinks").Output()
		if strings.Contains(string(b), "active port: <analog-output-headphones>") {
			return outputs.
				Pango(pango.Icon("mdi-headphones").Color(colors.Hex("#13ca91"))).
				OnClick(click.Left(func() {
					if err := exec.Command("pacmd", "set-sink-port", "@DEFAULT_SINK@", "analog-output-lineout").
						Run(); err != nil {
						return
					}
				}))
		} else {
			return outputs.
				Pango(pango.Icon("mdi-speaker").Color(colors.Hex("#13ca91"))).
				OnClick(click.Left(func() {
					if err := exec.Command("pacmd", "set-sink-port", "@DEFAULT_SINK@", "analog-output-headphones").
						Run(); err != nil {
						return
					}
				}))
		}
	})
}
