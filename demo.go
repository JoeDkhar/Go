// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {

	var (
		name   string  = "Lord of the Rings"
		rating float32 = 9.62
		cost   int     = 250
	)

	fmt.Printf("%v %T\t", name, name)
	fmt.Printf("%v %T\t", rating, rating)
	fmt.Printf("%v %T\n", cost, cost)
	//constants
	const (
		Sunday    = iota
		Monday    = iota
		Tuesday   = iota
		Wednesday = iota
		Thursday  = iota
		Friday    = iota
		Saturday  = iota
	)

	fmt.Println(Sunday)
	fmt.Println(Tuesday)
	fmt.Println(Saturday)
	//arrays
	var gari [3]string
	//gari := []string{"Mazada RX7"}
	gari[0] = "Mazada RX7"
	gari[1] = "AE86"
	gari[2] = "Ford GT"
	fmt.Println(gari)

	cars := [3]string{"Nissan Altima", "Toyota Celica", "Toyota FourRunner"}

	fmt.Println(cars)

	//maps : key value pairs

	employee := map[string]int{
		"John":  101,
		"Peter": 102,
		"Sam":   103,
	}

	fmt.Println(employee)
}
