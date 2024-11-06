package gui

import (
	"context"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/ShuaibKhan786/yours/cmd/yyd/channel"
	"github.com/ShuaibKhan786/yours/cmd/yyd/global"
	customlayout "github.com/ShuaibKhan786/yours/cmd/yyd/layout"
	customimage "github.com/ShuaibKhan786/yours/pkg/images"
	"github.com/kkdai/youtube/v2"
)

func (smp *sharedMainPage) videoMDtemplate(ctx context.Context, md *youtube.Video) *fyne.Container {
	//TODO: consider offload this image download to workerpool / load async
	//  size {500, 370}
	img, err := customimage.DecodeImageFromURI(ctx, md.Thumbnails[len(md.Thumbnails)-1].URL)
	if err != nil {
		log.Fatal(err)
	}
	smp.thumbnailImage.SetMinSize(fyne.NewSize(500, 370))
	smp.thumbnailImage.Image = img
	smp.thumbnailImage.FillMode = canvas.ImageFillContain
	smp.thumbnailImage.ScaleMode = canvas.ImageScaleSmooth

	smp.titleLabel.SetText(md.Title)
	smp.titleLabel.Alignment = fyne.TextAlignLeading
	smp.titleLabel.Wrapping = fyne.TextWrap(fyne.TextTruncateEllipsis)

	smp.durationLabel.SetText(md.Duration.String())

	smp.authorLabel.SetText(md.Author)

	smp.formatsSelect.ClearSelected()
	smp.formatsSelect.PlaceHolder = "select quality"
	smp.formatsSelect.Resize(fyne.NewSize(
		smp.formatsSelect.MinSize().Width+2*DefaultPaddingSize,
		smp.formatsSelect.MinSize().Height,
	))
	smp.formatsSelect.SetOptions(global.GetArrayOfFormatWithSize(md))
	smp.formatsSelect.OnChanged = func(s string) {
		if smp.formatsSelect.SelectedIndex() != -1 {
			smp.downloadButton.Enable()
			smp.downloadButton.SetIcon(theme.DownloadIcon())
			smp.downloadButton.SetText("Download")
		}
	}

	downloadChannel := &channel.DownloadChannel{}

	smp.downloadButton.Disable()

	smp.downloadButton.OnTapped = func() {
		switch smp.downloadButton.Text {
		case "Download":
			smp.linkEntry.Disable()

			smp.formatsSelect.Disable()
			smp.downloadButton.SetIcon(theme.CancelIcon())
			smp.downloadButton.SetText("Cancel Download")

			// send the video, itag, dir to backend channel for download
			downloadChannel.IsPlaylist = false
			downloadChannel.ItagNo = md.Formats[smp.formatsSelect.SelectedIndex()].ItagNo
			downloadChannel.Video = md
			downloadChannel.Dir = smp.FileSavedLocation

			smp.Channel.BackendChannel <- downloadChannel
		case "Cancel Download":
			smp.downloadButton.SetText("Cancelling...")
			smp.downloadButton.Disable()

			smp.Channel.BackendChannel <- channel.CancelDownload(md.ID)
		}
	}

	infoContainer := container.NewVBox(
		container.NewStack(
			smp.titleLabel,
		),
		container.NewHBox(
			smp.authorLabel,
			layout.NewSpacer(),
			smp.durationLabel,
		),
	)

	imageContainer := container.NewStack(
		smp.thumbnailImage,
		container.NewVBox(
			layout.NewSpacer(),
			container.New(
				customlayout.NewDynamicHBoxLayout(0, DefaultPaddingSize, DefaultPaddingSize, customlayout.CenterYAlignment),
				layout.NewSpacer(),
				smp.formatsSelect,
			),
		),
	)
	fmt.Println(imageContainer.Size())

	container := container.NewHBox(
		layout.NewSpacer(),
		container.NewVBox(
			infoContainer,
			imageContainer,
			SetYpadding(DefaultPaddingSize/2),
			smp.downloadButton,
			SetYpadding(DefaultPaddingSize),
		),
		layout.NewSpacer(),
	)

	return container
}
