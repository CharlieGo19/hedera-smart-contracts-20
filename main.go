package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/dcb9/go-ethereum/accounts/abi"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Errorf(err.Error()))
	}

	accId, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_TNET_ID"))
	if err != nil {
		panic(err.Error())
	}

	privKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_TNET_PRIVATE_KEY"))
	if err != nil {
		panic(err.Error())
	}

	var tnetClient *hedera.Client = hedera.ClientForTestnet()
	tnetClient.SetOperator(accId, privKey) // Operator is the default account to pay txn fee.

	var contractAddress string = "0.0.34117828"
	SetMapFunc("Hedera", "1234567890", contractAddress, tnetClient)
	//GetMapFunc("Hedera", contractAddress, tnetClient)

}

func SetMapFunc(name, number, contractId string, client *hedera.Client) {
	toBigInt, _ := new(big.Int).SetString(number, 10)
	u256 := abi.U256(toBigInt)

	cid, err := hedera.ContractIDFromString(contractId)
	if err != nil {
		fmt.Println("Failed to parse ContractID: ", err.Error())
		os.Exit(1)
	}

	// Contract function params.
	cfp := hedera.NewContractFunctionParameters().AddString(name).AddUint256(u256)
	// Contract Transaction Setup.
	cet := hedera.NewContractExecuteTransaction()
	// Set Contract ID in Transaction.
	cet.SetContractID(cid)
	cet.SetGas(100000)
	cet.SetMaxTransactionFee(hedera.NewHbar(1))
	// Set the function to call & associated params.
	cet.SetFunction("setMobNo", cfp)
	// Execute & get response from Hedera.
	rsp, err := cet.Execute(client)
	if err != nil {
		fmt.Println("Error executing contract call: ", err.Error())
		os.Exit(1)
	}

	rspRcrd, err := rsp.GetRecord(client)
	if err != nil {
		fmt.Println("Error getting receipt: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Gas used in setting: %s was %d\n", name, rspRcrd.CallResult.GasUsed)
	fmt.Println("Transaction Status: ", rspRcrd.Receipt.Status)

}

func GetMapFunc(name, contractId string, client *hedera.Client) {
	cid, err := hedera.ContractIDFromString(contractId)
	if err != nil {
		fmt.Println("Failed to parse ContractID: ", err.Error())
		os.Exit(1)
	}

	// Contract function params.
	cfp := hedera.NewContractFunctionParameters().AddString(name)

	// Contract Query Setup.
	ccq := hedera.NewContractCallQuery()
	ccq.SetContractID(cid)
	ccq.SetGas(100000)
	ccq.SetQueryPayment(hedera.NewHbar(1))
	// Set function to be called & associated params.
	ccq.SetFunction("getMobNo", cfp)
	// Execute & get response from Hedera.
	rsp, err := ccq.Execute(client)
	if err != nil {
		fmt.Println("Error executing contract call: ", err.Error())
		os.Exit(1)
	}

	rspNo := new(big.Int)
	rspNo.SetBytes(rsp.GetUint256(0))

	fmt.Printf("We found that %v belonged to %s\n", rspNo, name)
	fmt.Printf("Gas used to query %s was: %d\n", name, rsp.GasUsed)

}
