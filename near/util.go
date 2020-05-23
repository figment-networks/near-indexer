package near

import (
	"net/http"
	"time"
)

func reqWithTiming(c *http.Client, req *http.Request) (*http.Response, time.Duration, error) {
	ts := time.Now()
	resp, err := c.Do(req)
	te := time.Since(ts)

	return resp, te, err
}
