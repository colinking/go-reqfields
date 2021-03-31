package main

import "fmt"

type Hat struct {
	Style  string `required:"true"`
	OnHead bool
}

func main() {
	// This should error.
	fmt.Printf("hat w/out style: %+v\n", Hat{})

	// These should not error.
	fmt.Printf("hat w/ snazzy style: %+v\n", Hat{
		Style: "snazzy",
	})
	fmt.Printf("today's hat w/ snazzy style: %+v\n", Hat{
		Style:  "snazzy",
		OnHead: true,
	})
}
