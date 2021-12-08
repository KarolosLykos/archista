package main

import (
	"barista.run"
	"barista.run/group/modal"
	"barista.run/pango/icons/mdi"
	"github.com/KarolosLykos/archista/cpu"
	"github.com/KarolosLykos/archista/date"
	medias "github.com/KarolosLykos/archista/media"
	"github.com/KarolosLykos/archista/sound"
	"github.com/KarolosLykos/archista/utils"
)

func main() {

	mdi.Load(utils.Home("Downloads/MaterialDesign-Webfont"))

	localDate := date.GetLocalDate()
	localTime := date.GetLocalTime()
	volume := sound.GetVolume()
	source := sound.GetSource()
	temp := cpu.GetCPUTemp()
	mediaSummary, mediaDetail := medias.Player()

	mainModal := modal.New()
	mainModal.Mode("media").
		SetOutput(utils.MakeIconOutput("mdi-music")).
		Add(volume, mediaSummary).
		Detail(mediaDetail)

	mm, _ := mainModal.Build()
	panic(barista.Run(source, mm, temp, localDate, localTime))
}
