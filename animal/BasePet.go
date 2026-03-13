package animal

import "fmt"

type BasePet struct {
	name string
}

func (p *BasePet) Speak() {
	fmt.Println("speak:nil")
}

func (p *BasePet) Die() {
	fmt.Println("die")
}
