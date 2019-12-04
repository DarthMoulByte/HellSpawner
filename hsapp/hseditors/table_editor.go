package hseditors

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/OpenDiablo2/HellSpawner/hsproj"
	"github.com/OpenDiablo2/HellSpawner/hsutil"

	"github.com/inkyblackness/imgui-go"
)

type TableEditor struct {
	FileName       string
	MpqPath        hsutil.MpqPath
	Table          Table
	unsavedChanges bool
	renderName     string
}

type ColumnHeader struct {
	// Label represents the label of the column header
	Label string
	// Size represents the size of the column header field. Set it to -1 to calculate the size
	Size imgui.Vec2

	SyncOffset float32
}

type Table struct {
	header []ColumnHeader
	rows   [][]string
}

func CreateExcelEditor(filename string, mpqpath hsutil.MpqPath) *TableEditor {
	ed := TableEditor{}
	ed.FileName = filename
	ed.MpqPath = mpqpath
	ed.renderName = "TableEditor" + ed.MpqPath.MpqName + ":" + ed.MpqPath.FilePath + ":"

	// load the text itself
	bytes, err := hsproj.ActiveProject.MpqList.LoadFile(mpqpath)

	if err != nil {
		hsutil.PopupError(err)
	}

	tableData, err := ed.parseCSV(bytes)

	if err != nil {
		hsutil.PopupError(err)
	}

	ed.Table.header = ed.convertCSVHeader(tableData)
	ed.Table.rows = tableData[1:]

	return &ed
}

func (v *TableEditor) Name() string {
	return v.FileName
}

func (v *TableEditor) LongName() string {
	return v.MpqPath.MpqName + " :: " + v.MpqPath.FilePath
}

func (v *TableEditor) RenderName() string {
	return v.renderName
}

func (v *TableEditor) HasUnsavedChanges() bool {
	return v.unsavedChanges
}

func (v *TableEditor) Save() {
	// TODO
	v.unsavedChanges = false
}

func (v *TableEditor) editCallback(d imgui.InputTextCallbackData) int32 {
	v.unsavedChanges = true
	return 0
}

func (v *TableEditor) Render(size imgui.Vec2) {
	v.drawColumnHeader()
	v.beginColumnHeaderSync()

	// Draw the rows here

	for i := range v.Table.rows {
		for j := range v.Table.header {
			imgui.Selectable(v.Table.rows[i][j])

			if imgui.IsItemHovered() {
				imgui.BeginTooltip()
				imgui.SetTooltip(v.Table.rows[i][j])
				imgui.EndTooltip()
			}

			imgui.NextColumn()
		}
	}

	v.endColumnHeaderSync()
}

func (v *TableEditor) parseCSV(content []byte) ([][]string, error) {
	buf := bytes.NewReader(content)

	reader := csv.NewReader(buf)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	return reader.ReadAll()
}

func (v *TableEditor) cellSize(idx int, table [][]string) imgui.Vec2 {
	word := ""

	for i := range table {
		if idx <= len(table[i]) {
			if len(table[i][idx]) > len(word) {
				word = table[i][idx]
			}
		}
	}

	cellSize := imgui.CalcTextSize(word, false, 500)
	//cellSize.X += 20 // give it some extra space

	return cellSize
}

func (v *TableEditor) headerSize() imgui.Vec2 {
	rc := imgui.Vec2{X: 0, Y: 0}

	for _, v := range v.Table.header {
		if rc.Y < v.Size.Y {
			rc.Y = v.Size.Y
		}

		rc.X += v.Size.X
	}

	return rc
}

func (v *TableEditor) convertCSVHeader(table [][]string) []ColumnHeader {
	header := table[0]
	result := make([]ColumnHeader, len(header))

	for i := range header {
		result[i] = ColumnHeader{
			Label: header[i],
			Size:  v.cellSize(i, table),
		}
	}

	return result
}

func (v *TableEditor) drawColumnHeader() {
	columnCount := len(v.Table.header)

	if columnCount <= 0 {
		return
	}

	style := imgui.CurrentStyle()
	firstSize := imgui.CalcTextSize(v.Table.header[0].Label, false, 500)
	id := fmt.Sprintf("header_%s", v.Name())

	imgui.BeginChildV(id, imgui.Vec2{X: 0, Y: firstSize.Y + 2*style.ItemInnerSpacing().Y}, true, 0)
	imgui.Columns(columnCount, id)

	var offset float32 = 0.0

	for i := 0; i < columnCount; i++ {
		header := &v.Table.header[i]

		if header.SyncOffset < 0.0 {
			imgui.SetColumnOffset(i, offset)
			offset += header.Size.X
		} else {
			imgui.SetColumnOffset(i, header.SyncOffset)
		}

		header.SyncOffset = imgui.ColumnOffsetV(i)

		imgui.Selectable(header.Label) // Make header selectable
		imgui.NextColumn()
	}

	imgui.Columns(1, v.Name())
	imgui.EndChild()
}

func (v *TableEditor) beginColumnHeaderSync() {
	columnCount := len(v.Table.header)

	if columnCount <= 0 {
		return
	}

	id := fmt.Sprintf("rows_%s", v.Name())

	imgui.BeginChildV(id, imgui.Vec2{X: 0, Y: 0}, true, 0)
	imgui.Columns(columnCount, id)

	for i := range v.Table.header {
		header := &v.Table.header[i]
		header.SyncOffset = imgui.ColumnOffsetV(i)
	}
}

func (v *TableEditor) endColumnHeaderSync() {
	columnCount := len(v.Table.header)

	if columnCount <= 0 {
		return
	}

	id := fmt.Sprintf("rows_%s", v.Name())

	imgui.Columns(1, id)
	imgui.EndChild()
}
