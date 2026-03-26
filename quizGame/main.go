package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	fileName := flag.String("file", "problems.csv", "csv file name")
	quizTimeLimit := flag.Int("timer", 30, "quiz time limit")
	shouldShuffle := flag.Bool("shuffle", false, "shuffle the problems")
	flag.Parse()
	pf, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}

	defer pf.Close()

	r := csv.NewReader(pf)

	problems, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	if *shouldShuffle {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
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
				printScore(totalCorrect, len(problems))
				return
			}
			if strings.TrimSpace(userAnswer) == strings.TrimSpace(answer) {
				totalCorrect += 1
			}
		case <-timer.C:
			printScore(totalCorrect, len(problems))
			fmt.Println("Time up for quiz!")
			return
		}
	}

	printScore(totalCorrect, len(problems))
}

func printScore(totalCorrect, totalProblems int) {
	fmt.Printf("%d/%d got right\n", totalCorrect, totalProblems)
}
