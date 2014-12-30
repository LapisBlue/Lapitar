package cache

import (
	"github.com/LapisBlue/Lapitar/lapitar/mc"
	"log"
	"sync"
	"time"
)

type memorySkinCache struct {
	uuidLoader     map[string]chan SkinMeta
	uuidLoaderLock sync.RWMutex

	skins          map[string]SkinMeta
	skinsLock      sync.RWMutex
	skinLoader     map[string]chan SkinMeta
	skinLoaderLock sync.RWMutex
}

func Memory() SkinCache {
	result := &memorySkinCache{
		uuidLoader: make(map[string]chan SkinMeta),
		skins:      make(map[string]SkinMeta),
		skinLoader: make(map[string]chan SkinMeta),
	}

	result.skins["steve"] = Steve()
	result.skins["char"] = Steve()
	result.skins["alex"] = Alex()
	return result
}

type memorySkinMeta struct {
	mc.SkinMeta
	fallback  SkinMeta
	profile   mc.Profile
	timestamp time.Time
	skin      mc.Skin
	lock      sync.RWMutex
}

func (cache *memorySkinCache) FetchByName(realName string) (skin SkinMeta) {
	name := mc.ToLower(realName)
	skin = cache.pullSkin(name)
	if skin != nil {
		return
	}

	loader := cache.pullUUIDLoader(name)
	if loader != nil {
		skin = <-loader
		loader <- skin
		return skin
	}

	loader = make(chan SkinMeta, 1)
	go cache.pushUUIDLoader(name, loader)
	defer func() {
		go cache.pushUUIDLoader(name, nil)
	}()

	profile, err := mc.FetchProfile(name)
	if profile == nil || err != nil {
		log.Println("Failed to fetch UUID for", realName, err)
		skin = FallbackByName(realName)
		loader <- skin
		go cache.pushFallbackSkin(name, skin)
		return
	}

	skin = Fetch(profile.UUID())
	loader <- skin
	return
}

func (cache *memorySkinCache) Fetch(uuid string) (skin SkinMeta) {
	uuid = mc.ParseUUID(uuid)

	skin = cache.pullSkin(uuid)
	if skin != nil {
		return
	}

	loader := cache.pullSkinLoader(uuid)

	if loader != nil {
		skin = <-loader
		loader <- skin
		return
	}

	loader = make(chan SkinMeta, 1)
	go cache.pushSkinLoader(uuid, loader)
	defer func() {
		go cache.pushSkinLoader(uuid, nil)
	}()

	sk, err := mc.FetchSkin(uuid)
	if sk == nil || err != nil {
		log.Println("Failed to skin profile for", uuid, err)
		skin = FallbackByUUID(uuid)
		loader <- skin
		go cache.pushFallbackSkin(uuid, skin)
		return
	}

	skin = &memorySkinMeta{
		SkinMeta:  sk.Skin(),
		profile:   sk.Profile(),
		timestamp: time.Now(),
	}
	loader <- skin
	go cache.pushSkin(sk.Profile(), skin)
	return
}

func (cache *memorySkinCache) pullUUIDLoader(name string) chan SkinMeta {
	cache.uuidLoaderLock.RLock()
	defer cache.uuidLoaderLock.RUnlock()
	return cache.uuidLoader[name]
}

func (cache *memorySkinCache) pushUUIDLoader(name string, loader chan SkinMeta) {
	cache.uuidLoaderLock.Lock()
	defer cache.uuidLoaderLock.Unlock()
	cache.uuidLoader[name] = loader
}

func (cache *memorySkinCache) pullSkin(uuid string) SkinMeta {
	cache.skinsLock.RLock()
	defer cache.skinsLock.RUnlock()
	return cache.skins[uuid]
}

func (cache *memorySkinCache) pushFallbackSkin(id string, meta SkinMeta) {
	cache.skinsLock.Lock()
	defer cache.skinsLock.Unlock()
	cache.skins[id] = meta
}

func (cache *memorySkinCache) pushSkin(profile mc.Profile, meta SkinMeta) {
	cache.skinsLock.Lock()
	defer cache.skinsLock.Unlock()
	cache.skins[mc.ToLower(profile.UUID())] = meta
	cache.skins[mc.ToLower(profile.Name())] = meta
}

func (cache *memorySkinCache) pullSkinLoader(uuid string) chan SkinMeta {
	cache.skinLoaderLock.RLock()
	defer cache.skinLoaderLock.RUnlock()
	return cache.skinLoader[uuid]
}

func (cache *memorySkinCache) pushSkinLoader(uuid string, loader chan SkinMeta) {
	cache.skinLoaderLock.Lock()
	defer cache.skinLoaderLock.Unlock()
	cache.skinLoader[uuid] = loader
}

func (meta *memorySkinMeta) Profile() mc.Profile {
	return meta.profile
}

func (meta *memorySkinMeta) Fetch() (m SkinMeta, sk mc.Skin) {
	m = meta

	meta.lock.RLock()
	if meta.fallback == nil {
		sk = meta.skin
	} else {
		m, sk = meta.fallback.Fetch()
	}
	meta.lock.RUnlock()

	if sk != nil {
		return
	}

	// We need to download the skin first
	meta.lock.Lock()
	defer meta.lock.Unlock()

	if meta.fallback == nil {
		sk = meta.skin
	} else {
		m, sk = meta.fallback.Fetch()
	}

	if sk != nil {
		return
	}

	sk, err := meta.SkinMeta.Download()
	if err != nil {
		// Meh, we can't download the skin right now
		log.Println("Failed to fetch skin from", meta.URL(), err)
		m, sk = Fallback(meta.Profile()).Fetch()
		meta.fallback = m
		return
	}

	meta.skin = sk
	return
}

func (meta *memorySkinMeta) Download() (mc.Skin, error) {
	_, sk := meta.Fetch()
	return sk, nil
}

func (meta *memorySkinMeta) Timestamp() time.Time {
	return meta.timestamp
}
