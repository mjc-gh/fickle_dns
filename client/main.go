package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	// API flags
	var host string
	var zoneId string
	var recordId string
	var provider string

	// mTLS flags
	var certFile string
	var keyFile string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "Server host name running server component",
				Destination: &host,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "zone",
				Usage:       "Provider Zone ID",
				Destination: &zoneId,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "record",
				Usage:       "Provider Record ID",
				Destination: &recordId,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "provider",
				Value:       "cloudflare",
				Usage:       "Provider Type",
				Destination: &provider,
			},
			&cli.StringFlag{
				Name:        "cert",
				Value:       "cert/fickle_dns.crt",
				Usage:       "Client TLS certificate file",
				Destination: &certFile,
			},
			&cli.StringFlag{
				Name:        "key",
				Value:       "cert/fickle_dns.key",
				Usage:       "Client TLS private key file",
				Destination: &keyFile,
			},
		},
		Action: func(cCtx *cli.Context) error {
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return err
			}

			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						Certificates: []tls.Certificate{cert},
					},
				},
			}

			url := fmt.Sprintf("https://%s/providers/%s/zones/%s/records/%s", host, provider, zoneId, recordId)
			log.Printf("PATCH %s\n", url)

			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer([]byte{}))
			if err != nil {
				return err
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			log.Printf("response status: %s", resp.Status)
			log.Printf("response: %s", string(body))

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
