package main // ① 主包

import "fmt" // ② 打印包

type Player struct {
	Name string
	HP   int
}

type Monster struct {
	Name string
	HP   int
}

type Boss struct {
	Name string
	HP   int
}

type Damageable interface { //定义接口
	TakeDamage(dmg int) //只要"挨打"就行
}

func (b *Boss) TakeDamage(dmg int) {
	b.HP -= dmg
	if b.HP < 0 {
		b.HP = 0
	}
}

func (p *Player) TakeDamage(dmg int) {
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
}

func (m *Monster) TakeDamage(dmg int) {
	m.HP -= dmg
	if m.HP < 0 {
		m.HP = 0
	}
}

func CastFireball(t Damageable) {
	t.TakeDamage(100)
}

func main() {
	p := &Player{Name: "zhou", HP: 150}
	m := &Monster{Name: "史莱姆", HP: 80}
	boss := &Boss{Name: "深渊领主", HP: 500}

	CastFireball(p)
	CastFireball(m)
	CastFireball(boss)

	fmt.Println(p.Name, p.HP)
	fmt.Println(m.Name, m.HP)
	fmt.Println(boss.Name, boss.HP)

}
