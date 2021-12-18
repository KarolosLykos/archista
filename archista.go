package main

import (
	"flag"
	"fmt"

	"barista.run"
	"barista.run/pango/icons/mdi"

	"github.com/KarolosLykos/archista/config"
	"github.com/KarolosLykos/archista/cpu"
	"github.com/KarolosLykos/archista/date"
	"github.com/KarolosLykos/archista/lights"
	"github.com/KarolosLykos/archista/logo"
	medias "github.com/KarolosLykos/archista/media"
	"github.com/KarolosLykos/archista/sound"
	"github.com/KarolosLykos/archista/utils"
	"github.com/KarolosLykos/archista/yay"
)

func main() {
	mdi.Load(utils.Home("Downloads/MaterialDesign-Webfont"))

	var path string
	flag.StringVar(&path, "config", "", "configuration file path")
	flag.Parse()

	config, err := config.Load(path)
	if err != nil {
		fmt.Println(err)

		return
	}

	localDate := date.GetLocalDate()
	localTime := date.GetLocalTime()
	volume := sound.GetVolume()
	source := sound.GetSource()
	temp := cpu.GetCPUTemp()
	mediaMM := medias.GetMedia(volume)
	lightsCM := lights.GetLights(config)
	// updates := update.GetUpdates()
	arch := logo.GetLogo()

	y := yay.GetUpdates()

	panic(barista.Run(y, lightsCM, source, mediaMM, temp, localDate, localTime, arch))
}
