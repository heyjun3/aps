package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

func main() {
	testRequest()
}

func checkDefaultValue() {
	type res struct {
		r *http.Request
	}
	r := res{}
	fmt.Println(r.r)
}

func testNewRequest() {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		// 空のstringでもエラーにはならないみたいね。
		fmt.Println(err)
	}
	fmt.Println(req)
}

func testRequest() {
	form := url.Values{}
	form.Add("shopcode", "")
	form.Add("categorycode", "0")
	form.Add("hasStock", "1")
	form.Add("currentPage", "5")

	body := strings.NewReader(form.Encode())
	fmt.Println(form.Encode())

	req, err := http.NewRequest("POST", "https://kaago.com/ajax/catalog/list/init", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	fmt.Println(string(b))
}

func person() {
	p1 := Person{
		name:   "John",
		colors: []string{"red"},
		age:    1,
		email:  "test",
	}
	err := ValidatePerson(p1)
	fmt.Println(err)
}

func ValidatePerson(p Person) (err error) {
	structType := reflect.TypeOf(p)
	if structType.Kind() != reflect.Struct {
		return errors.New("input param should be a struct")
	}

	value := reflect.ValueOf(p)
	fields := value.NumField()

	for i := 0; i < fields; i++ {
		field := value.Field(i)
		fieldName := structType.Field(i).Name

		isSet := field.IsValid() && !field.IsZero()

		if !isSet {
			err = errors.New("error" + fieldName)
		}
	}
	return err
}

type Person struct {
	name   string
	age    int
	email  string
	colors []string
}

func send() {
	c1 := make(chan string, 10)
	c2 := make(chan string, 10)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go receive(c1, c2, &wg)

	for i := 0; i < 10; i++ {
		c1 <- fmt.Sprint(i)
	}

	time.Sleep(time.Second * 2)
	close(c1)
	close(c2)
	wg.Wait()
}

func receive(ch1, ch2 chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case s1, ok := <-ch1:
			if ok {
				fmt.Println(s1)
			} else {
				return
			}
		case s2 := <-ch2:
			fmt.Println(s2)
		default:
			fmt.Println("default")
			time.Sleep(time.Second * 1)
		}
	}
}
