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
	inset      = layout.UniformInset(unit.Dp(0))
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

			for quitButton.Clicked(gtx) {
				os.Exit(0)
			}

			widgets := []layout.Widget{
				makeLabelWidget(theme, InfoErrorMessage, text.Middle, colorBlack),
				makeLabelWidget(theme, ui.Error, text.Middle, colorBlack),
				makeButtonWidget(theme),
			}

			material.List(theme, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i])
			})

			e.Frame(gtx.Ops)
		}
	}
}

func makeLabelWidget(theme *material.Theme, in string, align text.Alignment, textColor color.NRGBA) layout.Widget {
	label := makeText(theme, in, align, textColor)
	return func(gtx C) D {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return inset.Layout(gtx, label.Layout)
				}),
			)
		})
	}
}

func makeButtonWidget(theme *material.Theme) layout.Widget {
	return func(gtx C) D {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return inset.Layout(gtx, material.Button(theme, quitButton, BtnQuit).Layout)
				}),
			)
		})
	}
}

func makeText(theme *material.Theme, in string, align text.Alignment, textColor color.NRGBA) material.LabelStyle {
	str := material.Body1(theme, in)
	str.Color = textColor
	str.Alignment = align

	return str
}
