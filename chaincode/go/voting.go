/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// Save each voter's data
type Voter struct {
	//ID        string `json:"user_id"`
	Vote_choice  string `json:"vote_choice"`
	Voted bool `json:"voted"`
}
// Save each voting item's data
type VoteItem struct {
	Name     string `json:"name"`
	Votesnum int    `json:"votesnum"`
}

// QueryResult structures used for handling result of query
type QueryResult struct {
	Key    string `json:"KeyVoter"`
	Record *Voter
}
type QueryResultItem struct {
	Key    string `json:"Key"`
	Record *VoteItem
}

// Create a voter given an user ID
func (s *SmartContract) CreateVoter(ctx contractapi.TransactionContextInterface) string {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "Failed to get user's ID"
	}
	voter_exists, err := s.QueryVoter(ctx, clientID)
	if err == nil {
		return "You have registered as a voter already!"
	}
	if voter_exists != nil {
		return "You have registered as a voter already!"
	}
	voter := Voter{
		Voted:   false,
		Vote_choice: "none",
	}
	voterAsBytes, err := json.Marshal(voter)
	if err != nil {
		return "failed to marshal voter object"

	}

	err = ctx.GetStub().PutState(clientID, voterAsBytes)
	if err != nil {
		return "failed to put state in public data"
	}
	return "CreateVoter ok"
}

// InitLedger with simple data
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	items := []VoteItem{
		VoteItem{Name: "Choice1", Votesnum: 0},
		VoteItem{Name: "Choice2", Votesnum: 6},
		VoteItem{Name: "Choice3", Votesnum: 30},
	}

	for i, item := range items {
		itemAsBytes, _ := json.Marshal(item)
		err := ctx.GetStub().PutState(strconv.Itoa(i), itemAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}
	return nil
}

// Get an item from world state
func (s *SmartContract) QueryItem(ctx contractapi.TransactionContextInterface, ID string) (*VoteItem, error) {
	itemAsBytes, err := ctx.GetStub().GetState(ID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if itemAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", ID)
	}

	item := new(VoteItem)
	_ = json.Unmarshal(itemAsBytes, item)

	return item, nil
}

// Get a voter from world state
func (s *SmartContract) QueryVoter(ctx contractapi.TransactionContextInterface, VoterID string) (*Voter, error) {
	voterAsBytes, err := ctx.GetStub().GetState(VoterID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if voterAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", VoterID)
	}

	voter := new(Voter)
	_ = json.Unmarshal(voterAsBytes, voter)

	return voter, nil
}

// Get all items from the world state
func (s *SmartContract) QueryAllItems(ctx contractapi.TransactionContextInterface) ([]QueryResultItem, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResultItem{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		item := new(VoteItem)
		_ = json.Unmarshal(queryResponse.Value, item)
		//check if object is VoteItem or Voter
		if item.Name != ""{
		queryResult := QueryResultItem{Key: queryResponse.Key, Record: item}
		results = append(results, queryResult)
		}
	}

	return results, nil
}

// Vote for item
func (s *SmartContract) Vote(ctx contractapi.TransactionContextInterface, ID string) string {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "Failed to get user's ID"
	}

	voter, err := s.QueryVoter(ctx, clientID)
	if err != nil {
		//return fmt.Errorf("Voter does not exist! Register as a voter first! %s", clientID)
		return "Voter does not exist! Register as a voter first!"
	}
	item, err := s.QueryItem(ctx, ID)
	if err != nil {
		return "Item does not exist! Check the item's ID again!"
	}

	if !voter.Voted {
		item.Votesnum += 1
		itemAsBytes, _ := json.Marshal(item)
		ctx.GetStub().PutState(ID, itemAsBytes)
		return "Voted successfully"
	}

	return "You have voted already!"
}

// Update the voter struct after voting in a separate function 
// to avoid unexpected behaviour when posting to the world state
func (s *SmartContract) UpdateVoter(ctx contractapi.TransactionContextInterface, ID string) string {
	clientID, err := ctx.GetClientIdentity().GetID()
	voter, err := s.QueryVoter(ctx, clientID)
	if err != nil {
		 return "err"
	}
	voter.Voted = true
	voter.Vote_choice = ID
	voterAsBytes, _ := json.Marshal(voter)
	ctx.GetStub().PutState(clientID, voterAsBytes)
	return "UpdateVoter successful"
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}



func (s *SmartContract) QueryAllUsers(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		voter := new(Voter)
		_ = json.Unmarshal(queryResponse.Value, voter)
		if voter.Vote_choice != ""{
			queryResult := QueryResult{Key: queryResponse.Key, Record: voter}
			results = append(results, queryResult)
		}
		
	
	}

	return results, nil
}