package server

import (
	"encoding/json"
	"github.com/LapisBlue/lapitar/face"
	"github.com/LapisBlue/lapitar/render"
	"github.com/LapisBlue/lapitar/util"
	"github.com/disintegration/imaging"
	"io"
)

const (
	indent = "    "
)

type config struct {
	Address  string          `json:"address"`
	Proxy    bool            `json:"proxy"`
	Defaults *renderDefaults `json:"defaults"`
}

type renderDefaults struct {
	Head     *renderConfig `json:"head"`
	Portrait *renderConfig `json:"portrait"`
	Body     *renderConfig `json:"body"`
	Face     *faceConfig   `json:"face"`
}

type faceConfig struct {
	Size  *limitedInt `json:"size"`
	Scale *scaling    `json:"scale"`
}

type renderConfig struct {
	*faceConfig
	Overlay       bool    `json:"overlay"`
	Angle         float32 `json:"angle"`
	Tilt          float32 `json:"tilt"`
	Zoom          float32 `json:"zoom"`
	SuperSampling int     `json:"supersampling"`
	Shadow        bool    `json:"shadow"`
	Lighting      bool    `json:"lighting"`
}

type limitedInt struct {
	Def int `json:"default"`
	Max int `json:"max"`
}

func defaultConfig() *config {
	return &config{
		":8088",
		false,
		&renderDefaults{
			newRenderConfig(-35, 20, -4.5),
			newRenderConfig(25, 10, -4.5),
			newRenderConfig(25, 10, -6),
			&faceConfig{
				&limitedInt{128, 512},
				&scaling{face.DefaultScale},
			},
		},
	}
}

func newRenderConfig(angle, tilt, zoom float32) *renderConfig {
	return &renderConfig{
		&faceConfig{
			&limitedInt{128, 512},
			&scaling{render.DefaultScale},
		},
		true,
		angle,
		tilt,
		zoom,
		4,
		true,
		true,
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

type scaling struct {
	*imaging.ResampleFilter
}

func (scale *scaling) Get() *imaging.ResampleFilter {
	return scale.ResampleFilter
}

func (scale *scaling) MarshalText() ([]byte, error) {
	return []byte(util.ScaleName(scale.ResampleFilter)), nil
}

func (scale *scaling) UnmarshalText(text []byte) (err error) {
	scale.ResampleFilter, err = util.ParseScale(string(text))
	return
}
