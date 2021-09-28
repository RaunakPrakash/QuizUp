package customui

import (
	"context"
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
	"log"
	"quiz/dataManager"
	"quiz/model"
	"quiz/rabbitmq"
	"quiz/redisDriver"
	"quiz/userController"
	"strconv"
	"time"
)




func createView() {
	view := ui.AddWindow(0, 0, 30, 7, "Login dialog")
	view.SetPack(ui.Vertical)
	view.SetGaps(2, 1)
	view.SetPaddings(2, 2)

	frmOpts := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmOpts.SetPack(ui.Horizontal)


	frmCred := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmCred.SetPack(ui.Horizontal)
	frmCred.SetGaps(1, 0)
	ui.CreateLabel(frmCred, ui.AutoSize, ui.AutoSize, "Username", ui.Fixed)
	edUser := ui.CreateEditField(frmCred, 8, "", 1)
	ui.CreateLabel(frmCred, ui.AutoSize, ui.AutoSize, "Password", ui.Fixed)
	edPass := ui.CreateEditField(frmCred, 8, "", 1)
	edPass.SetPasswordMode(true)

	lbRes := ui.CreateLabel(view, ui.AutoSize, ui.AutoSize, "", ui.Fixed)

	frmButtons := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmButtons.SetPack(ui.Horizontal)
	btnDlg := ui.CreateButton(frmButtons, ui.AutoSize, 4, "Login", ui.Fixed)
	btnDsn := ui.CreateButton(frmButtons, ui.AutoSize, 4, "Signup", ui.Fixed)
	btnQuit := ui.CreateButton(frmButtons, ui.AutoSize, 4, "Quit", ui.Fixed)
	ui.CreateFrame(frmButtons, 1, 1, ui.BorderNone, 1)

	ui.ActivateControl(view, edUser)
	flag:=false
	controller := userController.Mongo{}
	redisController := redisDriver.RedisDriver{}

	btnDlg.OnClick(func(ev ui.Event) {
		if edUser.Title() =="" || edPass.Title()==""{
			lbRes.SetTitle("Invalid input")
		}else{

			redisController.Init(edUser.Title(),edPass.Title())


			ctx := context.Background()
			user, ok := controller.Get(ctx,"quiz","user",edUser.Title())
			if !ok {
				lbRes.SetTitle("Invalid user")
			}else {
				if !controller.CheckPasswordHash(user.Password,edPass.Title()){
					lbRes.SetTitle("Incorrect password or user")
				}else {
					err := controller.SetCollection(ctx, "quiz", "userinfo")
					if err != nil {
						return
					}
					data, err := controller.GetUserData(ctx,edUser.Title())
					if err != nil {
						return
					}
					redisController.Init(edUser.Title(),edPass.Title())
					redisController.PutUser(ctx)
					SetupStartUI(&data)
					ui.Stop()
				}
			}
		}

	})

	btnDsn.OnClick(func(ev ui.Event) {

		if edUser.Title() =="" || edPass.Title()==""{
			lbRes.SetTitle("Invalid input")
		}else{
			ctx := context.Background()
			_ , ok := controller.Get(ctx,"quiz","user",edUser.Title())
			if ok {
				lbRes.SetTitle("Username Taken")
			}else {
				err := controller.SetCollection(ctx, "quiz", "user")
				if err != nil {
					return
				}
				controller.SetUser(model.User{Username: edUser.Title(),Password: edPass.Title()})
				_, err = controller.Put(ctx)
				if err != nil {
					fmt.Println(err)
				}
				userData := model.Score{
					Username: edUser.Title(),
					Level: 1,
					Points: []int{},
					Date: time.Now(),
				}

				err = controller.SetCollection(ctx, "quiz", "userinfo")
				if err != nil {
					return
				}
				controller.SetUserData(userData)
				_, err = controller.PutUserData(ctx)
				if err != nil {
					return
				}
				flag = true
				edPass.SetTitle("")
				edUser.SetTitle("")

				lbRes.SetTitle("Successfully registered")

			}
		}
	})

	btnQuit.OnClick(func(ev ui.Event) {
		ui.Stop()
	})


}


func StartGame() {

	ui.InitLibrary()
	defer ui.DeinitLibrary()

	createView()
	ui.MainLoop()
}





func SetupStartUI(user *model.Score)  {

	ui.InitLibrary()
	defer ui.DeinitLibrary()
	r := redisDriver.RedisDriver{}

	r.Init(user.Username,"")
	view := ui.AddWindow(0, 0, ui.AutoSize, ui.AutoSize, "QUIZ GAME")
	frmChk := ui.CreateFrame(view, 8, 5, ui.BorderNone, ui.Fixed)
	frmChk.SetPack(ui.Vertical)
	frmChk.SetGaps(0,1)


	ui.CreateLabel(frmChk,50,5,"Welcome to QuizUp",ui.Fixed).SetAlign(ui.AlignCenter)
	ui.CreateLabel(frmChk,5,1,"Name : "+user.Username,ui.Fixed).SetAlign(ui.AlignLeft)

	ui.CreateLabel(frmChk,5,1,"Level : "+strconv.Itoa(user.Level-1),ui.Fixed).SetAlign(ui.AlignLeft)
	ui.CreateLabel(frmChk,5,2,"Total Points : "+strconv.Itoa(user.Total),ui.Fixed).SetAlign(ui.AlignLeft)

	btnStart := ui.CreateButton(frmChk, 4, 4, "START", ui.Fixed)
	btnStart.OnClick(func(ev ui.Event) {
		if r.GetUser(context.Background()) != "" {
			showLevels(user)
			ui.Stop()
		}
		StartGame()
		ui.Stop()
	})

	btnScore := ui.CreateButton(frmChk,15,4,"HIGH SCORE",ui.Fixed)
	btnScore.OnClick(func(ev ui.Event) {
		if r.GetUser(context.Background()) != "" {
			showScore(user)
			ui.Stop()
		}
		StartGame()
		ui.Stop()

	})

	btnReset := ui.CreateButton(frmChk, 4, 4, "RESET", ui.Fixed)
	btnReset.OnClick(func(ev ui.Event) {
		// delete Data on MongoDB
		if r.GetUser(context.Background()) != "" {
			user.Level = 1
			user.Points = []int{}
			user.Total = 0
			user.Date = time.Now()
			reset(user)
			resetDone(user)
		}
		StartGame()
		ui.Stop()
	})

	btnQuit := ui.CreateButton(frmChk, 15, 4, "Logout", ui.Fixed)
	btnQuit.OnClick(func(ev ui.Event) {
		r.Del(context.Background(),user.Username)
		StartGame()
		ui.Stop()
	})
	ui.MainLoop()
}

func resetDone(user *model.Score)  {
	SetupStartUI(user)
}

func startGame(user *model.Score,l int)  {

	question := dataManager.FetchQuestions(l)
	k:=0
	for i,v := range question.Questions {
		if !setupQuesAndAns(i,v) {
			break
		}
		k++
	}
	if playGame(user,k,l) {
		if l >= 5 {
			showOver(user)
		}else if k>=3{
			l++
			startGame(user,l)
		}else {
			startGame(user,l)
		}
	}else {
		SetupStartUI(user)
	}
}

func showOver(user *model.Score) {
	ui.InitLibrary()
	defer ui.DeinitLibrary()
	view := ui.AddWindow(0, 0, ui.AutoSize, ui.AutoSize, "QUIZ GAME")
	view.SetPack(ui.Vertical)
	ui.CreateLabel(view,60,10,"CONGRATULATIONS !\nYOU HAVE REACHED THE END OF THE GAME.",1)
	btnOpt2 := ui.CreateButton(view, 4, 4, "HOME", 1)
	btnOpt2.OnClick(func(ev ui.Event) {
		SetupStartUI(user)
		ui.Stop()
	})
	ui.MainLoop()
}



func setupQuesAndAns(i int, v model.Quiz) bool {
	ui.InitLibrary()
	defer ui.DeinitLibrary()
	flag := false
	view := ui.AddWindow(0, 0, ui.AutoSize, ui.AutoSize, "QUIZ GAME")
	view.SetPack(ui.Vertical)
	view.SetGaps(0,1)
	ui.CreateLabel(view,50,10,"Problem #"+strconv.Itoa(i+1)+" : "+v.Question,1)
	btnOpt1 := ui.CreateButton(view, 15, 4, v.Options[0], 1)
	btnOpt1.OnClick(func(ev ui.Event) {
		if btnOpt1.Title() == v.Answer {
			flag = true
		}
		ui.Stop()
	})
	btnOpt2 := ui.CreateButton(view, 15, 4, v.Options[1], 1)
	btnOpt2.OnClick(func(ev ui.Event) {
		if btnOpt2.Title() == v.Answer {
			flag = true
		}
		ui.Stop()
	})
	btnOpt3 := ui.CreateButton(view, 15, 4, v.Options[2], 1)
	btnOpt3.OnClick(func(ev ui.Event) {
		if btnOpt3.Title() == v.Answer {
			flag = true
		}
		ui.Stop()
	})
	btnOpt4 := ui.CreateButton(view, 15, 4, v.Options[3], 1)
	btnOpt4.OnClick(func(ev ui.Event) {
		if btnOpt4.Title() == v.Answer {
			flag = true
		}
		ui.Stop()
	})
	ui.MainLoop()
	return flag
}

func playGame(user *model.Score,k int,l int) bool {

	ui.InitLibrary()
	defer ui.DeinitLibrary()
	if k>=3{
		if user.Level == l {
			user.Level++
			if len(user.Points) == l {
				diff := k-user.Points[l-1]
				user.Points[l-1]=k
				user.Total += diff
			}else {
				user.Points = append(user.Points,k)
				user.Total += k
			}

			user.Date = time.Now()
			rabbitmq.PublishToRabbit(user)
		}
	}else {
		if len(user.Points) >= l{
			if user.Points[l-1] < k {
				diff:= k-user.Points[l-1]
				user.Points[l-1] = k
				user.Total += diff
				user.Date = time.Now()
				rabbitmq.PublishToRabbit(user)
			}
		}else {
			user.Points = append(user.Points,k)
			user.Total += k
			user.Date = time.Now()
			rabbitmq.PublishToRabbit(user)
		}
	}
	flag := false
	view := ui.AddWindow(0, 0, ui.AutoSize, ui.AutoSize, "QUIZ GAME")
	view.SetPack(ui.Vertical)
	view.SetGaps(0,1)
	ui.CreateLabel(view,50,10,"Your score is "+ strconv.Itoa(k) +" out of 5.",1)

	btnOpt1 := ui.CreateButton(view, 15, 4,"Continue", 1)
	btnOpt1.OnClick(func(ev ui.Event) {
		flag = true
		go ui.Stop()
	})
	btnOpt2 := ui.CreateButton(view, 15, 4, "Back", 1)
	btnOpt2.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})
	ui.MainLoop()
	return flag
}


func showScore(user *model.Score) {

	ui.InitLibrary()
	defer ui.DeinitLibrary()

	u:=userController.Mongo{}
	err := u.SetCollection(context.Background(), "quiz", "userinfo")
	if err != nil {
		fmt.Println(err)
		return
	}
	users:=u.GetHighScore(context.Background())
	view := ui.AddWindow(0, 0, 20, 10, "High Score")
	b := ui.CreateTableView(view, 30, 15, 1)
	view.SetPack(ui.Vertical)
	b.SetRowCount(len(users))
	view.SetAlign(ui.AlignCenter)
	cols := []ui.Column{
		{Title: "S.No.", Width: 5, Alignment: ui.AlignLeft},
		{Title: "Username", Width: 10, Alignment: ui.AlignCenter},
		{Title: "Level", Width: 6, Alignment: ui.AlignCenter},
		{Title: "Points", Width: 6, Alignment: ui.AlignCenter},
	}
	b.SetColumns(cols)
	colCount := len(cols)

	values := make([]string, 10*colCount)
	for r := 0; r < len(users); r++ {
		for c := 0; c < 1; c++ {
			values[r*colCount+0] = strconv.Itoa(r+1)
			values[r*colCount+1] = users[r].Username
			l:= strconv.Itoa(users[r].Level-1)
			if users[r].Level > 5 {
				l="5"
			}
			values[r*colCount+2] = l
			values[r*colCount+3] = strconv.Itoa(users[r].Total)
		}
	}

	b.OnDrawCell(func(info *ui.ColumnDrawInfo) {
		info.Text = values[info.Row*colCount+info.Col]
	})

	btnQuit := ui.CreateButton(view, 15, 4, "BACK", ui.Fixed)
	btnQuit.OnClick(func(ev ui.Event) {
		SetupStartUI(user)
		go ui.Stop()
	})


	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func showLevels(user *model.Score)  {
	ui.InitLibrary()
	defer ui.DeinitLibrary()
	view := ui.AddWindow(0, 0, 30, 10, "Level")
	view.SetPack(ui.Vertical)
	view.SetAlign(ui.AlignCenter)
	view.SetGaps(0,1)
	l1:= checkLevel(*user,1)
	l2:=checkLevel(*user,2)
	l3:=checkLevel(*user,3)
	l4:=checkLevel(*user,4)
	l5:=checkLevel(*user,5)
	btn1 := ui.CreateButton(view, 15, 4, "LEVEL 1 ("+l1+" pts)", ui.Fixed)
	btn1.OnClick(func(ev ui.Event) {
		startGame(user,1)
		ui.Stop()
	})
	btn2 := ui.CreateButton(view, 30, 4, "LEVEL 2 ("+l2+" pts)", ui.Fixed)
	btn2.OnClick(func(ev ui.Event) {
		if user.Level>=2 {

			startGame(user,2)
			ui.Stop()
		}
	})
	btn3 := ui.CreateButton(view, 30, 4, "LEVEL 3 ("+l3+" pts)", ui.Fixed)
	btn3.OnClick(func(ev ui.Event) {
		if user.Level>=3 {
			startGame(user,3)
			ui.Stop()
		}
	})
	btn4 := ui.CreateButton(view, 30, 4, "LEVEL 4 ("+l4+" pts)", ui.Fixed)
	btn4.OnClick(func(ev ui.Event) {
		if user.Level>=4 {
			startGame(user,4)
			ui.Stop()
		}
	})
	btn5 := ui.CreateButton(view, 30, 4, "LEVEL 5 ("+l5+" pts)", ui.Fixed)
	btn5.OnClick(func(ev ui.Event) {
		if user.Level>=5 {
			startGame(user,5)
			ui.Stop()
		}
	})
	buttonQuit := ui.CreateButton(view, 30, 4, "BACK", ui.Fixed)
	buttonQuit.OnClick(func(ev ui.Event) {
		SetupStartUI(user)
		ui.Stop()
	})
	ui.MainLoop()
}

func reset(score *model.Score)  {
	u := userController.Mongo{}
	ctx := context.Background()
	e := u.SetCollection(ctx,"quiz","userinfo")
	if e != nil {
		fmt.Println(e)
	}
	u.SetUserData(*score)
	_, e = u.Reset(ctx)
	if e != nil {
		log.Fatalln(e)
	}
}

func checkLevel(score model.Score,l int) string {
	if len(score.Points) >= l {
		return strconv.Itoa(score.Points[l-1])
	}else {
		return "0"
	}
}