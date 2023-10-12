package excelStyle

import (
	_ "embed"

	"github.com/xuri/excelize/v2"
)

// 只有文字
func GetFontStyle(f *excelize.File, fontFamily string) (int, error) {
	fontStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 11, Family: fontFamily},
	})
	if err != nil {
		return 0, err
	}
	return fontStyle, nil
}

// 右邊誆線+文字
func GetRightBorderAndFontStyle(f *excelize.File, fontFamily string) (int, error) {
	rightBorderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "right",
				Color: "#000000",
				Style: 1,
			}},
		Font: &excelize.Font{Size: 11, Family: fontFamily},
	})
	if err != nil {
		return 0, err
	}
	return rightBorderStyle, nil
}

// 金額格式+右邊誆線+文字
func GetPriceFormatAndRightBorderAndFontStyle(f *excelize.File, fontFamily string) (int, error) {
	numberFormat := "#,##0 "
	priceFormatAndRightBorderAndFontStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "right",
				Color: "#000000",
				Style: 1,
			}},
		Font:         &excelize.Font{Size: 11, Family: fontFamily},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return 0, err
	}
	return priceFormatAndRightBorderAndFontStyle, nil
}

// 金額格式+左邊誆線+文字
func GetPriceFormatAndLeftBorderAndFontStyle(f *excelize.File) (int, error) {
	numberFormat := "#,##0 "
	priceFormatAndRightBorderAndFontStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "#000000",
				Style: 1,
			}},
		Font:         &excelize.Font{Size: 11, Family: "儷宋 Pro"},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return 0, err
	}
	return priceFormatAndRightBorderAndFontStyle, nil
}

// 金額格式+兩種誆線+文字
func GetPriceFormatAndTwoBorderAndFontStyle(f *excelize.File, fontSize float64, types []string, fontFamily string, borderStyle []int) (int, error) {
	numberFormat := "#,##0 "
	priceFormatAndRightBorderAndFontStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  types[0],
				Color: "#000000",
				Style: borderStyle[0],
			}, {
				Type:  types[1],
				Color: "#000000",
				Style: borderStyle[1],
			}},
		Font:         &excelize.Font{Size: fontSize, Family: fontFamily},
		CustomNumFmt: &numberFormat,
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return 0, err
	}
	return priceFormatAndRightBorderAndFontStyle, nil
}

// 金額格式+文字
func GetPriceFormatAndFontStyle(f *excelize.File, fontFamily string) (int, error) {
	numberFormat := "#,##0 "
	priceFormatAndRightBorderAndFontStyle, err := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 11, Family: fontFamily},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return 0, err
	}
	return priceFormatAndRightBorderAndFontStyle, nil
}

// 文字置右+文字
func GetAlignmentAndFontStyle(f *excelize.File) (int, error) {
	alignmentStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 11, Family: "儷宋 Pro"},
		Alignment: &excelize.Alignment{Horizontal: "right"},
	})
	if err != nil {
		return 0, err
	}
	return alignmentStyle, nil
}

// 中右下外框+置中(自動換行)+文字
func GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle(f *excelize.File, fontSize float64) (int, error) {
	classStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "#000000",
				Style: 1,
			}},
		Font:      &excelize.Font{Size: fontSize, Family: "Calibri (本文)"},
		Alignment: &excelize.Alignment{Horizontal: "center", WrapText: true, Vertical: "center"},
	})
	if err != nil {
		return 0, err
	}
	return classStyle, nil
}

// 右下外框+置中(自動換行)+文字
func GetRightAndBottomBorderAndLeftAlignmentAndFontStyle(f *excelize.File) (int, error) {
	classStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "#000000",
				Style: 1,
			}},
		Font:      &excelize.Font{Size: 11, Family: "Calibri (本文)"},
		Alignment: &excelize.Alignment{Horizontal: "left", WrapText: true, Vertical: "center"},
	})
	if err != nil {
		return 0, err
	}
	return classStyle, nil
}

// 外框
func GetAllBorderStyle(f *excelize.File) (int, error) {
	allBorderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "#000000",
				Style: 1,
			}, {
				Type:  "top",
				Color: "#000000",
				Style: 1,
			}, {
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			}, {
				Type:  "right",
				Color: "#000000",
				Style: 1,
			}}})

	if err != nil {
		return 0, err
	}
	return allBorderStyle, nil
}

// 底色(通常用在副標題)
func GetFillColorStyle(f *excelize.File) (int, error) {
	fillColor, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D8D8D8"},
			Pattern: 1,
		}})

	if err != nil {
		return 0, err
	}
	return fillColor, nil
}
