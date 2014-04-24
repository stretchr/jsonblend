package main

import (
	"bufio"
	"fmt"
	"github.com/stretchr/jsonblend/blend"
	"os"
)

// Consume JSON data from STDIN and blend it to
// STDOUT

func main() {

	dest := make(map[string]interface{})

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		json := scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Printf("Encountered an error while reading from STDIN: %s\n", err.Error())
		} else {
			source, err := blend.JsonToMSI(json)
			if err != nil {
				fmt.Printf("Encountered an error while unmarshalling JSON from STDIN: %s\n", err.Error())
			} else {
				blend.Blend(source, dest)
			}
		}
	}

	destString, err := blend.MSIToJson(dest)
	if err != nil {
		fmt.Printf("Encountered an error while marshalling blended data to JSON: %s\n", err.Error())
	}

	fmt.Println(destString)

}
