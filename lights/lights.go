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

func GetLights(cfg *config.Config) bar.Module {
	return funcs.OnClick(func(sink bar.Sink) {
		b, err := huego.DiscoverContext(context.Background())
		if err != nil {
			//nolint:errcheck,gosec // just a notification
			exec.Command("notify-send", "-i", "cancel", "Hue Error", fmt.Sprintf("Error: %v", err.Error())).Run()

			sink.Output(outputs.Pango(
				utils.Spacer,
				pango.Icon("mdi-home-lightbulb-outline").Color(colors.Hex("#4f4f4f")),
				pango.Icon("mdi-restart-alert").Color(colors.Hex("#ff0000")),
				utils.Spacer,
			))

			return
		}

		light1 := getLight(b.Host, cfg.HUE.User, 1)
		light2 := getLight(b.Host, cfg.HUE.User, 2)
		light3 := getLight(b.Host, cfg.HUE.User, 3)
		light5 := getLight(b.Host, cfg.HUE.User, 5)
		light6 := getLight(b.Host, cfg.HUE.User, 6)
		light7 := getLight(b.Host, cfg.HUE.User, 7)
		light8 := getLight(b.Host, cfg.HUE.User, 8)
		light9 := getLight(b.Host, cfg.HUE.User, 9)

		collapsingModule, g := collapsing.Group(light1, light2, light3, light5, light6, light7, light8, light9)
		g.ButtonFunc(collapsingButtons)

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
	case "LCT015":
		return "Bulb(" + strconv.Itoa(id) + ")"
	case "LCL001":
		return "Strip(" + strconv.Itoa(id) + ")"
	default:
		return "Light(" + strconv.Itoa(id) + ")"
	}
}

func collapsingButtons(c collapsing.Controller) (start, end bar.Output) {
	if c.Expanded() {
		return outputs.Pango(pango.Icon("mdi-menu-left-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse)),
			outputs.Pango(pango.Icon("mdi-menu-right-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse))
	}
	return outputs.Pango(pango.Icon("mdi-home-lightbulb-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Expand)), nil
}
