package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// Room :  Define the room structure, with 4 properties.  Structure tags are used by encoding/json library
type Room struct {
	Warranty   string `json:"warranty"`
	Rent  string `json:"rent"`
	Tenant string `json:"tenant"`
	Owner  string `json:"owner"`
}

type roomPrivateDetails struct {
	Owner string `json:"owner"`
	Price string `json:"price"`
}

// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("fabroom_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "queryRoom":
		return s.queryRoom(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "createRoom":
		return s.createRoom(APIstub, args)
	case "queryAllRooms":
		return s.queryAllRooms(APIstub)
	case "changeRoomOwner":
		return s.changeRoomOwner(APIstub, args)
	case "getHistoryForAsset":
		return s.getHistoryForAsset(APIstub, args)
	case "queryRoomsByOwner":
		return s.queryRoomsByOwner(APIstub, args)
	case "restictedMethod":
		return s.restictedMethod(APIstub, args)
	case "test":
		return s.test(APIstub, args)
	case "createPrivateRoom":
		return s.createPrivateRoom(APIstub, args)
	case "readPrivateRoom":
		return s.readPrivateRoom(APIstub, args)
	case "updatePrivateData":
		return s.updatePrivateData(APIstub, args)
	case "readRoomPrivateDetails":
		return s.readRoomPrivateDetails(APIstub, args)
	case "createPrivateRoomImplicitForOrg1":
		return s.createPrivateRoomImplicitForOrg1(APIstub, args)
	case "createPrivateRoomImplicitForOrg2":
		return s.createPrivateRoomImplicitForOrg2(APIstub, args)
	case "queryPrivateDataHash":
		return s.queryPrivateDataHash(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

func (s *SmartContract) queryRoom(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	roomAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(roomAsBytes)
}

func (s *SmartContract) readPrivateRoom(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	roomAsBytes, err := APIstub.GetPrivateData(args[0], args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[1] + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if roomAsBytes == nil {
		jsonResp := "{\"Error\":\"Room private details does not exist: " + args[1] + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(roomAsBytes)
}

func (s *SmartContract) readPrivateRoomIMpleciteForOrg1(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	roomAsBytes, _ := APIstub.GetPrivateData("_implicit_org_Org1MSP", args[0])
	return shim.Success(roomAsBytes)
}

func (s *SmartContract) readRoomPrivateDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	roomAsBytes, err := APIstub.GetPrivateData("collectionRoomPrivateDetails", args[0])

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[0] + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if roomAsBytes == nil {
		jsonResp := "{\"Error\":\"Marble private details does not exist: " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(roomAsBytes)
}

func (s *SmartContract) test(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	roomAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(roomAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	rooms := []Room{
		Room{Warranty: "850", Rent: "300", Tenant: "Michiel Verbeke", Owner: "Tomoko Satiko"},
		Room{Warranty: "700", Rent: "350", Tenant: "Milan Lemaire", Owner: "Tomoko Satiko"},
		Room{Warranty: "650", Rent: "500", Tenant: "Thibo Cuveele", Owner: "Tomoko Satiko"},
		Room{Warranty: "420", Rent: "365", Tenant: "Louis Mylle", Owner: "Tomoko Satiko"},
		Room{Warranty: "1200", Rent: "654", Tenant: "Ramzi Salhi", Owner: "Tomoko Satiko"},
	}

	i := 0
	for i < len(rooms) {
		roomAsBytes, _ := json.Marshal(rooms[i])
		APIstub.PutState("ROOM"+strconv.Itoa(i), roomAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createPrivateRoom(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	type roomTransientInput struct {
		Warranty  string `json:"arranty"` //the fieldtags are needed to keep case from bouncing around
		Rent string `json:"price"`
		Tenant string `json:"tenant"`
		Owner string `json:"owner"`
		Price string `json:"price"`
		Key   string `json:"key"`
	}
	if len(args) != 0 {
		return shim.Error("1111111----Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	logger.Infof("11111111111111111111111111")

	transMap, err := APIstub.GetTransient()
	if err != nil {
		return shim.Error("222222 -Error getting transient: " + err.Error())
	}

	roomDataAsBytes, ok := transMap["room"]
	if !ok {
		return shim.Error("room must be a key in the transient map")
	}
	logger.Infof("********************8   " + string(roomDataAsBytes))

	if len(roomDataAsBytes) == 0 {
		return shim.Error("333333 -marble value in the transient map must be a non-empty JSON string")
	}

	logger.Infof("2222222")

	var roomInput roomTransientInput
	err = json.Unmarshal(roomDataAsBytes, &roomInput)
	if err != nil {
		return shim.Error("44444 -Failed to decode JSON of: " + string(roomDataAsBytes) + "Error is : " + err.Error())
	}

	logger.Infof("3333")

	if len(roomInput.Key) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(roomInput.Warranty) == 0 {
		return shim.Error("color field must be a non-empty string")
	}
	if len(roomInput.Rent) == 0 {
		return shim.Error("model field must be a non-empty string")
	}
	if len(roomInput.Tenant) == 0 {
		return shim.Error("color field must be a non-empty string")
	}
	if len(roomInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}
	if len(roomInput.Price) == 0 {
		return shim.Error("price field must be a non-empty string")
	}

	logger.Infof("444444")

	// ==== Check if room already exists ====
	roomAsBytes, err := APIstub.GetPrivateData("collectionRooms", roomInput.Key)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	} else if roomAsBytes != nil {
		fmt.Println("This room already exists: " + roomInput.Key)
		return shim.Error("This room already exists: " + roomInput.Key)
	}

	logger.Infof("55555")

	var room = Room{Warranty: roomInput.Warranty, Rent: roomInput.Rent, Tenant: roomInput.Tenant, Owner: roomInput.Owner}

	roomAsBytes, err = json.Marshal(room)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutPrivateData("collectionRooms", roomInput.Key, roomAsBytes)
	if err != nil {
		logger.Infof("6666666")
		return shim.Error(err.Error())
	}

	roomPrivateDetails := &roomPrivateDetails{Owner: roomInput.Owner, Price: roomInput.Price}

	roomPrivateDetailsAsBytes, err := json.Marshal(roomPrivateDetails)
	if err != nil {
		logger.Infof("77777")
		return shim.Error(err.Error())
	}

	err = APIstub.PutPrivateData("collectionRoomPrivateDetails", roomInput.Key, roomPrivateDetailsAsBytes)
	if err != nil {
		logger.Infof("888888")
		return shim.Error(err.Error())
	}

	return shim.Success(roomAsBytes)
}

func (s *SmartContract) updatePrivateData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	type roomTransientInput struct {
		Owner string `json:"owner"`
		Price string `json:"price"`
		Key   string `json:"key"`
	}
	if len(args) != 0 {
		return shim.Error("1111111----Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	logger.Infof("11111111111111111111111111")

	transMap, err := APIstub.GetTransient()
	if err != nil {
		return shim.Error("222222 -Error getting transient: " + err.Error())
	}

	roomDataAsBytes, ok := transMap["room"]
	if !ok {
		return shim.Error("room must be a key in the transient map")
	}
	logger.Infof("********************8   " + string(roomDataAsBytes))

	if len(roomDataAsBytes) == 0 {
		return shim.Error("333333 -marble value in the transient map must be a non-empty JSON string")
	}

	logger.Infof("2222222")

	var roomInput roomTransientInput
	err = json.Unmarshal(roomDataAsBytes, &roomInput)
	if err != nil {
		return shim.Error("44444 -Failed to decode JSON of: " + string(roomDataAsBytes) + "Error is : " + err.Error())
	}

	roomPrivateDetails := &roomPrivateDetails{Owner: roomInput.Owner, Price: roomInput.Price}

	roomPrivateDetailsAsBytes, err := json.Marshal(roomPrivateDetails)
	if err != nil {
		logger.Infof("77777")
		return shim.Error(err.Error())
	}

	err = APIstub.PutPrivateData("collectionRoomPrivateDetails", roomInput.Key, roomPrivateDetailsAsBytes)
	if err != nil {
		logger.Infof("888888")
		return shim.Error(err.Error())
	}

	return shim.Success(roomPrivateDetailsAsBytes)

}

func (s *SmartContract) createRoom(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var room = Room{Warranty: args[1], Rent: args[2], Tenant: args[3], Owner: args[4]}

	roomAsBytes, _ := json.Marshal(room)
	APIstub.PutState(args[0], roomAsBytes)

	indexName := "owner~key"
	colorNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{room.Owner, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	APIstub.PutState(colorNameIndexKey, value)

	return shim.Success(roomAsBytes)
}

func (S *SmartContract) queryRoomsByOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	owner := args[0]

	ownerAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey("owner~key", []string{owner})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer ownerAndIdResultIterator.Close()

	var i int
	var id string

	var rooms []byte
	bArrayMemberAlreadyWritten := false

	rooms = append([]byte("["))

	for i = 0; ownerAndIdResultIterator.HasNext(); i++ {
		responseRange, err := ownerAndIdResultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		id = compositeKeyParts[1]
		assetAsBytes, err := APIstub.GetState(id)

		if bArrayMemberAlreadyWritten == true {
			newBytes := append([]byte(","), assetAsBytes...)
			rooms = append(rooms, newBytes...)

		} else {
			// newBytes := append([]byte(","), roomsAsBytes...)
			rooms = append(rooms, assetAsBytes...)
		}

		fmt.Printf("Found a asset for index : %s asset id : ", objectType, compositeKeyParts[0], compositeKeyParts[1])
		bArrayMemberAlreadyWritten = true

	}

	rooms = append(rooms, []byte("]")...)

	return shim.Success(rooms)
}

func (s *SmartContract) queryAllRooms(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "ROOM0"
	endKey := "ROOM999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllRooms:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) restictedMethod(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// get an ID for the client which is guaranteed to be unique within the MSP
	//id, err := cid.GetID(APIstub) -

	// get the MSP ID of the client's identity
	//mspid, err := cid.GetMSPID(APIstub) -

	// get the value of the attribute
	//val, ok, err := cid.GetAttributeValue(APIstub, "attr1") -

	// get the X509 certificate of the client, or nil if the client's identity was not based on an X509 certificate
	//cert, err := cid.GetX509Certificate(APIstub) -

	val, ok, err := cid.GetAttributeValue(APIstub, "role")
	if err != nil {
		// There was an error trying to retrieve the attribute
		shim.Error("Error while retriving attributes")
	}
	if !ok {
		// The client identity does not possess the attribute
		shim.Error("Client identity doesnot posses the attribute")
	}
	// Do something with the value of 'val'
	if val != "approver" {
		fmt.Println("Attribute role: " + val)
		return shim.Error("Only user with role as APPROVER have access this method!")
	} else {
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
		}

		roomAsBytes, _ := APIstub.GetState(args[0])
		return shim.Success(roomAsBytes)
	}

}

func (s *SmartContract) changeRoomOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	roomAsBytes, _ := APIstub.GetState(args[0])
	room := Room{}

	json.Unmarshal(roomAsBytes, &room)
	room.Owner = args[1]

	roomAsBytes, _ = json.Marshal(room)
	APIstub.PutState(args[0], roomAsBytes)

	return shim.Success(roomAsBytes)
}

func (t *SmartContract) getHistoryForAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	roomName := args[0]

	resultsIterator, err := stub.GetHistoryForKey(roomName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) createPrivateRoomImplicitForOrg1(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect arguments. Expecting 5 arguments")
	}

	var room = Room{Warranty: args[1], Rent: args[2], Tenant: args[3], Owner: args[4]}

	roomAsBytes, _ := json.Marshal(room)

	err := APIstub.PutPrivateData("_implicit_org_Org1MSP", args[0], roomAsBytes)
	if err != nil {
		return shim.Error("Failed to add asset: " + args[0])
	}
	return shim.Success(roomAsBytes)
}

func (s *SmartContract) createPrivateRoomImplicitForOrg2(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect arguments. Expecting 5 arguments")
	}

	var room = Room{Warranty: args[1], Rent: args[2], Tenant: args[3], Owner: args[4]}

	roomAsBytes, _ := json.Marshal(room)
	APIstub.PutState(args[0], roomAsBytes)

	err := APIstub.PutPrivateData("_implicit_org_Org2MSP", args[0], roomAsBytes)
	if err != nil {
		return shim.Error("Failed to add asset: " + args[0])
	}
	return shim.Success(roomAsBytes)
}

func (s *SmartContract) queryPrivateDataHash(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	roomAsBytes, _ := APIstub.GetPrivateDataHash(args[0], args[1])
	return shim.Success(roomAsBytes)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
