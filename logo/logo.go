package logo

import (
	"os/exec"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/group/collapsing"
	"barista.run/modules/static"
	"barista.run/outputs"
	"barista.run/pango"

	"github.com/KarolosLykos/archista/utils"
)

func GetLogo() bar.Module {
	shutdown := getShutdown()
	restart := getRestart()

	collapsingModule, g := collapsing.Group(restart, shutdown)
	g.ButtonFunc(collapsingButtons)

	return collapsingModule
}

func getShutdown() bar.Module {
	return static.New(outputs.Pango(
		utils.Spacer,
		pango.Icon("mdi-power").Color(colors.Hex("#13ca91")),
		utils.Spacer,
	).OnClick(click.Left(func() {
		//nolint:errcheck // just a notification
		exec.Command("poweroff").Run()
	})))
}

func getRestart() bar.Module {
	return static.New(outputs.Pango(
		utils.Spacer,
		pango.Icon("mdi-restart").Color(colors.Hex("#13ca91")),
		utils.Spacer,
	).OnClick(click.Left(func() {
		//nolint:errcheck // just a notification
		exec.Command("reboot").Run()
	})))
}

func collapsingButtons(c collapsing.Controller) (start, end bar.Output) {
	if c.Expanded() {
		return outputs.Pango(pango.Icon("mdi-menu-left-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse)),
			outputs.Pango(pango.Icon("mdi-menu-right-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse))
	}

	return outputs.Pango(
		utils.Spacer,
		pango.Icon("mdi-arch").Color(colors.Hex("#13ca91")),
		utils.Spacer,
	).OnClick(click.Left(c.Expand)), nil
}
