package global

import (
	"fmt"

	"github.com/ShuaibKhan786/yours/pkg/utils"
	"github.com/kkdai/youtube/v2"
)

func GetArrayOfFormatWithSize(md *youtube.Video) []string {
	arrFormat := make([]string, 0, len(md.Formats)-1)
	for _, format := range md.Formats {
		if format.ItagNo == 140 { //ACC 128k audio
			arrFormat = append(
				arrFormat,
				fmt.Sprintf("%-25s%s ", "128K", utils.ConvertBytesIntoSmart(float64(format.ContentLength))),
			)
		} else {
			var qualityLabel string
			switch len(format.QualityLabel) {
			case 5:
				qualityLabel = fmt.Sprintf("%-23s", format.QualityLabel)
			case 6:
				qualityLabel = fmt.Sprintf("%-22s", format.QualityLabel)
			case 7:
				qualityLabel = fmt.Sprintf("%-21s", format.QualityLabel)
			case 11:
				qualityLabel = fmt.Sprintf("%-16s", format.QualityLabel)
			default:
				qualityLabel = fmt.Sprintf("%-25s", format.QualityLabel)
			}
			size := utils.ConvertBytesIntoSmart(float64(format.ContentLength))
			arrFormat = append(arrFormat, fmt.Sprintf("%s%s", qualityLabel, size))
		}
	}

	return arrFormat
}
