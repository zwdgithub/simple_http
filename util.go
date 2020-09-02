package simple_http

import "net/url"

func BuildUrl(URL string, params url.Values) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	if params != nil && len(params) > 0 {
		u.RawQuery = params.Encode()
	}
	return u.String(), nil
}
