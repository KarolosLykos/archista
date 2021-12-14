package lights

import (
	"strconv"

	"barista.run/bar"
	"barista.run/colors"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/KarolosLykos/archista/hue"
	"github.com/KarolosLykos/archista/utils"
	"github.com/amimof/huego"
)

func GetLight(host, user string, ID int) *hue.Module {
	m := hue.New(host, user, ID)
	return m.Output(func(l *huego.Light) bar.Output {
		return lightFormatFunc(l)
	})
}

func lightFormatFunc(l *huego.Light) bar.Output {
	icon := GetIcon(l.IsOn(), l.ModelID)

	switch l.State.Reachable {
	case false:
		return outputs.Pango(
			utils.Spacer,
			pango.Icon("mdi-power-plug-off"),
			pango.Text(GetName(l.ModelID, l.ID)),
			utils.Spacer,
		).Color(colors.Hex("#4f4f4f"))
	default:
		if l.IsOn() {
			return outputs.Pango(
				utils.Spacer,
				pango.Icon(icon),
				pango.Text(GetName(l.ModelID, l.ID)),
				utils.Spacer,
			).Color(colors.Hex("#6c57299c"))
		} else {
			return outputs.Pango(
				utils.Spacer,
				pango.Icon(icon),
				pango.Text(GetName(l.ModelID, l.ID)),
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

func GetIcon(state bool, modelID string) string {
	if state {
		return lightIcons[modelID].on
	}

	return lightIcons[modelID].off
}

func GetName(modelID string, id int) string {
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
