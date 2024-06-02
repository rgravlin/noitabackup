package lib

import (
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
	exploreButton = new(widget.Clickable)
	launchButton  = new(widget.Clickable)
	backupButton  = new(widget.Clickable)
	restoreButton = new(widget.Clickable)
	phase         = stopped
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

			for exploreButton.Clicked(gtx) {
				err := LaunchExplorer()
				if err != nil {
					log.Printf("error launching explorer: %v", err)
				}
			}

			for launchButton.Clicked(gtx) {
				err := LaunchNoita()
				if err != nil {
					log.Printf("failed to launch noita: %v", err)
				}
			}

			for restoreButton.Clicked(gtx) {
				if !isNoitaRunning() {
					if err := RestoreNoita(""); err != nil {
						log.Printf("failed to restore latest noita save: %v", err)
					} else {
						log.Print("successful restore request")
					}
				} else {
					log.Printf("noita.exe cannot be running to restore")
				}
			}

			for backupButton.Clicked(gtx) {
				BackupNoita()
			}

			cl := clip.Rect{Max: e.Size}.Push(gtx.Ops)
			// TODO: make this not run every frame!
			if isNoitaRunning() {
				paint.ColorOp{Color: color.NRGBA{A: 0xff, R: 0xff}}.Add(gtx.Ops)
			} else {
				paint.ColorOp{Color: color.NRGBA{A: 0xff, G: 0xff}}.Add(gtx.Ops)
			}
			paint.PaintOp{}.Add(gtx.Ops)
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				in := layout.UniformInset(unit.Dp(8))
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
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
			cl.Pop()

			e.Frame(gtx.Ops)
		}
	}
}
