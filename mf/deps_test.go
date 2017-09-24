package mf_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/magefile/mage/mf"
	"github.com/pkg/errors"
)

func TestDeps(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	start := time.Now()
	Build()
	log.Printf("Test done, duration: %v", time.Since(start))
}

func Build() {
	mf.Deps(Root1, Root2)
	log.Println("Build done.")
}

func Root1() {
	mf.Deps(Leaf1, Leaf2)
	time.Sleep(time.Millisecond * 100)
	log.Println("Root1 done")
}

func Root2() {
	mf.Deps(Leaf1, Leaf3)
	time.Sleep(time.Millisecond * 100)
	log.Println("Root2 done")
}
func Leaf1() {
	time.Sleep(time.Millisecond * 100)
	log.Println("Leaf1 done")
}
func Leaf2() error {
	time.Sleep(time.Millisecond * 100)
	log.Println("Leaf2 done")
	return errors.New("bah!")
}
func Leaf3() error {
	time.Sleep(time.Millisecond * 100)
	log.Println("Leaf3 done")
	return errors.New("boo!")
}
