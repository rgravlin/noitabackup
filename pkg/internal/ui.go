package internal

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/spf13/viper"
	"image/color"
)

const (
	stopped int = iota
	started
)

var (
	logList           = list.List
	exploreButton     = new(widget.Clickable)
	launchButton      = new(widget.Clickable)
	backupButton      = new(widget.Clickable)
	restoreButton     = new(widget.Clickable)
	debugLog          = new(widget.Bool)
	autoLaunch        = new(widget.Bool)
	numBackups        = new(widget.Float)
	autoLaunchChecked = false
	debugLogChecked   = false
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

type UI struct {
	backup            *Backup
	restore           *Restore
	Logger            *LogRing
	autoLaunchChecked bool
}

func NewUI(autoLaunch bool) *UI {
	return &UI{
		Logger:            NewLogRing(16),
		autoLaunchChecked: autoLaunch,
	}
}

// Run handles all the events and rendering for the application window.
// It takes a *app.Window as a parameter and returns an error if any.
func (ui *UI) Run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops
	autoLaunch.Value, autoLaunchChecked = ui.autoLaunchChecked, ui.autoLaunchChecked

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

			if debugLog.Update(gtx) {
				debugLogChecked = !debugLogChecked
				ui.Logger.LogAndAppend(fmt.Sprintf("debug log set to %t", debugLogChecked))
			}

			if autoLaunch.Update(gtx) {
				autoLaunchChecked = !autoLaunchChecked
				ui.Logger.LogAndAppend(fmt.Sprintf("autolaunch set to %t", autoLaunchChecked))
			}

			for exploreButton.Clicked(gtx) {
				err := LaunchExplorer()
				if err != nil {
					ui.Logger.LogAndAppend(fmt.Sprintf("error launching explorer: %v", err))
				}
			}

			for launchButton.Clicked(gtx) {
				err := LaunchNoita(true)
				if err != nil {
					ui.Logger.LogAndAppend(fmt.Sprintf("error launching noita: %v", err))
				}
			}

			for restoreButton.Clicked(gtx) {
				ui.restore = NewRestore(
					"latest",
					NewBackup(
						true,
						autoLaunchChecked,
						viper.GetInt("num-backups"),
						viper.GetString("source-path"),
						viper.GetString("destination-path"),
					),
				)
				ui.restore.Backup.LogRing = ui.Logger
				ui.Logger.LogAndAppend("starting restore")
				ui.restore.RestoreNoita()
			}

			for backupButton.Clicked(gtx) {
				ui.backup = NewBackup(
					true,
					autoLaunchChecked,
					viper.GetInt("num-backups"),
					viper.GetString("source-path"),
					viper.GetString("destination-path"),
				)
				ui.backup.LogRing = ui.Logger
				ui.Logger.LogAndAppend("starting backup")
				ui.backup.BackupNoita()
			}

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
					if ui.isOperationRunning() {
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
					} else {
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
				func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return logList.Layout(gtx, ui.Logger.Len(), func(gtx layout.Context, i int) D {
								return layout.UniformInset(unit.Dp(1)).Layout(gtx, material.Label(theme, unit.Sp(14), ui.Logger.Print()[i]).Layout)
							})
						}),
					)
				},
			}

			material.List(theme, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
			})

			e.Frame(gtx.Ops)
		}
	}
}

func (ui *UI) isOperationRunning() bool {
	if ui.backup != nil {
		if ui.backup.phase == started {
			return true
		}
	}

	if ui.restore != nil {
		if ui.restore.Backup.phase == started {
			return true
		}
	}

	return false
}
