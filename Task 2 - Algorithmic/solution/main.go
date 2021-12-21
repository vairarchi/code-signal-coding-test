package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ElementStore struct {
	elements map[string]int
}

func (store *ElementStore) GetElement(writer http.ResponseWriter, element string) {
	if _, found := store.elements[element]; !found {

		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer)

	} else {
		fmt.Fprint(writer, strconv.Itoa(store.elements[element]))
	}

}

func (store *ElementStore) DeleteElement(writer http.ResponseWriter, element string) {
	if _, found := store.elements[element]; !found {

		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer)

	} else {
		delete(store.elements, element)
		fmt.Fprint(writer)
	}

}

func StartHttpServer(wg *sync.WaitGroup, store *ElementStore) *http.Server {
	server := &http.Server{Addr: ":7777"}

	http.HandleFunc("/elements", func(w http.ResponseWriter, r *http.Request) {
		keys, _ := r.URL.Query()["id"]
		element := keys[0]
		switch r.Method {
		case http.MethodDelete:
			store.DeleteElement(w, element)
		case http.MethodGet:
			store.GetElement(w, element)
		}
	})

	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe(): %v\n", err)
		}
	}()

	return server
}

func RetrieveResponse(requestBuilder func() (*http.Request, error)) string {
	request, _ := requestBuilder()
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Unable to make a request: %v", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		if len(body) > 0 {
			return strconv.Itoa(resp.StatusCode) + " " + string(body)
		} else {
			return strconv.Itoa(resp.StatusCode)
		}
	} else {
		fmt.Println("Unable to make a request: %v", err)
		return ""
	}
}

func main() {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	store := &ElementStore{map[string]int{"a": 2, "b": 10, "aba": 7}}
	server := StartHttpServer(httpServerExitDone, store)
	time.Sleep(300 * time.Millisecond)

	file, err := os.Open("/root/data/requests.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if parts[0] == "GET" {
			fmt.Println(RetrieveResponse(func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "http://localhost:7777/elements?id="+parts[1], nil)
			}))
		} else if parts[0] == "DELETE" {
			fmt.Println(RetrieveResponse(func() (*http.Request, error) {
				return http.NewRequest(http.MethodDelete, "http://localhost:7777/elements?id="+parts[1], nil)
			}))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
	if err := server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}

	httpServerExitDone.Wait()
}
