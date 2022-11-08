package docker

import (
	"os/exec"
	"time"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/modules/shell"
	"barista.run/outputs"
	"barista.run/pango"

	"github.com/KarolosLykos/archista/utils"
)

const serviceStatus = "active"

type docker struct {
	m *shell.Module
	s string
}

func New() bar.Module {
	d := &docker{
		m: shell.New("systemctl", "show", "docker.service", "-p", "ActiveState", "--value").
			Every(10 * time.Minute),
	}

	d.formatFunc()

	return d.m
}

func (d *docker) formatFunc() *docker {
	d.m.Output(func(status string) bar.Output {
		color := "#13ca91"
		if status != serviceStatus {
			color = "#a04f4f"
		}

		d.s = status

		return outputs.Pango(
			utils.Spacer,
			pango.Icon("mdi-docker").Color(colors.Hex(color)),
			utils.Spacer,
		).OnClick(click.Middle(d.m.Refresh)).OnClick(click.Left(d.toggle))
	})

	return d
}

func (d *docker) toggle() {
	notificationStatus := "Stopping..."
	switch d.s {
	case serviceStatus:
		//nolint:errcheck // just a notification
		exec.Command("notify-send", "-i", "chronometer", "-h", "int:value:20", "Docker service", notificationStatus).Run()

		_ = exec.Command("sudo", "systemctl", "stop", "docker.service").Run()
		_ = exec.Command("sudo", "systemctl", "stop", "docker.socket").Run()
	default:
		notificationStatus = "Starting..."
		//nolint:errcheck // just a notification
		exec.Command("notify-send", "-i", "chronometer", "Docker service", notificationStatus).Run()
		_ = exec.Command("sudo", "systemctl", "start", "docker.service").Run()
	}

	//nolint:errcheck // just a notification
	exec.Command("notify-send", "-i", "chronometer", "-h", "int:value:50", "Docker service", notificationStatus).Run()

	d.m.Refresh()

	//nolint:errcheck // just a notification
	exec.Command("notify-send", "-t", "1000", "-i", "chronometer", "-h", "int:value:100", "Docker service", notificationStatus).Run()
}
