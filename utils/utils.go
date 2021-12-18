package utils

import (
	"os/user"
	"path/filepath"

	"barista.run/bar"
	"barista.run/colors"
	"barista.run/outputs"
	"barista.run/pango"
)

func Home(path ...string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	args := append([]string{usr.HomeDir}, path...)
	return filepath.Join(args...)
}

func Truncate(in string, l int) string {
	fromStart := false
	if l < 0 {
		fromStart = true
		l = -l
	}
	inLen := len([]rune(in))
	if inLen <= l {
		return in
	}
	if fromStart {
		return "⋯" + string([]rune(in)[inLen-l+1:])
	}
	return string([]rune(in)[:l-1]) + "⋯"
}

var Spacer = pango.Text(" ").XSmall()

func MakeIconOutput(key string) *bar.Segment {
	return outputs.Pango(Spacer, pango.Icon(key).Color(colors.Hex("#13ca91")), Spacer)
}

func Threshold(out *bar.Segment, urgent bool, color ...bool) *bar.Segment {
	if urgent {
		return out.Urgent(true)
	}
	colorKeys := []string{"#FF0000", "#0000FF", "#00FF00"}
	for i, c := range colorKeys {
		if len(color) > i && color[i] {
			return out.Color(colors.Hex(c))
		}
	}
	return out
}
