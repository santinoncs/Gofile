
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	
	var counter = map[string]int{
		"info.log": 0,
		"warning.log": 0,
		"error.log": 0,
	}

	infoChan := make(chan string)
	warningChan := make(chan string)
	errorChan := make(chan string)

	wg.Add(4)

	go process(infoChan, &wg, "info.log", counter, mutex)
	go process(warningChan, &wg, "warning.log",counter,  mutex)
	go process(errorChan, &wg, "error.log", counter, mutex)



	file, err := os.Open("log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
		reLogLevel := regexp.MustCompile(`\[(.*)\].*`)

		loglevel := reLogLevel.FindStringSubmatch(scanner.Text())[1]

		if loglevel == "info" {

			infoChan <- scanner.Text()

		}

		if loglevel == "warning" {

			warningChan <- scanner.Text()

		}

		if loglevel == "error" {

			errorChan <- scanner.Text()

		}

	}
	close(infoChan)
	close(warningChan)
	close(errorChan)
}()

	wg.Wait()
	fmt.Println(counter)

}

func process(requests chan string, wg *sync.WaitGroup, file string, counter map[string]int, mutex *sync.Mutex) {
	var count int
	defer wg.Done()
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0666);
	if err != nil {
		fmt.Println(err)
		return
	}
	w := bufio.NewWriter(f)

	for item := range requests {

		_, err := w.WriteString(item+"\n")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		count ++
		mutex.Lock()
		counter[file] = count
		mutex.Unlock()
		// flush every N lines
		if count%2000 == 0 {
			w.Flush()
		}
	}

	w.Flush()

	defer f.Close()


}
