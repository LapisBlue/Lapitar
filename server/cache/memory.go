package cache

import (
	"github.com/LapisBlue/Lapitar/mc"
	"sync"
	"time"
)

type memorySkinCache struct {
	uuidLoader     map[string]chan skinMetaResult
	uuidLoaderLock sync.RWMutex

	skins          map[string]SkinMeta
	skinsLock      sync.RWMutex
	skinLoader     map[string]chan skinMetaResult
	skinLoaderLock sync.RWMutex
}

func Memory() SkinCache {
	return &memorySkinCache{
		uuidLoader: make(map[string]chan skinMetaResult),
		skins:      make(map[string]SkinMeta),
		skinLoader: make(map[string]chan skinMetaResult),
	}
}

type memorySkinMeta struct {
	profile   mc.Profile
	meta      mc.SkinMeta
	timestamp time.Time
	skin      mc.Skin
	loader    chan skinResult
}

type skinMetaResult struct {
	meta SkinMeta
	err  error
}

type skinResult struct {
	skin mc.Skin
	err  error
}

func (cache *memorySkinCache) FetchByName(name string) (skin SkinMeta, err error) {
	name = mc.ToLower(name)
	skin = cache.pullSkin(name)
	if skin != nil {
		return
	}

	loader := cache.pullUUIDLoader(name)

	if loader != nil {
		result := <-loader
		return result.meta, result.err
	}

	loader = make(chan skinMetaResult, 1)
	go cache.pushUUIDLoader(name, loader)

	profile, err := mc.FetchProfile(name)
	if err != nil {
		loader <- skinMetaResult{nil, err}
		return
	}

	skin, err = Fetch(profile.UUID())
	loader <- skinMetaResult{skin, err}

	go cache.pushUUIDLoader(name, nil)
	return
}

func (cache *memorySkinCache) Fetch(uuid string) (skin SkinMeta, err error) {
	uuid = mc.ParseUUID(uuid)

	skin = cache.pullSkin(uuid)
	if skin != nil {
		return
	}

	loader := cache.pullSkinLoader(uuid)

	if loader != nil {
		result := <-loader
		return result.meta, result.err
	}

	loader = make(chan skinMetaResult, 1)
	go cache.pushSkinLoader(uuid, loader)

	sk, err := mc.FetchSkin(uuid)
	if err != nil {
		loader <- skinMetaResult{nil, err}
		return
	}

	skin = &memorySkinMeta{sk.Profile(), sk.Skin(), time.Now(), nil, nil}
	loader <- skinMetaResult{skin, nil}

	go cache.pushSkin(sk.Profile(), skin)
	go cache.pushSkinLoader(uuid, nil)
	return
}

func (cache *memorySkinCache) pullUUIDLoader(name string) chan skinMetaResult {
	cache.uuidLoaderLock.RLock()
	defer cache.uuidLoaderLock.RUnlock()
	return cache.uuidLoader[name]
}

func (cache *memorySkinCache) pushUUIDLoader(name string, loader chan skinMetaResult) {
	cache.uuidLoaderLock.Lock()
	defer cache.uuidLoaderLock.Unlock()
	cache.uuidLoader[name] = loader
}

func (cache *memorySkinCache) pullSkin(uuid string) SkinMeta {
	cache.skinsLock.RLock()
	defer cache.skinsLock.RUnlock()
	return cache.skins[uuid]
}

func (cache *memorySkinCache) pushSkin(profile mc.Profile, meta SkinMeta) {
	cache.skinsLock.Lock()
	defer cache.skinsLock.Unlock()
	cache.skins[mc.ToLower(profile.UUID())] = meta
	cache.skins[mc.ToLower(profile.Name())] = meta
}

func (cache *memorySkinCache) pullSkinLoader(uuid string) chan skinMetaResult {
	cache.skinLoaderLock.RLock()
	defer cache.skinLoaderLock.RUnlock()
	return cache.skinLoader[uuid]
}

func (cache *memorySkinCache) pushSkinLoader(uuid string, loader chan skinMetaResult) {
	cache.skinLoaderLock.Lock()
	defer cache.skinLoaderLock.Unlock()
	cache.skinLoader[uuid] = loader
}

func (meta *memorySkinMeta) Profile() mc.Profile {
	return meta.profile
}

func (meta *memorySkinMeta) ID() string {
	return meta.meta.ID()
}

func (meta *memorySkinMeta) URL() string {
	return meta.meta.URL()
}

func (meta *memorySkinMeta) Download() (sk mc.Skin, err error) {
	if meta.skin != nil {
		return meta.skin, nil
	}

	if meta.loader != nil {
		result := <-meta.loader
		return result.skin, result.err
	}
	loader := make(chan skinResult, 1)
	meta.loader = loader

	sk, err = meta.meta.Download()
	if err != nil {
		loader <- skinResult{nil, err}
		return
	}

	meta.skin = sk
	loader <- skinResult{sk, nil}
	return
}

func (meta *memorySkinMeta) Timestamp() time.Time {
	return meta.timestamp
}
