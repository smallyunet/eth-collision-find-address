package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

var totalFile = "total.txt"
var accountsFile = "accounts.txt"
var speedFile = "speed.txt"

func main() {
	msg := make(chan *big.Int)
	for i := 0; i < 4; i++ {
		go generateAccountJob(msg)
	}
	totalStr := readFile(totalFile)
	n := new(big.Int)
	total, ok := n.SetString(totalStr, 10)
	if !ok {
		total = big.NewInt(0)
	}
	lastTotal := total
	tick := time.Tick(1 * time.Hour)
	for {
		select {
		case <-tick:
			speed := total.Sub(total, lastTotal)
			lastTotal = total
			addresses, err := fileCountLine(accountsFile)
			if err != nil {
				log.Println(err)
			}
			text := fmt.Sprintf("Total: %d\nSpeed: %d/h\nAddresses: %d\n", total, speed, addresses)
			appendFile(speedFile, text)
			sendMsgText(text)
		case count := <-msg:
			total = total.Add(total, count)
			writeFile(totalFile, total.String())
		}
	}
}

func generateAccountJob(msg chan *big.Int) {
	count := big.NewInt(0)
	tick := time.Tick(1 * time.Minute)
	for {
		select {
		case <-tick:
			msg <- count
			count = big.NewInt(0)
		default:
			generateAccount()
			count = count.Add(count, big.NewInt(1))
		}
	}
}

func generateAccount() {
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Println(err)
	}
	privateKey := hex.EncodeToString(key.D.Bytes())
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	handleAccount(privateKey, address)
}

func checkAddress(address string) bool {
	if strings.HasPrefix(address, "0x8888") && strings.HasSuffix(address, "8888") {
		return true
	}
	return false
}

func handleAccount(privateKey string, address string) {
	if checkAddress(address) {
		log.Println("Found: ", privateKey, address)
		text := fmt.Sprintf("%s,%s\n", privateKey, address)
		appendFile(accountsFile, text)
	}
}
