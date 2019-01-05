package config

type Config struct {
	Datasource 		string
	LogFile			string
	ApiKey			string
	MaxHttpRequests		int
	PlacesQuerys	[]string
}