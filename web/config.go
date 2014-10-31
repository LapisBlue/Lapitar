package web

import (
	"encoding/json"
	"io"
)

const (
	indent = "    "
)

type config struct {
	Address string        `json:"address" xml:"address"`
	Head    *headDefaults `json:"head" xml:"head"`
	Face    *faceDefaults `json:"face" xml:"face"`
}

type headDefaults struct {
	Width         int     `json:"width" xml:"width"`
	Height        int     `json:"height" xml:"height"`
	Angle         float32 `json:"angle" xml:"angle"`
	SuperSampling int     `json:"supersampling" xml:"supersampling"`
	Helm          bool    `json:"helm" xml:"helm"`
	Shadow        bool    `json:"shadow" xml:"shadow"`
	Lighting      bool    `json:"lighting" xml:"lighting"`
}

type faceDefaults struct {
	Size int  `json:"size" xml:"size"`
	Helm bool `json:"helm" xml:"helm"`
}

func defaultConfig() *config {
	return &config{
		":8088",
		&headDefaults{
			256, 256,
			45,
			4,
			true,
			true,
			true,
		}, &faceDefaults{
			256,
			true,
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
