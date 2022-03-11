package main

import (
	"io"
	"net/http"
	"regexp"
)

func getCountry(ip string) string {
	resp, err := http.Get("http://api.geoiplookup.net/?query=" + ip)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err.Error()
		}
		bodyString := string(bodyBytes)
		re := regexp.MustCompile(`<countryname.*?>(.*)</countryname>`)
		submatchall := re.FindAllStringSubmatch(bodyString, -1)
		for _, element := range submatchall {
			return element[1]
		}
	}

	return "no country"
}