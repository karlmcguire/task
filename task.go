package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	PATH      = "/Users/karl/.tasks/"
	DATA_PATH = PATH + "data.json"
	MARK_PATH = PATH + "readme.md"
)

type Task struct {
	Note     string        `json:"note"`
	Started  time.Time     `json:"started"`
	Duration time.Duration `json:"duration"`
}

func load(path string) ([]*Task, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			return make([]*Task, 0), nil
		} else {
			return nil, err
		}
	}
	tasks := make([]*Task, 0)
	return tasks, json.Unmarshal(data, &tasks)
}

func save(path string, tasks []*Task) error {
	data, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func mark(path string, tasks []*Task) error {
	page := ""
	prevDate := ""
	for i, task := range tasks {
		date := task.Started.Format("01-02-06")
		if date != prevDate {
			prevDate = date
			page += "\n"
			page += `## ` + date + "\n\n"
			page += `| start | duration | notes |` + "\n"
			page += `|:-----:|:--------:|:------|` + "\n"
		}
		page += fmt.Sprintf("| %s ", task.Started.Format("3:04pm"))
		if i == len(tasks)-1 {
			page += "| - "
		} else {
			page += fmt.Sprintf("| %d ", int(task.Duration.Minutes()))
		}
		page += fmt.Sprintf("| %s |\n", task.Note)
	}
	return ioutil.WriteFile(path, []byte(page), 0644)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: task [-pull -push] [notes...]")
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-pull":
			cmd := exec.Command("git", "-C", PATH, "pull")
			if err := cmd.Run(); err != nil {
				panic(err)
			}
			return
		case "-push":
			// add
			cmd := exec.Command("git", "-C", PATH, "add", "-A")
			if err := cmd.Run(); err != nil {
				panic(err)
			}
			// commit
			cmd = exec.Command("git", "-C", PATH, "commit", "-m",
				fmt.Sprintf(`"%v"`, time.Now().Format("01-02-06 3:04pm")))
			if err := cmd.Run(); err != nil {
				panic(err)
			}
			// push
			cmd = exec.Command("git", "-C", PATH, "push")
			if err := cmd.Run(); err != nil {
				panic(err)
			}
			return
		}
	}

	// load data.json
	tasks, err := load(DATA_PATH)
	if err != nil {
		panic(err)
	}

	tasks[len(tasks)-1].Duration = time.Since(tasks[len(tasks)-1].Started)
	tasks = append(tasks, &Task{
		Note:    strings.Join(os.Args[1:], " "),
		Started: time.Now(),
	})

	// save data.json
	if err = save(DATA_PATH, tasks); err != nil {
		panic(err)
	}

	// save readme.md
	if err = mark(MARK_PATH, tasks); err != nil {
		panic(err)
	}
}
