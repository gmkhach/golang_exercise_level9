package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type Person struct {
	First string
	Last  string
	Age   int
}

type Human interface {
	Speak()
}

func main() {
	/*
		Exercise 1
		1. In addition to the main goroutine, launch two additional goroutines that print something out.
		2. Use waitgroups to make sure each goroutine finishes before your program exists.
	*/
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		fmt.Println("thing one")
		fmt.Println("Goroutines:", runtime.NumGoroutine())
		wg.Done()
	}()

	go func() {
		fmt.Println("thing two")
		fmt.Println("Goroutines:", runtime.NumGoroutine())
		wg.Done()
	}()

	wg.Wait()

	/*
		Exercise 2
		This exercise will reinforce our understanding of method sets:
			1. Create a type person struct
			2. Attach a method speak to type person using a pointer receiver *person
			3. Create a type human interface
				(Hint: To implicitly implement the interface, a human must have the speak method)
			4. Create func “saySomething”
				a. Have it take in a human as a parameter
				b. Have it call the speak method
			5. Show the following in your code
				a. You CAN pass a value of type *person into saySomething
				b. You CANNOT pass a value of type person into saySomething
	*/
	p1 := Person{
		First: "Dr",
		Last:  "Seuss",
		Age:   87,
	}

	saySomething(&p1)
	// The following call will not work because p1 is the value, not the pointer, to the p1 variable.
	//	saySomething(p1)

	/*
		Exercise 3
		1. Using goroutines, create an incrementer program
			a. Have a variable to hold the incrementer value
			b. Launch a bunch of goroutines. Each goroutine should:
				- Read the incrementer value
				- Store it in a new variable
				- Yield the processor with runtime.Gosched()
				- Increment the new variable
				- Write the value in the new variable back to the incrementer variable
		2. Use waitgroups to wait for all of your goroutines to finish
		3. The above will create a race condition. Prove that it is a race condition by using the -race flag
	*/
	counter := 0
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(){
			fmt.Println("Goroutines:", runtime.NumGoroutine())	
			fmt.Println("Counter before:", counter)
			
			v := counter
			runtime.Gosched()
			v++
			counter = v
			
			fmt.Println("Counter after:", counter)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("Counter final:", counter)

	/*
		Exercise 4
		Fix the race condition you created in the previous exercise by using a mutex
			(Hint: It makes sense to remove runtime.Gosched())
	*/
	var mu sync.Mutex
	counter = 0
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(){
			fmt.Println("Goroutines:", runtime.NumGoroutine())	
			
			mu.Lock()
			fmt.Println("Counter before:", counter)
			
			v := counter
			runtime.Gosched()
			v++
			counter = v
			
			fmt.Println("Counter after:", counter)
			mu.Unlock()
			
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("Counter final:", counter)


	/*
		Exercise 5
		Fix the race condition you created in exercise #3 by using package atomic
	*/
	var incrementer int64
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(){
			fmt.Println("Goroutines:", runtime.NumGoroutine())	
			fmt.Println("Counter before:", atomic.LoadInt64(&incrementer))
			
			atomic.AddInt64(&incrementer, 1)
			fmt.Println("Counter after:", atomic.LoadInt64(&incrementer))
			
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("Counter final:", incrementer)


	/*
		Exercise 6
		1. Create a program that prints out your OS and ARCH.
		2. Use the following commands to run it
			a. go run
			b. go build
			c. go install
	*/
	fmt.Println("OS:", runtime.GOOS)
	fmt.Println("Architecture:", runtime.GOARCH)
}

func (p *Person) Speak() {
	fmt.Println("Hello. I'm a human.")
}

func saySomething(h Human) {
	h.Speak()
}