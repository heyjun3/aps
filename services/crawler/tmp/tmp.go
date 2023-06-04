package main

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

func main() {
	p1 := Person{
		name: "John",
		colors: []string{"red"},
		age: 1,
		email: "test",
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
			err = errors.New("error" + fieldName )
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
