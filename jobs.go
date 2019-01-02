package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jcuerdo/mymarket-app-go/database"
	"github.com/jcuerdo/mymarket-app-go/model"
	"github.com/polds/imgbase64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Place struct {
	PlaceId     	string  	`json:"place_id"`
	Name     		string  	`json:"name"`
	Address     	string  	`json:"formatted_address"`
	Geometry     	Geometry  	`json:"geometry"`
	Photos     		[]Photo  	`json:"photos"`
	Image     		string
}

type Geometry struct {
	Location     Location  `json:"location"`
}

type Location struct {
	Lat     	float32  `json:"lat"`
	Lng     	float32  `json:"lng"`
}

type Photo struct {
	Reference     	string  `json:"photo_reference"`
}

type Places struct {
	Places  []Place `json:"results"`
	NextUrlToken  string `json:"next_page_token"`
}

const PLACES_URL = "https://maps.googleapis.com/maps/api/place/textsearch/json?key=AIzaSyDlRrMhhZXm-uhLM6XYAa4EWKdqgDSPPQk&query=%s"
const PLACES_URL_NEXT = "https://maps.googleapis.com/maps/api/place/textsearch/json?key=AIzaSyDlRrMhhZXm-uhLM6XYAa4EWKdqgDSPPQk&pagetoken=%s"
const PHOTO_URL = "https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&key=AIzaSyDlRrMhhZXm-uhLM6XYAa4EWKdqgDSPPQk&photoreference=%s"

var maxImports = 500

func main() {

	searches := []string{
		"mercadillo%20en%20espa%C3%B1a",
		"mercadillo%20en%20barcelona",
		"mercadillo%20en%20madrid",
		"mercadillo%20en%asturias",
		"mercadillo%20in%20spain",
		"mercadillo%20in%20asturias",
		"mercadillo%20in%20barcelona",
		"mercadillo%20in%asturias",
		"mercado%20en%20espa%C3%B1a",
		"mercado%20en%20barcelona",
		"mercado%20en%20madrid",
		"mercado%20en%asturias",
		"mercado%20in%20spain",
		"mercado%20in%20barcelona",
		"mercado%20in%20madrid",
		"mercado%20in%20asturias",
		"mercadillo%20solidario%en%españa",
		"mercadillo%20solidario%en%barcelona",
		"mercadillo%20solidario%en%madrid",
		"mercadillo%20solidario%en%asturias",
		"mercadillo%20benefico%en%españa",
		"mercadillo%20benefico%en%barcelona",
		"mercadillo%20benefico%en%madrid",
		"mercadillo%20benefico%en%asturias",
	}

	wg := sync.WaitGroup{}
	for _ , search := range searches {
		search := url.PathEscape(search)
		fmt.Println("Searching url: " + fmt.Sprintf(PLACES_URL, search))
		importPlaces(fmt.Sprintf(PLACES_URL, search), &wg)
	}
	wg.Wait()

	fmt.Printf("Maximum is now %d", maxImports)

}

func importPlaces(url string, wg *sync.WaitGroup) {
	places := getPlacesFromUrl(url)
	importMorePlaces(fmt.Sprintf(PLACES_URL_NEXT, places.NextUrlToken), wg)
	storePlaces(places)
}

func importMorePlaces(url string, wg *sync.WaitGroup) {
	if maxImports > 0 {
		maxImports--
	} else {
		return
	}

	wg.Add(1)
	places := getPlacesFromUrl(url)
	importMorePlaces(fmt.Sprintf(PLACES_URL_NEXT, places.NextUrlToken), wg)
	storePlaces(places)
	wg.Done()
}

func storePlaces(places Places) {
	for _, place := range places.Places {

		if (existsPlace(place)) {
			fmt.Printf("Place %s already exists \n", place.PlaceId)
			continue
		}

		savePlace(place)
	}
}

func getPlacesFromUrl(url string) Places {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	places := Places{}
	json.Unmarshal(body, &places)
	for key, place := range places.Places {
		if len(place.Photos) <= 0 {
			continue
		}
		places.Places[key].Image = imgbase64.FromRemote(fmt.Sprintf(PHOTO_URL, place.Photos[0].Reference))
	}
	return places
}

func savePlace(place Place) {
	marketRepository := database.GetMarketRepository()
	photoRepository := database.GetPhotoRepository()
	market := model.MarketExportable{
		UserId:        2,
		Name:          place.Name,
		Description:   place.Name + " " + place.Address,
		Type:          "PUBLIC",
		Place:         "PUBLIC",
		Flexible:      true,
		Date:          time.Now().Format("2006-01-02 15:04:05"),
		Lat:           place.Geometry.Location.Lat,
		Lon:           place.Geometry.Location.Lng,
		GooglePlaceId: place.PlaceId,
	}
	fmt.Println("Import market")
	fmt.Println(market)
	id := marketRepository.Create(market)
	if id > 0 {
		fmt.Printf("Market imported with id %d \n", id)

		if place.Image != "" {
			photo := model.Photo{
				Content: place.Image,
			}
			fmt.Println("Import photo")
			created := photoRepository.Create(photo, id)
			if created {
				fmt.Printf("Photo imported \n")
			} else {
				fmt.Println("Fail importing market")
			}
		}

	} else {
		fmt.Println("Fail importing market")
	}

}
func existsPlace(place Place) (bool) {
	marketRepository := database.GetMarketRepository()

	return marketRepository.ExistsGooglePlaceId(place.PlaceId)

}
