package collection

import (
	"fmt"
	"strconv"
)

func (collection *Collection) Print() {
	for key, s := range collection.Structs {
		printStruct(key, s)
		fmt.Printf("    struct.Fields:\n\n")
		printFields(s)

	}
	fmt.Printf("\n\t\tEND.\n\n")
}

func printStruct(key int, s Struct) {
	fmt.Printf("collection.structs:\n")
	fmt.Printf(" %v. struct.Name: %v\n", strconv.Itoa(key), s.Name)

	fmt.Printf("-------------------------------------------------------------------------\n")
}

func printFields(s Struct) {
	if len(s.Fields) == 0 {
		fmt.Printf("\t\tStruct contains no fields!\n")
		return
	}
	for fkey, field := range s.Fields {
		if field.Name != "" {
			fmt.Printf("\t%v. field.Name: %v\n", strconv.Itoa(fkey), field.Name)
			fmt.Printf("\t   field.Type: %v\n", field.Type)
			fmt.Printf("\n")
		}
	}
}
