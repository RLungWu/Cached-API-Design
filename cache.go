package main

import "sync"

type CacheAd struct {
	ad       Ad
	gender   map[string]bool
	country  map[string]bool
	platform map[string]bool
}

type Cache struct {
	genderIndex   map[string](map[*CacheAd]bool)
	countryIndex  map[string](map[*CacheAd]bool)
	platformIndex map[string](map[*CacheAd]bool)
	ageIndex      []([]*CacheAd)
	ads           []*CacheAd
	lock          sync.RWMutex
}


