package main

import (
	"flag"
	"fmt"

	"barista.run"
	"barista.run/group/collapsing"
	"barista.run/group/modal"
	"barista.run/pango/icons/mdi"

	"github.com/KarolosLykos/archista/config"
	"github.com/KarolosLykos/archista/cpu"
	"github.com/KarolosLykos/archista/date"
	"github.com/KarolosLykos/archista/lights"
	medias "github.com/KarolosLykos/archista/media"
	"github.com/KarolosLykos/archista/sound"
	"github.com/KarolosLykos/archista/utils"
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
	mediaSummary, mediaDetail := medias.Player()

	light1 := lights.GetLight(config.HUE.Host, config.HUE.User, 1)
	light2 := lights.GetLight(config.HUE.Host, config.HUE.User, 2)
	light3 := lights.GetLight(config.HUE.Host, config.HUE.User, 3)
	light5 := lights.GetLight(config.HUE.Host, config.HUE.User, 5)
	light6 := lights.GetLight(config.HUE.Host, config.HUE.User, 6)
	light7 := lights.GetLight(config.HUE.Host, config.HUE.User, 7)
	light8 := lights.GetLight(config.HUE.Host, config.HUE.User, 8)
	light9 := lights.GetLight(config.HUE.Host, config.HUE.User, 9)

	collapingModule, g := collapsing.Group(light1, light2, light3, light5, light6, light7, light8, light9)

	g.ButtonFunc(utils.CollapsingButtons)

	mainModal := modal.New()
	mainModal.
		Mode("media").
		SetOutput(utils.MakeIconOutput("mdi-music")).
		Add(volume, mediaSummary).
		Detail(mediaDetail)

	mm, _ := mainModal.Build()
	if err := barista.Run(collapingModule, source, mm, temp, localDate, localTime); err != nil {
		fmt.Println(err)
	}
}

// package main

// import (
// 	"flag"
// 	"fmt"

// 	"github.com/KarolosLykos/archista/config"
// )

// func main() {
// 	var path string
// 	flag.StringVar(&path, "config", "", "configuration file path")
// 	flag.Parse()

// 	fmt.Println(path)

// 	config, err := config.Load(path)
// 	if err != nil {
// 		fmt.Println(err)

// 		return
// 	}

// 	fmt.Println(config)
// }
