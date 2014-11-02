package web

import (
	"encoding/json"
	"github.com/LapisBlue/Lapitar/face"
	"github.com/LapisBlue/Lapitar/head"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/disintegration/imaging"
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

type faceConfig struct {
	Size  *limitedInt `json:"size" xml:"size" schema:"size"`
	Helm  bool        `json:"helm" xml:"helm" schema:"helm"`
	Scale *scaling    `json:"scale" xml:"scale"`
}

type headConfig struct {
	*faceConfig
	Angle         float32 `json:"angle" xml:"angle" schema:"angle"`
	SuperSampling int     `json:"supersampling" xml:"supersampling"`
	Shadow        bool    `json:"shadow" xml:"shadow"`
	Lighting      bool    `json:"lighting" xml:"lighting"`
}

type limitedInt struct {
	Def int `json:"default" xml:"default"`
	Max int `json:"max", xml:"max"`
}

func defaultConfig() *config {
	return &config{
		":8088",
		&headConfig{
			&faceConfig{
				&limitedInt{256, 512},
				true,
				&scaling{head.DefaultScale},
			},
			-35,
			4,
			true,
			true,
		}, &faceConfig{
			&limitedInt{256, 512},
			false,
			&scaling{face.DefaultScale},
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
