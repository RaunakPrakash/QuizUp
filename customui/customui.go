package customui

import (
	ui "github.com/VladimirMarkelov/clui"
	"quiz/dataManager"
	"quiz/model"
	"strconv"
)

func SetupStartUI()  {

	ui.InitLibrary()
	defer ui.DeinitLibrary()
	view := ui.AddWindow(0, 0, ui.AutoSize, ui.AutoSize, "QUIZ GAME")
	view.SetPack(ui.Vertical)
	ui.CreateLabel(view,50,10,"Welcome to Maths Quiz Game",1).SetAlign(ui.AlignCenter)

	btnStart := ui.CreateButton(view, 4, 4, "START", 1)
	btnStart.OnClick(func(ev ui.Event) {
		startGame()
		go ui.Stop()
	})
	btnScore := ui.CreateButton(view,15,4,"HIGH SCORE",1)
	btnScore.OnClick(func(ev ui.Event) {
		//show()
	})
	btnQuit := ui.CreateButton(view, 15, 4, "EXIT", 1)
	btnQuit.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})
	ui.MainLoop()
}

func startGame()  {
	question := dataManager.FetchQuestions()
	k:=0
	for i,v := range question.Questions {
		if !setupQuseAndAns(i,v) {
			break
		}
		k++
	}
	if playGame(k) {
		startGame()
	}else {
		SetupStartUI()
	}
}



func setupQuseAndAns(i int, v model.Quiz) bool {
	ui.InitLibrary()
	defer ui.DeinitLibrary()
	flag := false
	view := ui.AddWindow(0, 0, ui.AutoSize, ui.AutoSize, "QUIZ GAME")
	view.SetPack(ui.Vertical)
	ui.CreateLabel(view,50,10,"Problem #"+strconv.Itoa(i+1)+" : "+v.Question,1)
	ans := v.Answer
	intAns, _ := strconv.Atoi(ans)
	btnOptn1 := ui.CreateButton(view, 15, 4, strconv.Itoa(intAns+1), 1)
	btnOptn1.OnClick(func(ev ui.Event) {
		if btnOptn1.Title() == v.Answer {
			flag = true
		}
		go ui.Stop()
	})
	btnOptn2 := ui.CreateButton(view, 15, 4, ans, 1)
	btnOptn2.OnClick(func(ev ui.Event) {
		if btnOptn2.Title() == v.Answer {
			flag = true
		}
		go ui.Stop()
	})
	btnOptn3 := ui.CreateButton(view, 15, 4, strconv.Itoa(intAns-1), 1)
	btnOptn3.OnClick(func(ev ui.Event) {
		if btnOptn3.Title() == v.Answer {
			flag = true
		}
		go ui.Stop()
	})
	btnOptn4 := ui.CreateButton(view, 15, 4, strconv.Itoa(intAns-2), 1)
	btnOptn4.OnClick(func(ev ui.Event) {
		if btnOptn4.Title() == v.Answer {
			flag = true
		}
		go ui.Stop()
	})
	ui.MainLoop()
	return flag
}

func playGame(k int) bool {
	ui.InitLibrary()
	defer ui.DeinitLibrary()
	flag := false
	view := ui.AddWindow(0, 0, ui.AutoSize, ui.AutoSize, "QUIZ GAME")
	view.SetPack(ui.Vertical)
	ui.CreateLabel(view,50,10,"Your score is "+ strconv.Itoa(k) +" out of 3." +
		"\n\n\nDo you want to play again",1)

	btnOptn1 := ui.CreateButton(view, 15, 4,"Yes", 1)
	btnOptn1.OnClick(func(ev ui.Event) {
		flag = true
		go ui.Stop()
	})
	btnOptn2 := ui.CreateButton(view, 15, 4, "No", 1)
	btnOptn2.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})
	ui.MainLoop()
	return flag
}
