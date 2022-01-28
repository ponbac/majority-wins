package game

type Question struct {
	Type        string
	Description string
	Answers     map[*Player]int
}

type JSONQuestion struct {
	Type        string        `json:"type"`
	Description string        `json:"description"`
	GroupOne    []*JSONPlayer `json:"group_one"`
	GroupTwo    []*JSONPlayer `json:"group_two"`
}

func (q *Question) ToJSONQuestion() *JSONQuestion {
	jsonQuestion := &JSONQuestion{Type: q.Type, Description: q.Description, GroupOne: []*JSONPlayer{}, GroupTwo: []*JSONPlayer{}}

	for player, value := range q.Answers {
		if value == 1 {
			jsonQuestion.GroupOne = append(jsonQuestion.GroupOne, player.ToJSONPlayer())
		} else {
			jsonQuestion.GroupTwo = append(jsonQuestion.GroupTwo, player.ToJSONPlayer())
		}
	}

	return jsonQuestion
}
