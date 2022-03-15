package clearing

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func convertSliceToMap(slice []string) map[string]struct{} {
	res := make(map[string]struct{})
	for _, val := range slice {
		res[strings.ToLower(val)] = struct{}{}
	}
	return res
}

const pathLangsToExsts = "../../configs/language_extensions.json"

type langEntry struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

func convertLangsToMap(langs []string) map[string]struct{} {
	langsMap := convertSliceToMap(langs)
	res := make(map[string]struct{})
	fl, err := os.Open(pathLangsToExsts)
	if err != nil {
		return nil
	}
	decoder := json.NewDecoder(fl)
	_, err = decoder.Token()
	if err != nil {
		return nil
	}

	for decoder.More() {
		var entry langEntry
		err = decoder.Decode(&entry)
		if err != nil {
			return nil
		}
		if _, ok := langsMap[strings.ToLower(entry.Name)]; ok {
			for _, ext := range entry.Extensions {
				res[ext] = struct{}{}
			}
		}
	}
	_, err = decoder.Token()
	if err != nil {
		return nil
	}
	return res
}

func TakeFilesToExtensions(files, langs, extensions []string) []string {
	emptyLangs, emplyExtensions := len(langs) == 0, len(extensions) == 0
	if emptyLangs && emplyExtensions {
		return files
	}
	mapLangs, mapExtensions := convertLangsToMap(langs), convertSliceToMap(extensions)
	var result []string

	for _, file := range files {
		extension := filepath.Ext(file)
		_, langInsert := mapLangs[extension]
		_, extensionsInsert := mapExtensions[extension]
		if (emptyLangs || langInsert) && (emplyExtensions || extensionsInsert) {
			result = append(result, file)
		}
	}
	return result
}

func TakeFilesToPaths(files, excludePaths, restrictPaths []string) []string {
	emptyExclude, emptyRestrict := len(excludePaths) == 0, len(restrictPaths) == 0
	if emptyExclude && emptyRestrict {
		return files
	}
	mapExclude, mapRestrict := convertSliceToMap(excludePaths), convertSliceToMap(restrictPaths)
	var result []string
out:
	for _, file := range files {
		for exclude := range mapExclude {
			if match, _ := filepath.Match(exclude, file); match {
				continue out
			}
		}

		if emptyRestrict {
			result = append(result, file)
		}

		for restrict := range mapRestrict {
			if match, _ := filepath.Match(restrict, file); match {
				result = append(result, file)
				continue out
			}
		}
	}
	return result
}
