// frontier.go - domain fronting http.RoundTriper wrapper
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of frontier, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package frontier

import (
	"crypto/tls"
	"net/http"
)

type Frontier struct {
	transport http.RoundTripper
	front     string
	addr      string
}

// New creates a Frontier that sets SNI of requests to front
// and network address to call to addr (or front if addr is empty).
// Frontier roundtrips all requests through t.
func New(t http.RoundTripper, front, addr string) *Frontier {
	fr := &Frontier{
		transport: t,
		front:     front,
		addr:      addr,
	}
	return fr
}

func (fr *Frontier) RoundTrip(r *http.Request) (*http.Response, error) {
	authority := r.URL.Host
	r.Host = authority
	if fr.front != "" {
		r.URL.Host = fr.front
	}
	if fr.addr != "" {
		r.URL.Host = fr.addr
		if r.URL.Scheme == "https" && fr.front != "" {
			t, ok := fr.transport.(*http.Transport)
			if ok {
				if t.TLSClientConfig == nil {
					t.TLSClientConfig = &tls.Config{}
				}
				t.TLSClientConfig.ServerName = fr.front

			}
		}
	}
	return fr.transport.RoundTrip(r)
}
