package hsapp

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/inkyblackness/imgui-go"

	"github.com/OpenDiablo2/HellSpawner/hsapp/hseditors"
	"github.com/OpenDiablo2/HellSpawner/hsinterface"
	"github.com/OpenDiablo2/HellSpawner/hsproj"
	"github.com/OpenDiablo2/HellSpawner/hsutil"
)

var WindowWidth float32
var WindowHeight float32

type MainWindow struct {
	aboutDialog      *AboutDialog
	openFolderDialog *OpenFolderDialog
	projectTreeView  *ProjectTreeView
	dynamicWindows   []hsinterface.UIWindow

	doClose    bool
	menuHeight int
}

func CreateMainWindow() *MainWindow {
	icons := CreateMpqTreeIcons()

	result := &MainWindow{
		aboutDialog:    CreateAboutDialog(),
		dynamicWindows: make([]hsinterface.UIWindow, 0),
	}

	startdir, _ := os.UserHomeDir()
	result.openFolderDialog = CreateOpenFolderDialog("Select Project Folder", startdir, icons, func(folderPath string) {
		loadProjectFolder(folderPath)
		result.RefreshProjectLoaded()
	})

	result.projectTreeView = CreateProjectTreeView("ProjectTreeView", icons, result)

	latestFolder := AppState.LatestFolder()
	if latestFolder != "" {
		loadProjectFolder(latestFolder)
		result.RefreshProjectLoaded()
	}

	return result
}

func loadProjectFolder(folderPath string) {
	AppState.SetLatestFolder(folderPath)
	hsproj.ActiveProject.PromptUnsavedChanges()
	hsproj.ActiveProject.Close()

	newproj, err := hsproj.LoadProjectStateFromFolder(folderPath)
	if err != nil {
		hsutil.PopupError(err)
		return
	}

	hsproj.ActiveProject = newproj
}

func (v MainWindow) DoClose() bool {
	return v.doClose
}

func (v *MainWindow) Render() {
	v.renderMainMenu()
	v.renderProjectTreeView()

	for windowIdx := range v.dynamicWindows {
		v.dynamicWindows[windowIdx].Render()
	}

	for windowIdx := 0; windowIdx < len(v.dynamicWindows); {
		if v.dynamicWindows[windowIdx].IsClosed() {
			v.dynamicWindows = append(v.dynamicWindows[:windowIdx], v.dynamicWindows[windowIdx+1:]...)
			continue
		}
		windowIdx++
	}

	v.aboutDialog.Render()
	v.openFolderDialog.Render()
}

func (v *MainWindow) renderMainMenu() {
	showAboutPopup := false
	showOpenFolderDialog := false
	if imgui.BeginMainMenuBar() {

		if imgui.BeginMenu("File") {
			if imgui.MenuItem("Open Folder") {
				v.openFolderDialog.Show()
				showOpenFolderDialog = true
			}
			if len(AppState.RecentFolders) > 0 && imgui.BeginMenu("Recent") {
				for _, folder := range AppState.RecentFolders {
					if imgui.MenuItem(folder) {
						loadProjectFolder(folder)
						v.RefreshProjectLoaded()
					}
				}

				imgui.EndMenu()
			}
			if imgui.MenuItem("Exit") {
				v.doClose = true
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Help") {
			if imgui.MenuItem("About HellSpawner...") {
				showAboutPopup = true
			}
			imgui.EndMenu()
		}
		imgui.EndMainMenuBar()
	}

	// popups
	if showAboutPopup {
		imgui.OpenPopup(AboutDialogPopupName)
	}
	if showOpenFolderDialog {
		imgui.OpenPopup(v.openFolderDialog.PopupName)
	}
}

func (v *MainWindow) renderProjectTreeView() {
	imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 0.0)
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: 20})
	imgui.SetNextWindowSize(imgui.Vec2{X: 300, Y: WindowHeight - 20})
	if imgui.BeginV("Project Treeview", nil, imgui.WindowFlagsNoMove|imgui.WindowFlagsNoResize|imgui.WindowFlagsNoCollapse|imgui.WindowFlagsNoSavedSettings) {
		v.projectTreeView.Render()
		imgui.End()
	}
	imgui.PopStyleVar()
	/*imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 0.0)
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: 20})
	imgui.SetNextWindowSize(imgui.Vec2{X: 300, Y: WindowHeight - 20})
	if imgui.BeginV("Project Treeview", nil, imgui.WindowFlagsNoMove|imgui.WindowFlagsNoResize|imgui.WindowFlagsNoCollapse|imgui.WindowFlagsNoSavedSettings) {
		if imgui.TreeNode("Project") {
			imgui.TreePop()
		}

		imgui.End()
	}
	imgui.PopStyleVar()*/
}

// should call sub components and ask them to refresh because a new ActiveProject was loaded
func (v *MainWindow) RefreshProjectLoaded() {
	v.projectTreeView.Refresh()
}

func (v *MainWindow) OnFileSelected(filename string, mpqpath hsutil.MpqPath) {
	// try to open an editor based on the file extension
	var ed hsinterface.UIEditor

	log.Printf("Loading file '%v' from MPQ '%v'", filename, mpqpath)

	ext := strings.ToLower(filepath.Ext(mpqpath.FilePath))
	switch ext {
	case ".txt":
		ed = hseditors.CreateTextEditor(filename, mpqpath)
	}

	if ed != nil {
		v.dynamicWindows = append(v.dynamicWindows, hseditors.CreateEditorWindow(&ed))
	}
}

func (v *MainWindow) OnViewMpqFileDetails(mpqpath string) {
	v.dynamicWindows = append(v.dynamicWindows, CreateMpqDetailsDialog(hsproj.ActiveProject.FolderPath+mpqpath))
}
