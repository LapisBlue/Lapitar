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
	Size          *limitedInt `json:"size" xml:"size" schema:"size"`
	Angle         float32     `json:"angle" xml:"angle" schema:"angle"`
	SuperSampling int         `json:"supersampling" xml:"supersampling"`
	Helm          bool        `json:"helm" xml:"helm" schema:"helm"`
	Shadow        bool        `json:"shadow" xml:"shadow"`
	Lighting      bool        `json:"lighting" xml:"lighting"`
}

type faceConfig struct {
	Size *limitedInt `json:"size" xml:"size" schema:"size"`
	Helm bool        `json:"helm" xml:"helm" schema:"helm"`
}

type limitedInt struct {
	Def int `json:"default" xml:"default"`
	Max int `json:"max", xml:"max"`
}

func defaultConfig() *config {
	return &config{
		":8088",
		&headConfig{
			&limitedInt{256, 512},
			45,
			4,
			true,
			true,
			true,
		}, &faceConfig{
			&limitedInt{256, 512},
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
