package pkg

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

const main = "Main"

var tableStyle int
var linkStyle int
var redStyle int
var yellowStyle int

func PrepareReport(testObjects []TestObject, filename string) {
	f := excelize.NewFile()
	defineTableStyle(f)
	defineHiperLinkStyle(f)
	defineWarnStyle(f)

	createSuccessRatePage(f, len(testObjects))

	for i, testObject := range testObjects {

		createTestRunsStatistic(f, testObject, len(testObject.Extra.History.Items), i)
		var passed = 0
		total := len(testObject.Extra.History.Items)
		for j, item := range testObject.Extra.History.Items {

			if item.Status == "passed" {
				passed++
			}
			if item.Status == "skipped" {
				total--
			}

			f.SetCellStr(testObject.Uid, fmt.Sprint("B", j+4), time.UnixMilli(int64(item.Time.Start)).Format("2006-01-02 15:04:05"))
			f.SetCellStyle(testObject.Uid, fmt.Sprint("B", j+4), fmt.Sprint("B", j+4), linkStyle)
			f.SetCellHyperLink(testObject.Uid, fmt.Sprint("B", j+4), item.ReportUrl, "External")

			f.SetCellInt(testObject.Uid, fmt.Sprint("C", j+4), int(item.Time.Duration))
			f.SetCellStr(testObject.Uid, fmt.Sprint("D", j+4), item.Status)
			f.SetCellStr(testObject.Uid, fmt.Sprint("E", j+4), item.StatusDetails)
		}

		// fullfill global statistic page
		service, class := GetServiceAndClass(testObject.FullName)
		f.SetCellStr(main, fmt.Sprintf("B%d", i+4), service)
		f.SetCellStr(main, fmt.Sprintf("C%d", i+4), class)
		f.SetCellStyle(main, fmt.Sprint("C", i+4), fmt.Sprint("C", i+4), linkStyle)
		f.SetCellHyperLink(main, fmt.Sprint("C", i+4), GetUrlToGithub(testObject.FullName), "External")
		f.SetCellStr(main, fmt.Sprintf("D%d", i+4), testObject.Name)

		f.SetCellInt(main, fmt.Sprintf("E%d", i+4), testObject.Extra.History.Statistic.Failed)
		f.SetCellInt(main, fmt.Sprintf("F%d", i+4), testObject.Extra.History.Statistic.Broken)
		f.SetCellInt(main, fmt.Sprintf("G%d", i+4), testObject.Extra.History.Statistic.Skipped)
		f.SetCellInt(main, fmt.Sprintf("H%d", i+4), testObject.Extra.History.Statistic.Passed)
		f.SetCellInt(main, fmt.Sprintf("I%d", i+4), testObject.Extra.History.Statistic.Unknown)
		f.SetCellInt(main, fmt.Sprintf("J%d", i+4), testObject.Extra.History.Statistic.Total)
		setRateCell(f, main, testObject.Extra.History.Statistic.Passed, testObject.Extra.History.Statistic.Total, fmt.Sprintf("K%d", i+4))

		//fullfill success rate page
		f.SetCellInt(main, fmt.Sprintf("L%d", i+4), passed)
		f.SetCellInt(main, fmt.Sprintf("M%d", i+4), total)
		setRateCell(f, main, passed, total, fmt.Sprintf("N%d", i+4))
		f.SetCellStr(main, fmt.Sprintf("O%d", i+4), "link")
		f.SetCellStyle(main, fmt.Sprint("O", i+4), fmt.Sprint("O", i+4), linkStyle)
		f.SetCellHyperLink(main, fmt.Sprintf("O%d", i+4), fmt.Sprintf("%s!A1", testObject.Uid), "Location")
	}

	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
	}
}

func createTestRunsStatistic(f *excelize.File, object TestObject, size, sourceIndex int) {
	f.NewSheet(object.Uid)

	f.SetCellStr(object.Uid, "A1", "<back")
	f.SetCellStyle(object.Uid, "A1", "A1", linkStyle)
	f.SetCellHyperLink(object.Uid, "A1", fmt.Sprintf("%s!O%d", main, sourceIndex), "Location")

	f.SetCellStr(object.Uid, "B1", "Test name")
	f.SetCellStr(object.Uid, "C1", object.Name)
	f.MergeCell(object.Uid, "C1", "E1")
	f.SetCellStr(object.Uid, "B2", "Full name")
	f.SetCellStr(object.Uid, "C2", object.FullName)
	f.MergeCell(object.Uid, "C2", "E2")

	f.SetCellStr(object.Uid, "B3", "Started at")
	f.SetColWidth(object.Uid, "B", "B", 20)
	f.SetCellStr(object.Uid, "C3", "Duration")
	f.SetCellStr(object.Uid, "D3", "Status")
	f.SetCellStr(object.Uid, "E3", "Status Details")
	f.SetColWidth(object.Uid, "E", "E", 70)

	err := f.SetCellStyle(object.Uid, "B3", fmt.Sprint("E", size+3), tableStyle)
	if err != nil {
		fmt.Println(err)
	}

}

func createSuccessRatePage(f *excelize.File, size int) {

	f.SetSheetName("Sheet1", main)

	f.SetCellStr(main, "B3", "Service")
	f.SetCellStr(main, "C3", "Class")
	f.SetCellStr(main, "D3", "Name")

	f.SetCellStr(main, "E2", "Global statistic")
	f.MergeCell(main, "E2", "K2")

	f.SetCellStr(main, "E3", "Failed")
	f.SetCellStr(main, "F3", "Broken")
	f.SetCellStr(main, "G3", "Skipped")
	f.SetCellStr(main, "H3", "Passed")
	f.SetCellStr(main, "I3", "Unknown")
	f.SetCellStr(main, "J3", "Total")
	f.SetCellStr(main, "K3", "Success rate %")

	f.SetCellStr(main, "L2", "Observed statistic")
	f.MergeCell(main, "L2", "O2")

	f.SetCellStr(main, "L3", "Passed")
	f.SetCellStr(main, "M3", "Total")
	f.SetCellStr(main, "N3", "Success rate %")
	f.SetCellStr(main, "O3", "Details")

	f.SetColWidth(main, "B", "B", 10)
	f.SetColWidth(main, "C", "C", 20)
	f.SetColWidth(main, "D", "D", 30)
	f.SetColWidth(main, "E", "J", 5)

	err := f.SetCellStyle(main, "B2", fmt.Sprint("O", size+3), tableStyle)
	if err != nil {
		fmt.Println(err)
	}

}

func defineTableStyle(f *excelize.File) {
	tableStyle, _ = f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
}

func defineHiperLinkStyle(f *excelize.File) {

	linkStyle, _ = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#1265BE", Underline: "single"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
}

func defineWarnStyle(f *excelize.File) {
	redStyle, _ = f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type: "pattern", Color: []string{"FF5733"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	yellowStyle, _ = f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type: "pattern", Color: []string{"FFF233"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
}

func setRateCell(f *excelize.File, page string, passed, total int, cell string) {

	if total == 0 {
		f.SetCellStr(page, cell, "N/A")
		return
	} else {
		f.SetCellStr(page, cell, fmt.Sprintf("%.2f %%", (float32)((passed*100)/total)))
	}
	rate := (float32)(passed) / (float32)(total)
	if rate < 0.8 {
		f.SetCellStyle(page, cell, cell, redStyle)
	} else if rate < 1 {
		f.SetCellStyle(page, cell, cell, yellowStyle)
	}

}
