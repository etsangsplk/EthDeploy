package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var TokenABI = `[
    {
      "constant": false,
      "inputs": [],
      "name": "kill",
      "outputs": [],
      "payable": false,
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "index",
          "type": "uint256"
        }
      ],
      "name": "remove",
      "outputs": [],
      "payable": false,
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "sshPublicKeys",
      "outputs": [
        {
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "sshkey",
          "type": "string"
        }
      ],
      "name": "addSshKey",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [],
      "name": "mortal",
      "outputs": [],
      "payable": false,
      "type": "function"
    },
  {
    "inputs":[
      {
        "name":"someinput",
        "type":"uint256"
      },
      {
        "name":"tokenName",
        "type":"string"
      },
      {
        "name":"decimalUnits",
        "type":"uint8"
      },
      {
        "name":"tokenSymbol",
        "type":"string"
      }
    ],
    "type":"constructor"
  }

  ]`

/*
   {
     "inputs": [],
     "payable": false,
     "type": "constructor"
   }*/
var TokenBin = `0x6060604052341561000c57fe5b5b60018054600160a060020a03191633600160a060020a03161790555b5b6105a1806100396000396000f300606060405263ffffffff60e060020a60003504166341c0e1b5811461004d5780634cc822151461005f57806353918ad114610074578063690252c814610107578063f1eae25c14610171575bfe5b341561005557fe5b61005d610183565b005b341561006757fe5b61005d6004356101ab565b005b341561007c57fe5b610087600435610281565b6040805160208082528351818301528351919283929083019185019080838382156100cd575b8051825260208311156100cd57601f1990920191602091820191016100ad565b505050905090810190601f1680156100f95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561010f57fe5b61015d600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284375094965061032b95505050505050565b604080519115158252519081900360200190f35b341561017957fe5b61005d61036b565b005b60015433600160a060020a03908116911614156101a857600154600160a060020a0316ff5b5b565b60008054821015806101bd5750600082105b156101c75761027d565b50805b6000546000190181101561023e5760008054600183019081106101e957fe5b906000526020600020900160005b50600080548390811061020657fe5b906000526020600020900160005b508154610234929060026000196101006001841615020190911604610389565b505b6001016101ca565b60008054600019810190811061025057fe5b906000526020600020900160005b610268919061040f565b600080549061027b906000198301610457565b505b5050565b600080548290811061028f57fe5b906000526020600020900160005b508054604080516020601f6002600019610100600188161502019095169490940493840181900481028201810190925282815293508301828280156103235780601f106102f857610100808354040283529160200191610323565b820191906000526020600020905b81548152906001019060200180831161030657829003601f168201915b505050505081565b6000805481906001810161033f8382610457565b916000526020600020900160005b50835161035f919060208601906104ab565b5050600190505b919050565b60018054600160a060020a03191633600160a060020a03161790555b565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106103c257805485556103fe565b828001600101855582156103fe57600052602060002091601f016020900482015b828111156103fe5782548255916001019190600101906103e3565b5b5061040b92915061052a565b5090565b50805460018160011615610100020316600290046000825580601f106104355750610453565b601f016020900490600052602060002090810190610453919061052a565b5b50565b81548183558181151161027b5760008381526020902061027b91810190830161054b565b5b505050565b81548183558181151161027b5760008381526020902061027b91810190830161054b565b5b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106104ec57805160ff19168380011785556103fe565b828001600101855582156103fe579182015b828111156103fe5782518255916020019190600101906104fe565b5b5061040b92915061052a565b5090565b61054891905b8082111561040b5760008155600101610530565b5090565b90565b61054891905b8082111561040b576000610565828261040f565b50600101610551565b5090565b905600a165627a7a72305820f6628617348853643d8514ec291cc930cf9cee748ab26ab56b917a9a115d60ed0029`

const keydata = `bddcbcf81af22330478ab2977b35e10dacb5ad38eae8263ca03e3b0ae8e009c7`

func main() {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	keyBytes := common.FromHex(keydata)
	key := crypto.ToECDSAUnsafe(keyBytes)

	auth := bind.NewKeyedTransactor(key)

	// Deploy a new awesome contract for the binding demo
	//	address, tx, token, err := DeployToken(auth, conn, new(big.Int), "Contracts in Go!!!", 0, "Go!")
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	inputs := parsed.Constructor.Inputs
	fmt.Printf("inputs -%v", inputs)

	var dataInputs []interface{} = make([]interface{}, 0)

	//Default all the parameters
	for _, i := range inputs {
		t := reflect.Zero(i.Type.Type)
		v := t.Interface()
		fmt.Printf("param name-%s -%s\n", i.Name, t.String())
		if t.Kind() == reflect.Ptr && t.IsNil() {
			fmt.Printf("its a pointer!\n")
			elem := reflect.TypeOf(v).Elem()
			v2 := reflect.New(elem)
			fmt.Printf("param name-v2-%s  %v-%v\n", i.Name, v2.Type(), v2.Pointer())
			//			v = (*big.Int)(unsafe.Pointer(v2.Pointer())
			v = v2.Interface()
		}
		dataInputs = append(dataInputs, v)
	}
	fmt.Printf("dataInputs -%v\n", dataInputs)

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenBin), conn, dataInputs...)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())
	fmt.Printf("Contract Object-%v", contract)

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P
	/*
		name, err := token.Name(&bind.CallOpts{Pending: true})
		if err != nil {
			log.Fatalf("Failed to retrieve pending name: %v", err)
		}
		fmt.Println("Pending name:", name)
	*/
}
