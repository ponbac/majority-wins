package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"

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

	var jsonQuestions JSONQuestions
	byteValue := readJson(QUESTIONS_PATH)
	json.Unmarshal(byteValue, &jsonQuestions)

	for _, q := range jsonQuestions.Questions {
		questions = append(questions, &game.Question{Type: q.Type, Description: q.Description, Choices: q.Choices, Reward: q.Reward, Answers: make(map[*game.Player]int)})
	}
	log.Debug().Msg("Fetched " + strconv.Itoa(len(questions)) + " questions")

	return questions
}

func readJson(path string) []byte {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Error().Err(err)
	}
	defer jsonFile.Close()

	// read our opened file as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}
