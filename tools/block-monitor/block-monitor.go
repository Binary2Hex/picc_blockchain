/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ChainResponse struct {
	Height            int    `json:"height"`
	CurrentBlockHash  string `json:"currentBlockHash"`
	PreviousBlockHash string `json:"previousBlockHash"`
}

const TICKER = 1 * time.Second

func main() {
	var eventAddress0 string
	var eventAddress1 string
	var eventAddress2 string
	var eventAddress3 string
	var restAddress0 string
	var restAddress1 string
	var restAddress2 string
	var restAddress3 string

	flag.StringVar(&eventAddress0, "events-address0", "0.0.0.0:31315", "address of events server")
	flag.StringVar(&eventAddress1, "events-address1", "0.0.0.1:31315", "address of events server")
	flag.StringVar(&eventAddress2, "events-address2", "0.0.0.2:31315", "address of events server")
	flag.StringVar(&eventAddress3, "events-address3", "0.0.0.3:31315", "address of events server")

	flag.StringVar(&restAddress0, "rest-address0", "0.0.0.0:5000", "address of rest server")
	flag.StringVar(&restAddress1, "rest-address1", "0.0.0.1:5000", "address of rest server")
	flag.StringVar(&restAddress2, "rest-address2", "0.0.0.2:5000", "address of rest server")
	flag.StringVar(&restAddress3, "rest-address3", "0.0.0.3:5000", "address of rest server")

	flag.Parse()

	fmt.Printf("Event Address0: %s\n", eventAddress0)
	fmt.Printf("Event Address1: %s\n", eventAddress1)
	fmt.Printf("Event Address2: %s\n", eventAddress2)
	fmt.Printf("Event Address3: %s\n", eventAddress3)

	fmt.Printf("REST Address0: %s\n", restAddress0)
	fmt.Printf("REST Address1: %s\n", restAddress1)
	fmt.Printf("REST Address2: %s\n", restAddress2)
	fmt.Printf("REST Address3: %s\n", restAddress3)

	events0 := 0
	events1 := 0
	events2 := 0
	events3 := 0

	printHeight(&events0, &events1, &events2, &events3)

	timer := time.NewTicker(TICKER)

	for {
		select {
		case <-timer.C:
			{
				initialHeight(&events0, restAddress0)
				initialHeight(&events1, restAddress1)
				initialHeight(&events2, restAddress2)
				initialHeight(&events3, restAddress3)
				printHeight(&events0, &events1, &events2, &events3)
			}
		}

	}

}

func printHeight(events0, events1, events2, events3 *int) {
	fmt.Printf("peers height: |%12d|, |%12d| ,|%12d|, |%12d|\n", *events0, *events1, *events2, *events3)
}

func initialHeight(events *int, addr string) (bool, error) {
	*events = 0
	res, err := http.Get("http://" + addr + "/chain")
	if err != nil {
		return false, errors.New("Couldn't get response from:" + addr)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, errors.New("Couldn't read response body")
	}
	var response ChainResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	*events = response.Height
	return true, nil
}
