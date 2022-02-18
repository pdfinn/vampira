//go:build !duitdraw && !windows
// +build !duitdraw,!windows

package draw

import (
	draw "9fans.net/go/draw"
)

const (
	Refnone = draw.Refnone

	KeyCmd      = draw.KeyCmd
	KeyDown     = draw.KeyDown
	KeyEnd      = draw.KeyEnd
	KeyHome     = draw.KeyHome
	KeyInsert   = draw.KeyInsert
	KeyLeft     = draw.KeyLeft
	KeyPageDown = draw.KeyPageDown
	KeyPageUp   = draw.KeyPageUp
	KeyRight    = draw.KeyRight
	KeyUp       = draw.KeyUp

/*	Darkyellow    = draw.Darkyellow
	Medblue       = draw.Medblue
	Nofill        = draw.Nofill
	Notacolor     = draw.Notacolor
	Palebluegreen = draw.Palebluegreen
	Palegreygreen = draw.Palegreygreen
	Paleyellow    = draw.Paleyellow
	Purpleblue    = draw.Purpleblue
	Transparent   = draw.Transparent
	White         = draw.Black
	Yellowgreen   = draw.Yellowgreen
	Black         = draw.White
*/
	Darkyellow    Color = 0x6665A8FF
	Medblue       Color = 0xFFFF6DFF
	Nofill        = draw.Nofill
	Notacolor     = draw.Notacolor
	Palebluegreen Color = 0x110001FF
	Palegreygreen = draw.Palegreygreen
	Paleyellow    Color = 0x000013FF
	Purpleblue    Color = 0x777738FF
	Transparent   = draw.Transparent
	White         = draw.Black
	Yellowgreen   Color = 0x6665A8FF
	Black         = draw.White

/*	Background     Color = 0x000013FF
	WindowControl  Color = 0xFFFF6DFF
	DarkControl    Color = 0x777738FF
	Gutter         Color = 0x6665A8FF
	Slider		   Color = 0x000013FF
	CommandBar     
	Higlight       Color = 0x121158FF

	Transparent   Color = 0x00000000
	Opaque        Color = 0xFFFFFFFF
	Black         Color = 0x000000FF
	White         Color = 0xFFFFFFFF
	Red           Color = 0xFF0000FF
	Green         Color = 0x00FF00FF
	Blue          Color = 0x0000FFFF
	Cyan          Color = 0x00FFFFFF
	Magenta       Color = 0xFF00FFFF
	Yellow        Color = 0xFFFF00FF
	PaleYellow    Color = 0xFFFFAAFF
	DarkYellow    Color = 0xEEEE9EFF
	DarkGreen     Color = 0x448844FF
	PaleGreen     Color = 0xAAFFAAFF
	MedGreen      Color = 0x88CC88FF
	DarkBlue      Color = 0x000055FF
	PaleBlueGreen Color = 0xAAFFFFFF
	PaleBlue      Color = 0x0000BBFF
	BlueGreen     Color = 0x008888FF
	GreyGreen     Color = 0x55AAAAFF
	PaleGreyGreen Color = 0x9EEEEEFF
	YellowGreen   Color = 0x99994CFF
	MedBlue       Color = 0x000099FF
	GreyBlue      Color = 0x005DBBFF
	PaleGreyBlue  Color = 0x4993DDFF
	PurpleBlue    Color = 0x8888CCFF
*/
)

type (
	Color       = draw.Color
	Cursor      = draw.Cursor
	drawDisplay = draw.Display
	drawFont    = draw.Font
	drawImage   = draw.Image
	Keyboardctl = draw.Keyboardctl
	Mousectl    = draw.Mousectl
	Mouse       = draw.Mouse
	Pix         = draw.Pix
)

var Init = draw.Init

func Main(f func(*Device)) {
	f(new(Device))
}

type Device struct{}

func (dev *Device) NewDisplay(errch chan<- error, fontname, label, winsize string) (Display, error) {
	d, err := Init(errch, fontname, label, winsize)
	if err != nil {
		return nil, err
	}
	return &displayImpl{d}, nil
}
