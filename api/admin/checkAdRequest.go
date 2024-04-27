package admin

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RLungWu/Dcard-Backend-HW/api/ad"
	"github.com/biter777/countries"
)

func checkAdRequest(ad *ad.AdRequest) error {
	if err := checkTitle(ad.Title); err != nil {
		return fmt.Errorf("title is incorrect with : %w", err)
	}

	if err := checkStartAndEnd(ad.StartAt, ad.EndAt); err != nil {
		return fmt.Errorf("start_at and end_at is incorrect with : %w", err)
	}

	if err := checkConditions(ad.Conditions); err != nil {
		return fmt.Errorf("conditions is incorrect with : %w", err)

	}

	return nil
}

func checkTitle(title string) error {
	if title == "" {
		return errors.New("title can't be null")
	}

	if !strings.HasPrefix(title, "Ad") {
		return errors.New("title must start with 'Ad'")
	}

	return nil
}

func checkStartAndEnd(startAt, endAt time.Time) error {
	if startAt.After(endAt) {
		return errors.New("start_at must be before end_at")
	}

	return nil
}

func checkConditions(conditions ad.Contition) error {
	if conditions.AgeStart != nil {
		if *conditions.AgeStart < 1 {
			return errors.New("ageStart must be greater than 0")
		}
		if *conditions.AgeStart > 100 {
			return errors.New("ageStart must be less than 100")
		}
		if *conditions.AgeStart > *conditions.AgeEnd {
			return errors.New("ageStart must be less than ageEnd")
		}
	}

	if conditions.AgeEnd != nil {
		if *conditions.AgeEnd < 1 {
			return errors.New("ageEnd must be greater than 0")
		}
		if *conditions.AgeEnd > 100 {
			return errors.New("ageEnd must be less than 100")
		}
	}

	if conditions.Gender != nil {
		for _, value := range conditions.Gender {
			if value != "M" && value != "F" {
				return errors.New("gender must be M or F")
			}
		}
	}

	if conditions.Country != nil {
		for _, value := range conditions.Country {
			country := countries.ByName(value)
			if country == countries.Unknown {
				return errors.New("country is invalid")
			}
		}
	}

	if conditions.Platform != nil {
		for _, value := range conditions.Platform{
			if value != "iOS" && value != "Android" {
				return errors.New("platform must be iOS or Android")
			}
		}
	}

	return nil
}
