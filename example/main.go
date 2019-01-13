package main

import (
	"fmt"

	"github.com/ieee0824/go-tel-num-parser-jp"
)

var telList = map[string]string{
	"東京都庁":           "03-5321-1111",
	"東京都庁2":          "03(5321)1111",
	"国土交通省":          "03-5253-4150",
	"株式会社壱番屋":        "0586-76-7545",
	"株式会社壱番屋 北海道営業所": "011-896-4081",
	"株式会社壱番屋 宮城営業所":  "022-381-0215",
	"株式会社壱番屋 埼玉営業所":  "0480-93-1221",
	"株式会社壱番屋 東京営業所":  "042-735-5331",
	"電話番号ではない":       "004-0031",
}

func main() {
	for k, v := range telList {
		fmt.Print(k + ": ")
		fmt.Println(tnp.IsTelNumber(v))
	}
}
