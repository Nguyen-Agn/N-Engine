//go:build ignore

// MathDemo tests all functionalities of the nmath module without running Ebitengine.
//
// Run: go run .\tests\simulation\MathDemo.go
package main

import (
	"fmt"

	"autoworld/modules/napi"
)

func main() {
	fmt.Println("=====================================")
	fmt.Println("      N-MATH CONSOLE TEST DEMO       ")
	fmt.Println("=====================================")

	// 1. Nhóm kiểm soát giá trị
	fmt.Println("\n--- 1. Value Control ---")
	fmt.Printf("Clamp(150, 0, 100) = %v\n", napi.Math.Clamp(150, 0, 100))
	fmt.Printf("Lerp(0, 100, 0.5)  = %v\n", napi.Math.Lerp(0, 100, 0.5))
	fmt.Printf("Approach(10, 100, 5) = %v\n", napi.Math.Approach(10, 100, 5))
	fmt.Printf("Wrap(370, 0, 360)  = %v\n", napi.Math.Wrap(370, 0, 360))
	fmt.Printf("Snap(22.5, 16)     = %v\n", napi.Math.Snap(22.5, 16))

	// 2. Nhóm khoảng cách & góc độ
	fmt.Println("\n--- 2. Geometry & Angles ---")
	fmt.Printf("Distance(0,0, 3,4) = %v\n", napi.Math.Distance(0, 0, 3, 4))

	// Góc từ (0,0) đến (1,1) -> 45 độ
	fmt.Printf("Angle(0,0, 1,1)    = %v\n", napi.Math.Angle(0, 0, 1, 1))

	// Khoảng cách góc ngắn nhất từ 350 độ đến 10 độ -> Phải là +20 độ
	fmt.Printf("AngleDif(350, 10)  = %v\n", napi.Math.AngleDif(350, 10))

	// Cộng góc 350 thêm 30 độ -> Phải là 20 độ (tự wrap vòng tròn)
	fmt.Printf("AngleAdd(350, 30)  = %v\n", napi.Math.AngleAdd(350, 30))

	// LengthDirX / Y: Từ (0,0) đi 100px theo góc 0 độ
	lx := napi.Math.LengthDirX(100, 0)
	ly := napi.Math.LengthDirY(100, 0)
	fmt.Printf("LengthDir(100px, 0 deg) = (X: %v, Y: %v)\n", lx, ly)

	// 3. Nhóm Vector chuyên sâu (Vec2)
	fmt.Println("\n--- 3. Vec2 & Slerp ---")
	dir1 := napi.Math.DirFromDeg(90) // Hướng nhìn xuống (nếu trục Y hướng xuống)
	dir2 := napi.Math.DirFromDeg(0)  // Hướng sang phải
	difV := napi.Math.AngleDifV(dir1, dir2)
	fmt.Printf("AngleDifV(dir90, dir0)  = %v\n", difV)

	slerpDir := napi.Math.SlerpDirDeg(0, 90, 0.5)
	fmt.Printf("SlerpDirDeg(0, 90, 50%%) = %v\n", slerpDir)

	// 4. Nhóm Random & Generic Choose
	fmt.Println("\n--- 4. Random & Choose ---")
	// Random float và int
	fmt.Printf("RandomRangeDouble(0, 1) = %v\n", napi.Math.RandomRangeDouble(0, 1))
	fmt.Printf("RandomRangeInt(1, 10)   = %v\n", napi.Math.RandomRangeInt(1, 10))

	// Chance (Xác suất)
	successCount := 0
	for i := 0; i < 100; i++ {
		if napi.Math.RandomChance(0.25) { // 25% tỷ lệ
			successCount++
		}
	}
	fmt.Printf("RandomChance(25%%) run 100 times -> True %d times\n", successCount)

	// Test các hàm Generic Wrapper bên napi
	item := napi.Choose("Kiếm", "Cung", "Trượng", "Khiên")
	fmt.Printf("napi.Choose(Vũ khí) = %s\n", item)

	// ChooseFromArray
	deck := []int{10, 20, 30, 40, 50}
	card, ok := napi.ChooseFromArray(deck)
	fmt.Printf("napi.ChooseFromArray(Deck) = %v (Valid: %v)\n", card, ok)

	// Shuffle
	napi.Shuffle(deck)
	fmt.Printf("napi.Shuffle(Deck) = %v\n", deck)

	// ChooseWeighted
	drop := napi.ChooseWeighted([]float32{70, 25, 5}, "Common", "Rare", "Legendary")
	fmt.Printf("napi.ChooseWeighted(Loot) = %s\n", drop)

	fmt.Println("\n=====================================")
	fmt.Println("             TEST HOÀN TẤT             ")
	fmt.Println("=====================================")
}
