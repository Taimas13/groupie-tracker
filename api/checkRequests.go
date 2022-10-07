package api

import (
	"groupie-tracker/model"
	"net/http"
	"strconv"
	"strings"
)

func checkIndexRequest(r *http.Request) int {
	if r.URL.Path != "/" {
		return http.StatusNotFound
	} else if r.Method != http.MethodGet && r.Method != http.MethodOptions && r.Method != http.MethodHead {
		return http.StatusMethodNotAllowed
	}
	return 0
}

func checkArtistRequest(r *http.Request) int {
	id := strings.TrimPrefix(r.URL.Path, "/artists/")
	if _, err := strconv.Atoi(id); err != nil {
		return http.StatusNotFound
	}

	if r.Method != http.MethodGet && r.Method != http.MethodOptions && r.Method != http.MethodHead {
		return http.StatusMethodNotAllowed
	}
	return 0
}

func checkSearchRequest(r *http.Request) int {
	if r.Method != http.MethodGet && r.Method != http.MethodOptions && r.Method != http.MethodHead {
		return http.StatusMethodNotAllowed
	}
	searchinput := r.URL.Query().Get("searchinput")
	searchinput = strings.Replace(searchinput, " ", "", -1)
	searchinput = strings.Replace(searchinput, "\n", "", -1)

	if searchinput == "" {
		return http.StatusNotFound
	}
	return 0
}

func checkFilterRequest(r *http.Request, filterInput *model.FilterRanges) int {
	if r.Method != http.MethodGet && r.Method != http.MethodOptions && r.Method != http.MethodHead {
		return http.StatusMethodNotAllowed
	}
	var err error
	if r.URL.Query().Get("creationDateFrom") == "" {
		filterInput.MinCreationDate = 0
	} else {
		filterInput.MinCreationDate, err = strconv.Atoi(r.URL.Query().Get("creationDateFrom"))
	}
	if err != nil {
		return http.StatusNotFound
	}

	if r.URL.Query().Get("creationDateTo") == "" {
		filterInput.MaxCreationDate = 0
	} else {
		filterInput.MaxCreationDate, err = strconv.Atoi(r.URL.Query().Get("creationDateTo"))
	}
	if err != nil {
		return http.StatusNotFound
	}
	filterInput.MinFirstAlbumDate = r.URL.Query().Get("firstAlbumDateFrom")
	filterInput.MaxFirstAlbumDate = r.URL.Query().Get("firstAlbumDateTo")

	if r.URL.Query().Get("numberOfMembers") == "" {
		filterInput.MembersNums = nil
	} else {
		filterInput.MembersNums = r.URL.Query()["numberOfMembers"]
	}
	for _, value := range filterInput.MembersNums {
		if _, err = strconv.Atoi(value); err != nil {
			return http.StatusNotFound
		}
	}

	filterInput.Location = r.URL.Query().Get("location")

	return 0
}
