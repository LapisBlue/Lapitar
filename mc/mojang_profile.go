package mc

import (
	"encoding/json"
	"github.com/LapisBlue/Lapitar/util/lhttp"
)

const (
	profileURL = "https://api.mojang.com/users/profiles/minecraft/"
)

type mojangProfile struct {
	ProfileId   string `json:"id"`
	ProfileName string `json:"name"`
}

func (p mojangProfile) Name() string {
	return p.ProfileName
}

func (p mojangProfile) UUID() string {
	return p.ProfileId
}

func FetchProfile(player string) (p Profile, err error) {
	// TODO: Is validation necessary here?

	req, err := lhttp.Get(profileURL + player)
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

	profile := mojangProfile{}
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil || profile.ProfileName == "" || profile.ProfileId == "" {
		return
	}

	p = profile
	return
}
