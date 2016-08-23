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
	"strconv"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var milesIndexStr = "_milesindex"				//name for the key/value that will store a list of all known miless
//var openTradesStr = "_opentrades"				//name for the key/value that will store all open trades
var A string

type AirMiles struct{
	
	upliftingAirline string `json:"upliftingAirline"`					
	flightNo string `json:"flightNo"`
	fromSector string `json:"fromSector"`
	toSector string `json:"toSector"`
	bookingClass string `json:"bookingClass"`
	bookingMiles string `json:"bookingMiles"`
	fFP string `json:"fFP"`
	rewardingAirline string `json:"rewardingAirline"`
	rewardedMiles string `json:"rewardedMiles"`
	passengerName string `json:"passengerName"`
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

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(milesIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Init Miles - create a new miles, store into chaincode state
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

	fmt.Println("- start init miles")
	
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
	err = stub.PutState(args[0], []byte(str))								//store miles with id as key
	if err != nil {
		return nil, err
	}
		
	//get the miles index
	var empty []string
	jsonAsBytes, _ = json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(milesIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	milesAsBytes, err := stub.GetState(milesIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get miles index")
	}
	var milesindex []string
	json.Unmarshal(milesAsBytes, &milesindex)							//un stringify it aka JSON.parse()
	
	//append
	milesindex = append(milesindex, args[0])								//add miles name to index list
	fmt.Println("! miles index: ", milesindex)
	jsonAsBytes, _ = json.Marshal(milesindex)
	err = stub.PutState(milesIndexStr, jsonAsBytes)						//store name of miles

	fmt.Println("- end init miles")
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
	
	var err error
	fmt.Println("running write()")
//var jsonAsBytes []byte
	//var upliftingAirline, flightNo, fromSector, toSector, bookingClass, fFP, rewardingAirline, passengerName string
	//var bookingMiles string
	//   0       			1     		  2     		3				4			5				6			7				8
	// "upliftingAirline", "flightNo", "bookingClass", "fromSector"		"toSector"	bookingMiles	fFP		rewardingAirline	passengerName	
	if len(args) != 9 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("- start init miles")
	res := AirMiles{}
	res.upliftingAirline = strings.ToLower(args[0])
	res.flightNo = strings.ToLower(args[1])
	res.bookingClass = strings.ToLower(args[2])
	res.fromSector = strings.ToLower(args[3])
	res.toSector = strings.ToLower(args[4])
	res.bookingMiles =strings.ToLower(args[5])
	res.fFP = strings.ToLower(args[6])
	res.rewardingAirline = strings.ToLower(args[7])
	res.passengerName = strings.ToLower(args[8])

	///res := AirMiles{}
	//json.Unmarshal(milesAsBytes, &res)										//un stringify it aka JSON.parse()
	//res.User = args[1]	
													//change the user

	var bookingClass = res.bookingClass
	var bookingMiles float64
	bookingMiles, err   = strconv.ParseFloat(res.bookingMiles,32)
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	switch bookingClass {
		case "W":
			res.rewardedMiles = strconv.FormatFloat((bookingMiles*0.25), 'f', -1, 64)
			
		case "V":
			res.rewardedMiles = strconv.FormatFloat((bookingMiles*0.50), 'f', -1, 64) 
		case "K":
			res.rewardedMiles = strconv.FormatFloat((bookingMiles*0.75), 'f', -1, 64) 
		case "P":
			res.rewardedMiles = strconv.FormatFloat(bookingMiles, 'f', -1, 64)
		case "U":
			res.rewardedMiles = strconv.FormatFloat((bookingMiles*1.25), 'f', -1, 64)
		default:
			panic("unrecognized escape character")
		}
	//res.RewardingMiles = rewardingMile
	//var rewardedMiles string
	//jsonAsBytes, _ = json.Marshal(res)
	err = stub.PutState(res.rewardingAirline, []byte(res.rewardedMiles))								//rewrite the marble with id as key
	//err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	//var valAsbytes [] byte
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	//var temp string = "10"
	fmt.Println("READ - Key :"+key)
	valAsbytes, err := stub.GetState(key)
	//valAsbytes := [] byte (temp)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// ============================================================================================================================
// Delete - remove a key/value pair from state
// ============================================================================================================================
func (t *SimpleChaincode) Delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	name := args[0]
	err := stub.DelState(name)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the miles index
	milesAsBytes, err := stub.GetState(milesIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get miles index")
	}
	var milesIndex []string
	json.Unmarshal(milesAsBytes, &milesIndex)								//un stringify it aka JSON.parse()
	
	//remove miles from index
	for i,val := range milesIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name{															//find the correct miles id
			fmt.Println("found miles id")
			milesIndex = append(milesIndex[:i], milesIndex[i+1:]...)			//remove it
			for x:= range milesIndex{											//debug prints...
				fmt.Println(string(x) + " - " + milesIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(milesIndex)									//save new index
	err = stub.PutState(milesIndexStr, jsonAsBytes)
	return nil, nil
}
