package opsman

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/schema"
	omnet "github.com/pivotal-cf/om/network"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type omWriter struct{}

func (o omWriter) Write(p []byte) (n int, err error) {
	log.Printf("[DEBUG] OM - %s", string(p[:]))
	return len(p), nil
}

func createAuthedClient(d *schema.ResourceData) (httpClient, error) {
	var authedClient httpClient
	authedClient, err := omnet.NewOAuthClient(
		d.Get("address").(string),
		d.Get("username").(string),
		d.Get("password").(string),
		"",    // client id
		"",    // client secret
		true,  // skip ssl validation
		false, // don't include cookies
		time.Duration(30)*time.Second)
	if err != nil {
		return nil, err
	}
	if logging.IsDebugOrHigher() {
		authedClient = omnet.NewTraceClient(authedClient, omWriter{})
	}
	return authedClient, nil
}

func createUnauthedClient(d *schema.ResourceData) omnet.UnauthenticatedClient {
	return omnet.NewUnauthenticatedClient(
		d.Get("address").(string),
		true, // skip ssl validation
		time.Duration(2)*time.Minute)
}
