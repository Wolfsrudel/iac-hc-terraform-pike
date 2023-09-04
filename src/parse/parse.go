package parse

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type provider struct {
	Resources   []string `json:"resources"`
	DataSources []string `json:"dataSources"`
}

func Parse(codebase string, name string) error {
	var err error

	myProvider := provider{}

	match := `"(` + strings.ToLower(name) + `_.*?)"`
	myProvider.Resources, err = GetMatches(codebase, match, "go")
	if err != nil {
		return err
	}

	myProvider.DataSources, err = GetMatches(codebase, `# Data Source:(.*)`, "markdown")
	if err != nil {
		return err
	}

	jsonOut, err := json.MarshalIndent(myProvider, "", "    ")

	if err != nil {
		return err
	}

	err = os.WriteFile(name+"-members.json", jsonOut, 0700)

	if err != nil {
		return err
	}

	return nil
}

func GetMatches(source string, match string, extension string) ([]string, error) {
	files, err := GetGoFiles(source, extension)

	if err != nil {
		return nil, err
	}

	var (
		m = make(map[string]bool)
		a []string
	)

	for _, file := range files {
		contents, _ := os.ReadFile(file)

		re := regexp.MustCompile(match)
		match := re.FindStringSubmatch(string(contents))

		for _, item := range match {
			if strings.Contains(item, "%s") {
				continue
			}

			matched := strings.TrimSpace(strings.ReplaceAll(item, "\"", ""))
			matched = strings.TrimSpace(strings.ReplaceAll(matched, "# Data Source: ", ""))
			a, m = add(matched, m, a)
		}
	}

	keys := GetKeys(m)

	return keys, nil
}

func GetGoFiles(path string, extension string) ([]string, error) {
	libRegEx, err := regexp.Compile("^.+\\." + extension + "$")
	if err != nil {
		return nil, err
	}

	var files []string

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			if strings.Contains(path, "_test") || strings.Contains(path, ".ci") || info.IsDir() {

			} else {
				files = append(files, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func GetKeys(m map[string]bool) []string {
	var keys []string

	for k := range m {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	return keys
}

func add(s string, m map[string]bool, a []string) ([]string, map[string]bool) {
	if m[s] {
		return a, m // Already in the map
	}
	a = append(a, s)

	m[s] = true

	return a, m
}
