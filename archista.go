package main

import (
	"flag"
	"fmt"

	"barista.run"
	"barista.run/pango/icons/mdi"

	"github.com/KarolosLykos/archista/config"
	"github.com/KarolosLykos/archista/cpu"
	"github.com/KarolosLykos/archista/date"
	"github.com/KarolosLykos/archista/docker"
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
	temp := cpu.GetCPUTemp()
	mediaMM := medias.GetMedia(volume)
	lightsCM := lights.GetLights(cfg)
	arch := logo.GetLogo()
	update := updates.GetUpdates()
	dockerModule := docker.New()

	panic(barista.Run(dockerModule, update, lightsCM, source, mediaMM, temp, localDate, localTime, arch))
}
