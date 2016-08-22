/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	//"strconv"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var milesIndexStr = "_milesindex"				//name for the key/value that will store a list of all known marbles
//var openTradesStr = "_opentrades"				//name for the key/value that will store all open trades

type AirMiles struct{
	
	UpliftingAirline string `json:"upliftingAirline"`					
	FlightNo string `json:"flightNo"`
	FromSector string `json:"fromSector"`
	ToSector string `json:"toSector"`
	BookingClass string `json:"bookingClass"`
	BookingMiles float32 `json:"bookingMiles"`
	FFP string `json:"fFP"`
	RewardingAirline string `json:"rewardingAirline"`
	RewardedMiles float32 `json:"rewardedMiles"`
	PassengerName string `json:"passengerName"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	fmt.Println("inside init function ")
	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Init Marble - create a new marble, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) init_miles(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
	var jsonAsBytes []byte
	var upliftingAirline, flightNo, fromSector, toSector, bookingClass, fFP, rewardingAirline, passengerName string
	var bookingMiles string
	//   0       			1     		  2     		3				4			5				6			7				8
	// "upliftingAirline", "flightNo", "bookingClass", "fromSector"		"toSector"	bookingMiles	fFP		rewardingAirline	passengerName	
	if len(args) != 9 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("- start init marble")
	
	upliftingAirline = strings.ToLower(args[0])
	flightNo = strings.ToLower(args[1])
	bookingClass = strings.ToLower(args[2])
	fromSector = strings.ToLower(args[3])
	toSector = strings.ToLower(args[4])
	bookingMiles =strings.ToLower(args[5])
	/////bookingMiles, err = strconv.Atoi(args[5])
	//if err != nil {
	//	return nil, errors.New("5 argument must be a numeric string")
	//}
	fFP = strings.ToLower(args[6])
	rewardingAirline = strings.ToLower(args[7])
	passengerName = strings.ToLower(args[8])

	str := `{"upliftingAirline": "` + upliftingAirline + `", "flightNo": "` + flightNo + `", "bookingClass": "` + bookingClass + `", "fromSector": "` + fromSector + `" , "toSector": "` + toSector + `", "bookingMiles": "` + bookingMiles + `", "fFP": "` + fFP + `", "rewardingAirline": "` + rewardingAirline + `", "passengerName": "` + passengerName+	`"}`
	err = stub.PutState(args[0], []byte(str))								//store marble with id as key
	if err != nil {
		return nil, err
	}
		
	//get the marble index
	var empty []string
	jsonAsBytes, _ = json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(milesIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	milesAsBytes, err := stub.GetState(milesIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var milesindex []string
	json.Unmarshal(milesAsBytes, &milesindex)							//un stringify it aka JSON.parse()
	
	//append
	milesindex = append(milesindex, args[0])								//add marble name to index list
	fmt.Println("! marble index: ", milesindex)
	jsonAsBytes, _ = json.Marshal(milesindex)
	err = stub.PutState(milesIndexStr, jsonAsBytes)						//store name of marble

	fmt.Println("- end init marble")
	return nil, nil
}
// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//  W       V		K		P		U
	// "25%", "50%"		"75%"	100%	"125%"
	//var key, value string
	var err error
	fmt.Println("running write()")

	
	
	
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	
	//key = args[0] //rename for funsies
	//fmt.Println("- start set user")
	//fmt.Println(args[0] + " - " + args[1])
	milesAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	res := AirMiles{}
	json.Unmarshal(milesAsBytes, &res)										//un stringify it aka JSON.parse()
	//res.User = args[1]	
													//change the user

	var bookingClass string  = res.BookingClass
	var bookingMiles float32  = res.BookingMiles
	
	switch bookingClass {
		case "W":
			res.RewardedMiles = bookingMiles*0.25
		case "V":
			res.RewardedMiles = bookingMiles*0.50
		case "K":
			res.RewardedMiles = bookingMiles*0.75
		case "P":
			res.RewardedMiles = bookingMiles
		case "U":
			res.RewardedMiles = bookingMiles*1.25	
		default:
			panic("unrecognized escape character")
		}
	//res.RewardingMiles = rewardingMile
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes)								//rewrite the marble with id as key
	//value = args[1]
	//value = "fixed value"
	//err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
