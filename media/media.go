package medias

import (
	"fmt"
	"strings"
	"time"

	"barista.run/bar"
	"barista.run/colors"
	"barista.run/group/modal"
	"barista.run/modules/media"
	"barista.run/modules/meta/split"
	"barista.run/modules/volume"
	"barista.run/outputs"
	"barista.run/pango"

	"github.com/KarolosLykos/archista/utils"
)

func New(vol *volume.Module) bar.Module {
	mediaSummary, mediaDetail := player()

	mainModal := modal.New()
	mainModal.
		Mode("media").
		SetOutput(utils.MakeIconOutput("mdi-music")).
		Add(vol, mediaSummary).
		Detail(mediaDetail)

	mm, _ := mainModal.Build()

	return mm
}

func player() (bar.Module, bar.Module) {
	return split.New(media.Auto().Output(mediaFormatFunc), 1)
}

func mediaFormatFunc(m media.Info) bar.Output {
	if m.PlaybackStatus == media.Stopped || m.PlaybackStatus == media.Disconnected {
		return nil
	}
	artist := utils.Truncate(m.Artist, 35)
	title := utils.Truncate(m.Title, 70-len(artist))
	if len(title) < 35 {
		artist = utils.Truncate(m.Artist, 70-len(title))
	}
	var iconAndPosition bar.Output
	if m.PlaybackStatus == media.Playing {
		iconAndPosition = outputs.Repeat(func(time.Time) bar.Output {
			return makeMediaIconAndPosition(m)
		}).Every(time.Second)
	} else {
		iconAndPosition = makeMediaIconAndPosition(m)
	}

	return outputs.Group(iconAndPosition, outputs.Pango(artist, " - ", title))
}

func makeMediaIconAndPosition(m media.Info) (iconAndPosition *pango.Node) {
	if m.PlaybackStatus == media.Playing {
		iconAndPosition = pango.Icon("mdi-pause").Color(colors.Hex("#13ca91"))
	} else {
		iconAndPosition = pango.Icon("mdi-play").Color(colors.Hex("#13ca91"))
	}

	if m.PlaybackStatus == media.Playing {
		iconAndPosition.Append(
			utils.Spacer,
			pango.Textf("%s/", formatMediaTime(m.Position())).Color(colors.Hex("#4f4f4f")),
		)
	}

	if m.PlaybackStatus == media.Paused || m.PlaybackStatus == media.Playing && strings.Contains(m.PlayerName, "chromium") {
		iconAndPosition.Append(
			utils.Spacer,
			pango.Textf("%s", formatMediaTime(m.Length)).Color(colors.Hex("#4f4f4f")),
		)
	}
	return iconAndPosition
}

func formatMediaTime(d time.Duration) string {
	h, m, s := hms(d)
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

func hms(d time.Duration) (h int, m int, s int) {
	h = int(d.Hours())
	m = int(d.Minutes()) % 60
	s = int(d.Seconds()) % 60
	return
}
