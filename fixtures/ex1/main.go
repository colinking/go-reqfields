package main

import "fmt"

type Hat struct {
	Style  string `require:"true"`
	OnHead bool
}

func main() {
	// Required fields error when not set.
	fmt.Printf("hat w/out style: %+v\n", Hat{})

	// Required fields don't error if set.
	// Optional fields don't error when not set.
	fmt.Printf("hat w/ snazzy style: %+v\n", Hat{
		Style: "snazzy",
	})

	// Optional fields can be set without error.
	fmt.Printf("today's hat w/ snazzy style: %+v\n", Hat{
		Style:  "snazzy",
		OnHead: true,
	})

	// Required fields error when not set in pointer structs.
	fmt.Printf("today's hat w/ snazzy style: %+v\n", &Hat{})

	// Empty struct lists do not error.
	fmt.Printf("today's hat w/ snazzy style: %+v\n", []Hat{})
}
