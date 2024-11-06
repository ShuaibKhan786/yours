package gui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ShuaibKhan786/yours/cmd/yyd/channel"
	"github.com/ShuaibKhan786/yours/cmd/yyd/global"
	customlayout "github.com/ShuaibKhan786/yours/cmd/yyd/layout"
	customimage "github.com/ShuaibKhan786/yours/pkg/images"
	"github.com/ShuaibKhan786/yours/pkg/utils"
	"github.com/ShuaibKhan786/yours/pkg/yt"
	"github.com/kkdai/youtube/v2"
)

type PlaylistBindingData struct {
	SelectIndex binding.Int
	ButtonText  binding.String
	Video       *youtube.Video
}

func (smp *sharedMainPage) playlistMDtemplate(ctx context.Context, md *yt.PlaylistMetadata) *fyne.Container {
	smp.playlistBindingMap = make(map[string]*PlaylistBindingData)

	smp.titleLabel.SetText(md.Title)
	smp.titleLabel.Alignment = fyne.TextAlignLeading
	smp.titleLabel.Truncation = fyne.TextTruncateEllipsis

	videosNumberLabel := widget.NewLabel(
		fmt.Sprintf(" %d videos", len(md.Videos)),
	)

	smp.authorLabel.SetText(md.Author)

	audioFullSize := calculateFullAudioSizes(md)
	smp.formatsSelect.ClearSelected()
	smp.formatsSelect.Resize(fyne.NewSize(
		smp.formatsSelect.MinSize().Width+2*DefaultPaddingSize,
		smp.formatsSelect.MinSize().Height,
	))
	smp.formatsSelect.SetOptions([]string{
		fmt.Sprintf("128K  %s", utils.ConvertBytesIntoSmart(float64(audioFullSize))),
	})
	smp.formatsSelect.OnChanged = func(s string) {
		if smp.formatsSelect.SelectedIndex() != -1 {
			smp.playlistDownloadDone = 0
			smp.playlistDownloadProgressDoneBind.Set(smp.playlistDownloadDone)
			smp.playlistDownloadProgressLabel.Hidden = true
			smp.playlistDownloadAllButton.SetText("Download All")
			smp.playlistDownloadAllButton.SetIcon(theme.DownloadIcon())
			smp.playlistDownloadAllButton.Enable()
		}
	}

	smp.playlistDownloadProgressLabel.Hidden = true

	smp.playlistDownloadAllButton.Disable()
	smp.playlistDownloadAllButton.Resize(
		fyne.NewSize(
			smp.playlistDownloadAllButton.MinSize().Width+18,
			smp.playlistDownloadAllButton.MinSize().Height,
		),
	)
	smp.playlistDownloadAllButton.OnTapped = func() {
		switch smp.playlistDownloadAllButton.Text {
		case "Download All":
			smp.linkEntry.Disable()

			smp.formatsSelect.Disable()
			smp.playlistDownloadAllButton.SetIcon(theme.CancelIcon())
			smp.playlistDownloadAllButton.SetText("Cancel All")

			smp.playlistDownloadNeeded = len(md.Videos)
			smp.playlistDownloadProgressLabel.Hidden = false
			smp.playlistDownloadProgressNeededBind.Set(smp.playlistDownloadNeeded)
			smp.Channel.BackendChannel <- &channel.DownloadAllPlaylist{
				PlaylistMD: md,
				ItagNo:     140,
				Dir:        smp.FileSavedLocation,
			}
		case "Cancel All":
			smp.playlistDownloadAllButton.SetText("Cancelling...")
			smp.playlistDownloadAllButton.Disable()

			smp.Channel.BackendChannel <- channel.CancelDownload(md.ID)
		}
	}

	playlistInfoContainer := container.NewVBox(
		container.NewBorder(
			nil, nil,
			SetXpadding(2),
			videosNumberLabel,
			container.NewStack(smp.titleLabel),
		),
		container.New(
			customlayout.NewDynamicHBoxLayout(theme.InnerPadding(), DefaultPaddingSize, DefaultPaddingSize, customlayout.CenterYAlignment),
			smp.authorLabel,
			layout.NewSpacer(),
			smp.playlistDownloadProgressLabel,
			smp.formatsSelect,
			smp.playlistDownloadAllButton,
		),
	)

	//**** Playlist List ****
	smp.playlistBinding = binding.NewUntypedList()
	smp.playlistBindingMap = make(map[string]*PlaylistBindingData)
	for _, video := range md.Videos {
		item := &PlaylistBindingData{
			SelectIndex: binding.NewInt(),
			ButtonText:  binding.NewString(),
			Video:       video,
		}
		item.ButtonText.Set("Download")
		item.SelectIndex.Set(-1)
		smp.playlistBinding.Append(item)
	}

	smp.playlistList = widget.NewListWithData(
		smp.playlistBinding,
		func() fyne.CanvasObject {
			thumbnailImage := canvas.NewImageFromImage(nil)
			thumbnailImage.FillMode = canvas.ImageFillContain
			thumbnailImage.ScaleMode = canvas.ImageScaleSmooth
			thumbnailImage.SetMinSize(fyne.NewSize(120, 90))

			titleLabel := widget.NewLabel("template title")
			titleLabel.Wrapping = fyne.TextWrap(fyne.TextTruncateEllipsis)
			titleLabel.Alignment = fyne.TextAlignLeading
			durationLabel := widget.NewLabel("template duration")

			authorLabel := widget.NewLabel("template author")
			formatSelect := widget.NewSelect(nil, nil)
			formatSelect.PlaceHolder = "select quality"
			downloadButton := widget.NewButtonWithIcon("Download", theme.DownloadIcon(), nil)
			downloadButton.Resize(fyne.NewSize(
				smp.playlistDownloadAllButton.MinSize().Width,
				smp.playlistDownloadAllButton.MinSize().Height,
			))
			downloadButton.Disable()

			return container.NewBorder(
				//paddedContainer
				SetYpadding(DefaultPaddingSize), SetYpadding(DefaultPaddingSize),
				SetXpadding(DefaultPaddingSize), SetXpadding(DefaultPaddingSize),
				container.NewBorder(
					nil, nil,
					thumbnailImage,
					nil,
					//firstContainer
					container.NewVBox(
						//secondContainer
						container.NewBorder(
							//thirdContainer
							nil, nil, nil,
							durationLabel,
							titleLabel,
						),
						layout.NewSpacer(),
						container.NewHBox(
							//fourthContainer
							authorLabel,
							layout.NewSpacer(),
							formatSelect,
							downloadButton,
						),
					),
				),
			)
		},
		func(di binding.DataItem, co fyne.CanvasObject) {

			paddedContainer := co.(*fyne.Container)
			diItem, _ := di.(binding.Untyped).Get()
			diActualItem := diItem.(*PlaylistBindingData)

			firstContainer := paddedContainer.Objects[0].(*fyne.Container)
			thumbnailImage := firstContainer.Objects[1].(*canvas.Image)

			secondContainer := firstContainer.Objects[0].(*fyne.Container)
			thirdContainer := secondContainer.Objects[0].(*fyne.Container)
			fourthContainer := secondContainer.Objects[2].(*fyne.Container)

			titleLabel := thirdContainer.Objects[0].(*widget.Label)
			durationLabel := thirdContainer.Objects[1].(*widget.Label)

			authorLabel := fourthContainer.Objects[0].(*widget.Label)
			formatSelect := fourthContainer.Objects[2].(*widget.Select)
			downloadButton := fourthContainer.Objects[3].(*widget.Button)

			var ok bool
			var err error

			d, ok := smp.PlaylistMap.Details[diActualItem.Video.ID]
			if ok {
				thumbnailImage.Image = d.Thumbnail
				formatSelect.SetOptions(d.Formats)
			} else {
				thumbnailImage.Image, err = customimage.DecodeImageFromURI(
					ctx,
					diActualItem.Video.Thumbnails[0].URL,
				)
				if err != nil {
					thumbnailImage.Resource = resourceDefaultimageSvg
				}

				formatSelect.SetOptions(global.GetArrayOfFormatWithSize(diActualItem.Video))
			}

			titleLabel.SetText(diActualItem.Video.Title)
			durationLabel.SetText(diActualItem.Video.Duration.String())
			authorLabel.SetText(diActualItem.Video.Author)

			index, _ := diActualItem.SelectIndex.Get()
			formatSelect.ClearSelected()
			formatSelect.SetSelectedIndex(index)

			text, _ := diActualItem.ButtonText.Get()
			downloadButton.SetText(text)
			switch text {
			case "Download":
				downloadButton.SetIcon(theme.DownloadIcon())
			case "Cancel":
				formatSelect.Disable()
				downloadButton.SetIcon(theme.CancelIcon())
			case "Downloaded":
				downloadButton.SetIcon(theme.ConfirmIcon())
				downloadButton.Enable()
				formatSelect.Enable()
			case "Cancelled":
				downloadButton.SetIcon(theme.ConfirmIcon())
				downloadButton.Enable()
				formatSelect.Enable()
			}

			diActualItem.ButtonText.AddListener(binding.NewDataListener(func() {
				text, _ := diActualItem.ButtonText.Get()
				switch text {
				case "Download":
					downloadButton.SetIcon(theme.DownloadIcon())
				case "Cancel":
					formatSelect.Disable()
					downloadButton.SetIcon(theme.CancelIcon())
				case "Downloaded":
					downloadButton.SetIcon(theme.ConfirmIcon())
					downloadButton.Enable()
					formatSelect.Enable()
				case "Cancelled":
					downloadButton.SetIcon(theme.ConfirmIcon())
					downloadButton.Enable()
					formatSelect.Enable()
				}
				downloadButton.SetText(text)
			}))

			diActualItem.SelectIndex.AddListener(binding.NewDataListener(func() {
				index, _ := diActualItem.SelectIndex.Get()
				formatSelect.SetSelectedIndex(index)
			}))

			//dynamic changes item will occur here
			formatSelect.OnChanged = func(s string) {
				if formatSelect.SelectedIndex() != 1 {
					diActualItem.SelectIndex.Set(formatSelect.SelectedIndex())
					diActualItem.ButtonText.Set("Download")
					downloadButton.Enable()
				}
			}

			downloadButton.OnTapped = func() {
				switch downloadButton.Text {
				case "Download":
					if formatSelect.SelectedIndex() != -1 {
						smp.linkEntry.Disable()

						diActualItem.ButtonText.Set("Cancel")
						diActualItem.SelectIndex.Set(formatSelect.SelectedIndex())
						smp.playlistBindingMap[diActualItem.Video.ID] = diActualItem

						// send something to backend for download
						smp.Channel.BackendChannel <- &channel.DownloadChannel{
							IsPlaylist: true,
							ItagNo:     diActualItem.Video.Formats[formatSelect.SelectedIndex()].ItagNo,
							Video:      diActualItem.Video,
							Dir:        smp.FileSavedLocation,
						}
					} else {
						formatSelect.ClearSelected()
						downloadButton.Disable()
					}
				case "Cancel":
					//send to appropriate channel for cancellation
					smp.Channel.BackendChannel <- channel.CancelDownload(diActualItem.Video.ID)
				}
			}
		},
	)

	return container.NewBorder(
		playlistInfoContainer,
		nil, nil, nil,
		smp.playlistList,
	)
}

func calculateFullAudioSizes(md *yt.PlaylistMetadata) int64 {
	// TODO: figure out one good thing for video also
	var audioSize int64

	for _, format := range md.Videos {
		audioSize += format.Formats.Itag(140)[0].ContentLength
	}

	return audioSize
}
