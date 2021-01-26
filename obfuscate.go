package main

import (
	"fmt"
	"math/rand"
	"time"
	"bufio"
	"os"
	"strings"
	"log"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// Construct Obfucated Payload func
type obfuscate func(string, map[string]string) string

// Struct to Build Payload
type payload struct {
	cmd string
	mapping map[string]string
	obfuscated obfuscate
}

func main() {
	// Seed Random
	rand.Seed(time.Now().UnixNano())

	// Map Characters, Get Input, Create Payload
	mapping := mapChars()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Command> ")
	cmd, _ := reader.ReadString('\n')
	cmd = strings.ReplaceAll(cmd, "\r\n", "")
	cmd = strings.ReplaceAll(cmd, "\n", "")
	p := payload{
		cmd: cmd,
		mapping: mapping,
		obfuscated: func(cmd string, mapping map[string]string) string {
			var initial, getLetters, command []string
			var initialPart, letterPart, commandPart string

			// Initial Setup
			initial = append(initial, fmt.Sprintf("set %s=set", mapping["set"]))
			initial = append(initial, fmt.Sprintf("%%%s%% %s= ", mapping["set"], mapping[" "]))
			initial = append(initial, fmt.Sprintf("%%%s%%%%%s%%%s==", mapping["set"], mapping[" "], mapping["="]))
			initialPart =  strings.Join(initial, "\n")

			// Letters Setup
			for i := range letters {
				if strings.Contains(cmd, string(letters[i])){
					getLetters = append(getLetters, fmt.Sprintf("%%%s%%%%%s%%%s%%%s%%%s", mapping["set"],mapping[" "], string(letters[i]),mapping["="], mapping[string(letters[i])]))
				}
			}
			letterPart = strings.Join(getLetters, " & ")

			// Command Setup
			for i := range cmd {
				command = append(command, fmt.Sprintf("%%%s%%", mapping[string(cmd[i])]))
			}
			commandPart = strings.Join(command, "")

			// Combine Initial + Letters + Command
			return fmt.Sprintf("%s\n%s\n%s", initialPart, letterPart, commandPart)
		},
	}

	// Constructed Payload
	constructed := p.obfuscated(p.cmd, p.mapping)

	// Output Information
	fmt.Printf("\n[+] Command: %s\n[+] Payload Size: %d Characters\n[+] Payload:\n%s\n\n[+] Writing to payload.bat...", p.cmd, len(constructed), constructed)

	// Write to a File
	f, err := os.Create("payload.bat")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(constructed)
	if err != nil {
		log.Fatal(err)
	}

}

// Take Length, Return Random String of Given Length
func genRandomString(length int) string {
	str := make([]rune, length)
	for i := range str {
		// Choose Random Letters
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}

// Map Characters A-Z, a-z, "=", " " and "set" to Unique Random Strings 
func mapChars() map[string]string {
	mapping := make(map[string]string)
	var randNum int
	var attempt string
	// A-Z, a-z
	for i := range letters {
		randNum = rand.Intn(6) + 4
		attempt = genRandomString(randNum)
		for {
			// Check if Value Exists
			if uniq := checkUniq(mapping, attempt); uniq {
				randNum = rand.Intn(6) + 4
				attempt = genRandomString(randNum)
				mapping[string(letters[i])] = attempt
				break
			}
		}
	}
	// Get Random String for "=", " ", and "set"
	other := [3]string{"=", " ", "set"}
	for i := range other {
		randNum = rand.Intn(6) + 4
		attempt = genRandomString(randNum)
		for {
			// Check if Value Exists
			if uniq := checkUniq(mapping, attempt); uniq {
				randNum = rand.Intn(6) + 4
				attempt = genRandomString(randNum)
				mapping[other[i]] = attempt
				break
			}
		}
	}

	return mapping
}

// Function to Check if Value is Unique in a Map
func checkUniq(mapping map[string]string, value string) bool {
	for i := range letters {
		if val, ok := mapping[string(letters[i])]; ok {
			if val == value {
				return false
			}
		}
	}
	return true
}
