package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ShuaibKhan786/yours/cmd/yyd/navigation"
	"github.com/ShuaibKhan786/yours/pkg/yt"
)

func (y *YYD) getStartedPage(sv *navigation.StackNavigation) *fyne.Container {
	yydLabel := widget.NewLabel("Yours Youtube Downloader")
	yydLabel.TextStyle.Bold = true
	// yydLabel.TextStyle.Italic = true

	yydLogoImage := canvas.NewImageFromResource(resourceYydLogoSvg)
	yydLogoImage.SetMinSize(fyne.NewSize(24, 24))
	yydLogoImage.FillMode = canvas.ImageFillContain

	infoLabel := widget.NewLabel("An app that lets you download non-copyrighted YouTube content")
	infoLabel.TextStyle.Italic = true

	nextButton := widget.NewButton("Get Started", func() {
		y.setIsSetup()
		sv.Push()
	})
	nextButton.Disable()

	var ffmpegInstalledState bool
	var fileSavedLocationState bool

	ffmpegLabel := widget.NewLabel("")
	ffmpegCheckButton := widget.NewButtonWithIcon("Check", theme.ConfirmIcon(), nil)
	ffmpegConfirmationImage := canvas.NewImageFromResource(nil)
	ffmpegConfirmationImage.SetMinSize(fyne.NewSize(16, 16))
	ffmpegConfirmationImage.FillMode = canvas.ImageFillContain
	if yt.IsFFmpegInstalled() {
		ffmpegInstalledState = true
		ffmpegLabel.SetText("FFmpeg Installed")
		ffmpegConfirmationImage.Resource = resourceConfirmSvg
		ffmpegCheckButton.Hidden = true
	} else {
		ffmpegLabel.SetText("FFmpeg Not Installed")
		ffmpegConfirmationImage.Resource = resourceWrongSvg
		ffmpegCheckButton.Hidden = false
	}
	ffmpegCheckButton.OnTapped = func() {
		ffmpegCheckButton.SetText("Checking...")
		ffmpegCheckButton.Disable()

		if yt.IsFFmpegInstalled() {
			ffmpegInstalledState = true
			ffmpegLabel.SetText("FFmpeg Installed")
			ffmpegConfirmationImage.Resource = resourceConfirmSvg
			ffmpegCheckButton.Hidden = true
		} else {
			ffmpegLabel.SetText("FFmpeg Not Installed")
			ffmpegConfirmationImage.Resource = resourceWrongSvg
			ffmpegCheckButton.SetIcon(theme.ConfirmIcon())
			ffmpegCheckButton.SetText("Check")
			ffmpegCheckButton.Hidden = false
		}

		if ffmpegInstalledState && fileSavedLocationState {
			nextButton.Enable()
		}
	}

	fileLocationLabel := widget.NewLabel("")
	fileLocationButton := widget.NewButtonWithIcon("Open", theme.FolderIcon(), nil)
	fileLocationConfirmationImage := canvas.NewImageFromResource(nil)
	fileLocationConfirmationImage.SetMinSize(fyne.NewSize(16, 16))
	fileLocationConfirmationImage.FillMode = canvas.ImageFillContain
	if y.FileSavedLocation != "" {
		fileLocationLabel.SetText("File Location Saved")
		fileLocationConfirmationImage.Resource = resourceConfirmSvg
		fileLocationButton.Hidden = true
	} else {
		fileLocationLabel.SetText("File Location not yet Saved")
		fileLocationConfirmationImage.Resource = resourceWrongSvg
		fileLocationButton.Hidden = false
	}
	fileLocationButton.OnTapped = func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
			if err != nil {
				return
			}
			if lu == nil || lu.Path() == "" {
				return
			}
			y.setFileSavedLocation(lu.Path())
			fileSavedLocationState = true

			if y.FileSavedLocation != "" {
				fileLocationLabel.SetText("File Location Saved")
				fileLocationConfirmationImage.Resource = resourceConfirmSvg
				fileLocationButton.Hidden = true
			} else {
				fileLocationLabel.SetText("File Location not yet Saved")
				fileLocationConfirmationImage.Resource = resourceWrongSvg
				fileLocationButton.Hidden = false
			}

			if ffmpegInstalledState && fileSavedLocationState {
				nextButton.Enable()
			}
		}, y.Window)
	}

	if fileSavedLocationState && ffmpegInstalledState {
		nextButton.Enable()
	}

	authorLabel := widget.NewLabel("Developed by Md Shuaib Khan")
	authorLabel.TextStyle.Italic = true

	addLabel := widget.NewLabel("and")
	addLabel.TextStyle.Italic = true

	usingLabel := widget.NewLabel("using")
	usingLabel.TextStyle.Italic = true

	fyneImage := canvas.NewImageFromResource(resourceFyneLogoTransparentPng)
	fyneImage.SetMinSize(fyne.NewSize(32, 32))
	fyneImage.FillMode = canvas.ImageFillContain

	goLogoImage := canvas.NewImageFromResource(resourceGoLogoBlueSvg)
	goLogoImage.SetMinSize(fyne.NewSize(32, 32))
	goLogoImage.FillMode = canvas.ImageFillContain

	return container.NewCenter(
		container.NewVBox(
			container.NewHBox(
				yydLogoImage,
				yydLabel,
			),
			infoLabel,
			SetYpadding(4*DefaultPaddingSize),
			container.NewHBox(
				ffmpegLabel, ffmpegConfirmationImage, layout.NewSpacer(), ffmpegCheckButton,
			),
			container.NewHBox(
				fileLocationLabel, fileLocationConfirmationImage, layout.NewSpacer(), fileLocationButton,
			),
			SetYpadding(4*DefaultPaddingSize),
			container.NewHBox(
				layout.NewSpacer(),
				authorLabel,
			),
			container.NewHBox(
				layout.NewSpacer(),
				usingLabel, goLogoImage, addLabel, fyneImage,
			),
			nextButton,
		),
	)
}
