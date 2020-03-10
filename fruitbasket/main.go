package main

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"fmt"
)

type FruitBasket struct {
    Capacity int `yaml:"capacity"`
    Fruits []Fruit
}

func NewFruitBasket() *FruitBasket {
	fb := new(FruitBasket)

	return fb
}

type Fruit interface {
	GetFruitName() string
	GetNumber() int
}

type Apple struct {
	Name string `yaml:"name"`
	Number int `yaml:"number"`
}

func (apple *Apple) GetFruitName() string {
	return apple.Name
}

func (apple *Apple) GetNumber() int {
	return apple.Number
}

type tmpFruitBasket struct {
	Capacity int `yaml:"capacity"`
	Fruits []map[string]yaml.Node
}

func (fruitBasket *FruitBasket) UnmarshalYAML(value *yaml.Node) error {
    var tmpFruitBasket tmpFruitBasket

    if err := value.Decode(&tmpFruitBasket); err != nil {
        return err
	}
	
	fruitBasket.Capacity = tmpFruitBasket.Capacity

	fruits := make([]Fruit, 0, len(tmpFruitBasket.Fruits))

    for i := 0; i < len(tmpFruitBasket.Fruits); i++ {
        for tag, node := range tmpFruitBasket.Fruits[i] {
            switch tag {
            case "Apple":
                apple := &Apple{}
                if err := node.Decode(apple); err != nil {
                    return err
                }

                fruits = append(fruits, apple)
            default:
                return errors.New("Failed to interpret the fruit of type: \"" + tag + "\"")
            }
        }
    }

    fruitBasket.Fruits = fruits

    return nil
}

func main() {
	data := []byte(`
capacity: 2
Apple:
- name: "apple1"
  number: 1
- name: "apple2"
  number: 1
`)

	fruitBasket := NewFruitBasket()

	err := yaml.Unmarshal(data, &fruitBasket)
    
    if err != nil {
        log.Fatalf("error: %v", err)
	}

	fmt.Println(fruitBasket.Capacity)

    for i := 0; i < len(fruitBasket.Fruits); i++ {
        switch fruit := fruitBasket.Fruits[i].(type) {
        case *Apple:
			fmt.Println(fruit.Name)
			fmt.Println(fruit.Number)
		}
	}
}