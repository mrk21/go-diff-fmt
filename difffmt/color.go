package difffmt

import (
	"io"

	"github.com/fatih/color"
)

type ColorMode uint8

const (
	ColorNone ColorMode = iota
	ColorAlways
	ColorTerminalOnly
)

type Colorizer struct {
	color *color.Color
}

func NewColorizer(mode ColorMode, attrs ...color.Attribute) *Colorizer {
	c := &Colorizer{
		color: color.New(attrs...),
	}
	c.SetMode(mode)
	return c
}

func (c *Colorizer) SetMode(mode ColorMode) {
	var enabled bool
	switch mode {
	case ColorNone:
		enabled = false
	case ColorAlways:
		enabled = true
	case ColorTerminalOnly:
		enabled = !color.NoColor
	}

	if enabled {
		c.color.EnableColor()
	} else {
		c.color.DisableColor()
	}
}

func (c *Colorizer) Fprint(w io.Writer, a ...interface{}) (int, error) {
	return c.color.Fprint(w, a...)
}

func (c *Colorizer) Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
	return c.color.Fprintf(w, format, a...)
}
