package disk

import (
	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/format"
	"barista.run/group/collapsing"
	"barista.run/modules/diskio"
	"barista.run/modules/diskspace"
	"barista.run/outputs"
	"barista.run/pango"

	"github.com/KarolosLykos/archista/config"
	"github.com/KarolosLykos/archista/utils"
)

func New(cfg *config.Config) bar.Module {
	modules := make([]bar.Module, 0)

	dm := diskspace.New("/home").Output(func(i diskspace.Info) bar.Output {
		return outputs.Pango(
			pango.Icon("mdi-harddisk").Color(colors.Hex("#13ca91")),
			utils.Spacer,
			pango.Textf("%s", format.IBytesize(i.Available)).Color(colors.Hex("#4f4f4f")))
	})

	modules = append(modules, dm)

	for _, disk := range cfg.Disk.Disks {
		m := diskio.New(disk).Output(func(i diskio.IO) bar.Output {
			return outputs.Pango(
				pango.Icon("mdi-arrow-up-down").Color(colors.Hex("#13ca91")).Rise(-1),
				utils.Spacer,
				pango.Textf("%s", format.IByterate(i.Total())),
			).Color(colors.Hex("#4f4f4f"))
		})

		modules = append(modules, m)
	}

	collapsingModule, g := collapsing.Group(modules...)
	g.ButtonFunc(collapsingButtons)

	return collapsingModule
}

func collapsingButtons(c collapsing.Controller) (start bar.Output, end bar.Output) {
	if c.Expanded() {
		return outputs.Pango(pango.Icon("mdi-menu-left-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse)),
			outputs.Pango(pango.Icon("mdi-menu-right-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse))
	}
	return outputs.Pango(pango.Icon("mdi-harddisk").Color(colors.Hex("#13ca91"))).
		OnClick(click.Left(c.Expand)), nil
}
