package util

import (
	"bytes"
	"io/ioutil"
	"time"

	apiclient "github.com/Optum/dce-cli/client"
	"github.com/Optum/dce-cli/client/operations"
	"github.com/Optum/dce-cli/internal/observation"
	"github.com/aws/aws-sdk-go/aws/credentials"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"net/http"
	"net/http/httputil"
)

// Adapted from https://stackoverflow.com/questions/39527847/is-there-middleware-for-go-http-client
type Sig4RoundTripper struct {
	Proxied http.RoundTripper
	Creds   *credentials.Credentials
	Region  string
	Logger  observation.Logger
}

func (srt Sig4RoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	log := srt.Logger
	dumpedReq, err := httputil.DumpRequest(req, true)
	if err != nil {
		srt.Logger.Fatalf(err.Error())
	}
	log.Debugln("V4 Signing Request:\n", string(dumpedReq))

	signer := sigv4.NewSigner(srt.Creds)
	now := time.Now().Add(time.Duration(30) * time.Second)

	// If there's a json provided, add it when signing
	// Body does not matter if added before the signing, it will be overwritten

	executeAPI := "execute-api"
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatalln("Error reading payload for v4 signing. ", err)
		}

		if err != nil {
			log.Fatalln("Error marshaling payload. ", err)
		}
		req.Header.Set("Content-Type", "application/json")
		_, err = signer.Sign(req, bytes.NewReader(body),
			executeAPI, srt.Region, now)

	} else {
		_, err := signer.Sign(req, nil,
			executeAPI, srt.Region, now)
		if err != nil {
			log.Fatalln("Error while v4 signing request. ", err)
		}
	}

	res, e = srt.Proxied.RoundTrip(req)

	log.Debugln("Response: ", res)
	return res, e
}

func (u *APIUtil) InitApiClient() *operations.Client {

	sig4RoundTripper := Sig4RoundTripper{
		Proxied: http.DefaultTransport,
		Creds: credentials.NewStaticCredentials(
			*u.Config.System.MasterAccount.Credentials.AwsAccessKeyID,
			*u.Config.System.MasterAccount.Credentials.AwsSecretAccessKey,
			*u.Config.System.MasterAccount.Credentials.AwsSessionToken,
		),
		Region: *u.Config.Region,
		Logger: log,
	}
	sig4HTTTPClient := http.Client{Transport: &sig4RoundTripper}
	httpTransport := httptransport.NewWithClient(*u.Config.API.Host, *u.Config.API.BasePath, nil, &sig4HTTTPClient)
	client := apiclient.New(httpTransport, strfmt.Default)
	return client.Operations
}
