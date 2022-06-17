package data

//package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	//"strconv"

	"github.com/ponbac/majority-wins/game"
)

type JSONQuestions struct {
	Questions []JSONQuestion `json:"questions"`
}

type JSONQuestion struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Choices     []string `json:"choices"`
	Reward      int      `json:"reward"`
}

const QUESTIONS_PATH = "questions.json"

func FetchQuestions() []*game.Question {
	var questions []*game.Question

	byteValue := readJson(QUESTIONS_PATH)

	// we initialize our Users array
	var jsonQuestions JSONQuestions

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &jsonQuestions)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for _, q := range jsonQuestions.Questions {
		questions = append(questions, &game.Question{Type: q.Type, Description: q.Description, Choices: q.Choices, Reward: q.Reward, Answers: make(map[*game.Player]int)})
	}

	return questions
}

func readJson(path string) []byte {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

// func main() {
// 	questions := FetchQuestions()
// 	for _, q := range questions {
// 		fmt.Println(q.Description)
// 	}
// }
