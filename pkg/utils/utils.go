package utils

import (
	"fmt"
)

const (
	_          = iota
	KB float64 = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func ConvertBytesIntoSmart(b float64) string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.0fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.0fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.0fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.0fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.0fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.0fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.0fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.0fKB", b/KB)
	}
	return fmt.Sprintf("%.0fB", b)
}
