package constants

import "github.com/xuri/excelize/v2"

// ExcelStyle ExcelStyle
type ExcelStyle struct {
	Title          int
	Subtitle       int
	TitleCenter    int
	SubtitleCenter int
	Bold           int
	Heading        int
	BoldCenter     int
	Normal         int
	Center         int
	TextRight      int
	TextRightBold  int
}

// NewExcelStyle NewExcelStyle
func NewExcelStyle(file *excelize.File) *ExcelStyle {
	normal, _ := file.NewStyle(&excelize.Style{})
	title, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Italic: false, Size: 16}})
	subtitle, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Italic: false, Size: 14}})
	bold, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Italic: false}})
	heading, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Italic: false, Size: 14}})
	boldCenter, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Italic: false}, Alignment: &excelize.Alignment{Horizontal: "center"}})
	center, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: false, Italic: false}, Alignment: &excelize.Alignment{Horizontal: "center"}})
	textRight, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: false, Italic: false}, Alignment: &excelize.Alignment{Horizontal: "right"}})
	textRightBold, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Italic: false}, Alignment: &excelize.Alignment{Horizontal: "right"}})
	titleCenter, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Italic: false, Size: 16}, Alignment: &excelize.Alignment{Horizontal: "center"}})
	subtitleCenter, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Italic: false, Size: 14}, Alignment: &excelize.Alignment{Horizontal: "center"}})
	return &ExcelStyle{
		Title:          title,
		TitleCenter:    titleCenter,
		Subtitle:       subtitle,
		SubtitleCenter: subtitleCenter,
		Bold:           bold,
		Heading:        heading,
		BoldCenter:     boldCenter,
		Normal:         normal,
		Center:         center,
		TextRight:      textRight,
		TextRightBold:  textRightBold,
	}
}
