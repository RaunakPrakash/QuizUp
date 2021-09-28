package model

import "time"

type Questions struct {
	Questions []Quiz `json:"quiz"`
}

type Quiz struct {
	Question string `json:"q"`
	Options []string `json:"options"`
	Answer string `json:"a"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Score struct {
	Username string `json:"username"`
	Level int `json:"level"`
	Points []int `json:"points"`
	Total int `json:"total"`
	Date time.Time `json:"date"`
}