package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var reason string
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "reason, r",
			Value:       "generic reason",
			Usage:       "what distracted you?",
			Destination: &reason,
		},
	}

	app.Action = func(c *cli.Context) error {
		t := time.Now()

		year, month, day := t.Date()
		formattedDate := fmt.Sprintf("%d-%02d-%02d", year, month, day)

		workingDir := "/tmp/distraction_tracker"

		_ = os.Mkdir(workingDir, os.ModePerm)

		fileName := fmt.Sprintf("%s/%s.log", workingDir, formattedDate)

		if _, err := os.Stat(fileName); err == nil {
			f, openError := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			check(openError)

			writeErr := WriteToFile(f, t, reason)
			check(writeErr)

			defer f.Close()
		} else if os.IsNotExist(err) {
			fmt.Println("Creating file: ", fileName)

			f, fileError := os.Create(fileName)
			check(fileError)

			writeErr := WriteToFile(f, t, reason)
			check(writeErr)

			defer f.Close()
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteToFile(f *os.File, t time.Time, reason string) error {
	reasonText := fmt.Sprintf("%s -- %s\n", t.Format(time.RFC3339), reason)

	_, writeError := f.WriteString(reasonText)

	return writeError
}
