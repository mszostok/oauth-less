package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Config handles application configuration
type Config struct {
	ClientID     string
	ClientSecret string
	TokenURL     string

	RedirectBaseURL string
}

func main() {
	var cfg Config
	err := envconfig.InitWithPrefix(&cfg, "APP")
	fatalOnErr(err, "while init cfg")

	clientCredCfg := clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
	}
	oauth2.RegisterBrokenAuthHeaderProvider(cfg.TokenURL)
	client := clientCredCfg.Client(context.Background())

	target, err := url.Parse(cfg.RedirectBaseURL)
	fatalOnErr(err, "while parsing RedirectBaseURL")

	customDirector := reverseProxyDirector(target)
	proxy := &httputil.ReverseProxy{
		Director:  customDirector,
		Transport: client.Transport,
	}

	err = http.ListenAndServe("localhost:9090", proxy)
	fatalOnErr(err, "while executing ListenAndServe")
}

func fatalOnErr(err error, ctx string) {
	if err != nil {
		log.Fatal(errors.Wrap(err, ctx).Error())
	}
}

func reverseProxyDirector(target *url.URL) func(req *http.Request) {
	singleJoiningSlash := func(a, b string) string {
		aslash := strings.HasSuffix(a, "/")
		bslash := strings.HasPrefix(b, "/")
		switch {
		case aslash && bslash:
			return a + b[1:]
		case !aslash && !bslash:
			return a + "/" + b
		}
		return a + b
	}

	// https://github.com/golang/go/issues/10342#issuecomment-89458313
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.Host = target.Host

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return director
}
