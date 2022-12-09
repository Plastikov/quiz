package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Challenge struct {
	question string
	answer   string
}

func parseCSV(data string) ([][]string, error) {
	file, err := os.Open(data)
	if err != nil {
		log.Fatalf("Unable to open file. Tried to open: %s\n", data)
	}

	reader := csv.NewReader(file)
	fileData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return fileData, nil
}

func populateChallenges(data [][]string) []Challenge {
	for _, row := range data {
		challenges = append(challenges, Challenge{
			question: row[0],
			answer:   row[1],
		})
	}
	return challenges
}

func run(cs []Challenge) {
	var (
		scanner = bufio.NewScanner(os.Stdin)
		timer = time.NewTimer(time.Duration(*stopwatch) * time.Second)
	)
ChallengeLoop:
	for i, c := range cs {
		fmt.Printf("Question #%d: what is %s? ", i+1, c.question)
		scoreCh := make(chan string)
		go func() {
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())
			scoreCh <- input
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nTime out!\n")
			break ChallengeLoop
		case input := <- scoreCh:
			if input == c.answer {
			score++
		}
		answeredQuestions++
		}
	}
}

func reportScore() {
	fmt.Printf("You answered %d questions and got %d questions right out of %d\n",
	answeredQuestions, 
	score,
	len(challenges))
}

func shuffler() {
	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(challenges), func(i, j int) {
			challenges[i], challenges[j] = challenges[j], challenges[i]
		})
	}
}

var (
	score int
	answeredQuestions int
	challenges []Challenge
	shuffle = flag.Bool("shuffle", false, "Pass a true value to this flag if you want the quiz questions to be shuffled")
	stopwatch = flag.Int64("timer", 30, "Allotted time to run the quiz in seconds")
	
)

func main() {
	csvFileName := flag.String("filename", "./problem.csv", "file containing all the questions to be presented to the user")
	flag.Parse()

	rawData, err := parseCSV(*csvFileName)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	quiz := populateChallenges(rawData)
	shuffler()
	run(quiz)
	reportScore()
	
}
