package internal

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"os"
)

type ErrorUI struct {
	Error string
}

func NewErrorUI(error string) *ErrorUI {
	return &ErrorUI{
		Error: error,
	}
}

var (
	quitButton = new(widget.Clickable)
	colorBlack = color.NRGBA{A: 255}
	colorWhite = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
)

// Run handles all the events and rendering for the application window.
// It takes a *app.Window as a parameter and returns an error if any.
func (ui *ErrorUI) Run(window *app.Window) error {
	var ops op.Ops
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.ColorOp{Color: color.NRGBA{A: 0xff, R: 0xff}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			errorString := material.Body1(theme, ui.Error)
			errorString.Color = colorBlack
			errorString.Alignment = text.Middle

			for quitButton.Clicked(gtx) {
				os.Exit(0)
			}

			widgets := []layout.Widget{
				func(gtx C) D {
					in := layout.UniformInset(unit.Dp(0))
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return in.Layout(gtx, errorString.Layout)
							}),
						)
					})
				},
				func(gtx C) D {
					in := layout.UniformInset(unit.Dp(0))
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return in.Layout(gtx, material.Button(theme, quitButton, BtnQuit).Layout)
							}),
						)
					})
				},
			}

			material.List(theme, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i])
			})

			e.Frame(gtx.Ops)
		}
	}
}
