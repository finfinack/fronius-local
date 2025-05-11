package fronius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/icholy/digest"

	"github.com/finfinack/fronius-local/pkg/config"
	"github.com/finfinack/logger/logging"
)

const (
	batteriesPath = "/api/config/batteries"
	timeOfUsePath = "/api/config/timeofuse"
)

type Fronius struct {
	cfg    *config.Fronius
	log    *logging.Logger
	client *http.Client
}

func New(cfg *config.Fronius, tz *time.Location) *Fronius {
	timezone = tz
	return &Fronius{
		cfg: cfg,
		log: logging.NewLogger("FRON"),
		client: &http.Client{
			Transport: &digest.Transport{
				Username: cfg.Username,
				Password: cfg.Password,
				FindChallenge: func(h http.Header) (*digest.Challenge, error) {
					value := h.Get("X-WWW-Authenticate")
					if value == "" {
						return nil, digest.ErrNoChallenge
					}
					return digest.ParseChallenge(value)
				},
			},
		},
	}
}

func (f *Fronius) request(method, path string, body []byte) ([]byte, error) {
	u, err := url.Parse(f.cfg.Backend)
	if err != nil {
		return nil, fmt.Errorf("unable to parse base URL %q: %s", f.cfg.Backend, err)
	}
	u.Path = path

	req, err := http.NewRequest(method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %s", err)
	}

	res, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch data: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d for %q", res.StatusCode, u.String())
	}

	return io.ReadAll(res.Body)
}

func (f *Fronius) FetchTimeOfUse() ([]*TimeOfUse, error) {
	f.log.Infoln("fetching time of use data")
	data, err := f.request(http.MethodGet, timeOfUsePath, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to run request: %s", err)
	}

	resp := &TimeOfUseResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body: %s", err)
	}

	return resp.TimeOfUse, nil
}

func (f *Fronius) FetchBattery() (*Battery, error) {
	f.log.Infoln("fetching battery data")
	data, err := f.request(http.MethodGet, batteriesPath, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to run request: %s", err)
	}

	resp := &Battery{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body: %s", err)
	}

	return resp, nil
}
