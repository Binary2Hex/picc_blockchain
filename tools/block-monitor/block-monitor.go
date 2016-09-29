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
	"flag"
	"fmt"
	"github.com/hyperledger/fabric/events/consumer"
	pb "github.com/hyperledger/fabric/protos"
	"net/http"
	"errors"
	"io/ioutil"
	"encoding/json"
	"os"
)

const PEER0 = 0
const PEER1 = 1
const PEER2 = 2
const PEER3 = 3

const CHAN_SIZE = 100

type adapter struct {
	name int
	addr string
	notfy chan *pb.Event_Block
}

type BlockEvent struct {
	peer int
	event *pb.Event_Block
}

type ChainResponse struct {
	Height int  `json:"height"`
	CurrentBlockHash string  `json:"currentBlockHash"`
	PreviousBlockHash string  `json:"previousBlockHash"`

}

//GetInterestedEvents implements consumer.EventAdapter interface for registering interested events
func (a *adapter) GetInterestedEvents() ([]*pb.Interest, error) {
	return []*pb.Interest{{EventType: pb.EventType_BLOCK}}, nil
}

//Recv implements consumer.EventAdapter interface for receiving events
func (a *adapter) Recv(msg *pb.Event) (bool, error) {
	switch msg.Event.(type) {
	case *pb.Event_Block:
		a.notfy <- msg.Event.(*pb.Event_Block)
		return true, nil
	default:
		//a.notfy <- nil
		return false, nil
	}
}

//Disconnected implements consumer.EventAdapter interface for disconnecting
func (a *adapter) Disconnected(err error) {
	//fmt.Println(a.name)
	//fmt.Println(a.addr)
	os.Exit(1)
}

func createEventClient(name int, eventAddress string) *adapter {
	var obcEHClient *consumer.EventsClient

	done := make(chan *pb.Event_Block, CHAN_SIZE)
	adapter := &adapter{name: name, addr: eventAddress, notfy: done}
	obcEHClient = consumer.NewEventsClient(eventAddress, adapter)
	if obcEHClient == nil {
		return nil
	}
	if err := obcEHClient.Start(); err != nil {
		fmt.Printf("could not start chat %s\n", err)
		//obcEHClient.Stop()
		return nil
	}
	return adapter
}

func main() {
	var eventAddress0 string
	var eventAddress1 string
	var eventAddress2 string
	var eventAddress3 string
	var restAddress0 string
	var restAddress1 string
	var restAddress2 string
	var restAddress3 string
	var events0 int
	var events1 int
	var events2 int
	var events3 int

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


	a := createEventClient(PEER0, eventAddress0)
	b := createEventClient(PEER1, eventAddress1)
	c := createEventClient(PEER2, eventAddress2)
	d := createEventClient(PEER3, eventAddress3)
	fan := make(chan *BlockEvent, 4*CHAN_SIZE)

	events0 = 0
	events1 = 0
	events2 = 0
	events3 = 0

	initialHeight(&events0, restAddress0)
	if a != nil {
		go recvEvent(PEER0, a.notfy, fan)
	}
	initialHeight(&events1, restAddress1)
	if b != nil {
		go recvEvent(PEER1,  b.notfy, fan)
	}
	initialHeight(&events2, restAddress2)
	if c != nil {
		go recvEvent(PEER2, c.notfy, fan)
	}
	initialHeight(&events3, restAddress3)
	if d != nil {
		go recvEvent(PEER3, d.notfy, fan)
	}


	printHeight(events0, events1, events2, events3)

	if a == nil || b == nil || c == nil || d == nil {
		fmt.Printf("Error creating event client\n")
		os.Exit(2)
	}

	for {
		event := <- fan
		if event.peer == PEER0{
			events0++
			printHeight(events0, events1, events2, events3)
			continue
		}
		if event.peer == PEER1{
			events1++
			printHeight(events0, events1, events2, events3)
			continue
		}
		if event.peer == PEER2{
			events2++
			printHeight(events0, events1, events2, events3)
			continue
		}
		if event.peer == PEER3{
			events3++
			printHeight(events0, events1, events2, events3)
			continue
		}

	}

}

func recvEvent(peer int, in chan *pb.Event_Block, out chan *BlockEvent) {
	for blockEvent := range in{
		if blockEvent.Block.NonHashData.TransactionResults == nil {
			//fmt.Printf("INVALID BLOCK ... NO TRANSACTION RESULTS %v\n", b)
		} else {
			out <- &BlockEvent{peer: peer}
		}
	}
	fmt.Println("Goodbye World")
}

func printHeight(events0, events1, events2, events3 int) {
	fmt.Printf("peers height: |%12d|, |%12d| ,|%12d|, |%12d|\n", events0, events1, events2, events3)
}

func initialHeight(events *int, addr string) (bool, error){
	res, err := http.Get("http://"+addr+"/chain")
	if err != nil {
		fmt.Errorf("Couldn't get response from:" + addr)
		return false, errors.New("Couldn't get response from:" + addr)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Couldn't read response body")
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
