package dataManager

import (
	"quiz/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)


func FetchQuestions() model.Questions {

	jsonFile,err:= os.Open("quiz.json")
	if err!=nil {
		fmt.Println(err)
	}

	byteValue,_ := ioutil.ReadAll(jsonFile)

	var questions model.Questions

	json.Unmarshal(byteValue,&questions)

	return questions

	//for i,v :=  range questions.Questions {
	//	ans:=""
	//	fmt.Printf("Problem#%d ",i+1)
	//	fmt.Println(string(v.Question))
	//
	//	fmt.Scanf("%s",&ans)
	//	if ans==v.Answer {
	//		fmt.Println()
	//		fmt.Println("Correct!")
	//		fmt.Println()
	//	}else {
	//		fmt.Println()
	//		fmt.Println("Incorrect")
	//		fmt.Println()
	//		fmt.Printf("Do you want to play again? (y/n)")
	//		a:=""
	//		fmt.Scanf("%s",&a)
	//		if a=="y"{
	//			//startGame()
	//		}else if a=="n"{
	//			break
	//		}
	//	}
	//}
}
