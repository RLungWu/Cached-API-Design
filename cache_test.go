package main

import (
	"testing"
	"time"
)

func setupCache(cache *Cache) {
	int20 := 20
	int30 := 30
	int40 := 40
	ad1 := Ad{
		Title:   "Ad 1",
		StartAt: time.Now().Add(-time.Hour * 24),
		EndAt:   time.Now().Add(time.Hour * 24),
		Conditions: AdConditions{
			Gender:  &[]string{"M"},
			Country: &[]string{"TW", "JP"},
		},
	}
	ad2 := Ad{
		Title:   "Ad 2",
		StartAt: time.Now().Add(-time.Hour * 24),
		EndAt:   time.Now().Add(time.Hour * 24),
		Conditions: AdConditions{
			AgeStart: &int20,
			AgeEnd:   &int30,
			Gender:   &[]string{"F", "M"},
		},
	}
	ad3 := Ad{
		Title:   "Ad 3",
		StartAt: time.Now().Add(-time.Hour * 24),
		EndAt:   time.Now().Add(time.Hour * 24),
		Conditions: AdConditions{
			AgeStart: &int30,
			AgeEnd:   &int40,
			Country:  &[]string{"TW"},
			Platform: &[]string{"android", "ios"},
		},
	}
	ad4 := Ad{
		Title:   "Ad 4",
		StartAt: time.Now().Add(-time.Hour * 24),
		EndAt:   time.Now().Add(time.Hour * 24),
		Conditions: AdConditions{
			AgeStart: &int20,
			AgeEnd:   &int40,
			Country:  &[]string{"JP"},
			Platform: &[]string{"android", "ios", "web"},
		},
	}

	cache.Update([]Ad{ad1, ad2, ad3, ad4})
}

func TestCacheFilter(t *testing.T) {
	cache := Cache{}
	setupCache(&cache)

	testCases := []struct {
		name            string
		query           AdQuery
		expectedResults int
		expectedTitles  []string
	}{
		{"Filter by Limit", AdQuery{Offset: 0, Limit: 3}, 3, []string{"Ad 1", "Ad 2", "Ad 3"}},
		{"Filter by Offset", AdQuery{Offset: 1, Limit: 3}, 3, []string{"Ad 2", "Ad 3", "Ad 4"}},
		{"Filter by out of range Offset", AdQuery{Offset: 5, Limit: 3}, 0, nil},
		{"Filter by Age", AdQuery{Offset: 0, Limit: 3, Age: 25}, 2, []string{"Ad 2", "Ad 4"}},
		{"Filter by Country", AdQuery{Offset: 0, Limit: 3, Country: "TW"}, 2, []string{"Ad 1", "Ad 3"}},
		{"Filter by Platform", AdQuery{Offset: 0, Limit: 3, Platform: "android"}, 2, []string{"Ad 3", "Ad 4"}},
		{"Filter by Country and Gender", AdQuery{Offset: 0, Limit: 3, Country: "TW", Gender: "M"}, 1, []string{"Ad 1"}},
		{"Filter by Country and Platform", AdQuery{Offset: 0, Limit: 3, Country: "JP", Platform: "ios"}, 1, []string{"Ad 4"}},
		{"Filter by Age and Platform", AdQuery{Offset: 0, Limit: 3, Age: 30, Platform: "web"}, 1, []string{"Ad 4"}},
		{"Filter by Age, Country, and Platform", AdQuery{Offset: 0, Limit: 3, Age: 30, Country: "TW", Platform: "web"}, 0, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := cache.Filter(tc.query)
			if len(results) != tc.expectedResults {
				t.Errorf("Expected %d results, but got %d", tc.expectedResults, len(results))
			}
			for i, title := range tc.expectedTitles {
				if results[i].Title != title {
					t.Errorf("Expected title %s, but got %s", title, results[i].Title)
				}
			}
		})
	}
}
