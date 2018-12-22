package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/polds/imgbase64"
	"io/ioutil"
	"log"
	"net/http"
)

type Place struct {
	PlaceId     	string  `json:"place_id"`
	Name     		string  `json:"name"`
	Geometry     	Geometry  `json:"geometry"`
	Photos     		[]Photo  `json:"photos"`
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
	Places []Place `json:"results"`
}

func main() {

	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/textsearch/json?key=AIzaSyDlRrMhhZXm-uhLM6XYAa4EWKdqgDSPPQk&query=mercadillo%20in%20spain")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	places := Places{}
	json.Unmarshal(body, &places)

	for  key, place := range places.Places {
		if len(place.Photos) <= 0 {
			continue
		}


		places.Places[key].Image = imgbase64.FromRemote("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&key=AIzaSyDlRrMhhZXm-uhLM6XYAa4EWKdqgDSPPQk&photoreference="  + place.Photos[0].Reference)
	}



}
