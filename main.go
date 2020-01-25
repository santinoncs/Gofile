
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

	
	
	infoChan := make(chan string)
	warningChan := make(chan string)
	errorChan := make(chan string)

	wg.Add(4)
	
	go process(infoChan, &wg, "info.log")
	go process(warningChan, &wg, "warning.log")
	go process(errorChan, &wg, "error.log")



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

}

func process(requests chan string, wg *sync.WaitGroup, file string) {
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
		// flush every N lines
		if count%2000 == 0 {
			w.Flush()
			fmt.Println("flushing!", file)
		}
	}

	w.Flush()

	defer f.Close()


}
/*
func infoProcess(requests chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.OpenFile("info.log", os.O_CREATE|os.O_RDWR, 0666);
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
		// flush every N lines
		w.Flush()
	}

	defer f.Close()

}

func warningProcess(requests chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.OpenFile("warning.log", os.O_CREATE|os.O_RDWR, 0666);
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
		//fmt.Println(n4, "bytes written successfully")
		w.Flush()
	}

	defer f.Close()

}

func errorProcess(requests chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.OpenFile("error.log", os.O_CREATE|os.O_RDWR, 0666);
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
		//fmt.Println(n4, "bytes written successfully")
		w.Flush()
	}

	defer f.Close()

}
*/