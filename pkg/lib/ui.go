package lib

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"log"
)

const (
	stopped int = iota
	started
)

var (
	exploreButton     = new(widget.Clickable)
	launchButton      = new(widget.Clickable)
	backupButton      = new(widget.Clickable)
	restoreButton     = new(widget.Clickable)
	autoLaunch        = new(widget.Bool)
	numBackups        = new(widget.Float)
	autoLaunchChecked = false
	phase             = stopped
	list              = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
)

type (
	D = layout.Dimensions
	C = layout.Context
)

func Run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			window.Option(
				app.MaxSize(unit.Dp(640), unit.Dp(105)),
				app.MinSize(unit.Dp(640), unit.Dp(105)),
			)
			numBackups.Value = ConfigNumBackupsToKeep / 100.00

			if autoLaunch.Update(gtx) {
				autoLaunchChecked = !autoLaunchChecked
				log.Printf("autolaunch set to %t", autoLaunchChecked)
			}

			for exploreButton.Clicked(gtx) {
				err := LaunchExplorer()
				if err != nil {
					log.Printf("error launching explorer: %v", err)
				}
			}

			for launchButton.Clicked(gtx) {
				err := LaunchNoita(true)
				if err != nil {
					log.Printf("failed to launch noita: %v", err)
				}
			}

			for restoreButton.Clicked(gtx) {
				if !isNoitaRunning() {
					if err := RestoreNoita("", true); err != nil {
						log.Printf("failed to restore latest noita save: %v", err)
					} else {
						log.Print("successful restore request")
					}
				} else {
					log.Printf("noita.exe cannot be running to restore")
				}
			}

			for backupButton.Clicked(gtx) {
				BackupNoita(true)
			}

			cl := clip.Rect{Max: e.Size}.Push(gtx.Ops)
			// TODO: make this not run every frame!
			if isNoitaRunning() {
				paint.ColorOp{Color: color.NRGBA{A: 0xff, R: 0xff}}.Add(gtx.Ops)
			} else {
				paint.ColorOp{Color: color.NRGBA{A: 0xff, G: 0xff}}.Add(gtx.Ops)
			}
			paint.PaintOp{}.Add(gtx.Ops)

			widgets := []layout.Widget{
				func(gtx C) D {
					in := layout.UniformInset(unit.Dp(8))
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return in.Layout(gtx, material.Button(theme, launchButton, "Launch Noita").Layout)
							}),
							layout.Rigid(func(gtx C) D {
								return in.Layout(gtx, material.Button(theme, backupButton, "Backup Noita").Layout)
							}),
							layout.Rigid(func(gtx C) D {
								return in.Layout(gtx, material.Button(theme, restoreButton, "Restore Noita").Layout)
							}),
							layout.Rigid(func(gtx C) D {
								return in.Layout(gtx, material.Button(theme, exploreButton, "Explore Backups").Layout)
							}),
						)
					})
				},
				func(gtx C) D {
					in := layout.UniformInset(unit.Dp(8))
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, material.CheckBox(theme, autoLaunch, "Auto Launch").Layout)
						}),
						layout.Rigid(layout.Spacer{Width: 50}.Layout),
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, material.Label(theme, theme.TextSize, "Number backups to keep").Layout)
						}),
						layout.Flexed(1, material.Slider(theme, numBackups).Layout),
						layout.Rigid(func(gtx C) D {
							return layout.UniformInset(unit.Dp(8)).Layout(gtx,
								material.Body1(theme, fmt.Sprintf("%.0f", numBackups.Value*100)).Layout,
							)
						}),
					)
				},
			}

			material.List(theme, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
			})
			cl.Pop()

			e.Frame(gtx.Ops)
		}
	}
}
