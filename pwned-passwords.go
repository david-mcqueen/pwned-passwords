package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, password := range os.Args[1:] {
		hashedPwrd := hashInput(password)

		head, tail := hashedPwrd[:5], hashedPwrd[5:]
		url := "https://api.pwnedpasswords.com/range/" + head

		resp, err := http.Get(url)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get response from api")
			os.Exit(1)
		}

		if resp.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, "Got response other than 200. StatusCode: %d\r\n", resp.StatusCode)
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading response %s: %v\n", url, err)
			os.Exit(1)
		}

		bodyString := string(b)

		hashLines := strings.Split(bodyString, "\r\n")

		matchFound := false

		for _, line := range hashLines {
			hash := strings.Split(line, ":")
			if strings.EqualFold(hash[0], tail){
				fmt.Fprintf(os.Stdout, "FOUND! This password appears %s times\r\n", hash[1])
				matchFound = true
			}
		}

		if !matchFound {
			fmt.Println("Password doesn't appear in list")
		}
	}
}

func hashInput(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}