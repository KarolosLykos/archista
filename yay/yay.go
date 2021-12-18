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

// RateLimiter throttles state updates to once every ~20ms to avoid unexpected behaviour.
var RateLimiter = rate.NewLimiter(rate.Every(20*time.Millisecond), 1)

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
	updates        int
	packageDetails PackageDetails
	lastUpdated    time.Time
}

// New constructs an instance of the hue module.
func New() *Module {
	m := &Module{
		currentInfo: new(value.ErrorValue),
		scheduler:   timing.NewScheduler(),
	}

	m = update(m)

	m.RefreshInterval(1 * time.Hour)

	m.Output(defaultOutput)

	return m
}

// RefreshInterval configures the polling frequency.
func (m *Module) RefreshInterval(interval time.Duration) *Module {
	m.interval = interval
	m.scheduler.Every(interval)

	return m
}

// Output configures a module to display the output of a user-defined function.
func (m *Module) Output(outputFunc func(Yay) bar.Output) *Module {
	m.outputFunc.Set(outputFunc)
	return m
}

// defaultOutput configurea a default bar output
func defaultOutput(y Yay) bar.Output {
	return outputs.Textf("updates").Color(colors.Hex("#00FF00"))
}

// Stream starts the module.
func (m *Module) Stream(s bar.Sink) {
	i, err := m.currentInfo.Get()

	nextInfo, done := m.currentInfo.Subscribe()
	defer done()

	outf := m.outputFunc.Get().(func(Yay) bar.Output)

	nextOutputFunc, done := m.outputFunc.Subscribe()
	defer done()

	for {
		if s.Error(err) {
			return
		} else if info, ok := i.(Yay); ok {
			s.Output(outputs.Group(outf(info)).OnClick(defaultClickHandler(m, info)))
		}

		select {
		case <-nextOutputFunc:
			outf = m.outputFunc.Get().(func(Yay) bar.Output)
		case <-nextInfo:
			i, err = m.currentInfo.Get()
		case <-m.scheduler.C:
			update(m)
		}
	}
}

// defaultClickHandler provides a simple example of the click handler capabilities.
func defaultClickHandler(m *Module, y Yay) func(bar.Event) {
	return func(e bar.Event) {
		if m.err != nil {
			exec.Command("notify-send", "-i", "cancel", "Error", fmt.Sprintf("Error: %v", m.err)).Run()

			return
		}

		if !RateLimiter.Allow() {
			// Don't update the state if it was updated <20ms ago or the light is unreachable
			return
		}

		if y.lastUpdated.After(time.Now().Add(-m.interval)) {
			body := fmt.Sprintf("Last updated at: %s", y.lastUpdated.Format("15:04:05"))
			exec.Command("notify-send", "-i", "chronometer", "Up to date", body).Run()

			return
		}

		if e.Button == bar.ButtonLeft {
			m = update(m)
		}

		m.currentInfo.Set(y)
	}
}

func update(m *Module) *Module {
	y := Yay{
		updates:     0,
		lastUpdated: time.Now(),
	}

	if _, err := exec.Command("yay", "-Sy").CombinedOutput(); err != nil {
		m.currentInfo.Set(y)
		m.err = err

		return m
	}

	output, err := exec.Command("yay", "-Qu").CombinedOutput()
	if err != nil {
		m.currentInfo.Set(y)
		m.err = err

		return m
	}

	details, err := parsePackageDetails(output)
	if err != nil {
		m.currentInfo.Set(y)
		m.err = err

		return m
	}

	y.updates = len(details)
	y.packageDetails = details
	y.lastUpdated = time.Now()
	m.err = nil

	m.currentInfo.Set(y)

	return m
}

// PackageDetail contains information about a single package update.
type PackageDetail struct {
	// PackageName is the name of the package.
	PackageName string
	// CurrentVersion is the currently installed package version.
	CurrentVersion string
	// TargetVersion is the version of the available package update.
	TargetVersion string
}

// PackageDetails contains details about package updates.
type PackageDetails []PackageDetail

// ParsePackageDetails parses package details from pacman compatible output of
// the form "packageName currentVersion -> targetVersion" and returns the
// package details. Returns an error if raw contains malformed lines.
func parsePackageDetails(raw []byte) (PackageDetails, error) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))

	details := PackageDetails{}

	for scanner.Scan() {
		var detail PackageDetail

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

func GetUpdates() *Module {
	m := New().RefreshInterval(1 * time.Hour)
	return m.Output(func(y Yay) bar.Output {
		if y.updates == 0 {
			return outputs.Pango(pango.Icon("mdi-package-variant").Color(colors.Hex("#13ca91")), utils.Spacer)
		}

		return outputs.
			Pango(
				utils.Spacer,
				pango.Icon("mdi-package-variant-closed"),
				pango.Textf("%d", y.updates),
				utils.Spacer,
			).
			Color(colors.Hex("#a04f4f"))
	})
}