package api

import (
	"groupie-tracker/config"
	"groupie-tracker/model"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
)

func AppMux() http.Handler {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("ui"))
	mux.Handle("/ui/", http.StripPrefix("/ui", fs))

	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("/artists/", artistPageHandler)
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/filter/", filterHandler)
	return mux
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	artists := new(model.ArtistsList)
	relations := new(model.RelationList)
	if err := getAllArtistsAndRelations(artists, relations); err != nil {
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		executeErrorPage(w, http.StatusInternalServerError)
		return
	}

	artists.AddDatesLocations(relations.List["index"])
	relations = nil
	searchResultArtists := new(model.ArtistsList)
	artists.ChangeKeys()
	searchResultArtists.SeachDataList = artists.List
	searchResultArtists.FilterRanges = artists.FilterRanges

	filterInput := new(model.FilterRanges)
	if errStatus := checkFilterRequest(r, filterInput); errStatus != 0 {
		executeSearchErrorPage(w, errStatus, *searchResultArtists)
		errorLog.Printf(http.StatusText(errStatus))
		return
	}

	searchResultArtists.FilterInputInArtistsList(artists.List, *filterInput)
	if len(searchResultArtists.List) == 0 {
		executeSearchErrorPage(w, http.StatusNotFound, *searchResultArtists)
		errorLog.Printf(http.StatusText(http.StatusNotFound))
		return
	}

	if tmpl, err := template.ParseFiles(config.IndexTmplPath); err != nil {
		executeErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	} else {
		if err := tmpl.Execute(w, searchResultArtists); err != nil {
			executeErrorPage(w, http.StatusInternalServerError)
			errorLog.Printf(err.Error())
		}
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	searchinput := r.URL.Query().Get("searchinput")

	artists := new(model.ArtistsList)
	relations := new(model.RelationList)
	if err := getAllArtistsAndRelations(artists, relations); err != nil {
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		executeErrorPage(w, http.StatusInternalServerError)
		return
	}

	artists.AddDatesLocations(relations.List["index"])
	relations = nil

	searchResultArtists := new(model.ArtistsList)
	artists.ChangeKeys()
	searchResultArtists.SeachDataList = artists.List

	if errStatus := checkSearchRequest(r); errStatus != 0 {
		executeSearchErrorPage(w, errStatus, *searchResultArtists)
		errorLog.Printf(http.StatusText(errStatus))
		return
	}
	searchResultArtists.SearchInputInArtistsList(artists.List, searchinput)
	if len(searchResultArtists.List) == 0 {
		executeSearchErrorPage(w, http.StatusNotFound, *searchResultArtists)
		errorLog.Printf(http.StatusText(http.StatusNotFound))
		return
	}

	if tmpl, err := template.ParseFiles(config.IndexTmplPath); err != nil {
		executeErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	} else {
		if err := tmpl.Execute(w, searchResultArtists); err != nil {
			executeErrorPage(w, http.StatusInternalServerError)
			errorLog.Printf(err.Error())
		}
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	if errStatus := checkIndexRequest(r); errStatus != 0 {
		executeErrorPage(w, errStatus)
		errorLog.Printf(http.StatusText(errStatus))
		return
	}

	artists := new(model.ArtistsList)
	relations := new(model.RelationList)
	if err := getAllArtistsAndRelations(artists, relations); err != nil {
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		executeErrorPage(w, http.StatusInternalServerError)
		return
	}

	artists.AddDatesLocations(relations.List["index"])
	relations = nil
	artists.ChangeKeys()
	artists.SeachDataList = artists.List

	if tmpl, err := template.ParseFiles(config.IndexTmplPath); err != nil {
		executeErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	} else {
		if err := tmpl.Execute(w, artists); err != nil {
			executeErrorPage(w, http.StatusInternalServerError)
			errorLog.Printf(err.Error())
		}
	}
}

func artistPageHandler(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if errStatus := checkArtistRequest(r); errStatus != 0 {
		executeErrorPage(w, errStatus)
		errorLog.Printf(http.StatusText(errStatus))
		return
	}
	artId := strings.TrimPrefix(r.URL.Path, "/artists/")
	artist := new(model.Artist)

	if err := getArtistAndRelation(artist, artId); err != nil {
		executeErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	}

	if artist.Id == 0 {
		executeErrorPage(w, http.StatusNotFound)
		errorLog.Printf(http.StatusText(http.StatusNotFound))
		return
	}

	artist.ChangeKey()

	if tmpl, err := template.ParseFiles(config.ArtistTmplPath); err != nil {
		executeErrorPage(w, http.StatusInternalServerError)
		errorLog.Printf(http.StatusText(http.StatusInternalServerError), err)
		return
	} else {
		if err := tmpl.Execute(w, artist); err != nil {
			executeErrorPage(w, http.StatusInternalServerError)
			errorLog.Printf(err.Error())
		}
	}
}

func executeErrorPage(w http.ResponseWriter, errStatus int) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tmpl, err := template.ParseFiles(config.ErrorTmplPath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		errorLog.Printf(err.Error())
		return
	}
	w.WriteHeader(errStatus)
	if err := tmpl.Execute(w, http.StatusText(errStatus)); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		errorLog.Printf(err.Error())
		return
	}
}

func executeSearchErrorPage(w http.ResponseWriter, errStatus int, artists model.ArtistsList) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	artists.ErrorText = http.StatusText(errStatus)
	tmpl, err := template.ParseFiles(config.SearchErrorTmplPath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		errorLog.Printf(err.Error())
		return
	}
	w.WriteHeader(errStatus)
	if err := tmpl.Execute(w, artists); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		errorLog.Printf(err.Error())
		return
	}
}

func getAllArtistsAndRelations(artists *model.ArtistsList, relations *model.RelationList) error {
	var wg sync.WaitGroup
	var errGetArtists error
	var errGetRelations error
	wg.Add(2)
	go func() {
		defer wg.Done()
		errGetArtists = model.GetAllAtritst(&artists.List)
	}()

	go func() {
		defer wg.Done()
		errGetRelations = model.GetAllRelations(&relations.List)
	}()

	wg.Wait()
	if errGetArtists != nil {
		return errGetArtists
	}
	if errGetRelations != nil {
		return errGetRelations
	}
	return nil
}

func getArtistAndRelation(artist *model.Artist, artId string) error {
	var wg sync.WaitGroup
	var errGetArtist error
	var errGetRelation error
	wg.Add(2)
	go func() {
		defer wg.Done()
		errGetArtist = model.GetArtist(artId, artist)
	}()

	go func() {
		defer wg.Done()
		errGetRelation = model.GetRelation(artId, artist)
	}()

	wg.Wait()
	if errGetArtist != nil {
		return errGetArtist
	}
	if errGetRelation != nil {
		return errGetRelation
	}
	return nil
}
