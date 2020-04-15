package commands

// Weather is a subcommand that demonstrates using REST calls to capture informaiton which
// can be output using the formatting tools. This uses freely-available data from openweather.com
//

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"

	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/profile"
	"github.com/tucats/gopackages/app-cli/tables"
	"github.com/tucats/gopackages/app-cli/ui"
)

var stateNames = map[string]string{
	"nc": "north carolina",
	"va": "virginia",
	"sc": "south carolina",
	"dc": "district of columbia",
}

// LatLong defines the coordinates of a report.
type LatLong struct {
	Longitude float32 `json:"lon"`
	Lattitude float32 `json:"lat"`
}

// WeatherText contains human-readable descriptions of the conditions.
type WeatherText struct {
	Description string `json:"description"`
}

// WeatherWind summarizes wind conditions
type WeatherWind struct {
	Speed     float64 `json:"speed"`
	Direction float64 `json:"deg"`
}

// WeatherOverview summarizes temperature.
type WeatherOverview struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Minimum   float64 `json:"temp_min"`
	Maximum   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

// Weather contains a weather report.
type Weather struct {
	Coord LatLong         `json:"coord"`
	Text  []WeatherText   `json:"weather"`
	Main  WeatherOverview `json:"main"`
	Wind  WeatherWind     `json:"wind"`
}

// WeatherGrammar defines the subgrammar of the weather command.
var WeatherGrammar = []cli.Option{
	cli.Option{
		LongName:    "location",
		Description: "The location (city, state) for which the weather is displayed",
		OptionType:  cli.StringListType,
	},
}

// WeatherAction is the command line action for handling the weather subcommand.
func WeatherAction(c *cli.Context) error {

	var city string
	var state string

	location, found := c.GetStringList("location")

	if !found {
		city = profile.Get("weather-city")
		state = profile.Get("weather-state")
	} else {
		if len(location) != 2 {
			return errors.New("Missing location value")
		}

		city = strings.ToLower(location[0])
		state = strings.ToLower(location[1])

		if longName, found := stateNames[state]; found {
			state = longName
		}

		profile.Set("weather-city", city)
		profile.Set("weather-state", state)
	}

	keyValue := profile.Get("weather-api-key")
	if keyValue == "" {
		// Get your own darn key
		keyValue = "fbd457b51b56eddf1644edefd591f89c"
		profile.Set("weather-api-key", keyValue)
	}

	parms := url.QueryEscape(city+","+state) + "&appid=" + keyValue + "&units=imperial"
	url := "https://api.openweathermap.org/data/2.5/weather?q=" + parms

	ui.Debug("URL: %s", url)

	response, err := http.Get(url)

	if err != nil {
		return err
	}

	weather := Weather{}
	data, _ := ioutil.ReadAll(response.Body)
	ui.Debug("Reply: %s", string(data))

	err = json.Unmarshal(data, &weather)
	if err != nil {
		return err
	}

	t := tables.New([]string{"Item", "Value"})
	t.AddRowItems("Summary", weather.Text[0].Description)
	t.AddRowItems("Temperature", weather.Main.Temp)
	t.AddRowItems("  Feels Like", weather.Main.FeelsLike)
	t.AddRowItems("  Minimum", weather.Main.Minimum)
	t.AddRowItems("  Maximum", weather.Main.Maximum)
	t.AddRowItems("Wind Speed", weather.Wind.Speed)

	if weather.Wind.Direction > 0.0 {
		d := int((math.Round(weather.Wind.Direction / 22.5)))
		windDirections := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW", "N"}
		t.AddRowItems("Wind Direction", windDirections[d])
	}

	t.AddRowItems("Pressure", weather.Main.Pressure)
	t.AddRowItems("Humidity", weather.Main.Humidity)

	t.Print(profile.Get("output-format"))

	return nil
}
