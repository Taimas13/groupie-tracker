package model

import (
	"encoding/json"
	"errors"
	"groupie-tracker/config"
	"net/http"
	"time"
)

func GetAllAtritst(Data *[]Artist) error {
	netClient := http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := netClient.Get(config.ArtistURL)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorText := "invlid status of the api's response; url: " + config.ArtistURL + "; status: " + resp.Status
		return errors.New(errorText)
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(&Data)

	if err != nil {
		return err
	}

	return nil
}

func GetArtist(id string, Data *Artist) error {
	netClient := http.Client{}

	resp, err := netClient.Get(config.ArtistURL + "/" + id)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorText := "invlid status of the api's response; url: " + config.ArtistURL + "/" + id + "; status: " + resp.Status
		return errors.New(errorText)
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(&Data)

	if err != nil {
		return err
	}

	return nil
}

func GetAllRelations(Data *map[string][]DateLocation) error {
	netClient := http.Client{}

	resp, err := netClient.Get(config.RelationURL)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorText := "invlid status of the api's response; url: " + config.RelationURL + "; status: " + resp.Status
		return errors.New(errorText)
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(&Data)

	if err != nil {
		return err
	}

	return nil
}

func GetRelation(id string, Data *Artist) error {
	netClient := http.Client{}

	resp, err := netClient.Get(config.RelationURL + "/" + id)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorText := "invlid status of the api's response; url: " + config.RelationURL + "/" + id + "; status: " + resp.Status
		return errors.New(errorText)
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(&Data)

	if err != nil {
		return err
	}

	return nil
}
