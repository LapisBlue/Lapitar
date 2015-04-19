package mc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/LapisBlue/lapitar/util/lhttp"
	"path"
)

const skinProfileURL = "https://sessionserver.mojang.com/session/minecraft/profile/"

type mojangSkinProfile struct {
	mojangProfile
	ProfileSkin mojangSkinMeta `json:"properties"`
}

func (p mojangSkinProfile) Profile() Profile {
	return p.mojangProfile
}

func (p mojangSkinProfile) Skin() SkinMeta {
	return p.ProfileSkin
}

type mojangSkinMeta struct {
	id   string
	url  string
	alex bool
}

func (meta mojangSkinMeta) ID() string {
	return meta.id
}

func (meta mojangSkinMeta) URL() string {
	return meta.url
}

func FetchSkin(uuid string) (p SkinProfile, err error) {
	req, err := lhttp.Get(skinProfileURL + uuid)
	if err != nil {
		return
	}

	req.Header.Set("Accept", lhttp.TypeJSON)

	resp, err := lhttp.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if lhttp.IsNoContent(resp) {
		return
	}

	err = lhttp.ExpectSuccess(resp)
	if err != nil {
		return
	}

	err = lhttp.ExpectContent(resp, lhttp.TypeJSON)
	if err != nil {
		return
	}

	profile := mojangSkinProfile{}
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil || profile.ProfileId == "" || profile.ProfileName == "" {
		return
	}

	p = profile
	return
}

type mojangProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (meta *mojangSkinMeta) UnmarshalJSON(data []byte) (err error) {
	var properties []mojangProperty
	err = json.Unmarshal(data, &properties)
	if err != nil {
		return
	}

	for _, property := range properties {
		if property.Name == "textures" {
			data, err = base64.StdEncoding.DecodeString(property.Value)
			if err != nil {
				return
			}

			var i interface{}
			err = json.Unmarshal(data, &i)
			if err != nil {
				return
			}

			textures, ok := i.(map[string]interface{})
			if !ok {
				err = errors.New("Invalid texture properties: missing property content")
				return
			}

			textures, ok = textures["textures"].(map[string]interface{})
			if !ok {
				err = errors.New("Invalid texture properties: missing texture content")
				return
			}

			if skin, ok := textures["SKIN"].(map[string]interface{}); ok {
				if meta.url, ok = skin["url"].(string); ok {
					meta.id = path.Base(meta.url)
				} else {
					err = errors.New("Invalid skin: missing url")
					return
				}

				if i, ok = skin["metadata"]; ok {
					metadata, ok := i.(map[string]interface{})
					if !ok {
						err = errors.New("Invalid skin: invalid metadata")
						return
					}

					if model, ok := metadata["model"]; ok {
						meta.alex = model == "slim"
					}
				}
			} else {
				err = errors.New("Invalid texture property: missing skin property")
				return
			}

			break
		}
	}

	return nil
}
