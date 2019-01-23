package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"unicode"
)

type Processor struct {
	input string
	index int

	curChar  rune
	peekChar rune

	elementCounter map[string]int
}

func NewProcessor(input string) *Processor {
	p := &Processor{
		input:          input,
		index:          -2,
		elementCounter: make(map[string]int),
	}

	p.nextChar()
	p.nextChar()

	return p
}

func (p *Processor) nextChar() {
	p.index++
	p.curChar = p.peekChar
	if p.index+1 >= len(p.input) {
		p.peekChar = 0
	} else {
		p.peekChar = rune(p.input[p.index+1])
	}
}

func (p *Processor) Process() error {
	for p.curChar != 0 {
		curSym := ""
		curMul := ""

		if !unicode.IsUpper(p.curChar) {
			return fmt.Errorf("Expected upper case letter, but got, %c", p.curChar)
		}
		curSym += string(p.curChar)
		p.nextChar()

		if unicode.IsLower(p.curChar) {
			curSym += string(p.curChar)
			p.nextChar()
		}

		for unicode.IsDigit(p.curChar) {
			curMul += string(p.curChar)
			p.nextChar()
		}

		if _, ok := MapOfEls[curSym]; !ok {
			return fmt.Errorf("Could not find element with symbol %s", curSym)
		}

		curMulInt := int64(1)
		if curMul != "" {
			curMulInt, _ = strconv.ParseInt(curMul, 10, 64)
		}

		p.elementCounter[curSym] += int(curMulInt)
	}

	return nil
}

func (p *Processor) GetWeight() float64 {
	sum := float64(0)
	for el, amount := range p.elementCounter {
		element := MapOfEls[el]

		sum += element.Weight * float64(amount)
	}

	return sum
}

func (p *Processor) GetAtomicRelation() map[string]float64 {
	res := make(map[string]float64)

	sumOfAll := 0
	for _, amount := range p.elementCounter {
		sumOfAll += amount
	}

	for el, amount := range p.elementCounter {
		res[el] = float64(amount) / float64(sumOfAll)
	}

	return res
}

func (p *Processor) GetWeightRelation() map[string]float64 {
	res := make(map[string]float64)

	weightOfAll := float64(0)
	for el, amount := range p.elementCounter {
		weightOfAll += MapOfEls[el].Weight * float64(amount)
	}

	for el, amount := range p.elementCounter {
		res[el] = (MapOfEls[el].Weight * float64(amount)) / float64(weightOfAll)
	}

	return res
}

type Element struct {
	Name         string
	Symbol       string
	AtomicNumber int
	Weight       float64
}

var MapOfEls map[string]*Element

func main() {
	atomicRelation := flag.Bool("atomic-relation", false, "print relations of atoms")
	weightRelation := flag.Bool("weight-relation", false, "print relations of atomic weight")
	flag.Parse()

	input := flag.Arg(0)
	if input == "" {
		inputByte, _ := ioutil.ReadAll(os.Stdin)
		input = string(inputByte[:len(inputByte)-1])
	}

	p := NewProcessor(input)
	err := p.Process()
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println(p.GetWeight())

	if *atomicRelation {
		fmt.Println("atomic relation")
		fmt.Println(p.GetAtomicRelation())
	}

	if *weightRelation {
		fmt.Println("atomic weight relation")
		fmt.Println(p.GetWeightRelation())
	}
}
