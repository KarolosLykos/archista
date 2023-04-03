package main

import (
	"flag"
	"fmt"

	"barista.run"
	"barista.run/pango/icons/mdi"

	"github.com/KarolosLykos/archista/bt"
	"github.com/KarolosLykos/archista/config"
	"github.com/KarolosLykos/archista/cpu"
	"github.com/KarolosLykos/archista/date"
	"github.com/KarolosLykos/archista/disk"
	"github.com/KarolosLykos/archista/docker"
	"github.com/KarolosLykos/archista/language"
	"github.com/KarolosLykos/archista/lights"
	"github.com/KarolosLykos/archista/logo"
	medias "github.com/KarolosLykos/archista/media"
	"github.com/KarolosLykos/archista/sound"
	"github.com/KarolosLykos/archista/updates"
	"github.com/KarolosLykos/archista/utils"
)

func main() {
	if err := mdi.Load(utils.Home("Downloads/MaterialDesign-Webfont")); err != nil {
		panic(err)
	}

	var path string
	flag.StringVar(&path, "config", "", "configuration file path")
	flag.Parse()

	cfg, err := config.Load(path)
	if err != nil {
		fmt.Println(err)

		return
	}

	s, err := sound.New()
	if err != nil {
		panic(err)
	}

	localDate := date.GetLocalDate()
	localTime := date.GetLocalTime()
	volume := s.GetVolume()
	source := s.GetSource()
	temperatureModule := cpu.GetCPUTemp()
	mediaModule := medias.New(volume)
	logoModule := logo.New()
	lightsModule := lights.New(cfg)
	updateModule := updates.New()
	dockerModule := docker.New()
	bluetoothModule := bt.New(cfg)
	diskModule := disk.New(cfg)

	panic(barista.Run(
		dockerModule,
		updateModule,
		lightsModule,
		bluetoothModule,
		source,
		mediaModule,
		diskModule,
		temperatureModule,
		language.New(),
		localDate,
		localTime,
		logoModule,
	))
}
