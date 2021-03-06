// Package sbweather displays the current weather in the provided zip code.
// Currently supported for US only.
package sbweather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var colorEnd = "^d^"

// Routine is the main object for this package.
type Routine struct {
	// Error encountered along the way, if any.
	err error

	// HTTP client to reuse for all requests out.
	client http.Client

	// User-supplied zip code for where the temperature should reflect.
	zip string

	// NWS-provided URL for getting the temperature, as found during the init.
	url string

	// Current temperature for the provided zip code.
	temp int

	// Forecast high.
	high int

	// Forecast low.
	low int

	// Trio of user-provided colors for displaying various states.
	colors struct {
		normal  string
		warning string
		error   string
	}
}

// New runs a sanity check on the provided zip code and makes a new routine object. zip is the zip code to use for
// localizing the weather. colors is an optional triplet of hex color codes for colorizing the output based on these
// rules:
//   1. Normal color, used for printing the current temperature and forecast.
//   2. Warning color, currently unused.
//   3. Error color, used for error messages.
func New(zip string, colors ...[3]string) *Routine {
	var r Routine

	if len(zip) != 5 {
		r.err = fmt.Errorf("invalid zip code length")
		return &r
	}

	_, err := strconv.Atoi(zip)
	if err != nil {
		r.err = err
		return &r
	}
	r.zip = zip

	// Do a minor sanity check on the color codes.
	if len(colors) == 1 {
		for _, color := range colors[0] {
			if !strings.HasPrefix(color, "#") || len(color) != 7 {
				r.err = fmt.Errorf("invalid color")
				return &r
			}
		}
		r.colors.normal = "^c" + colors[0][0] + "^"
		r.colors.warning = "^c" + colors[0][1] + "^"
		r.colors.error = "^c" + colors[0][2] + "^"
	} else {
		// If a color array wasn't passed in, then we don't want to print this.
		colorEnd = ""
	}

	return &r
}

// Update gets the current hourly temperature.
func (r *Routine) Update() {
	r.err = nil

	// If this is the first run of the session, initialize the routine.
	if r.url == "" {
		// Get coordinates.
		lat, long, err := getCoords(r.client, r.zip)
		if err != nil {
			r.err = fmt.Errorf("no coordinates")
			return
		}

		// Get forecast URL.
		url, err := getURL(r.client, lat, long)
		if err != nil {
			r.err = fmt.Errorf("no forecast data")
			return
		}
		r.url = url
	}

	// Get hourly temperature.
	temp, err := getTemp(r.client, r.url+"/hourly")
	if err != nil {
		r.err = err
		return
	}
	r.temp = temp

	high, low, err := getForecast(r.client, r.url)
	if err != nil {
		r.err = err
		return
	}
	r.high = high
	r.low = low
}

// String formats and prints the current temperature.
func (r *Routine) String() string {
	var s string

	if r.err != nil {
		return r.colors.error + r.err.Error() + colorEnd
	}

	t := time.Now()
	if t.Hour() < 15 {
		s = "today"
	} else {
		s = "tom"
	}

	return fmt.Sprintf("%s%v °F (%s: %v/%v)%s", r.colors.normal, r.temp, s, r.high, r.low, colorEnd)
}

// getCoords gets the geographic coordinates for the provided zip code. It should receive a response in this format:
// {"status":1,"output":[{"zip":"90210","latitude":"34.103131","longitude":"-118.416253"}]}
func getCoords(client http.Client, zip string) (string, string, error) {
	type coords struct {
		Status int                 `json:"status"`
		Output []map[string]string `json:"output"`
	}

	url := "https://api.promaptools.com/service/us/zip-lat-lng/get/?zip=" + zip + "&key=17o8dysaCDrgv1c"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	c := coords{}
	err = json.Unmarshal(body, &c)
	if err != nil {
		return "", "", err
	}

	// Make sure the status is good.
	if c.Status != 1 {
		return "", "", fmt.Errorf("coordinates request failed")
	}

	// Make sure we got back just one dictionary.
	if len(c.Output) != 1 {
		return "", "", fmt.Errorf("received invalid coordinates array")
	}

	// TODO: reduce to 4 decimal points of precision
	lat := c.Output[0]["latitude"]
	long := c.Output[0]["longitude"]
	if lat == "" || long == "" {
		return "", "", fmt.Errorf("missing coordinates in response")
	}

	return lat, long, nil
}

// getURL queries the NWS to determine which URL we should be using for getting the weather forecast.
// Our value should be here: properties -> forecast.
func getURL(client http.Client, lat string, long string) (string, error) {
	type props struct {
		// Properties map[string]interface{} `json:"properties"`
		Properties struct {
			Forecast string `json: "temperature"`
		} `json:"properties"`
	}

	url := "https://api.weather.gov/points/" + lat + "," + long
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	p := props{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		return "", err
	}

	url = p.Properties.Forecast
	if url == "" {
		return "", fmt.Errorf("missing temperature URL")
	}

	return url, nil
}

// getTemp gets the current temperature from the NWS database.
// Our value should be here: properties -> periods -> (latest period) -> temperature.
func getTemp(client http.Client, url string) (int, error) {
	type temp struct {
		Properties struct {
			Periods []interface{} `json:"periods"`
		} `json:"properties"`
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, fmt.Errorf("temp: bad request")
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, fmt.Errorf("temp: bad client")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("temp: bad read")
	}

	t := temp{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		return -1, fmt.Errorf("temp: bad JSON")
	}

	// Get the list of weather readings.
	periods := t.Properties.Periods
	if len(periods) == 0 {
		return -1, fmt.Errorf("missing hourly temperature periods")
	}

	// Use the most recent reading.
	latest := periods[0].(map[string]interface{})
	if len(latest) == 0 {
		return -1, fmt.Errorf("missing current temperature")
	}

	// Get just the temperature reading.
	return tempConvert(latest["temperature"])
}

// getForecast gets the forecasted temperatures from the NWS database.
// Our values should be here: properties -> periods -> (chosen periods) -> temperature.
// We're going to use these rules to determine which day's forecast we want:
//   1. If it's before 3 pm, we'll use the current day.
//   2. If it's after 3 pm, we'll display the high/low for the next day.
func getForecast(client http.Client, url string) (int, int, error) {
	type forecast struct {
		Properties struct {
			Periods []map[string]interface{} `json:"periods"`
		} `json:"properties"`
	}

	// If it's before 3pm, we'll use the forecast of the current day.
	// After that, we'll use tomorrow's forecast.
	t := time.Now()
	if t.Hour() >= 15 {
		t = t.Add(time.Hour * 12)
	}
	timeStr := t.Format("2006-01-02T") + "18:00:00"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, -1, fmt.Errorf("forecast: bad request")
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, -1, fmt.Errorf("forecast: bad client")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, fmt.Errorf("forecast: bad read")
	}
	// TODO: handle expired grid.

	f := forecast{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return -1, -1, fmt.Errorf("forecast: bad JSON")
	}

	// Get the list of forecasts.
	periods := f.Properties.Periods
	if len(periods) == 0 {
		return -1, -1, fmt.Errorf("missing forecast periods")
	}

	// Iterate through the list until we find the forecast for tomorrow.
	var high int
	var low int
	for _, f := range periods {
		et := f["endTime"].(string)
		st := f["startTime"].(string)
		if strings.Contains(et, timeStr) {
			// We'll get the high from here.
			high, _ = tempConvert(f["temperature"])
		} else if strings.Contains(st, timeStr) {
			// We'll get the low from here.
			low, _ = tempConvert(f["temperature"])

			// This is all we need from the forecast, so we can exit now.
			return high, low, nil
		}
	}

	// If we're here, then we didn't find the forecast.
	return -1, -1, fmt.Errorf("failed to determine forecast")
}

// tempConvert converts the temperature from either a float or string into an int.
func tempConvert(val interface{}) (int, error) {
	switch val.(type) {
	case float64:
		return int(val.(float64)), nil
	case string:
		return strconv.Atoi(val.(string))
	}

	return -1, fmt.Errorf("unknown temperature format")
}
