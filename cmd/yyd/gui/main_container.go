package gui

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ShuaibKhan786/yours/cmd/yyd/channel"
	"github.com/ShuaibKhan786/yours/cmd/yyd/global"
	customlayout "github.com/ShuaibKhan786/yours/cmd/yyd/layout"
	"github.com/ShuaibKhan786/yours/pkg/yt"
	"github.com/kkdai/youtube/v2"

	xwidget "fyne.io/x/fyne/widget"
)

func (y *YYD) mainPage() *fyne.Container {

	smp := newSharedMainPage(y)

	dynamicContainer := container.NewStack(defaultPreview())
	mpContainer := container.NewBorder(
		smp.navbar(),
		nil, nil, nil,
		dynamicContainer,
	)

	//****ONE GO ROUTINE RUNNING HERE****
	go func() {
		for {
			select {
			case i := <-y.Channel.GUIChannel:
				switch ct := i.(type) {
				case *youtube.Video:
					smp.downloadButton.SetText("Download")
					smp.downloadButton.SetIcon(theme.DownloadIcon())
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					md := smp.videoMDtemplate(ctx, ct)

					dynamicContainer.RemoveAll()
					dynamicContainer.Add(md)
					dynamicContainer.Refresh()
					smp.linkEntry.SetText("")
					smp.linkEntry.Enable()
				case *yt.PlaylistMetadata:
					//just to track the downloaded and not downloaded in playlist
					smp.playlistDownloadNeeded = 0
					smp.playlistDownloadDone = 0
					smp.playlistDownloadError = 0
					y.Channel.BackendChannel <- channel.CancelPlaylistImageDownload{}
					y.PlaylistMap.Mu.Lock()
					y.PlaylistMap.Details = make(map[string]*global.MediaDetails)
					y.PlaylistMap.Mu.Unlock()
					dpi := &channel.DownloadPlaylistImages{
						PlaylistMap: y.PlaylistMap,
						MD:          ct,
					}
					y.Channel.BackendChannel <- dpi

					ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
					defer cancel()
					md := smp.playlistMDtemplate(ctx, ct)

					dynamicContainer.RemoveAll()
					dynamicContainer.Add(md)
					dynamicContainer.Refresh()

					smp.linkEntry.SetText("")
					smp.linkEntry.Enable()
				case channel.DownloadDone: //used both for playlistVideo and video
					smp.linkEntry.Enable()
					if ct == "" { //it means video
						smp.formatsSelect.Enable()
						smp.downloadButton.Enable()
						smp.downloadButton.SetIcon(theme.ConfirmIcon())
						smp.downloadButton.SetText("Successfully Downloaded")
					} else {
						if state, ok := smp.playlistBindingMap[string(ct)]; ok {
							state.ButtonText.Set("Downloaded")
						}
					}
				case channel.DownloadCancelDone: //used both for playlistVideo and video
					smp.linkEntry.Enable()
					if ct == "" { //it means video
						smp.formatsSelect.Enable()
						smp.downloadButton.Enable()
						smp.downloadButton.SetIcon(theme.ConfirmIcon())
						smp.downloadButton.SetText("Successfully Canceled")
					} else {
						if state, ok := smp.playlistBindingMap[string(ct)]; ok {
							state.ButtonText.Set("Cancelled")
						}
					}
				case channel.ErrorChannel:
					dynamicContainer.RemoveAll()
					dynamicContainer.Add(displayError(string(ct)))
					dynamicContainer.Refresh()
					smp.linkEntry.SetText("")
					smp.linkEntry.Enable()
				case channel.FetchMDStarted:
					dynamicContainer.RemoveAll()
					dynamicContainer.Add(displayFetching())
					dynamicContainer.Refresh()
				case channel.PlaylistDownloadAllDone:
					smp.playlistDownloadDone++
					smp.playlistDownloadProgressDoneBind.Set(smp.playlistDownloadDone)
					if smp.playlistDownloadNeeded == smp.playlistDownloadDone+smp.playlistDownloadError {
						smp.linkEntry.Enable()
						smp.formatsSelect.Enable()
						smp.playlistDownloadAllButton.Enable()
						smp.playlistDownloadAllButton.SetIcon(theme.ConfirmIcon())
						smp.playlistDownloadAllButton.SetText("All Downloaded")
					}
				case channel.PlaylistDownloadAllCancelDone:
					smp.linkEntry.Enable()
					smp.formatsSelect.Enable()
					smp.playlistDownloadAllButton.Enable()
					smp.playlistDownloadAllButton.SetIcon(theme.ConfirmIcon())
					smp.playlistDownloadAllButton.SetText("All Canceled")
				case channel.PlaylistDownloadAllError:
					smp.playlistDownloadError++
					if smp.playlistDownloadNeeded == smp.playlistDownloadDone+smp.playlistDownloadError {
						smp.linkEntry.Enable()
						smp.formatsSelect.Enable()
						smp.playlistDownloadAllButton.Enable()
						smp.playlistDownloadAllButton.SetIcon(theme.ConfirmIcon())
						smp.playlistDownloadAllButton.SetText("All Downloaded")
					}
				}
			case <-y.RootCtx.Done():
				return
			}
		}
	}()

	return mpContainer
}

func defaultPreview() *fyne.Container {
	previewLabel := widget.NewLabel("Youtube preview will be shown here")

	previewImage := canvas.NewImageFromResource(resourcePreviewSvg)
	previewImage.FillMode = canvas.ImageFillContain
	previewImage.SetMinSize(fyne.NewSize(24, 24))

	return container.New(
		customlayout.NewDynamicHBoxLayout(
			theme.Size(theme.SizeNameInnerPadding),
			DefaultPaddingSize,
			DefaultPaddingSize,
			customlayout.CenterYAlignment,
		),
		layout.NewSpacer(),
		previewLabel,
		previewImage,
		layout.NewSpacer(),
	)
}

func displayError(err string) *fyne.Container {
	errorImage := canvas.NewImageFromResource(resourceErrorSvg)
	errorImage.FillMode = canvas.ImageFillContain
	errorImage.SetMinSize(fyne.NewSize(24, 24))

	return container.New(
		customlayout.NewDynamicHBoxLayout(
			theme.Size(theme.SizeNameInnerPadding),
			DefaultPaddingSize,
			DefaultPaddingSize,
			customlayout.CenterYAlignment,
		),
		layout.NewSpacer(),
		widget.NewLabel(err),
		errorImage,
		layout.NewSpacer(),
	)
}

func displayFetching() *fyne.Container {
	var c *fyne.Container

	gif, err := xwidget.NewAnimatedGifFromResource(resourceSpinningGif)
	if err != nil {
		timeImage := canvas.NewImageFromResource(resourceTimeSvg)
		timeImage.FillMode = canvas.ImageFillOriginal
		timeImage.Resize(fyne.NewSize(24, 24))

		c = container.New(
			customlayout.NewDynamicHBoxLayout(
				theme.Size(theme.SizeNameInnerPadding),
				DefaultPaddingSize,
				DefaultPaddingSize,
				customlayout.CenterYAlignment,
			),
			layout.NewSpacer(),
			widget.NewLabel("Fetching preview details"),
			timeImage,
			layout.NewSpacer(),
		)
	} else {
		gif.Start()
		gif.SetMinSize(fyne.NewSize(24, 24))
		c = container.New(
			customlayout.NewDynamicHBoxLayout(
				theme.Size(theme.SizeNameInnerPadding),
				DefaultPaddingSize,
				DefaultPaddingSize,
				customlayout.CenterYAlignment,
			),
			layout.NewSpacer(),
			widget.NewLabel("Fetching preview detials"),
			gif,
			layout.NewSpacer(),
		)
	}

	return c
}

type sharedMainPage struct {
	*YYD
	defaultImage *canvas.Image
	//video metadata template objects
	thumbnailImage     *canvas.Image
	fileLocationButton *widget.Button
	linkEntry          *widget.Entry
	previewButton      *widget.Button
	titleLabel         *widget.Label
	authorLabel        *widget.Label
	durationLabel      *widget.Label
	formatsSelect      *widget.Select
	downloadButton     *widget.Button
	//playlist metadata template objects
	playlistList                       *widget.List
	playlistBinding                    binding.UntypedList
	playlistDownloadAllButton          *widget.Button
	playlistDownloadProgressLabel      *widget.Label
	playlistDownloadProgressDoneBind   binding.Int
	playlistDownloadProgressNeededBind binding.Int
	playlistDownloadNeeded             int
	playlistDownloadDone               int
	playlistDownloadError              int
	playlistBindingMap                 map[string]*PlaylistBindingData
}

func newSharedMainPage(y *YYD) *sharedMainPage {

	smp := &sharedMainPage{
		YYD:                       y,
		defaultImage:              canvas.NewImageFromResource(resourceDefaultimageSvg),
		thumbnailImage:            canvas.NewImageFromImage(nil),
		fileLocationButton:        widget.NewButtonWithIcon("Edit", theme.FolderIcon(), nil),
		linkEntry:                 widget.NewEntry(),
		previewButton:             widget.NewButtonWithIcon("Preview", theme.VisibilityIcon(), nil),
		titleLabel:                widget.NewLabel(""),
		authorLabel:               widget.NewLabel(""),
		durationLabel:             widget.NewLabel(""),
		formatsSelect:             widget.NewSelect(nil, nil),
		downloadButton:            widget.NewButtonWithIcon("Download", theme.DownloadIcon(), nil),
		playlistDownloadAllButton: widget.NewButtonWithIcon("Download All", theme.DownloadIcon(), nil),
	}

	smp.playlistDownloadProgressDoneBind = binding.BindInt(&smp.playlistDownloadDone)
	smp.playlistDownloadProgressNeededBind = binding.BindInt(&smp.playlistDownloadNeeded)
	smp.playlistDownloadProgressLabel = widget.NewLabelWithData(
		binding.NewSprintf(
			"%d / %d",
			smp.playlistDownloadProgressDoneBind,
			smp.playlistDownloadProgressNeededBind,
		),
	)

	return smp
}

func (smp *sharedMainPage) navbar() *fyne.Container {
	bannerImage := canvas.NewImageFromResource(resourceYydLogoSvg)
	bannerImage.FillMode = canvas.ImageFillOriginal
	bannerImage.Resize(fyne.NewSize(36, 36))

	smp.fileLocationButton.Resize(fyne.NewSize(
		smp.fileLocationButton.MinSize().Width+18,
		smp.fileLocationButton.MinSize().Height,
	))

	smp.fileLocationButton.OnTapped = func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
			if err != nil {
				return
			}
			if lu != nil && lu.Path() != "" {
				smp.setFileSavedLocation(lu.Path())
			}
		}, smp.Window)
	}

	smp.linkEntry.Resize(fyne.NewSize(300, smp.linkEntry.MinSize().Height))
	smp.linkEntry.SetPlaceHolder("  Paste your youtube link here")

	smp.previewButton.Resize(fyne.NewSize(smp.previewButton.MinSize().Width+18, smp.previewButton.MinSize().Height))
	smp.previewButton.Disable()

	smp.linkEntry.OnChanged = func(s string) {
		if len(s) > 0 {
			smp.previewButton.Enable()
		} else {
			smp.previewButton.Disable()
		}
	}

	smp.previewButton.OnTapped = func() {
		// send link to backend channel
		link := channel.LinkChannel(smp.linkEntry.Text)
		if link != "" {
			smp.linkEntry.Disable()
			smp.previewButton.Disable()
			if len(smp.PlaylistMap.Details) > 0 {
				smp.playlistList.ScrollToTop()
				smp.playlistList.UnselectAll()
			}
			smp.Channel.BackendChannel <- link
		}
	}

	return container.NewBorder(
		nil,
		SetYpadding(DefaultPaddingSize),
		nil, nil,
		container.New(
			customlayout.NewDynamicHBoxLayout(theme.InnerPadding(), 0, 0, customlayout.CenterYAlignment),
			layout.NewSpacer(),
			SetYpadding(2*DefaultPaddingSize),
			bannerImage,
			smp.linkEntry,
			smp.previewButton,
			smp.fileLocationButton,
			layout.NewSpacer(),
		))
}
