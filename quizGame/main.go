package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	fileName := flag.String("file", "", "csv file name")
	quizTimeLimit := flag.Int("timer", 30, "quiz time limit")
	flag.Parse()
	pf, err := os.Open(fmt.Sprintf("%s.csv", *fileName))
	if err != nil {
		panic(err)
	}

	defer pf.Close()

	r := csv.NewReader(pf)

	problems, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	inputCh := make(chan string)
	done := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if !scanner.Scan() {
				close(inputCh)
				return
			}

			select {
			case inputCh <- scanner.Text():
			case <-done:
				return
			}
		}
	}()

	timer := time.NewTimer(time.Duration(*quizTimeLimit) * time.Second)
	defer timer.Stop()
	defer close(done)

	totalCorrect := 0
	for _, problem := range problems {
		question, answer := problem[0], problem[1]

		fmt.Println(question)

		select {
		case userAnswer, ok := <-inputCh:
			if !ok {
				fmt.Println("input closed")
				return
			}
			if userAnswer == answer {
				totalCorrect += 1
			}
		case <-timer.C:
			fmt.Printf("%d/%d got right", totalCorrect, len(problems))
			fmt.Println("Time up for quiz!")
			return
		}
	}

	fmt.Printf("%d/%d", totalCorrect, len(problems))
}
