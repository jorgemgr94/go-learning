package main

import "fmt"

type contactInfo struct {
	email   string
	zipCode int
}

// embedded struct
type person struct {
	firstName string
	lastName  string
	contact   contactInfo
}

func structs() {
	fmt.Println("// == Structs ===========================================")
	jack := person{
		"Jack", "Bauer",
		contactInfo{"jack@h.com", 67119},
	}
	alex := person{
		firstName: "Alex", lastName: "Anderson",
		contact: contactInfo{email: "alex@h.com", zipCode: 67119},
	}

	jack.print()

	// stills works because Go pointers shortcut
	alex.updateName("Alexis")
	alex.print()

	// NOTE: person created with zero values
	// string: ""
	// int: 0
	// float: 0
	// bool: false
	var john person
	john.firstName = "John"
	john.lastName = "Doe"

	// old fashioned way to update a struct by using a pointer
	johnPointer := &john
	johnPointer.updateName("Johnny")
	john.print()
}

// a receiver function with a pointer to a person
// *person is a type description, only means we're working with a pointer to a person
// *pointerToPerson is an operator, it means we want to manipulate the value the pointer is referencing
func (pointerToPerson *person) updateName(newFirstName string) {
	(*pointerToPerson).firstName = newFirstName
}

func (p person) print() {
	fmt.Printf("%+v\n", p)
}

/*
NOTE:
- Value types: creates a new copy of the value itself for these types: int, float, string, bool, structs
- Reference types: creates a reference, which  to the value for these types: slices, maps, channels, pointers, functions
which means that the reference (pointer) is pointing to the value in memory.
*/
