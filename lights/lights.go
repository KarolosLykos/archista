package lights

import (
	"strconv"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/group/collapsing"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/KarolosLykos/archista/config"
	"github.com/KarolosLykos/archista/hue"
	"github.com/KarolosLykos/archista/utils"
	"github.com/amimof/huego"
)

func GetLights(cfg *config.Config) bar.Module {
	light1 := getLight(cfg.HUE.Host, cfg.HUE.User, 1)
	light2 := getLight(cfg.HUE.Host, cfg.HUE.User, 2)
	light3 := getLight(cfg.HUE.Host, cfg.HUE.User, 3)
	light5 := getLight(cfg.HUE.Host, cfg.HUE.User, 5)
	light6 := getLight(cfg.HUE.Host, cfg.HUE.User, 6)
	light7 := getLight(cfg.HUE.Host, cfg.HUE.User, 7)
	light8 := getLight(cfg.HUE.Host, cfg.HUE.User, 8)
	light9 := getLight(cfg.HUE.Host, cfg.HUE.User, 9)

	collapsingModule, g := collapsing.Group(light1, light2, light3, light5, light6, light7, light8, light9)
	g.ButtonFunc(collapsingButtons)

	return collapsingModule
}

func getLight(host, user string, ID int) *hue.Module {
	m := hue.New(host, user, ID)
	return m.Output(func(l *huego.Light) bar.Output {
		return lightFormatFunc(l)
	})
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
