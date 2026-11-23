package main

import "fmt"

type Player struct {
	name string
	HP   int
}

func (p Player) Show() {
	fmt.Println(p.name, "的HP是", p.HP)
}

func (p *Player) Damage(dmg int) {
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
}

func main() {
	p := Player{
		name: "zhou",
		HP:   100,
	}
	p.Show()

	p.Damage(30)
	p.Show()

	p.Damage(100)
	p.Show()
}
