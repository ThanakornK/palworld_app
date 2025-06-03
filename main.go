package main

import (
	"bufio"
	"fmt"
	"os"
	"palworld_tools/services/scrapper"
	"strings"
	"time"
)

func main() {
	var functionName string

	fmt.Println("Enter function: ")
	fmt.Println("1. Update Data\n")
	fmt.Print("Function number: ")
	reader := bufio.NewReader(os.Stdin)
	functionName, _ = reader.ReadString('\n')
	functionName = strings.TrimSpace(functionName)

	switch functionName {
	case "1":
		fmt.Println("Update Data")

		// Spinner characters
		spinner := []string{"|", "/", "-", "\\"}
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return // Stop the spinner loop when done is received
				case <-time.After(100 * time.Millisecond):
					// Print the spinner with a carriage return to overwrite the last character
					fmt.Print("\rLoading... ", spinner[time.Now().Unix()%4])
				}
			}
		}()

		err := updateData()
		if err != nil {
			errMessage := fmt.Sprintf("Error: %v", err)
			panic(errMessage)
		}

		done <- true // Send a signal to stop the spinner loop

	default:
		fmt.Println("Invalid function number")
	}

}

func updateData() error {
	err := scrapper.ScrapperPalInfo()
	if err != nil {
		return err
	}

	// Wait for 2 seconds before running the next function
	time.Sleep(5 * time.Second)

	err = scrapper.ScrapperPassiveSkill()
	if err != nil {
		return err
	}

	// Wait for 2 seconds before running the next function
	time.Sleep(5 * time.Second)

	err = scrapper.BestComboPassiveSkill()
	if err != nil {
		return err
	}

	return nil
}
