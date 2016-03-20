// Copyright © 2016 Eduard Angold eddyhub@users.noreply.github.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sync

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/csv"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"net/http"
	"os"
	"time"
)

// set environment config
func NewClient(
	envAddress string,
	envCACert string,
	envInsecure bool,
	envTLSServerName string,
	token string) (*api.Client, error) {

	var newCertPool *x509.CertPool
	var config *api.Config
	var err error

	if envAddress != "" {
		config = &api.Config{
			Address:    envAddress,
			HttpClient: cleanhttp.DefaultClient(),
		}
		config.HttpClient.Timeout = time.Second * 60
		transport := config.HttpClient.Transport.(*http.Transport)
		transport.TLSHandshakeTimeout = 10 * time.Second
		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	} else {
		fmt.Println("Pleas specify the address to the vault server!")
		os.Exit(-1)
	}

	if envCACert != "" || envInsecure {
		//var err error
		if envCACert != "" {
			newCertPool, err = api.LoadCACert(envCACert)
			if err != nil {
				fmt.Errorf("Error setting up CACert: %s", err)
				return nil, nil
			}
		}

	}

	clientTLSConfig := config.HttpClient.Transport.(*http.Transport).TLSClientConfig
	if newCertPool != nil {
		clientTLSConfig.RootCAs = newCertPool
	}
	if envTLSServerName != "" {
		clientTLSConfig.ServerName = envTLSServerName
	}

	client, err := api.NewClient(config)
	client.SetToken(token)
	if err != nil {
		fmt.Errorf("err: %s", err)
		return nil, err
	}

	return client, nil
}

func WriteData(c *api.Client, path string, data map[string]interface{}) {
	c.Logical().Write(path, data)
}

func DataWriter(
	envAddress string,
	envCACert string,
	envInsecure bool,
	envTLSServerName string,
	token string) func(branch string, dbName string, schemaUser string, password string) {
	c, err := NewClient(envAddress, envCACert, envInsecure, envTLSServerName, token)
	if err != nil {
		panic("Error initializing the client!")
	}
	writer := func(branch string, dbName string, schemaUser string, password string) {
		c.Logical().Write("secret/"+branch+"/"+dbName+"/"+schemaUser, map[string]interface{}{"password": password})
	}
	return writer
}

func ReadCsv(path string) (header []string, data [][]string) {
	csvFile, err := os.Open(path)

	if err != nil {
		fmt.Println("Error reading csv file!")
		return
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	header, err = reader.Read()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data, err = reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}

//edi@edi vaultproject]$ ./vault init -address=https://localhost:8200 -ca-cert=ca.crt
//Key 1: 84914ffd0a27b6ee5822504b7542f7425b30832fcbbf0e7078d82564ed379ed701
//Key 2: fe8ec12295576c8505d00687b3e00a9e609aa95e383fb99f7fae02265d00b44802
//Key 3: 7a9311f4853a429ad04421a37f38935d02181ffb2a396ee36a74779c753f85b403
//Key 4: 7e689af79541fd2a62356912bf078981c93e971b39b06086a94d225fe3db66db04
//Key 5: fa754a21852cd335b7a14e3673df1042abbc21be2bb6b7fabc9757e5cbe4572705
//Initial Root Token: 63f56001-8746-2dcf-4bfe-6897b048cd4a
//
//Vault initialized with 5 keys and a key threshold of 3. Please
//securely distribute the above keys. When the Vault is re-sealed,
//restarted, or stopped, you must provide at least 3 of these keys
//to unseal it again.
//
//Vault does not store the master key. Without at least 3 keys,
//your Vault will remain permanently sealed.
