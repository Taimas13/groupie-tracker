package model

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (artist *Artist) ChangeKey() {
	newMap := make(map[string][]string)
	for key, value := range artist.DatesLocations {
		newKey := strings.ToUpper(key)
		reg := regexp.MustCompile(`-`)
		newKey = reg.ReplaceAllString(newKey, "$1 - $2")
		reg = regexp.MustCompile(`_`)
		newKey = reg.ReplaceAllString(newKey, "$1-$2")
		newMap[newKey] = value
	}
	artist.DatesLocations = newMap
}

func (artistsList *ArtistsList) ChangeKeys() {
	for i := 0; i < len(artistsList.List); i++ {
		artistsList.List[i].ChangeKey()
	}
}

func contains(s []string, e int) bool {
	for _, a := range s {
		b, _ := strconv.Atoi(a)
		if b == e {
			return true
		}
	}
	return false
}

func (artistsList *ArtistsList) AddDatesLocations(inputList []DateLocation) {
	artistsList.FilterRanges.MinCreationDate = artistsList.List[0].CreationDate
	artistsList.FilterRanges.MaxCreationDate = 0
	artistsList.FilterRanges.MinFirstAlbumDate = time.Now().Format("yyyy-mm-dd")
	artistsList.FilterRanges.MaxFirstAlbumDate = "0000-00-00"

	for i := 0; i < len(artistsList.List); i++ {
		if !contains(artistsList.FilterRanges.MembersNums, len(artistsList.List[i].Members)) {
			artistsList.FilterRanges.MembersNums = append(artistsList.FilterRanges.MembersNums, strconv.Itoa(len(artistsList.List[i].Members)))
		}
		if artistsList.List[i].CreationDate > artistsList.FilterRanges.MaxCreationDate {
			artistsList.FilterRanges.MaxCreationDate = artistsList.List[i].CreationDate
		}
		if artistsList.List[i].CreationDate < artistsList.FilterRanges.MinCreationDate {
			artistsList.FilterRanges.MinCreationDate = artistsList.List[i].CreationDate
		}
		firstAlbumDate := artistsList.List[i].FirstAlbum[6:] + "-" + artistsList.List[i].FirstAlbum[3:5] + "-" + artistsList.List[i].FirstAlbum[:2]
		if firstAlbumDate > artistsList.FilterRanges.MaxFirstAlbumDate {
			artistsList.FilterRanges.MaxFirstAlbumDate = firstAlbumDate
		}
		if firstAlbumDate < artistsList.FilterRanges.MinFirstAlbumDate {
			artistsList.FilterRanges.MinFirstAlbumDate = firstAlbumDate
		}
		for j := 0; j < len(inputList); j++ {
			if artistsList.List[i].Id == inputList[j].Id {
				artistsList.List[i].DatesLocations = inputList[j].DatesLocations
				break
			}
		}
	}
	sort.Strings(artistsList.FilterRanges.MembersNums)
}

func (artist *Artist) ContainsSearchInput(input string) bool {
	if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(input)) {
		return true
	} else if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(input)) {
		return true
	} else if strings.Contains(strconv.Itoa(artist.CreationDate), input) {
		return true
	}

	for i := 0; i < len(artist.Members); i++ {
		if strings.Contains(strings.ToLower(artist.Members[i]), strings.ToLower(input)) {
			return true
		}
	}

	for key, value := range artist.DatesLocations {
		if strings.Contains(strings.ToLower(key), strings.ToLower(input)) {
			return true
		}
		for i := 0; i < len(value); i++ {
			if strings.Contains(strings.ToLower(value[i]), strings.ToLower(input)) {
				return true
			}
		}
	}

	return false
}

func (searchResultArtists *ArtistsList) SearchInputInArtistsList(artists []Artist, input string) {
	for i := 0; i < len(artists); i++ {
		if artists[i].ContainsSearchInput(input) {
			searchResultArtists.List = append(searchResultArtists.List, artists[i])
		}
	}
}

func (searchResultArtists *ArtistsList) FilterInputInArtistsList(artists []Artist, filterInput FilterRanges) {
	if filterInput.Location == "" && filterInput.MaxCreationDate == 0 && filterInput.MinCreationDate == 0 && filterInput.MaxFirstAlbumDate == "" && filterInput.MinFirstAlbumDate == "" && len(filterInput.MembersNums) == 0 {
		searchResultArtists.List = artists
		return
	}

	for i := 0; i < len(artists); i++ {
		if artists[i].ContainsFilterInput(filterInput) {
			searchResultArtists.List = append(searchResultArtists.List, artists[i])
		}
	}
}

func (artist *Artist) ContainsFilterInput(filterInput FilterRanges) bool {
	if filterInput.MinCreationDate != 0 && filterInput.MinCreationDate > artist.CreationDate {
		return false
	}
	if filterInput.MaxCreationDate != 0 && filterInput.MaxCreationDate < artist.CreationDate {
		return false
	}

	firstAlbumDate := artist.FirstAlbum[6:] + "-" + artist.FirstAlbum[3:5] + "-" + artist.FirstAlbum[:2]
	if filterInput.MinFirstAlbumDate != "" && filterInput.MinFirstAlbumDate > firstAlbumDate {
		return false
	}
	if filterInput.MaxFirstAlbumDate != "" && filterInput.MaxFirstAlbumDate < firstAlbumDate {
		return false
	}
	if len(filterInput.MembersNums) > 0 {
		contains := false
		for _, value := range filterInput.MembersNums {
			if intValue, _ := strconv.Atoi(value); intValue == len(artist.Members) {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}

	if filterInput.Location != "" {
		contains := false
		for key := range artist.DatesLocations {
			if strings.Contains(strings.ToLower(key), strings.ToLower(filterInput.Location)) {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}
	return true
}
