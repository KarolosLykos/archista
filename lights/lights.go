package lights

import (
	"fmt"
	"os/exec"
	"strconv"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/group/collapsing"
	"barista.run/modules/funcs"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/amimof/huego"
	"golang.org/x/net/context"

	"github.com/KarolosLykos/archista/config"
	"github.com/KarolosLykos/archista/hue"
	"github.com/KarolosLykos/archista/utils"
)

func New(cfg *config.Config) bar.Module {
	return funcs.OnClick(func(sink bar.Sink) {
		b := huego.New(cfg.HUE.Host, cfg.HUE.User)
		_, err := b.GetConfig()
		if err != nil {
			//nolint:errcheck,gosec // just a notification
			exec.Command(
				"notify-send", "-i", "cancel", "Hue module error",
				fmt.Sprintf("Could not get bridge manually, error: %v", err)).
				Run()
			b, err = huego.DiscoverContext(context.Background())
			if err != nil {
				sink.Output(outputs.Pango(
					utils.Spacer,
					pango.Icon("mdi-home-lightbulb-outline").Color(colors.Hex("#a04f4f")),
					utils.Spacer,
				))

				return
			}
		}

		lights := make([]bar.Module, len(cfg.HUE.Lights))
		for i := 0; i < len(cfg.HUE.Lights); i++ {
			lights[i] = getLight(b.Host, cfg.HUE.User, cfg.HUE.Lights[i])
		}

		collapsingModule, g := collapsing.Group(lights...)
		g.ButtonFunc(collapsingButtons(lights))

		collapsingModule.Stream(sink)
	})
}

func getLight(host, user string, id int) *hue.Module {
	m := hue.New(host, user, id)
	return m.Output(lightFormatFunc)
}

func lightFormatFunc(l *huego.Light) bar.Output {
	icon := getIcon(l.IsOn(), l.ModelID)

	switch l.State.Reachable {
	case false:
		return outputs.Pango(
			utils.Spacer,
			pango.Icon("mdi-power-plug-off"),
			pango.Text(getName(l.ModelID, l.ID)),
			utils.Spacer,
		).Color(colors.Hex("#4f4f4f"))
	default:
		if l.IsOn() {
			return outputs.Pango(
				utils.Spacer,
				pango.Icon(icon),
				pango.Text(getName(l.ModelID, l.ID)),
				utils.Spacer,
			).Color(colors.Hex("#6c57299c"))
		} else {
			return outputs.Pango(
				utils.Spacer,
				pango.Icon(icon),
				pango.Text(getName(l.ModelID, l.ID)),
				utils.Spacer,
			).Color(colors.Hex("#4f4f4f"))
		}
	}
}

type state struct {
	on  string
	off string
}

var lightIcons = map[string]state{
	"LCG002": {
		on:  "mdi-spotlight-beam",
		off: "mdi-lightbulb-outline",
	},
	"LCT015": {
		on:  "mdi-lightbulb",
		off: "mdi-lightbulb-outline",
	},
	"LTA004": {
		on:  "mdi-lightbulb",
		off: "mdi-lightbulb-outline",
	},
	"LCL001": {
		on:  "mdi-led-strip",
		off: "mdi-led-strip-variant",
	},
}

func getIcon(state bool, modelID string) string {
	if state {
		return lightIcons[modelID].on
	}

	return lightIcons[modelID].off
}

func getName(modelID string, id int) string {
	switch modelID {
	case "LCG002":
		return "Spot(" + strconv.Itoa(id) + ")"
	case "LCT015", "LTA004":
		return "Bulb(" + strconv.Itoa(id) + ")"
	case "LCL001":
		return "Strip(" + strconv.Itoa(id) + ")"
	default:
		return "Light(" + strconv.Itoa(id) + ")"
	}
}

func collapsingButtons(lights []bar.Module) func(collapsing.Controller) (start bar.Output, end bar.Output) {
	return func(c collapsing.Controller) (start, end bar.Output) {
		if c.Expanded() {
			return outputs.Pango(pango.Icon("mdi-menu-left-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse)),
				outputs.Pango(pango.Icon("mdi-menu-right-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse))
		}
		return outputs.Pango(pango.Icon("mdi-home-lightbulb-outline").Color(colors.Hex("#13ca91"))).
			OnClick(clickHandler(c, lights)), nil
	}
}

func clickHandler(c collapsing.Controller, lights []bar.Module) func(bar.Event) {
	return func(e bar.Event) {
		switch e.Button {
		case bar.ButtonLeft:
			c.Expand()
		case bar.ButtonRight:
			dimLights(lights, 100)
		case bar.ButtonMiddle:
			dimLights(lights, 254)
		}
	}
}

func dimLights(lights []bar.Module, brightness uint8) {
	for _, l := range lights {
		lm, ok := l.(*hue.Module)
		if ok {
			light := lm.GetLight()
			if light != nil {
				if light.IsOn() {
					_ = light.Bri(brightness)
				}
			}
		}
	}
}
