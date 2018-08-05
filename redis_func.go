package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/SUN-XIN/know-your-friends/geo"
	"github.com/SUN-XIN/know-your-friends/types"
)

// return all significant places of the given user
// FAKE! TODO: get from redis DB
func GetPlacesByID(userID string) []*types.SignificantPlace {
	// fake
	// in home: 48.823305, 2.361281
	// in school: 48.847016, 2.355808
	// in work: 48.854095, 2.373171
	switch {
	case strings.Contains(userID, "places"):
		return []*types.SignificantPlace{
			&types.SignificantPlace{
				UserID: userID,
				ID:     fmt.Sprintf("%d", time.Now().Unix()),
				LLBox: geo.LLBox{
					N: 48.823534,
					E: 2.361528,
					S: 48.823058,
					W: 2.360991,
				},
				Name: types.PLACE_NAME_HOME,
			},
			&types.SignificantPlace{
				UserID: userID,
				ID:     fmt.Sprintf("%d", time.Now().Unix()),
				LLBox: geo.LLBox{
					N: 48.849049,
					E: 2.357621,
					S: 48.844919,
					W: 2.355068,
				},
				Name: types.PLACE_NAME_SCHOOL,
			},
			&types.SignificantPlace{
				UserID: userID,
				ID:     fmt.Sprintf("%d", time.Now().Unix()),
				LLBox: geo.LLBox{
					N: 48.854147,
					E: 2.373185,
					S: 48.853667,
					W: 2.372981,
				},
				Name: types.PLACE_NAME_WORK,
			},
		}
	case strings.Contains(userID, "home"):
		return []*types.SignificantPlace{
			&types.SignificantPlace{
				UserID: userID,
				ID:     fmt.Sprintf("%d", time.Now().Unix()),
				LLBox:  geo.LLBox{},
				Name:   types.PLACE_NAME_HOME,
			},
		}
	default:
		return []*types.SignificantPlace{
			/*
				&types.SignificantPlace{
					UserID: userID,
					ID:     fmt.Sprintf("%d", time.Now().Unix()),
					LLBox: geo.LLBox{
						N: 49.823473,
						E: 3.361238,
						S: 49.822901,
						W: 3.360728,
					},
					Name: types.PLACE_NAME_HOME,
				},
			*/
			&types.SignificantPlace{
				UserID: userID,
				ID:     fmt.Sprintf("%d", time.Now().Unix()),
				LLBox: geo.LLBox{
					N: 49.849049,
					E: 3.357621,
					S: 49.844919,
					W: 3.355068,
				},
				Name: types.PLACE_NAME_SCHOOL,
			},
			&types.SignificantPlace{
				UserID: userID,
				ID:     fmt.Sprintf("%d", time.Now().Unix()),
				LLBox: geo.LLBox{
					N: 49.854147,
					E: 3.373185,
					S: 49.853667,
					W: 3.372981,
				},
				Name: types.PLACE_NAME_WORK,
			},
		}
	}
}
