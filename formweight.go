package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
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

func (p *Processor) getTotalWeight() float64 {
	sumOfAll := float64(0)
	for el, amount := range p.elementCounter {
		ame := MapOfEls[el].Weight
		sumOfAll += ame * float64(amount)
	}

	return sumOfAll
}

func (p *Processor) GetSubTotalMass() map[string]float64 {
	res := make(map[string]float64)

	for el, amount := range p.elementCounter {
		res[el] = float64(amount) * MapOfEls[el].Weight
	}

	return res
}

func (p *Processor) GetSubTotalRelation() map[string]float64 {
	res := make(map[string]float64)

	sumOfAll := p.getTotalWeight()

	for el, amount := range p.elementCounter {
		res[el] = (MapOfEls[el].Weight * float64(amount)) / float64(sumOfAll)
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

type elements []*Element

func (els elements) Len() int {
	return len(els)
}

func (els elements) Less(x, y int) bool {
	if els[x].AtomicNumber < els[y].AtomicNumber {
		return true
	} else {
		return false
	}
}

func (els elements) Swap(x, y int) {
	els[x], els[y] = els[y], els[x]
}

func main() {
	verbose := flag.Bool("verbose", false, "print additional information")
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

	if !(*verbose) {
		fmt.Println(p.GetWeight())
	} else {
		fmt.Println("total mass", p.GetWeight())
		w := tabwriter.NewWriter(os.Stdout, 15, 4, 0, ' ', 0)
		fmt.Fprintf(w, "element\tsymbol\tame\tsubtotal_ame\tsubtotal_%%\n")
		subtotalMassMap := p.GetSubTotalMass()
		subtotalRelMap := p.GetSubTotalRelation()
		var els elements
		for elSym := range p.elementCounter {
			els = append(els, MapOfEls[elSym])
		}
		sort.Sort(els)

		for _, el := range els {
			name := el.Name
			symbol := el.Symbol
			ame := el.Weight
			subtotalMass := subtotalMassMap[el.Symbol]
			subtotalRel := subtotalRelMap[el.Symbol] * 100
			fmt.Fprintf(w, "%s\t%s\t%.2f\t%.2f\t%.2f\n", name, symbol, ame, subtotalMass, subtotalRel)
		}
		w.Flush()
	}
}
