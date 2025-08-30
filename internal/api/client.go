package groupieapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	ArtistsLink   = "https://groupietrackers.herokuapp.com/api/artists"
	LocationsLink = "https://groupietrackers.herokuapp.com/api/locations"
	DatesLink     = "https://groupietrackers.herokuapp.com/api/dates"
	RelationsLink = "https://groupietrackers.herokuapp.com/api/relation"
)

func GetArtist(id int) (Artist, error) {
	artist, err := getFromAPI[Artist](id)
	if ArtistNotFound(artist) {
		return artist, ErrArtistNotFound
	}
	return artist, err
}

func GetDates(id int) (ArtistDates, error) {
	data, err := getFromAPI[ArtistDates](id)
	for i, date := range data.Dates {
		data.Dates[i] = strings.TrimLeft(date, "*")
	}
	return data, err
}

func GetRelation(id int) (Relation, error) {
	data, err := getFromAPI[Relation](id)
	newMap := make(map[string][]string, len(data.DatesLocations))
	for location, dates := range data.DatesLocations {
		newMap[normalizeLocation(location)] = dates
	}
	data.DatesLocations = newMap
	return data, err
}

func GetLocations(id int) (ArtistLocations, error) {
	data, err := getFromAPI[ArtistLocations](id)
	for i, loc := range data.Locations {
		data.Locations[i] = normalizeLocation(loc)
	}
	return data, err
}

func IndexArtists() ([]Artist, error) {
	return getFromAPI[[]Artist](0)
}

func getFromAPI[T HasAPI](id int) (result T, err error) {
	var link string // link goes here
	switch any(result).(type) {
	case Artist:
		link = ArtistsLink
	case ArtistLocations:
		link = LocationsLink
	case ArtistDates:
		link = DatesLink
	case Relation:
		link = RelationsLink
	case []Artist:
		link = ArtistsLink
	default:
		return result, ErrUnsupportedAPI
	}
	if _, ok := any(result).([]Artist); !ok {
		link += "/" + strconv.Itoa(id)
	}
	data, err := getData(link)
	if err != nil {
		return result, fmt.Errorf("error getting from api: %w", err)
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("error getting from api: %w", err)
	}
	return
}

func ArtistNotFound(artist Artist) bool {
	return artist.ID == 0 || artist.CreationDate == 0 || artist.FirstAlbum == "" || len(artist.Members) == 0 || artist.Name == ""
}

func getData(link string) (result []byte, err error) {
	response, err := http.Get(link)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	return result, nil
}

func normalizeLocation(location string) string {
	split := strings.Split(location, "-")
	if len(split) != 2 {
		return location
	}
	city, country := split[0], split[1]
	fix := func(s string) string {
		return strings.ReplaceAll(s, "_", " ")
	}
	city, country = fix(city), fix(country)
	city = capitalize(city)
	country = capitalize(country)
	switch strings.ToLower(country) {
	case "usa", "uk":
		country = strings.ToUpper(country)
	default:
		country = strings.Title(country)
	}
	return city + " - " + country
}

func capitalize(str string) string {
	split := strings.Split(str, " ")
	for i, s := range split {
		split[i] = strings.Title(s)
	}
	return strings.Join(split, " ")
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type ArtistLocations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type ArtistDates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type HasAPI interface {
	Artist | ArtistLocations | ArtistDates | Relation | []Artist
}

var (
	ErrUnsupportedAPI = errors.New("unsupported type of api call")
	ErrArtistNotFound = errors.New("artist not found")
)
