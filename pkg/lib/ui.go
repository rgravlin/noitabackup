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
	"github.com/spf13/viper"
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

// Run handles all the events and rendering for the application window.
// It takes a *app.Window as a parameter and returns an error if any.
func Run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			// TODO: Forcing window size every frame results in memory leak,
			//       but sometimes it resizes for no reason
			// window.Option(
			// 	app.MaxSize(unit.Dp(640), unit.Dp(105)),
			// 	app.MinSize(unit.Dp(640), unit.Dp(105)),
			// )
			numBackups.Value = float32(viper.GetInt("num-backups")) / ConfigMaxNumBackupsToKeep

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
				restore := NewRestore(
					"latest",
					NewBackup(
						true,
						viper.GetBool("auto-launch"),
						viper.GetInt("num-backups"),
						viper.GetString("source-path"),
						viper.GetString("destination-path"),
					),
				)
				restore.RestoreNoita()
			}

			for backupButton.Clicked(gtx) {
				backup := NewBackup(
					true,
					viper.GetBool("auto-launch"),
					viper.GetInt("num-backups"),
					viper.GetString("source-path"),
					viper.GetString("destination-path"),
				)
				backup.BackupNoita()
				// BackupNoita(true, viper.GetInt("num-backups"))
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
					var loadFunc layout.FlexChild
					switch phase {
					case stopped:
						loadFunc = layout.Rigid(layout.Spacer{Width: 56}.Layout)
					case started:
						loadFunc = layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Top:    unit.Dp(4),
								Bottom: unit.Dp(4),
								Left:   unit.Dp(16),
								Right:  unit.Dp(16),
							}.Layout(gtx, func(gtx C) D {
								gtx.Constraints.Max.X = gtx.Dp(unit.Dp(24))
								gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(24))
								return material.Loader(theme).Layout(gtx)
							})
						})
					default:
						loadFunc = layout.Rigid(layout.Spacer{Width: 56}.Layout)
					}
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, material.CheckBox(theme, autoLaunch, "Auto Launch").Layout)
						}),
						loadFunc,
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, material.Label(theme, theme.TextSize, "Number backups to keep").Layout)
						}),
						layout.Flexed(1, material.Slider(theme, numBackups).Layout),
						layout.Rigid(func(gtx C) D {
							return layout.UniformInset(unit.Dp(8)).Layout(gtx,
								material.Body1(theme, fmt.Sprintf("%.0f", numBackups.Value*ConfigMaxNumBackupsToKeep)).Layout,
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
