package main

import (
	"fmt"
	"log"
	"os"

	. "github.com/brianw0924/go_assignment_3/raid"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	inputFile := os.Args[1]

	fmt.Printf("Block size: %d\n", BLOCK_SIZE)
	fmt.Printf("Block per disk: %d\n", BLOCK_PER_DISK)
	fmt.Printf("Disk size: %d\n", BLOCK_SIZE*BLOCK_PER_DISK)
	fmt.Printf("Number of disks: %d\n", STRIPE_WIDTH)
	fmt.Printf("Stripe size: %d\n", STRIPE_WIDTH*BLOCK_SIZE)
	fmt.Printf("Total storage: %d\n\n", BLOCK_SIZE*BLOCK_PER_DISK*STRIPE_WIDTH)

	data, err := os.ReadFile(inputFile)
	check(err)

	fmt.Printf("Input length: %d\n\n", len(data))
	fmt.Printf("Input string:  %s\n\n", string(data))

	raid := NewRaid10()
	err = raid.Write(data)
	check(err)

	s, err := raid.Read(len(data))
	check(err)
	fmt.Printf("Output string: %s\n", s)
}
