package model


type Questions struct {
	Questions []Quiz `json:"quiz"`
}

type Quiz struct {
	Question string `json:"q"`
	Answer string `json:"a"`
}