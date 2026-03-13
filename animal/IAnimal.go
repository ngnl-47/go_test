package animal

type IAnimal interface {
	Speak()
	Die()
}

func NewDog() IAnimal {
	return IAnimal(&Dog{})
}
