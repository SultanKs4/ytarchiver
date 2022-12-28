package fileutils

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

// Sanitize String for Filename
func SanitizeFilename(fileName string) string {
	// Characters not allowed on mac
	//	:/
	// Characters not allowed on linux
	//	/
	// Characters not allowed on windows
	//	<>:"/\|?*

	// Ref https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file#naming-conventions

	fileName = regexp.MustCompile(`[:/<>\:"\\|?*]`).ReplaceAllString(fileName, "")
	fileName = regexp.MustCompile(`\s+`).ReplaceAllString(fileName, " ")

	return fileName
}

// Get Base Root Project
func BasePath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	basepath = filepath.Join(basepath, "../../")
	return basepath
}

// Get Resource Path
func ResourcePath() string {
	basePath := filepath.Join(BasePath(), "resources")
	_ = os.MkdirAll(basePath, os.ModePerm)
	return basePath
}

// Scan where channelData.json inside path with one nested folder
func GetDataJson(path string) ([]string, error) {
	dirResource, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	listFile := []string{}
	for _, v := range dirResource {
		if !v.IsDir() {
			continue
		}
		nestedPath := filepath.Join(path, v.Name())
		nestedDir, _ := os.ReadDir(nestedPath)
		for _, v := range nestedDir {
			if v.Name() == "channelData.json" {
				fileJson := filepath.Join(nestedPath, v.Name())
				listFile = append(listFile, fileJson)
				break
			}
		}
	}
	return listFile, nil
}
