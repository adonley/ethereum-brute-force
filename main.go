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
	"strconv"
	"fmt"
)

var partitions int = 6
var count int64
var oldCount int64
var addressesMap map[string]float64 = make(map[string]float64)

func main() {
	count = int64(0)

	loadAddresses()

	value, _ := time.ParseDuration("1s")
	checkTimer := time.NewTimer(value)
	go func() {
		for {
			select {
			case <-checkTimer.C:
				log.Printf("Checked: %d, Speed: %d per second", count, count - oldCount)
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

func loadAddresses() {
	count := int64(0)
	f, _ := os.Open("./balances.csv")

	r := csv.NewReader(bufio.NewReader(f))
	first := true
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if first {
			first = !first
			continue
		}

		count++

		if string(record[2]) == "False" {
			val, err := strconv.ParseFloat(record[1], 64)
			if err != nil {
				log.Printf("Error parsing: %s - %v", record[1], err.Error())
				continue
			}
			addressesMap[record[0]] = val
		}
	}
	log.Printf("Number of addresses loaded: %d", count)
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
		if amount, ok := addressesMap[address.Hex()]; ok {
			log.Printf("Found address with ETH balance, priv: %s, addr: %s, ammount %v", priv.D, address.Hex(), amount)
			writeToFound(fmt.Sprintf("Private: %s, Address: %s, Balance: %v\n", priv.D, address.Hex(), amount))
		}
		count++
	}
}

func writeToFound(text string) {
	foundFileName := "./found.txt"
	if _, err := os.Stat(foundFileName); os.IsNotExist(err) {
		_, _ = os.Create(foundFileName)
	}
	f, err := os.OpenFile(foundFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Printf(err.Error())
	}
	_, err = f.WriteString(text)
	if err != nil {
		log.Printf(err.Error())
	}
	f.Close()
}

func incrementPrivKey(privKey []byte) {
	for i := 31; i > 0; i-- {
		if privKey[i] + 1 == 255 {
			privKey[i] = 0
		} else {
			privKey[i] += 1
			break
		}
	}
}

func convertToPrivateKey(privKey []byte) (*ecdsa.PrivateKey) {
	return crypto.ToECDSA(privKey)
}