package hseditors

import (
	"bufio"
	"github.com/OpenDiablo2/HellSpawner/hsinterface"
	"github.com/OpenDiablo2/HellSpawner/hsproj"
	"github.com/OpenDiablo2/HellSpawner/hsutil"
	"strings"

	"github.com/inkyblackness/imgui-go"
)

type TextEditor struct {
	FileName string
	MpqPath  hsutil.MpqPath

	Text *string

	unsavedChanges bool

	isTable bool

	tableView hsinterface.UIEditor
	tableWindow *EditorWindow

	renderName string
}

func CreateTextEditor(filename string, mpqpath hsutil.MpqPath) *TextEditor {
	ed := TextEditor{isTable: false, FileName:filename, MpqPath:mpqpath}
	ed.renderName = "TextEditor" + ed.MpqPath.MpqName + ":" + ed.MpqPath.FilePath + ":"

	// load the text itself
	bytes, err := hsproj.ActiveProject.MpqList.LoadFile(mpqpath)

	var str string

	if err != nil {
		hsutil.PopupError(err)
		str = "Failed to load file"
	} else {
		str = string(bytes)

	}

	ed.Text = &str
	ed.isTable = ed.isCSV(&str)

	return &ed
}

func (v *TextEditor) Name() string {
	return v.FileName
}

func (v *TextEditor) LongName() string {
	return v.MpqPath.MpqName + " :: " + v.MpqPath.FilePath
}

func (v *TextEditor) RenderName() string {
	return v.renderName
}

func (v *TextEditor) HasUnsavedChanges() bool {
	return v.unsavedChanges
}

func (v *TextEditor) Save() {
	// TODO
	v.unsavedChanges = false
}

func (v *TextEditor) editCallback(d imgui.InputTextCallbackData) int32 {
	v.unsavedChanges = true
	return 0
}

func (v *TextEditor) isCSV(data *string) bool {
	scanner := bufio.NewScanner(strings.NewReader(*data))
	scanner.Scan()
	return len(strings.Split(scanner.Text(), "\t")) > 1 // more than one header element
}

func (v *TextEditor) renderTableView() {
	if v.tableWindow == nil || v.tableView == nil {
		return
	}

	if v.tableWindow.IsClosed() {
		v.tableView = nil
		v.tableWindow = nil
		return
	}

	v.tableWindow.Render()
}

func (v *TextEditor) Render(size imgui.Vec2) {
	if v.isTable {
		if imgui.Button("Show in table view") && v.tableWindow == nil {
			v.tableView = CreateExcelEditor(v.FileName, v.MpqPath)
			v.tableWindow = CreateEditorWindow(&v.tableView)
		}
	}
	imgui.InputTextMultilineV(v.renderName + "InputText", v.Text,
		imgui.Vec2{X: size.X - 10, Y: size.Y - 30},
		imgui.InputTextFlagsAllowTabInput|imgui.InputTextFlagsCallbackCharFilter,
		v.editCallback)
	v.renderTableView()
}
