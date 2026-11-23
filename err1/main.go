package main

import (
	"errors"
	"fmt"
)

func BuyItem(gold int, price int) (int, error) {
	if gold < price {
		return gold, errors.New("金币不足")
	}
	return gold - price, nil
}

func main() {
	gold := 100
	price := 150

	rest, err := BuyItem(gold, price)
	if err != nil {
		fmt.Println("购买失败", err)
		fmt.Println("你还差", price-gold, "金币")
		return
	}
	fmt.Println("购买成功,剩余", rest)
}
