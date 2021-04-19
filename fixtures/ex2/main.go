package ex2

import (
	"fmt"
	"github.com/colinking/go-reqfields/fixtures/ex1"
)

func run() {
	// Required fields error when not set.
	fmt.Printf("hat w/out style: %+v\n", ex1.Hat{})

	// Required fields don't error if set.
	// Optional fields don't error when not set.
	fmt.Printf("hat w/ snazzy style: %+v\n", ex1.Hat{
		Style: "snazzy",
	})

	// Optional fields can be set without error.
	fmt.Printf("today's hat w/ snazzy style: %+v\n", ex1.Hat{
		Style:  "snazzy",
		OnHead: true,
	})

	// Required fields error when not set in pointer structs.
	fmt.Printf("today's hat w/ snazzy style: %+v\n", &ex1.Hat{})

	// Empty struct lists do not error.
	fmt.Printf("today's hat w/ snazzy style: %+v\n", []ex1.Hat{})
}
