package animal

import "fmt"

type Dog struct {
	BasePet
}

func (d *Dog) Speak() {
	fmt.Println("speak:wang!")
}

//func (d *Dog) Die() {
//	fmt.Println("dog die")
//}
