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
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var marbleIndexStr = "_marbleindex"	
var driverIndexStr = "_driverindex"				//name for the key/value that will store a list of all known marbles
var openTradesStr = "_opentrades"				//name for the key/value that will store all open trades

type Marble struct{
	Name string `json:"name"`					//the fieldtags are needed to keep case from bouncing around
	Color string `json:"color"`
	Size int `json:"size"`
	User string `json:"user"`
}

type Driver struct{
	Name string `json:"name"`
	DL string `json:"dl"`		
	DOB string `json:"dob"`	
	Email string `json:"email"`
	Mobile string `json:"mobile"`
	Password string `json:"password"`
	Address string `json:"address"`
	Status string `json:"status"`
	Modifyby string `json:"modifyby"`
	Adminemail string `json:"adminemail"`
	Rejectreason string `json:"rejectreason"`
	Anycomment string `json:"anycomment"`
	Bookingid  string `json:"bookingid"`
}

type Bookcar struct{
	Bookacarname string `json:"bookacarname"`
	Bookacaremail string `json:"bookacaremail"`		
	Bookacarclass string `json:"bookacarclass"`
	Bookacarlocation string `json:"bookacarlocation"`
	Bookacardroplocation string `json:"bookacardroplocation"`
	Bookacarpickupdate string `json:"bookacarpickupdate"`
	Bookacarpickuptime string `json:"bookacarpickuptime"`
	Bookacardropoffdate string `json:"bookacardropoffdate"`
	Bookacardropofftime string `json:"bookacardropofftime"`
	Bookingid string `json:"bookingid"`
}


type Description struct{
	Color string `json:"color"`
	Size int `json:"size"`
}

type AnOpenTrade struct{
	User string `json:"user"`					//user who created the open trade order
	Timestamp int64 `json:"timestamp"`			//utc timestamp of creation
	Want Description  `json:"want"`				//description of desired marble
	Willing []Description `json:"willing"`		//array of marbles willing to trade away
}

type AllTrades struct{
	OpenTrades []AnOpenTrade `json:"open_trades"`
}

type Adminlogin struct{
	Userid string `json:"userid"`					//User login for system Admin
	Password string `json:"password"`
}

// ============================================================================================================================
// Main
// =======================================================	=====================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	
//Changes for the Hertz Blockchain

   if len(args) != 2 {
	   return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
    //Write the User Id "mail Id" arg[0] and password arg[1]
	userid := args[0]															//argument for UserID
	password := args[1]  	//argument for password
	str := `{"userid": "` + userid+ `", "password": "` + password + `"}`
	
	err = stub.PutState(userid, []byte(str))								//Put the userid and password in blockchain
	if err != nil {
		return nil, err
	}
	
	
//End of Changes for the Hertz Blockchain}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(marbleIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	//var empty []string
	 //jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	 err = stub.PutState(driverIndexStr, jsonAsBytes)
	 if err != nil {
		return nil, err
	 }
	
	var trades AllTrades
	jsonAsBytes, _ = json.Marshal(trades)								//clear the open trade struct
	err = stub.PutState(openTradesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "delete" {										//deletes an entity from its state
		res, err := t.Delete(stub, args)
		cleanTrades(stub)													//lets make sure all open trades are still valid
		return res, err
	} else if function == "write" {											//writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "init_marble" {									//create a new marble
		return t.init_marble(stub, args)
	} else if function == "signup_driver" {									//create a new marble
		return t.signup_driver(stub, args)
	} else if function == "book_car" {									//create a new marble
		return t.book_car(stub, args)
	}else if function == "set_user" {										//change owner of a marble
		res, err := t.set_user(stub, args)
		cleanTrades(stub)													//lets make sure all open trades are still valid
		return res, err	
	}else if function == "set_status" {										//change owner of a marble
		res, err := t.set_status(stub, args)
		cleanTrades(stub)													//lets make sure all open trades are still valid
		return res, err	
	} else if function == "open_trade" {									//create a new trade order
		return t.open_trade(stub, args)
	} else if function == "perform_trade" {									//forfill an open trade order
		res, err := t.perform_trade(stub, args)
		cleanTrades(stub)													//lets clean just in case
		return res, err
	} else if function == "remove_trade" {									//cancel an open trade order
		return t.remove_trade(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	} else if function == "read_sysadmin" {									//Read system admin User id and password
		return t.read_sysadmin(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)	//get the var from chaincode 
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}
//=============================================================
// Read - query function to read key/value pair (System Admin read User id and Password)
//===============================================================================================================================
func (t *SimpleChaincode) read_sysadmin(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	userid := args[0]
	PassAsbytes, err := stub.GetState(userid)
	
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + userid + "\"}"
		return nil, errors.New(jsonResp)
	}
	
	res := Adminlogin{}
	json.Unmarshal(PassAsbytes,&res)
	
	if res.Userid == userid{
	   fmt.Println("Userid Password Matched: " +res.Userid + res.Password)
	  }else {
	   fmt.Println("Wrong ID Password: " +res.Userid + res.Password)
	   }
	
	return PassAsbytes, nil
}

// ============================================================================================================================
// Delete - remove a key/value pair from state
// ============================================================================================================================
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	name := args[0]
	err := stub.DelState(name)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the marble index
	marblesAsBytes, err := stub.GetState(marbleIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var marbleIndex []string
	json.Unmarshal(marblesAsBytes, &marbleIndex)								//un stringify it aka JSON.parse()
	
	//remove marble from index
	for i,val := range marbleIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name{															//find the correct marble
			fmt.Println("found marble")
			marbleIndex = append(marbleIndex[:i], marbleIndex[i+1:]...)			//remove it
			for x:= range marbleIndex{											//debug prints...
				fmt.Println(string(x) + " - " + marbleIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(marbleIndex)									//save new index
	err = stub.PutState(marbleIndexStr, jsonAsBytes)
	return nil, nil
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Init Marble - create a new marble, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) init_marble(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}
	
    fmt.Println("-Amit C code-Init_marble driver")
	
	//input sanitation
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	name := args[0]
	color := strings.ToLower(args[1])
	user := strings.ToLower(args[3])
	size, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}

	//check if marble already exists
	marbleAsBytes, err := stub.GetState(name)
	if err != nil {
		return nil, errors.New("Failed to get marble name")
	}
	res := Marble{}
	json.Unmarshal(marbleAsBytes, &res)
	if res.Name == name{
		fmt.Println("This marble arleady exists: " + name)
		fmt.Println(res);
		return nil, errors.New("This marble arleady exists")				//all stop a marble by this name exists
	}
	
	//build the marble json string manually
	str := `{"name": "` + name + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "user": "` + user + `"}`
	err = stub.PutState(name, []byte(str))									//store marble with id as key
	if err != nil {
		return nil, err
	}
		
	//get the marble index
	marblesAsBytes, err := stub.GetState(marbleIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var marbleIndex []string
	json.Unmarshal(marblesAsBytes, &marbleIndex)							//un stringify it aka JSON.parse()
	
	//append
	marbleIndex = append(marbleIndex, name)									//add marble name to index list
	fmt.Println("! marble index: ", marbleIndex)
	jsonAsBytes, _ := json.Marshal(marbleIndex)
	err = stub.PutState(marbleIndexStr, jsonAsBytes)						//store name of marble

	fmt.Println("- end init marble")
	return nil, nil
}


// ============================================================================================================================
// Init Marble - create a new marble, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) signup_driver(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       			2						 3
	// "Mainak", "Mandal", "mainakmandal@hotmail.com", "password"
	
	name := args[0]
	dl := args[1]
	dob := args[2]
	email := strings.ToLower(args[3])
	mobile := args[4]
	password := args[5]
	address := args[6]
	status := args[7]
	modifyby := args[8]
	adminemail := args[9]
	rejectreason := args[10]
	anycomment := args[11]

	
	//if err != nil {
		//return nil, errors.New("3rd argument must be a numeric string")
	//}

	//check if marble already exists
	driverAsBytes, err := stub.GetState(email)
	if err != nil {
		return nil, errors.New("Failed to get driver name")
	}
	res := Driver{}
	json.Unmarshal(driverAsBytes, &res)
	if res.Email == email{
		fmt.Println("This marble arleady exists: " + email)
		fmt.Println(res);
		return nil, errors.New("This driver arleady exists")				//all stop a marble by this name exists
	}
	
	//build the marble json string manually
	str := `{"name": "` + name + `", "dl": "` + dl + `", "dob": "` + dob + `", "email": "` + email + `",  "mobile": "` + mobile + `", "password": "` + password + `","address": "` + address + `","status": "` + status + `","modifyby": "` + modifyby + `" ,"adminemail": "` + adminemail + `" ,"rejectreason": "` + rejectreason + `" ,"anycomment": "` + anycomment + `"}`
	err = stub.PutState(email, []byte(str))									//store marble with id as key
	if err != nil {
		return nil, err
	}
		
	//get the driver index
	driversAsBytes, err := stub.GetState(driverIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get driver index")
	}
	var driverIndex []string
	json.Unmarshal(driversAsBytes, &driverIndex)							//un stringify it aka JSON.parse()
	
	//append
	 driverIndex = append(driverIndex, email)									//add marble name to index list
	 fmt.Println("! driver index: ", driverIndex)
	 jsonAsBytes, _ := json.Marshal(driverIndex)
	 err = stub.PutState(driverIndexStr, jsonAsBytes)						//store name of marble

	fmt.Println("- end signup driver")
	return nil, nil
}
// ============================================================================================================================
// Book_car - create a new booking, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) book_car(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	//   0       1       			2						 3
	// "Mainak", "Mandal", "mainakmandal@hotmail.com", "password"
	
	bookacarname := args[0]
	bookacaremail := args[1]
	bookacarclass := args[2]
	bookacarlocation := args[3]
        bookacardroplocation := args[4]
	bookacarpickupdate := args[5]
	bookacarpickuptime := args[6]
	bookacardropoffdate := args[7]
	bookacardropofftime := args[8]
	bookingid := args[9]	
	
	//This function is to check the login id of Driver
	  PassAsbytes, err := stub.GetState(bookacaremail)
	
	  fmt.Println("Wrong ID Password-1: " +bookacaremail)
	
	  if err != nil {
		 jsonResp := "{\"Error\":\"Failed to get state for " + bookacaremail + "\"}"
		 return nil, errors.New(jsonResp)
	     }
	
	  res := Driver{}
	  json.Unmarshal(PassAsbytes,&res)
	  fmt.Println("res.Email-2: " +res.Email)
	
	  if res.Email != bookacaremail{
	    fmt.Println("Wrong ID Password-3: " +res.Email)
	    
	  }else {
	    fmt.Println("Userid Password Matched-4: " +res.Email)
	  
	//-----------------------------------------------------------------
	//build the Bookid json string manually
	str := `{"bookacarname": "` + bookacarname + `", "bookacaremail": "` + bookacaremail + `", "bookacarclass": "` + bookacarclass + `","bookacarlocation": "` + bookacarlocation + `", "bookacardroplocation": "` + bookacardroplocation + `","bookacarpickupdate": "` + bookacarpickupdate + `", "bookacarpickuptime": "` + bookacarpickuptime + `", "bookacardropoffdate": "` + bookacardropoffdate + `", "bookacardropofftime": "` + bookacardropofftime + `","bookingid": "` + bookingid + `"}`
	err = stub.PutState(bookacaremail+bookingid, []byte(str))									 
	if err != nil {
		return nil, err
	}
	
	//-------------------------------------------------get the driver index
	driversAsBytes, err := stub.GetState(driverIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get driver index")
	}
	var driverIndex []string
	json.Unmarshal(driversAsBytes, &driverIndex)							 
	
	//append
	 driverIndex = append(driverIndex, bookingid)									 
	 fmt.Println("! driver index: ", driverIndex)
	 jsonAsBytes1, _ := json.Marshal(driverIndex)
	 err = stub.PutState(driverIndexStr, jsonAsBytes1)	
     //-------------------------------------------------Driver index end	 

    
	//++++++++Boooking id for driver details
	driverAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return nil, errors.New("Failed to get driver name")
	}
	res := Driver{}
	json.Unmarshal(driverAsBytes, &res)
	res.Bookingid = args[9]	 //change the user
 	jsonAsBytes, _ := json.Marshal(res)
 	err = stub.PutState(args[1], jsonAsBytes)
	
	if err != nil {
		return nil, err
	}
	
	//---------------------------get the driver index for Driver details
	driversAsBytes1, err := stub.GetState(driverIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get driver index")
	}
	var driverIndex1 []string
	json.Unmarshal(driversAsBytes1, &driverIndex1)							 
	
	//append
	 driverIndex1 = append(driverIndex1, bookacaremail)									 
	 fmt.Println("! driver index: ", driverIndex1)
	 jsonAsBytes3, _ := json.Marshal(driverIndex1)
	 err = stub.PutState(driverIndexStr, jsonAsBytes3)		
	 //----------------------------------------------------------------------
	 
	fmt.Println("- end signup driver")
	}
	return PassAsbytes, nil
}
// ============================================================================================================================
// Set User Permission on Marble
// ============================================================================================================================
  func (t *SimpleChaincode) set_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  	var err error
	//   0       1
	// "name", "bob"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	 fmt.Println("- start set user")
	 fmt.Println(args[0] + " - " + args[1])
	 marbleAsBytes, err := stub.GetState(args[0])
 	if err != nil {
	 	return nil, errors.New("Failed to get thing")
	 }
	 res := Marble{}
	 json.Unmarshal(marbleAsBytes, &res)										//un stringify it aka JSON.parse()
	 res.User = args[1]														//change the user
	
 	jsonAsBytes, _ := json.Marshal(res)
 	err = stub.PutState(args[0], jsonAsBytes)								//rewrite the marble with id as key
	if err != nil {
		return nil, err
	}
	
 	fmt.Println("- end set user")
	return nil, nil
  }
  
 // ============================================================================================================================
// Set User Permission on Marble
// ============================================================================================================================
  func (t *SimpleChaincode) set_status(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  	var err error
	//   0       1
	// "email", "status"
	if len(args) < 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	 fmt.Println("- start set user")
	 fmt.Println(args[0] + " - " + args[1])
	 
	driverAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get driver name")
	}
	res := Driver{}
	json.Unmarshal(driverAsBytes, &res) 
	 res.Name = args[1]	 //change the user
	 res.DL = args[2]
	 res.DOB = args[3]
	 res.Mobile = args[4]
	 res.Password = args[5]
	 res.Address = args[6]
	 res.Status = args[7]
	 res.Modifyby = args[8]
	 res.Adminemail  = args[9]
	 res.Rejectreason   = args[10]
	 res.Anycomment   = args[11]
	 
 	jsonAsBytes, _ := json.Marshal(res)
 	err = stub.PutState(args[0], jsonAsBytes)								//rewrite the user status with email-id as key
	if err != nil {
		return nil, err
	}
	
 	fmt.Println("- end set user")
	return nil, nil
  } 
// ============================================================================================================================
// Open Trade - create an open trade for a marble you want with marbles you have 
// ============================================================================================================================
func (t *SimpleChaincode) open_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var will_size int
	var trade_away Description
	
	//	0        1      2     3      4      5       6
	//["bob", "blue", "16", "red", "16"] *"blue", "35*
	if len(args) < 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting like 5?")
	}
	if len(args)%2 == 0{
		return nil, errors.New("Incorrect number of arguments. Expecting an odd number")
	}

	size1, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}

	open := AnOpenTrade{}
	open.User = args[0]
	open.Timestamp = makeTimestamp()											//use timestamp as an ID
	open.Want.Color = args[1]
	open.Want.Size =  size1
	fmt.Println("- start open trade")
	jsonAsBytes, _ := json.Marshal(open)
	err = stub.PutState("_debug1", jsonAsBytes)

	for i:=3; i < len(args); i++ {												//create and append each willing trade
		will_size, err = strconv.Atoi(args[i + 1])
		if err != nil {
			msg := "is not a numeric string " + args[i + 1]
			fmt.Println(msg)
			return nil, errors.New(msg)
		}
		
		trade_away = Description{}
		trade_away.Color = args[i]
		trade_away.Size =  will_size
		fmt.Println("! created trade_away: " + args[i])
		jsonAsBytes, _ = json.Marshal(trade_away)
		err = stub.PutState("_debug2", jsonAsBytes)
		
		open.Willing = append(open.Willing, trade_away)
		fmt.Println("! appended willing to open")
		i++;
	}
	
	//get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)										//un stringify it aka JSON.parse()
	
	trades.OpenTrades = append(trades.OpenTrades, open);						//append to open trades
	fmt.Println("! appended open to trades")
	jsonAsBytes, _ = json.Marshal(trades)
	err = stub.PutState(openTradesStr, jsonAsBytes)								//rewrite open orders
	if err != nil {
		return nil, err
	}
	fmt.Println("- end open trade")
	return nil, nil
}

// ============================================================================================================================
// Perform Trade - close an open trade and move ownership
// ============================================================================================================================
func (t *SimpleChaincode) perform_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	//	0		1					2					3				4					5
	//[data.id, data.closer.user, data.closer.name, data.opener.user, data.opener.color, data.opener.size]
	if len(args) < 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}
	
	fmt.Println("- start close trade")
	timestamp, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, errors.New("1st argument must be a numeric string")
	}
	
	size, err := strconv.Atoi(args[5])
	if err != nil {
		return nil, errors.New("6th argument must be a numeric string")
	}
	
	//get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)															//un stringify it aka JSON.parse()
	
	for i := range trades.OpenTrades{																//look for the trade
		fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))
		if trades.OpenTrades[i].Timestamp == timestamp{
			fmt.Println("found the trade");
			
			
			marbleAsBytes, err := stub.GetState(args[2])
			if err != nil {
				return nil, errors.New("Failed to get thing")
			}
			closersMarble := Marble{}
			json.Unmarshal(marbleAsBytes, &closersMarble)											//un stringify it aka JSON.parse()
			
			//verify if marble meets trade requirements
			if closersMarble.Color != trades.OpenTrades[i].Want.Color || closersMarble.Size != trades.OpenTrades[i].Want.Size {
				msg := "marble in input does not meet trade requriements"
				fmt.Println(msg)
				return nil, errors.New(msg)
			}
			
			marble, e := findMarble4Trade(stub, trades.OpenTrades[i].User, args[4], size)			//find a marble that is suitable from opener
			if(e == nil){
				fmt.Println("! no errors, proceeding")

				t.set_user(stub, []string{args[2], trades.OpenTrades[i].User})						//change owner of selected marble, closer -> opener
				t.set_user(stub, []string{marble.Name, args[1]})									//change owner of selected marble, opener -> closer
			
				trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)		//remove trade
				jsonAsBytes, _ := json.Marshal(trades)
				err = stub.PutState(openTradesStr, jsonAsBytes)										//rewrite open orders
				if err != nil {
					return nil, err
				}
			}
		}
	}
	fmt.Println("- end close trade")
	return nil, nil
}

// ============================================================================================================================
// findMarble4Trade - look for a matching marble that this user owns and return it
// ============================================================================================================================
func findMarble4Trade(stub shim.ChaincodeStubInterface, user string, color string, size int )(m Marble, err error){
	var fail Marble;
	fmt.Println("- start find marble 4 trade")
	fmt.Println("looking for " + user + ", " + color + ", " + strconv.Itoa(size));

	//get the marble index
	marblesAsBytes, err := stub.GetState(marbleIndexStr)
	if err != nil {
		return fail, errors.New("Failed to get marble index")
	}
	var marbleIndex []string
	json.Unmarshal(marblesAsBytes, &marbleIndex)								//un stringify it aka JSON.parse()
	
	for i:= range marbleIndex{													//iter through all the marbles
		//fmt.Println("looking @ marble name: " + marbleIndex[i]);

		marbleAsBytes, err := stub.GetState(marbleIndex[i])						//grab this marble
		if err != nil {
			return fail, errors.New("Failed to get marble")
		}
		res := Marble{}
		json.Unmarshal(marbleAsBytes, &res)										//un stringify it aka JSON.parse()
		//fmt.Println("looking @ " + res.User + ", " + res.Color + ", " + strconv.Itoa(res.Size));
		
		//check for user && color && size
		if strings.ToLower(res.User) == strings.ToLower(user) && strings.ToLower(res.Color) == strings.ToLower(color) && res.Size == size{
			fmt.Println("found a marble: " + res.Name)
			fmt.Println("! end find marble 4 trade")
			return res, nil
		}
	}
	
	fmt.Println("- end find marble 4 trade - error")
	return fail, errors.New("Did not find marble to use in this trade")
}

// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

// ============================================================================================================================
// Remove Open Trade - close an open trade
// ============================================================================================================================
func (t *SimpleChaincode) remove_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	//	0
	//[data.id]
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	fmt.Println("- start remove trade")
	timestamp, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, errors.New("1st argument must be a numeric string")
	}
	
	//get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)																//un stringify it aka JSON.parse()
	
	for i := range trades.OpenTrades{																	//look for the trade
		//fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))
		if trades.OpenTrades[i].Timestamp == timestamp{
			fmt.Println("found the trade");
			trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)				//remove this trade
			jsonAsBytes, _ := json.Marshal(trades)
			err = stub.PutState(openTradesStr, jsonAsBytes)												//rewrite open orders
			if err != nil {
				return nil, err
			}
			break
		}
	}
	
	fmt.Println("- end remove trade")
	return nil, nil
}

// ============================================================================================================================
// Clean Up Open Trades - make sure open trades are still possible, remove choices that are no longer possible, remove trades that have no valid choices
// ============================================================================================================================
func cleanTrades(stub shim.ChaincodeStubInterface)(err error){
	var didWork = false
	fmt.Println("- start clean trades")
	
	//get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)																		//un stringify it aka JSON.parse()
	
	fmt.Println("# trades " + strconv.Itoa(len(trades.OpenTrades)))
	for i:=0; i<len(trades.OpenTrades); {																		//iter over all the known open trades
		fmt.Println(strconv.Itoa(i) + ": looking at trade " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10))
		
		fmt.Println("# options " + strconv.Itoa(len(trades.OpenTrades[i].Willing)))
		for x:=0; x<len(trades.OpenTrades[i].Willing); {														//find a marble that is suitable
			fmt.Println("! on next option " + strconv.Itoa(i) + ":" + strconv.Itoa(x))
			_, e := findMarble4Trade(stub, trades.OpenTrades[i].User, trades.OpenTrades[i].Willing[x].Color, trades.OpenTrades[i].Willing[x].Size)
			if(e != nil){
				fmt.Println("! errors with this option, removing option")
				didWork = true
				trades.OpenTrades[i].Willing = append(trades.OpenTrades[i].Willing[:x], trades.OpenTrades[i].Willing[x+1:]...)	//remove this option
				x--;
			}else{
				fmt.Println("! this option is fine")
			}
			
			x++
			fmt.Println("! x:" + strconv.Itoa(x))
			if x >= len(trades.OpenTrades[i].Willing) {														//things might have shifted, recalcuate
				break
			}
		}
		
		if len(trades.OpenTrades[i].Willing) == 0 {
			fmt.Println("! no more options for this trade, removing trade")
			didWork = true
			trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)					//remove this trade
			i--;
		}
		
		i++
		fmt.Println("! i:" + strconv.Itoa(i))
		if i >= len(trades.OpenTrades) {																	//things might have shifted, recalcuate
			break
		}
	}

	if(didWork){
		fmt.Println("! saving open trade changes")
		jsonAsBytes, _ := json.Marshal(trades)
		err = stub.PutState(openTradesStr, jsonAsBytes)														//rewrite open orders
		if err != nil {
			return err
		}
	}else{
		fmt.Println("! all open trades are fine")
	}

	fmt.Println("- end clean trades")
	return nil
}
