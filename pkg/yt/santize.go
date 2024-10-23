package yt

import (
	"regexp"
	"strings"

	"github.com/kkdai/youtube/v2"
)

type supportedMimetype struct {
	index int8
	state bool
}

var supportedMimetypes map[string]*supportedMimetype

func init() {
	supportedMimetypes = map[string]*supportedMimetype{
		"hd2160": {index: -1},
		"hd1440": {index: -1},
		"hd1080": {index: -1},
		"hd720":  {index: -1},
		"large":  {index: -1},
		"medium": {index: -1},
		"small":  {index: -1},
		"tiny":   {index: -1},
	}
}
func sanitizedMP4FormatsOnly(formats youtube.FormatList) youtube.FormatList {
	sanitizedFormats := make(youtube.FormatList, 0, 9)
	index := int8(0)

	for _, format := range formats {
		mimetypeArr := strings.Split(format.MimeType, ";")

		switch mimetypeArr[0] {
		case "video/mp4":
			codec := strings.Split(mimetypeArr[1], "=")
			actualCodec := codec[1]

			if smt, ok := supportedMimetypes[format.Quality]; ok {
				if smt.state {
					if smt.index != -1 {
						if strings.Contains(actualCodec, "avc1") { //h.264 codec
							sanitizedFormats[smt.index] = format
							smt.index = -1
						}
					}
				} else {
					if !strings.Contains(actualCodec, "avc1") {
						smt.index = index
					}
					sanitizedFormats = append(sanitizedFormats, format)
					smt.state = true
					index++
				}
			}
		case "audio/mp4":
			if format.ItagNo == 140 { //AAC 128k
				sanitizedFormats = append(sanitizedFormats, format)
				index++
			}
		}
	}
	return sanitizedFormats
}
func SanitizeFilename(fileName string) string {
	invalidChars := regexp.MustCompile(`[<>:"/\\|?* ]`)
	safeFileName := invalidChars.ReplaceAllString(fileName, "")
	return strings.TrimSpace(safeFileName)
}
