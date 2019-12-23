package main

import (
	"fmt"
	"math"

	lsdp "github.com/deltam/go-lsd-parametrized"
)

func main() {
	fmt.Println(CalScore("瀬戸大橋線", "name"))
	fmt.Println(CalScore("東武宇都宮", "宇都宮"))
	fmt.Println(CalScore("神戸高速南北", "神戸高速"))
	fmt.Println(CalScore("しなの鉄道しなの鉄道", "しなの鉄道"))
	fmt.Println(CalScore("神戸新交通ポートライナー", "ポートライナー"))
	fmt.Println(CalScore("神戸新交通ポートライナー", "新交通ポートライナー"))
	fmt.Println(CalScore("神戸新交通ポートライナー", "交通ポートライナー"))
	fmt.Println(CalScore("名古屋地下鉄上飯田", "飯田"))
	fmt.Println(CalScore("神戸新交通ポートライナー", "神戸新交通ポートライナー"))
	fmt.Println(CalScore("飯田", "飯田"))
	fmt.Println(CalScore("関電トンネルトロリーバス", "東海道・山陽新幹線"))
	fmt.Println(CalScore("えちごトキめき鉄道日本海ひすいライン", "中央本線支線"))
	fmt.Println(CalScore("えびの", "上ノ国"))
}

// Equal はレーベンシュタイン距離を用いてスコア化する
func CalScore(x1 string, x2 string) int {

	wd := lsdp.Weights{Insert: 0.1, Delete: 1, Replace: 1}
	nd := lsdp.Normalized(wd)

	return round(100.0 - nd.Distance(x1, x2)*100.0)
}

// 丸め
func round(f float64) int {
	return int(math.Floor(f + .5))
}
