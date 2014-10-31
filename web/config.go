package web

import (
	"encoding/json"
	"io"
)

const (
	indent = "    "
)

type config struct {
	Address string      `json:"address" xml:"address"`
	Head    *headConfig `json:"head" xml:"head"`
	Face    *faceConfig `json:"face" xml:"face"`
}

type headConfig struct {
	Size          int     `json:"size" xml:"size" schema:"size"`
	Angle         float32 `json:"angle" xml:"angle" schema:"angle"`
	SuperSampling int     `json:"supersampling" xml:"supersampling"`
	Helm          bool    `json:"helm" xml:"helm" schema:"helm"`
	Shadow        bool    `json:"shadow" xml:"shadow"`
	Lighting      bool    `json:"lighting" xml:"lighting"`
}

type faceConfig struct {
	Size int  `json:"size" xml:"size" schema:"size"`
	Helm bool `json:"helm" xml:"helm" schema:"helm"`
}

func defaultConfig() *config {
	return &config{
		":8088",
		&headConfig{
			256,
			45,
			4,
			true,
			true,
			true,
		}, &faceConfig{
			256,
			false,
		},
	}
}

func parseConfig(reader io.Reader) (conf *config, err error) {
	conf = defaultConfig()
	err = json.NewDecoder(reader).Decode(conf)
	return
}

func writeConfig(writer io.Writer, conf *config) (err error) {
	// We can't use an encoder here because it cannot print with indentation
	buf, err := json.MarshalIndent(conf, "", indent)
	if err != nil {
		return
	}
	_, err = writer.Write(buf)
	if err != nil {
		return
	}
	_, err = io.WriteString(writer, "\n")
	return
}
