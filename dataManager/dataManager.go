package dataManager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"quiz/model"
	"strconv"
)


func FetchQuestions(k int) model.Questions {

	jsonFile,err:= os.Open("questions/quiz"+strconv.Itoa(k)+".json")
	if err!=nil {
		fmt.Println(err)
	}

	byteValue,_ := ioutil.ReadAll(jsonFile)

	var questions model.Questions

	err = json.Unmarshal(byteValue, &questions)
	if err != nil {
		fmt.Println(err)
	}

	return questions

}
