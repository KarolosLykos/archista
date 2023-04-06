package bt

import (
	"fmt"
	"os/exec"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/group/collapsing"
	"barista.run/modules/bluetooth"
	"barista.run/outputs"
	"barista.run/pango"

	"github.com/KarolosLykos/archista/config"
)

func New(cfg *config.Config) bar.Module {
	adapter := getAdapter()
	device := getHeadSet(cfg.Bluetooth.Adapter, cfg.Bluetooth.HeadsetMacAddress)

	collapsingModule, g := collapsing.Group(device, adapter)
	g.ButtonFunc(collapsingButtons)

	return collapsingModule
}

func getAdapter() *bluetooth.AdapterModule {
	return bluetooth.DefaultAdapter().
		Output(func(info bluetooth.AdapterInfo) bar.Output {
			if info.Powered {
				switch {
				case info.Discovering:
					return outputs.Pango(pango.Icon("mdi-bluetooth-connect").Color(colors.Hex("#2a52be")))
				default:
					return outputs.
						Pango(pango.Icon("mdi-bluetooth-audio").Color(colors.Hex("#1d3985"))).
						OnClick(click.Left(func() { powerOnOffAdapter("off") }))
				}
			} else {
				return outputs.
					Pango(pango.Icon("mdi-bluetooth-off").Color(colors.Hex("#a04f4f"))).
					OnClick(click.Left(func() { powerOnOffAdapter("on") }))
			}
		})
}

func getHeadSet(adapter, headsetMacAddress string) *bluetooth.DeviceModule {
	return bluetooth.Device(adapter, headsetMacAddress).
		Output(func(info bluetooth.DeviceInfo) bar.Output {
			if info.Paired {
				if info.Connected {
					return outputs.Pango(pango.Icon("mdi-headphones-bluetooth").Color(colors.Hex("#13ca91")))
				} else {
					return outputs.Pango(pango.Icon("mdi-headphones-bluetooth").Color(colors.Hex("#a04f4f")))
				}
			} else {
				return outputs.Pango(pango.Icon("mdi-headphones-box").Color(colors.Hex("#4f4f4f")))
			}
		})
}

func collapsingButtons(c collapsing.Controller) (start, end bar.Output) {
	if c.Expanded() {
		return outputs.Pango(pango.Icon("mdi-menu-left-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse)),
			outputs.Pango(pango.Icon("mdi-menu-right-outline").Color(colors.Hex("#13ca91"))).OnClick(click.Left(c.Collapse))
	}

	return outputs.Pango(
		pango.Icon("mdi-bluetooth").Color(colors.Hex("#13ca91")),
	).OnClick(click.Left(c.Expand)), nil
}

func powerOnOffAdapter(onOff string) {
	if err := exec.Command("bluetoothctl", "power", onOff).Run(); err != nil {
		//nolint:errcheck,gosec // just a notification
		exec.Command(
			"notify-send",
			"-i",
			"cancel",
			fmt.Sprintf("Bluetoothctl power %s error", onOff),
			fmt.Sprintf("Error: %v", err),
		).Run()
	}
}
