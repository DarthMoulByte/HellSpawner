package hsapp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
)

var AppState *State

type State struct {
	RecentFolders []string `json:"recentFolders"`
}

func (s *State) LatestFolder() string {
	if len(s.RecentFolders) == 0 {
		return ""
	}
	return s.RecentFolders[0]
}


func (s *State) SetLatestFolder(openedFolder string)  {
	s.setLatestFolder(openedFolder)
	_ = SaveAppState()
}

// setLatestFolder Insert or move folder to the top of the recent folders.
func (s *State) setLatestFolder(openedFolder string)  {
	maxLength := 8
	recent := make([]string, 0, maxLength)
	recent = append(recent, openedFolder)
	for i := 0; i < len(s.RecentFolders) && i < maxLength; i++ {
		if i >= maxLength {
			break
		}

		if s.RecentFolders[i] != openedFolder {
			recent = append(recent, s.RecentFolders[i])
		}
	}
	s.RecentFolders = recent
}

func LoadAppState() {
	f, err := os.Open(getConfigSavePath())
	if err != nil {
		if !os.IsNotExist(err) {
			log.Panicf("failed to open config: %s", err.Error())
		}

		AppState = &State{}
		_ = SaveAppState()
		return
	}
	defer f.Close()

	AppState = loadAppState(f)
}

func SaveAppState() error {
	data, err := json.Marshal(AppState)
	if err != nil {
		log.Printf("failed to marshal config: %s", err.Error())
		return err
	}

	filePath := getConfigSavePath()
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Printf("failed to create config: %s", err.Error())
		return err
	}
	return nil
}

func loadAppState(r io.Reader) *State {
	result := &State{
	}
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Printf("failed to read config: %s", err.Error())
	}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		log.Printf("failed to unmarshal config: %s", err.Error())
	}

	return result
}

func getConfigSavePath() string {
	return path.Join(getAppBaseSavePath(), "config.json")
}

func getAppBaseSavePath() string {
	if runtime.GOOS == "windows" {
		appDataPath := os.Getenv("APPDATA")
		basePath := path.Join(appDataPath, "HellSpawner")
		if err := os.MkdirAll(basePath, os.ModeDir); err != nil {
			log.Panicf(err.Error())
		}
		return basePath
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Panicf(err.Error())
	}
	basePath := path.Join(configDir, "HellSpawner")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Panicf(err.Error())
	}
	return basePath
}
