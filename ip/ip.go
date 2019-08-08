package ip

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const DefaultTimeout = 10 * time.Second

type IPGetter struct {
	ServiceURL string
	GetIP      func(url string, ctx context.Context) (string, error)
}

func HTTPGetBody(url string, ctx context.Context) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func HTTPGetString(url string, ctx context.Context) (string, error) {
	bz, err := HTTPGetBody(url, ctx)
	if err != nil {
		return "", err
	}
	return string(bz), err
}

func HTTPJSONGetField(field string) func(string, context.Context) (string, error) {
	return func(url string, ctx context.Context) (string, error) {
		body, err := HTTPGetBody(url, ctx)
		if err != nil {
			return "", err
		}
		m := map[string]string{}
		err = json.Unmarshal(body, &m)
		if err != nil {
			return "", err
		}
		return m[field], err
	}
}

func RunProviders(ipGetters []IPGetter, timeout time.Duration) (string, error) {
	ch := make(chan string, len(ipGetters))
	for _, ipGetter := range ipGetters {
		url := ipGetter.ServiceURL
		go func(url string, ipGetter IPGetter) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			ip, err := ipGetter.GetIP(url, ctx)
			if err == nil {
				ch <- ip
			} else {
				ch <- ""
			}
		}(url, ipGetter)
	}
	limit := len(ipGetters)/2 + 1
	ipCount := map[string]int{}
	for i := 0; i < len(ipGetters); i++ {
		ip := <-ch
		ipCount[ip]++
		if len(ip) > 0 && ipCount[ip] >= limit {
			return ip, nil
		}
	}
	return "", errors.New("cannot get majority agree on the same IP")
}
