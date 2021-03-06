// Package sbvolume displays the current volume as a percentage.
package sbvolume

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var colorEnd = "^d^"

// Routine is the main object for this package.
type Routine struct {
	// Error encountered along the way, if any.
	err error

	// Control to query, as provided by caller.
	control string

	// System volume, in multiples of ten, as percentage of max.
	vol int

	// True if volume is muted.
	muted bool

	// Trio of user-provided colors for displaying various states.
	colors struct {
		normal  string
		warning string
		error   string
	}
}

// New stores the provided control value and makes a new routine object. control is the mixer
// control to monitor. See the man pages for amixer for more information on that. colors is an
// optional triplet of hex color codes for colorizing the output based on these rules:
//   1. Normal color, for normal printing.
//   2. Warning color, for when the volume is muted.
//   3. Error color, for error messages.
func New(control string, colors ...[3]string) *Routine {
	var r Routine

	// Store the color codes. Don't do any validation.
	if len(colors) > 0 {
		r.colors.normal = "^c" + colors[0][0] + "^"
		r.colors.warning = "^c" + colors[0][1] + "^"
		r.colors.error = "^c" + colors[0][2] + "^"
	} else {
		// If a color array wasn't passed in, then we don't want to print this.
		colorEnd = ""
	}

	r.control = control
	return &r
}

// Update runs the 'amixer' command and parses the output for mute status and volume percentage.
func (r *Routine) Update() (bool, error) {
	if r == nil {
		return false, fmt.Errorf("bad routine")
	}

	// Handle error from New.
	if r.control == "" {
		if r.err == nil {
			r.err = fmt.Errorf("invalid control")
		}
		return false, r.err
	}

	r.muted = false
	r.vol = -1

	cmd := exec.Command("amixer", "get", r.control)
	out, err := cmd.Output()
	if err != nil {
		r.err = fmt.Errorf("error getting volume")
		return true, err
	}

	// Find the line that has the percentage volume and mute status in it.
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Playback") && strings.Contains(line, "%]") {
			// We found it. Check the mute status, then pull out the volume.
			fields := strings.Fields(line)
			for _, field := range fields {
				field = strings.Trim(field, "[]")
				if field == "off" {
					r.muted = true
				} else if strings.HasSuffix(field, "%") {
					s := strings.TrimRight(field, "%")
					vol, err := strconv.Atoi(s)
					if err != nil {
						r.err = fmt.Errorf("error parsing volume")
						return true, err
					}
					// Ensure that the volume is a multiple of 10 (so it looks nicer).
					r.vol = (vol + 5) / 10 * 10
				}
			}
			break
		}
	}

	if r.vol < 0 {
		r.err = fmt.Errorf("no volume found")
	}

	return true, nil
}

// String prints either an error, the mute status, or the volume percentage.
func (r *Routine) String() string {
	if r == nil {
		return "bad routine"
	}

	if r.muted {
		return r.colors.warning + "Vol mute" + colorEnd
	}

	return fmt.Sprintf("%sVol %v%%%s", r.colors.normal, r.vol, colorEnd)
}

// Error formats and returns an error message.
func (r *Routine) Error() string {
	if r == nil {
		return "bad routine"
	}

	if r.err == nil {
		r.err = fmt.Errorf("unknown error")
	}

	return r.colors.error + r.err.Error() + colorEnd
}

// Name returns the display name of this module.
func (r *Routine) Name() string {
	return "Volume"
}
