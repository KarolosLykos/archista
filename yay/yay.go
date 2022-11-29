package yay

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"barista.run/bar"
	"barista.run/base/value"
	"barista.run/colors"
	"barista.run/outputs"
	"barista.run/pango"
	"barista.run/timing"
	"golang.org/x/time/rate"

	"github.com/KarolosLykos/archista/utils"
)

// rateLimiter throttles state updates to once every ~20ms to avoid unexpected behaviour.
var rateLimiter = rate.NewLimiter(rate.Every(10*time.Minute), 1)

// Module represents a hue bar module. It supports setting the output
// format, click handler, and update frequency.
type Module struct {
	outputFunc  value.Value
	currentInfo *value.ErrorValue
	scheduler   *timing.Scheduler
	interval    time.Duration
	err         error
}

type Yay struct {
	Updates        int
	packageDetails packageDetails
	lastUpdated    time.Time
}

// New constructs an instance of the hue module.
func New() *Module {
	m := &Module{
		currentInfo: new(value.ErrorValue),
		scheduler:   timing.NewScheduler(),
	}

	m = update(m)

	m.refreshInterval(1 * time.Hour)

	m.Output(defaultOutput)

	return m
}

// refreshInterval configures the polling frequency.
func (m *Module) refreshInterval(interval time.Duration) *Module {
	m.interval = interval
	m.scheduler.Every(interval)

	return m
}

// Output configures a module to display the output of a user-defined function.
func (m *Module) Output(outputFunc func(Yay) bar.Output) *Module {
	m.outputFunc.Set(outputFunc)
	return m
}

// defaultOutput configures a default bar output.
func defaultOutput(_ Yay) bar.Output {
	return outputs.Textf("updates").Color(colors.Hex("#00FF00"))
}

// Stream starts the module.
func (m *Module) Stream(s bar.Sink) {
	var info Yay
	i, err := m.currentInfo.Get()

	nextInfo, done := m.currentInfo.Subscribe()
	defer done()

	outf, ok := m.outputFunc.Get().(func(Yay) bar.Output)
	if !ok {
		return
	}

	nextOutputFunc, done := m.outputFunc.Subscribe()
	defer done()

	for {
		if s.Error(err) {
			s.Output(outputs.Pango(pango.Icon("mdi-update").Color(colors.Hex("#13ca91")), utils.Spacer).
				OnClick(func(e bar.Event) {
					if e.Button == bar.ButtonLeft {
						update(m)
					}
					if e.Button == bar.ButtonRight {
						//nolint:errcheck,gosec // just a notification
						exec.Command("notify-send", "-i", "cancel", "Error", fmt.Sprintf("Error: %v", m.err)).Run()
					}
				}))

			return
		} else if info, ok = i.(Yay); ok {
			s.Output(outputs.Group(outf(info)).OnClick(defaultClickHandler(m, info)))
		}

		select {
		case <-nextOutputFunc:
			outf, ok = m.outputFunc.Get().(func(Yay) bar.Output)
			if !ok {
				return
			}
		case <-nextInfo:
			i, err = m.currentInfo.Get()
		case <-m.scheduler.C:
			update(m)
		}
	}
}

// defaultClickHandler provides a simple example of the click handler capabilities.
//
//nolint:gosec // notify-send required in docs.
func defaultClickHandler(m *Module, y Yay) func(bar.Event) {
	return func(e bar.Event) {
		if m.err != nil {
			//nolint:errcheck // just a notification
			exec.Command("notify-send", "-i", "cancel", "Error", fmt.Sprintf("Error: %v", m.err)).Run()

			return
		}

		if e.Button == bar.ButtonMiddle && y.Updates > 0 {
			s := ""
			for _, p := range y.packageDetails {
				s += fmt.Sprintf("%s(%s) -> %s", p.PackageName, p.CurrentVersion, p.TargetVersion)
				s += "\n"
			}

			s = strings.TrimSuffix(s, "\n")

			//nolint:errcheck // just a notification
			exec.Command("notify-send", "-i", "view-process-tree", "Packages", s).Run()

			return
		}

		if e.Button == bar.ButtonRight {
			body := fmt.Sprintf("Last updated at: %s", y.lastUpdated.Format("15:04:05"))

			//nolint:errcheck // just a notification
			exec.Command("notify-send", "-i", "chronometer", "Forced update", body).Run()
			m = update(m)

			return
		}

		if !rateLimiter.Allow() && m.err == nil {
			// Don't update the state if it was updated <10m ago.
			body := fmt.Sprintf("Last updated at: %s", y.lastUpdated.Format("15:04:05"))

			//nolint:errcheck // just a notification
			exec.Command("notify-send", "-i", "chronometer", "Rate limited", body).Run()

			return
		}

		if e.Button == bar.ButtonLeft {
			if y.lastUpdated.After(time.Now().Add(-m.interval / 2)) {
				body := fmt.Sprintf("Last updated at: %s", y.lastUpdated.Format("15:04:05"))
				//nolint:errcheck // just a notification
				exec.Command("notify-send", "-i", "chronometer", "Up to date", body).Run()

				return
			}

			m = update(m)
		}

		m.currentInfo.Set(y)
	}
}

func update(m *Module) *Module {
	body := fmt.Sprintf("Current update time: %s", time.Now().Format("15:04:05"))

	y := Yay{
		Updates:     0,
		lastUpdated: time.Now(),
	}

	//nolint:errcheck // just a notification
	exec.Command("notify-send", "-i", "chronometer", "-h", "int:value:20", "Updating...", body).Run()
	if _, err := exec.Command("yay", "-Sy").CombinedOutput(); err != nil {
		m.currentInfo.Set(y)
		m.err = err

		return m
	}

	//nolint:errcheck // just a notification
	exec.Command("notify-send", "-i", "chronometer", "-h", "int:value:40", "Updating...", body).Run()
	output, err := exec.Command("yay", "-Qu").CombinedOutput()
	if err != nil {
		m.currentInfo.Set(y)
		m.err = err

		return m
	}

	//nolint:errcheck // just a notification
	exec.Command("notify-send", "-i", "chronometer", "-h", "int:value:60", "Updating...", body).Run()
	details, err := parsePackageDetails(output)
	if err != nil {
		m.currentInfo.Set(y)
		m.err = err

		return m
	}

	y.Updates = len(details)
	y.packageDetails = details
	y.lastUpdated = time.Now()
	m.err = nil

	//nolint:errcheck // just a notification
	exec.Command("notify-send", "-i", "chronometer", "-h", "int:value:80", "Updating...", body).Run()
	m.currentInfo.Set(y)

	//nolint:errcheck // just a notification
	exec.Command("notify-send", "-i", "chronometer", "-t", "800", "-h", "int:value:100", "Updating...", body).Run()

	return m
}

// packageDetail contains information about a single package update.
type packageDetail struct {
	// PackageName is the name of the package.
	PackageName string
	// CurrentVersion is the currently installed package version.
	CurrentVersion string
	// TargetVersion is the version of the available package update.
	TargetVersion string
}

// packageDetails contains details about package updates.
type packageDetails []packageDetail

// ParsePackageDetails parses package details from pacman compatible output of
// the form "packageName currentVersion -> targetVersion" and returns the
// package details. Returns an error if raw contains malformed lines.
func parsePackageDetails(raw []byte) (packageDetails, error) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))

	details := packageDetails{}

	for scanner.Scan() {
		var detail packageDetail

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		_, err := fmt.Sscanf(line, "%s %s -> %s", &detail.PackageName, &detail.CurrentVersion, &detail.TargetVersion)
		if err != nil {
			return nil, err
		}

		details = append(details, detail)
	}

	return details, nil
}
