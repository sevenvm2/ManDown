package web

import (
	"net/http"
	"net/url"
	"strings"
)

// Health struct to contain Health of site
type Health struct {
	Site   string
	Status int
	Misc   string
}

// CheckBulk checks multiple domain healths at once
func CheckBulk(sites []string) []Health {
	ch := make(chan Health)
	var results []Health

	sites = deDupeStr(sites)
	for _, site := range sites {
		go GetStatus(site, ch)
	}

	for range sites {
		results = append(results, <-ch)
	}

	return results
}

// GetStatus gets the Status code of a single website
func GetStatus(site string, ch chan<- Health) {
	web, err := url.ParseRequestURI(site)
	if err != nil {
		ch <- Health{
			Site:   site,
			Misc:   err.Error(),
			Status: 0,
		}
		return
	}

	resp, err := http.Get(web.String())
	if err != nil {
		ch <- Health{
			Site:   site,
			Misc:   err.Error(),
			Status: 1,
		}
		return
	}

	status := resp.StatusCode
	ch <- Health{
		Site:   site,
		Misc:   "",
		Status: status,
	}
	return
}

// CheckHealth gets the Status code of a website
func CheckHealth(site string) (int, error) {
	web, err := url.ParseRequestURI(site)
	if err != nil {
		return 0, err
	}

	resp, err := http.Get(web.String())
	if err != nil {
		return 1, err
	}

	return resp.StatusCode, nil
}

// Sanitise makes sure only the domain name gets through
func Sanitise(site string) string {
	web, err := url.Parse(strings.ToLower(site))
	if err != nil {
		return ""
	}

	if web.Scheme != "http" {
		web.Scheme = "https"
	}

	web, _ = url.Parse(web.String())
	return web.Scheme + "://" + web.Hostname()
}

// deDupeStr takes an array of strings and returns only unique items strings
func deDupeStr(strs []string) []string {
	deDuped := []string{}
	for _, item := range strs {
		skip := false
		for _, el := range deDuped {
			if el == item {
				skip = true
				break
			}
		}
		if !skip && item != "" {
			deDuped = append(deDuped, item)
		}
	}
	return deDuped
}
