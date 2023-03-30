package sound

import (
	"strings"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/modules/volume"
	"barista.run/modules/volume/pulseaudio"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/KarolosLykos/pulse"
	"github.com/KarolosLykos/pulse/proto"

	"github.com/KarolosLykos/archista/utils"
)

type Sound struct {
	c          *pulse.Client
	activeSink int
	activePort map[string]int
	sinks      []string
	ports      map[string][]string
}

func New() (*Sound, error) {
	c, err := pulse.NewClient()
	if err != nil {
		return nil, err
	}

	sinks, err := c.ListSinks()
	if err != nil {
		return nil, err
	}

	sink, err := c.DefaultSink()
	if err != nil {
		return nil, err
	}

	availableSinks := make([]string, len(sinks))
	availablePorts := make(map[string][]string)
	activeSink := 0
	activePort := make(map[string]int)

	for i, s := range sinks {
		availableSinks[i] = s.ID()

		if s.ID() == sink.ID() {
			activeSink = i
		}

		for i, p := range s.Info().Ports {
			if p.Name == sink.Info().ActivePortName {
				activePort[sink.ID()] = i
			}
			availablePorts[s.ID()] = append(availablePorts[s.ID()], p.Name)
		}
	}

	return &Sound{c: c, sinks: availableSinks, ports: availablePorts, activePort: activePort, activeSink: activeSink}, nil
}

func (s *Sound) GetVolume() *volume.Module {
	return volume.New(pulseaudio.DefaultSink()).Output(func(v volume.Volume) bar.Output {
		if v.Mute {
			return outputs.
				Pango(pango.Icon("mdi-volume-off")).
				Color(colors.Hex("#FF0000")).
				OnClick(click.Left(func() {
					_ = s.c.SinkMuteToggle(proto.Undefined, s.sinks[s.activeSink])
				}))
		}
		iconName := "mute"
		pct := v.Pct()
		switch {
		case pct > 66:
			iconName = "high"
		case pct > 33:
			iconName = "medium"
		case pct > 1:
			iconName = "low"
		}

		return outputs.Pango(
			pango.Icon("mdi-volume-"+iconName).Color(colors.Hex("#13ca91")),
			utils.Spacer,
			pango.Textf("%2d%%", pct),
		)
	})
}

func (s *Sound) GetSource() *volume.Module {
	return volume.New(pulseaudio.DefaultSink()).Output(func(v volume.Volume) bar.Output {
		return outputs.
			Pango(s.getNode()).
			OnClick(click.Left(s.onClick))
	})
}

func (s *Sound) getNode() *pango.Node {
	sink := s.sinks[s.activeSink]
	switch {
	case strings.Contains(sink, "analog-stereo"):
		if s.ports[sink][s.activePort[sink]] == "analog-output-lineout" {
			return pango.Icon("mdi-speaker").Color(colors.Hex("#13ca91"))
		}

		return pango.Icon("mdi-headphones").Color(colors.Hex("#13ca91"))
	case strings.Contains(sink, "bluez_sink"):
		return pango.Icon("mdi-headphones-bluetooth").Color(colors.Hex("#13ca91"))
	default:
		return pango.Icon("mdi-television").Color(colors.Hex("#13ca91"))
	}
}

func (s *Sound) onClick() {
	sink := s.sinks[s.activeSink]
	switch {
	case strings.Contains(sink, "analog-stereo"):
		if s.ports[sink][s.activePort[sink]] == "analog-output-lineout" {
			_ = s.c.SetSinkPort(proto.Undefined, sink, s.nextPort())

			return
		}

		_ = s.c.SetDefaultSink(s.nextSink())
	default:
		nextSink := s.nextSink()
		_ = s.c.SetSinkPort(proto.Undefined, nextSink, s.nextPort())
		_ = s.c.SetDefaultSink(nextSink)
	}
}

func (s *Sound) nextSink() string {
	s.activeSink++
	if s.activeSink > len(s.sinks)-1 {
		s.activeSink = 0
	}

	return s.sinks[s.activeSink]
}

func (s *Sound) nextPort() string {
	activeSink := s.sinks[s.activeSink]

	s.activePort[activeSink]++
	if s.activePort[activeSink] > len(s.ports[activeSink])-1 {
		s.activePort[activeSink] = 0
	}

	return s.ports[activeSink][s.activePort[activeSink]]
}
