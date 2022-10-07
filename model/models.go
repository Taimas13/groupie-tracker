package model

type Artist struct {
	Id             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistsList struct {
	List          []Artist
	SeachDataList []Artist
	ErrorText     string
	FilterRanges  FilterRanges
}

type FilterRanges struct {
	MinCreationDate   int
	MaxCreationDate   int
	MinFirstAlbumDate string
	MaxFirstAlbumDate string
	Location          string
	MembersNums       []string
}

type RelationList struct {
	List map[string][]DateLocation `json:"index"`
}

type DateLocation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
