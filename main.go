package main

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
	"encoding/csv"
	"bufio"
	"math/rand"
	"io"
	"time"
	"sync"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"encoding/json"
	"strconv"
	"runtime"
	"flag"
)

type concurrentMap struct {
	sync.Mutex
	addresses map[string]bool
}

var partitions = int(7)
var count int64
var oldCount int64
var addressesMap = concurrentMap { addresses: make(map[string]bool), }
var provider string
var shouldUpdate bool
var semaphoreChan = make(chan int, 5)
var goRequest = gorequest.New()

func main() {

	processFlags()

	runtime.GOMAXPROCS(runtime.NumCPU() + 1)

	count = int64(0)

	processedBlocks := loadAddresses()

	if shouldUpdate {
		updateAddressList(processedBlocks)
	}


	value, _ := time.ParseDuration("1s")
	checkTimer := time.NewTimer(value)
	go func() {
		for {
			select {
			case <-checkTimer.C:
				log.Printf("Checked: %d, Speed: %d per second", count, count-oldCount)
				oldCount = count
				checkTimer.Reset(value)
			}
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < partitions; i++ {
		wg.Add(1)
		addr := generateSeedAddress()
		log.Printf("Seed addr: %v\n", addr)
		go generateAddresses(addr)
	}
	wg.Wait()
}

func processFlags() {
	shouldUpdate = *flag.Bool("update", false, "a boolean on whether or not to update the list of addresses.")
	provider = *flag.String("provider", "http://localhost:8545", "http location of an ethereum node to use updating the address list")
}

func getBlockNumber() int64 {
	var requestBody = `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}`

	_, body, errs := goRequest.Post(provider).
		Send(requestBody).
		End()
	if errs != nil {
		fmt.Println(errs)
		os.Exit(1)
	}

	var raw map[string]interface{}

	err :=  json.Unmarshal([]byte(body), &raw)
	if err != nil {
		panic(err)
	}

	resultString := raw["result"].(string)
	blockNum, err := strconv.ParseInt(resultString[2:], 16, 64)
	if err != nil {
		panic(err)
	}
	return blockNum
}

func getBlock(blockNumber int64) []string {
	var requestBody = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x` + fmt.Sprintf("%x", blockNumber)  + `", true],"id":1}`
	var addressList []string
	var failed = true

	for failed {
		var goreq = gorequest.New()
		goreq.
			Post(provider).
			Send(requestBody).
			End(func(resp gorequest.Response, body string, errs []error) {
			if errs != nil {
				log.Print(errs)
				return
			} else {
				failed = false
			}

			var raw map[string]interface{}

			err :=  json.Unmarshal([]byte(body), &raw)
			if err != nil {
				log.Panic(err)
			}

			// Get the miner of the block
			addressList = append(addressList, raw["result"].(map[string]interface{})["miner"].(string))

			// Get the addresses of the transactions
			transactions := raw["result"].(map[string]interface{})["transactions"]
			for _, transaction := range transactions.([]interface{}) {
				if val, ok := transaction.(map[string]interface{})["to"].(string); ok {
					addressList = append(addressList, val)
				}
				if val, ok := transaction.(map[string]interface{})["from"].(string); ok {
					addressList = append(addressList, val)
				}
			}
		})


	}

	// TODO: Uncles

	return addressList
}

func updateAddressList(processedBlocks int64) {
	currentBlock := getBlockNumber()
	log.Print("Using {} to get new addresses", provider)
	var wg sync.WaitGroup

	for i := processedBlocks; i < currentBlock; i++  {
		wg.Add(1)
		semaphoreChan <- 1
		go func () {
			defer func() {
				// Release a slot
				<-semaphoreChan
			}()
			addressesList := getBlock(i)

			addressesMap.Lock()
			// Can this fail within this loop?
			for _, address := range addressesList {
				addressesMap.addresses[address] = true
			}
			addressesMap.Unlock()

			if i % 100 == 0 {
				log.Print("Gathering addresses from block: ", i)
			}
			wg.Done()
		} ()
	}

	wg.Wait()

	f, err := os.OpenFile("./balances.csv", os.O_WRONLY, 0600)
	defer f.Close()

	if err != nil {
		log.Panic(err)
	}

	f.WriteString(strconv.FormatInt(currentBlock,10) + "\n")

	addressesMap.Lock()
	for k := range addressesMap.addresses {
		f.WriteString(k + "\n")
	}
	addressesMap.Unlock()
}

func loadAddresses() int64 {
	processedBlocks := int64(0)
	count := int64(0)
	f, _ := os.Open("./balances.csv")
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	first := true
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if first {
			processedBlocks, err = strconv.ParseInt(record[0], 10, 64)
			if err != nil {
				log.Panic(err)
			}
			first = !first
			continue
		}

		count++

		addressesMap.Lock()
		addressesMap.addresses[record[0]] = true
		addressesMap.Unlock()
	}
	log.Printf("Number of addresses loaded: %d", count)

	return processedBlocks
}

func generateSeedAddress() []byte {
	privKey := make([]byte, 32)
	for i := 0; i < 32; i++ {
		privKey[i] = byte(rand.Intn(256))
	}
	return privKey
}

func generateAddresses(seedPrivKey []byte) {
	for ; ; {
		incrementPrivKey(seedPrivKey)
		priv := convertToPrivateKey(seedPrivKey)
		address := crypto.PubkeyToAddress(priv.PublicKey)
		addressesMap.Lock()
		if amount, ok := addressesMap.addresses[address.Hex()]; ok {
			log.Printf("Found address with ETH balance, priv: %s, addr: %s, ammount %v", priv.D, address.Hex(), amount)
			writeToFound(fmt.Sprintf("Private: %s, Address: %s, Balance: %v\n", priv.D, address.Hex(), amount))
		}
		addressesMap.Unlock()
		count++
	}
}

func writeToFound(text string) {
	foundFileName := "./found.txt"
	if _, err := os.Stat(foundFileName); os.IsNotExist(err) {
		_, _ = os.Create(foundFileName)
	}
	f, err := os.OpenFile(foundFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	if err != nil {
		log.Printf(err.Error())
	}
	_, err = f.WriteString(text)
	if err != nil {
		log.Printf(err.Error())
	}
}

func incrementPrivKey(privKey []byte) {
	for i := 31; i > 0; i-- {
		if privKey[i]+1 == 255 {
			privKey[i] = 0
		} else {
			privKey[i] += 1
			break
		}
	}
}

func convertToPrivateKey(privKey []byte) (*ecdsa.PrivateKey) {
	return crypto.ToECDSAUnsafe(privKey)
}
