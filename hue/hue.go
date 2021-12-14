package hue

import (
	"time"

	"barista.run/bar"
	"barista.run/base/value"
	"barista.run/outputs"
	"barista.run/timing"
	"golang.org/x/time/rate"

	"github.com/amimof/huego"
)

// RateLimiter throttles state updates to once every ~20ms to avoid unexpected behaviour.
var RateLimiter = rate.NewLimiter(rate.Every(20*time.Millisecond), 1)

// Module represents a hue bar module. It supports setting the output
// format, click handler, and update frequency.
type Module struct {
	b           *huego.Bridge
	outputFunc  value.Value
	currentInfo *value.ErrorValue
	scheduler   *timing.Scheduler
	id          int
}

// New constructs an instance of the hue module.
func New(host, user string, id int) *Module {
	b := huego.New(host, user)

	m := &Module{
		b:           b,
		id:          id,
		currentInfo: new(value.ErrorValue),
		scheduler:   timing.NewScheduler(),
	}

	m = update(m)

	m.RefreshInterval(5 * time.Minute)

	m.Output(defaultOutput)

	return m
}

// RefreshInterval configures the polling frequency.
func (m *Module) RefreshInterval(interval time.Duration) {
	m.scheduler.Every(interval)
}

// Output configures a module to display the output of a user-defined function.
func (m *Module) Output(outputFunc func(*huego.Light) bar.Output) *Module {
	m.outputFunc.Set(outputFunc)
	return m
}

// defaultOutput configurea a default bar output
func defaultOutput(l *huego.Light) bar.Output {
	return outputs.Textf("id: %d, status: %t, reach: %t", l.ID, l.IsOn(), l.State.Reachable)
}

// Stream starts the module.
func (m *Module) Stream(s bar.Sink) {
	i, err := m.currentInfo.Get()

	nextInfo, done := m.currentInfo.Subscribe()
	defer done()

	outf := m.outputFunc.Get().(func(*huego.Light) bar.Output)

	nextOutputFunc, done := m.outputFunc.Subscribe()
	defer done()

	for {
		if s.Error(err) {
			return
		} else if info, ok := i.(*huego.Light); ok {
			s.Output(outputs.Group(outf(info)).OnClick(defaultClickHandler(m, info)))
		}

		select {
		case <-nextOutputFunc:
			outf = m.outputFunc.Get().(func(*huego.Light) bar.Output)
		case <-nextInfo:
			i, err = m.currentInfo.Get()
		case <-m.scheduler.C:
			update(m)
		}
	}
}

// defaultClickHandler provides a simple example of the click handler capabilities.
func defaultClickHandler(m *Module, light *huego.Light) func(bar.Event) {
	return func(e bar.Event) {
		if !RateLimiter.Allow() || !light.State.Reachable {
			// Don't update the state if it was updated <20ms ago or the light is unreachable
			return
		}

		// Set the light on
		if e.Button == bar.ButtonLeft {
			if light.IsOn() {
				light.Off()
			} else {
				light.On()
			}
		}

		if light.IsOn() {
			// Dim the lights
			if e.Button == bar.ScrollUp && light.State.Bri+10 < 254 {
				light.Bri(light.State.Bri + 10)
			}

			if e.Button == bar.ScrollDown && light.State.Bri-10 >= 1 {
				light.Bri(light.State.Bri - 10)
			}

			// Set maximum brightness
			if e.Button == bar.ButtonRight {
				light.Bri(255)
			}
		}

		m.currentInfo.Set(light)
	}
}

func update(m *Module) *Module {
	l, err := m.b.GetLight(m.id)
	if err != nil {
		m.currentInfo.Error(err)

		return m
	}

	m.currentInfo.Set(l)

	return m
}
