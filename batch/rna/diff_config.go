package rna

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// DiffConfig is a model of config file for Diff function.
type DiffConfig struct {
	ConfigFile   string
	ProjectPath  string       `json:"projectPath"`
	Groups       []Group      `json:"groups"`
	Comparisons  []Comparison `json:"comparisons"`
	CuffdiffExec string       `json:"cuffdiffExec"`
	GenomeBtwIdx string       `json:"genomeBtwIdx"`
	GtfPath      string       `json:"gtfPath"`
	TranxBtwIdx  string       `json:"tranxBtwIdx"`
}

// Parse is a method for DiffConfig, which we use to parse config file.
func (dc *DiffConfig) Parse() error {
	// Return error if `ConfigFile` not exists.
	if _, err := os.Stat(dc.ConfigFile); os.IsNotExist(err) {
		return fmt.Errorf("Error: %s", err)
	}

	// Store _Clean_ & absolute file path of `ConfigFile`
	if filepath.IsAbs(dc.ConfigFile) {
		dc.ConfigFile = filepath.Clean(dc.ConfigFile)
	} else {
		absPath, absErr := filepath.Abs(dc.ConfigFile)
		if absErr != nil {
			return fmt.Errorf("Error: %s", absErr)
		}
		dc.ConfigFile = absPath
	}

	content, fileErr := ioutil.ReadFile(dc.ConfigFile)
	if fileErr != nil {
		return fmt.Errorf("Error: %s", fileErr)
	}

	jsonErr := json.Unmarshal(content, &dc)
	if jsonErr != nil {
		return fmt.Errorf("Error: %s", jsonErr)
	}

	return nil
}

// GetGroupByName returns the first `Group` by `Name` matching `name`.
func (dc *DiffConfig) GetGroupByName(name string) (Group, error) {
	for _, g := range dc.Groups {
		if g.Name == name {
			return g, nil
		}
	}

	return Group{}, errors.New("Group not found")
}

// Group is a named set of samples.
type Group struct {
	Name    string   `json:"name"`
	Samples []string `json:"samples"`
}

// Comparison is a pair of `Case` and `Control` group names.
type Comparison struct {
	Case    string `json:"case"`
	Control string `json:"control"`
}
