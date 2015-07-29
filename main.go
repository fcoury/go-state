package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Executable interface {
	Execute(*User, string) Executable
}

type InitialMenu struct{}

func (a *InitialMenu) Execute(user *User, msg string) Executable {
	menu := WaitingMenu{
		Options: map[string]string{
			"1": "News",
			"2": "Sports",
		},
	}
	user.Send(menu.Prompt())
	return &menu
}

type WaitingMenu struct {
	Options map[string]string
}

func (a *WaitingMenu) Execute(user *User, msg string) Executable {
	opt := a.Options[msg]
	if opt == "" {
		user.Send("Unknown option: " + msg)
		user.Send(a.Prompt())
		return a
	}
	user.Send("You chose: " + opt)
	return nil
}

func (a *WaitingMenu) Prompt() string {
	s := "Menu:\n"
	for k := range a.Options {
		s += fmt.Sprintf("%s - %s\n", k, a.Options[k])
	}
	return s
}

type User struct {
	Name          string
	CurrentAction Executable
}

func (u *User) Send(msg string) {
	fmt.Printf("[%s] -> %s\n", u.Name, msg)
}

func (u *User) Execute(msg string) {
	if u.CurrentAction == nil {
		cmd := strings.Split(msg, " ")[0]
		fmt.Println("   Cmd:", cmd)
		fmt.Println("Action:", actions[cmd])
		action := actions[cmd]
		if action == nil {
			return
		}
		u.CurrentAction = action
	}

	u.CurrentAction = u.CurrentAction.Execute(u, msg)
}

func (u *User) Handle(msg string) {
	fmt.Printf("[%s] <- %s\n", u.Name, msg)
	u.Execute(msg)
}

type Question struct {
	Question string
	Answer   string
}

type Quiz struct {
	Questions []Question
}

func (a *Quiz) Execute(user *User, msg string) Executable {
	a.Questions = []Question{
		Question{"How old was Jesus?", "33"},
		Question{"What is the capital of USA?", "Washington"},
		Question{"Who is Brazil's president?", "Dilma"},
	}

	quiz := AskQuiz{Quiz: *a}
	return quiz.Execute(user, msg)
}

type AskQuiz struct {
	Quiz
	currentQuestion int
	score           int
}

func (a *AskQuiz) Execute(user *User, msg string) Executable {
	if a.currentQuestion > 0 {
		prev := a.Questions[a.currentQuestion-1]
		if msg == prev.Answer {
			a.score += 10
			user.Send(fmt.Sprintf("Great, you got it!\nScore: %d\n", a.score))
		} else {
			user.Send("That's not the right answer, sorry!\n")
		}
	}

	if a.currentQuestion >= len(a.Questions) {
		user.Send(fmt.Sprintf("Your final score: %d/%d\n", a.score, 10*len(a.Questions)))
		return nil
	}

	user.Send(a.Questions[a.currentQuestion].Question)
	a.currentQuestion++

	return a
}

var actions map[string]Executable

func initActions() {
	actions = map[string]Executable{}
	actions["menu"] = &InitialMenu{}
	actions["quiz"] = &Quiz{}
}

func main() {
	initActions()
	user := User{Name: "fcoury"}
	for true {
		reader := bufio.NewReader(os.Stdin)
		val := "<nil>"
		if user.CurrentAction != nil {
			val = reflect.TypeOf(user.CurrentAction).Elem().Name()
		}
		fmt.Printf("%s# ", val)
		text, _ := reader.ReadString('\n')
		user.Handle(strings.TrimSpace(text))
	}
}
