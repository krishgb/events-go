package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event map[string]map[int]map[int][]string

func serror(e error) {
	if e != nil {
		panic(e)
	}
}

func writeInFile(event string, decorate string, month int, day int, year string) {
	file, err := os.ReadFile("events.json")
	serror(err)

	events := Event{}
	json.Unmarshal(file, &events)
	

	if _, ok := events[event]; !ok {
		events[event] = map[int]map[int][]string{month: {day: {year, decorate}}}
	} else if _, ok = events[event][month]; !ok {
		events[event][month] = map[int][]string{day: {year, decorate}}
	} else {
		events[event][month][day] = append(events[event][month][day], year, decorate)
	}

	r, err := json.Marshal(events)
	serror(err)

	err = os.WriteFile("events.json", r, 0600)
	serror(err)
}

func readbyte() string {
	reader := bufio.NewReaderSize(os.Stdin, 64*1024)
	input, _, err := reader.ReadLine()
	serror(err)
	return strings.TrimSpace(string(input))
}

func decorate(event string) {
	fmt.Println("\n--------------------------------------------\n  ------------Added New Event-------------")
	fmt.Println("\n" + event)
	fmt.Println("---------------------------------------------")
}

func birthday() {
	year, month, day := time.Now().Date()

	fmt.Print("\nTell me whose :: ")
	name := strings.Title(readbyte())

	fmt.Print("\nAge(if you know or just skip) :: ")
	age := 0
	ageStr := ""

	if h := readbyte(); h != "" {
		age, _ = strconv.Atoi(h)
		if h[len(h)-1:] == "1" {
			ageStr = fmt.Sprintf("%dst", age)
		} else if h[len(h)-1:] == "2" {
			ageStr = fmt.Sprintf("%dnd", age)
		} else if h[len(h)-1:] == "3" {
			ageStr = fmt.Sprintf("%drd", age)
		} else if age == 0 {
			ageStr = ""
		} else {
			ageStr = fmt.Sprintf("%dth", age)
		}
	}

	decorateFormat := fmt.Sprintf("%s's %s Birthday", name, ageStr)
	writeInFile("birthdays", decorateFormat, int(month), day, strconv.Itoa(year))
	decorate(fmt.Sprintf("%02d-%02d-%d   %s", day, month, year, decorateFormat))

}

func otherEvent() {
	year, month, day := time.Now().Date()
	fmt.Print("\nWrite that down buddy ::  ")
	event := strings.TrimSpace(readbyte())
	event = strings.ToUpper(event[0:1]) + event[1:]
	writeInFile("events", event, int(month), day, strconv.Itoa(year))

	decorate(fmt.Sprintf("%02d-%02d-%d   %s", day, month, year, event))
}

func Creator() {
	fmt.Println("\n\n--------------A New Event? Great!-------------\nHere you go....\n ")

	fmt.Print("If '''BIRTHDAY🎂🎂🎂''' enter \"b\"  or press any key: ")
	b := readbyte()
	if strings.ToLower(b) == "b" {
		birthday()
	} else {
		otherEvent()
	}
}

func readFile(month int, day int, event string) []string {

	var birthday Event


	f, err := os.Open("events.json")

	if _, ok := err.(*fs.PathError); ok{
		f, _ = os.Create("events.json")
	}

	defer f.Close()

	file, err := os.ReadFile("events.json")

	serror(err)
	

	json.Unmarshal(file, &birthday)


	if (len(birthday) == 0) {return []string{}}
	return birthday[event][month][day]
}

func birthdays(year int, month int, day int) {
	b := readFile(month, day, "birthdays")
	if len(b) != 0 {
	for j, i := 1, 0; i <= len(b)/2; i, j = i+2, j+1 {

		y, err := strconv.Atoi(b[i])
		serror(err)
		
		str := strings.Split(b[i+1], " ")
		
		if strings.TrimSpace(str[len(str)-1]) != "" && strings.Contains(str[len(str) - 1], "'s"){

			age, err := strconv.Atoi(str[len(str)-2 : len(str)-1][0][0:2])
			age += year - y
			serror(err)

			ageInStr := fmt.Sprintf("%02s", strconv.Itoa(age))
			if ageInStr[1:] == "1" {
				ageInStr += "st"
			} else if ageInStr[1:] == "2" {
				ageInStr += "nd"
			} else if ageInStr[1:] == "3" {
				ageInStr += "rd"
			} else {
				ageInStr += "th"
			}
			str[len(str)-2 : len(str)-1][0] = ageInStr

		}
		fmt.Println(fmt.Sprintf("%d. ", j), strings.Join(str, " "))
	}
	}
}

func otherEvents(year int, month int, day int) {
	e := readFile(month, day, "events")
	for i, j := 0, 1; i < len(e)/2; i, j = i+2, j+1 {
		fmt.Printf("%d.  %s\n", j, e[i+1])
	}
}

func Teller() {
	year, month, day := time.Now().Date()
	fmt.Print("------------------\U0001f382\U0001f382\U0001f382 Birthdays\U0001f382\U0001f382\U0001f382------------------\n\n")
	birthdays(year, int(month), int(day))
	fmt.Print("\n\n\n\n------------------------- Events -------------------------\n")
	otherEvents(year, int(month), int(day))
	fmt.Print("\n\n\n")
}

func main() {
	Teller()
	fmt.Print("Want to add an event?(y/n): ")
	y := readbyte()
	if strings.ToLower(y) == "y" || strings.ToLower(y) == "yes" {
		Creator()
		fmt.Scanln()
	}

}
