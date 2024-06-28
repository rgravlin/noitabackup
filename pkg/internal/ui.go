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
	DefaultMinHeight = 155
	DefaultMaxHeight = 580
	DefaultWidth     = 640
	ErrorWidth       = DefaultWidth
	ErrorHeight      = 280
)

var (
	debugLogListFunc  = func(gtx C) D { return layout.Spacer{Width: 0}.Layout(gtx) }
	loadFunc          = layout.Rigid(layout.Spacer{Width: 56}.Layout)
	logList           = list.List
	exploreButton     = new(widget.Clickable)
	launchButton      = new(widget.Clickable)
	backupButton      = new(widget.Clickable)
	restoreButton     = new(widget.Clickable)
	debugLog          = new(widget.Bool)
	debugHeight       = 155
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
	theme             *material.Theme
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
	ui.theme = material.NewTheme()
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
			numBackups.Value = float32(viper.GetInt(ViperNumBackups)) / ConfigMaxNumBackupsToKeep

			if debugLog.Update(gtx) {
				ui.toggleDebugHeight()
				ui.adjustWindowSize(window, DefaultWidth, debugHeight)

				debugLogChecked = !debugLogChecked
				if debugLogChecked {
					ui.enableDebugLog()
				} else {
					ui.disableDebugLog()
				}

				ui.Logger.LogAndAppend(fmt.Sprintf("%s %t", InfoDebugLogSet, debugLogChecked))
			}

			if autoLaunch.Update(gtx) {
				autoLaunchChecked = !autoLaunchChecked
				ui.Logger.LogAndAppend(fmt.Sprintf("%s %t", InfoAutoLaunchSet, autoLaunchChecked))
			}

			for exploreButton.Clicked(gtx) {
				err := LaunchExplorer()
				if err != nil {
					ui.Logger.LogAndAppend(fmt.Sprintf("%s: %v", ErrLaunchingExplorer, err))
				}
			}

			for launchButton.Clicked(gtx) {
				if !ui.isOperationRunning() {
					err := LaunchNoita(true)
					if err != nil {
						ui.Logger.LogAndAppend(fmt.Sprintf("%s: %v", ErrLaunchingNoita, err))
					}
				} else {
					ui.Logger.LogAndAppend(ErrOperationAlreadyInProgress)
				}
			}

			for restoreButton.Clicked(gtx) {
				if !ui.isOperationRunning() {
					ui.runRestore()
				} else {
					ui.Logger.LogAndAppend(ErrOperationAlreadyInProgress)
				}
			}

			for backupButton.Clicked(gtx) {
				if !ui.isOperationRunning() {
					ui.runBackup()
				} else {
					ui.Logger.LogAndAppend(ErrOperationAlreadyInProgress)
				}
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
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							ui.makeButton(launchButton, BtnLaunch),
							ui.makeButton(backupButton, BtnBackup),
							ui.makeButton(restoreButton, BtnRestore),
							ui.makeButton(exploreButton, BtnExplore),
						)
					})
				},
				func(gtx C) D {
					ui.updateLoadFunc()

					in := layout.UniformInset(unit.Dp(8))
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, material.CheckBox(ui.theme, autoLaunch, ChkAutoLaunch).Layout)
						}),
						loadFunc,
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, material.Label(ui.theme, ui.theme.TextSize, SldNumBackupsToKeep).Layout)
						}),
						layout.Flexed(1, material.Slider(ui.theme, numBackups).Layout),
						layout.Rigid(func(gtx C) D {
							return layout.UniformInset(unit.Dp(8)).Layout(gtx,
								material.Body1(ui.theme, fmt.Sprintf("%.0f", numBackups.Value*ConfigMaxNumBackupsToKeep)).Layout,
							)
						}),
					)
				},
				func(gtx C) D {
					in := layout.UniformInset(unit.Dp(8))
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, material.CheckBox(ui.theme, debugLog, ChkDebugLog).Layout)
						}),
					)
				},
				debugLogListFunc,
			}

			material.List(ui.theme, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
			})

			e.Frame(gtx.Ops)
		}
	}
}

func (ui *UI) updateLoadFunc() {
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
				return material.Loader(ui.theme).Layout(gtx)
			})
		})
	} else {
		loadFunc = layout.Rigid(layout.Spacer{Width: 56}.Layout)
	}
}

func (ui *UI) makeButton(button *widget.Clickable, label string) layout.FlexChild {
	in := layout.UniformInset(unit.Dp(8))
	return layout.Rigid(func(gtx C) D {
		return in.Layout(gtx, material.Button(ui.theme, button, label).Layout)
	})
}

func (ui *UI) runRestore() {
	ui.restore = NewRestore(
		StrLatest,
		NewBackup(
			true,
			autoLaunchChecked,
			viper.GetInt(ViperNumBackups),
			viper.GetString(ViperSourcePath),
			viper.GetString(ViperDestinationPath),
		),
	)
	ui.restore.Backup.LogRing = ui.Logger
	ui.Logger.LogAndAppend(InfoStartingRestore)
	ui.restore.RestoreNoita()
}

func (ui *UI) runBackup() {
	ui.backup = NewBackup(
		true,
		autoLaunchChecked,
		viper.GetInt(ViperNumBackups),
		viper.GetString(ViperSourcePath),
		viper.GetString(ViperDestinationPath),
	)
	ui.backup.LogRing = ui.Logger
	ui.Logger.LogAndAppend(InfoStartingBackup)
	ui.backup.BackupNoita()
}

func (ui *UI) isOperationRunning() bool {
	running := false

	if ui.backup != nil {
		if ui.backup.phase == started {
			running = true
		}
	}

	if ui.restore != nil {
		if ui.restore.Backup.phase == started {
			running = true
		}
	}

	return running
}

func (ui *UI) enableDebugLog() {
	debugLogListFunc = func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return logList.Layout(gtx, ui.Logger.Len(), func(gtx layout.Context, i int) D {
					return layout.UniformInset(unit.Dp(1)).Layout(gtx, material.Label(ui.theme, unit.Sp(14), ui.Logger.Print()[i]).Layout)
				})
			}),
		)
	}
}

func (ui *UI) disableDebugLog() {
	debugLogListFunc = func(gtx C) D {
		return layout.Spacer{Width: 0}.Layout(gtx)
	}
}

func (ui *UI) toggleDebugHeight() {
	switch debugHeight {
	case DefaultMaxHeight:
		debugHeight = DefaultMinHeight
	case DefaultMinHeight:
		debugHeight = DefaultMaxHeight
	default:
	}
}

func (ui *UI) adjustWindowSize(window *app.Window, width, height int) {
	window.Option(
		app.MaxSize(unit.Dp(width), unit.Dp(height)),
		app.MinSize(unit.Dp(width), unit.Dp(height)),
	)
}
