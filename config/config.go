package config

import (
	"bytes"
	_ "embed"

	"gopkg.in/yaml.v3"
)

/*
	Eventually these will be user configurable rather than compiled time
	options.
*/

//go:embed defaults.yaml
var defaults []byte

func init() {
	err := Default()
	if err != nil {
		panic(err)
	}
}

func Default() error {
	yml := yaml.NewDecoder(bytes.NewReader(defaults))
	yml.KnownFields(true)
	return yml.Decode(&Config)
}

var Config configT

type configT struct {
	Shell struct {
		Default  []string `yaml:"Default"`
		Fallback []string `yaml:"Fallback"`
	} `yaml:"Shell"`

	Terminal struct {
		ScrollbackHistory       int  `yaml:"ScrollbackHistory"`
		ScrollbackCloseKeyPress bool `yaml:"ScrollbackCloseKeyPress"`
		JumpScrollLineCount     int  `yaml:"JumpScrollLineCount"`
		RefreshInterval         int  `yaml:"RefreshInterval"`
		LightMode               bool `yaml:"LightMode"`

		TypeFace struct {
			FontName         string `yaml:"FontName"`
			FontSize         int    `yaml:"FontSize"`
			Ligatures        bool   `yaml:"Ligatures"`
			DropShadow       bool   `yaml:"DropShadow"`
			AdjustCellWidth  int    `yaml:"AdjustCellWidth"`
			AdjustCellHeight int    `yaml:"AdjustCellHeight"`
		} `yaml:"TypeFace"`
	} `yaml:"Terminal"`

	Window struct {
		Opacity  int      `yaml:"Opacity"`
		Fallback []string `yaml:"Fallback"`
	} `yaml:"Window"`

	Tmux struct {
		Enabled bool `yaml:"Enabled"`
	} `yaml:"Tmux"`
}
