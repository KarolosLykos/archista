package main

import (
	"barista.run"
	"barista.run/pango/icons/mdi"
	"flag"

	"github.com/KarolosLykos/archista/cpu"
	"github.com/KarolosLykos/archista/date"
	"github.com/KarolosLykos/archista/docker"
	"github.com/KarolosLykos/archista/logo"
	medias "github.com/KarolosLykos/archista/media"
	"github.com/KarolosLykos/archista/sound"
	"github.com/KarolosLykos/archista/utils"
)

func main() {
	if err := mdi.Load(utils.Home("Downloads/MaterialDesign-Webfont")); err != nil {
		panic(err)
	}

	var path string
	flag.StringVar(&path, "config", "", "configuration file path")
	flag.Parse()

	//cfg, err := config.Load(path)
	//if err != nil {
	//	fmt.Println(err)
	//
	//	return
	//}

	localDate := date.GetLocalDate()
	localTime := date.GetLocalTime()
	volume := sound.GetVolume()
	source := sound.GetSource()
	temp := cpu.GetCPUTemp()
	mediaMM := medias.GetMedia(volume)
	//lightsCM := lights.GetLights(cfg)
	arch := logo.GetLogo()
	//update := updates.GetUpdates()
	dockerModule := docker.New()

	panic(barista.Run(dockerModule, source, mediaMM, temp, localDate, localTime, arch))
}
