package util

import (
	"fmt"
	"math"
	"os"
	"time"
)

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func GetAfterPercent(price float64, per int) float64 {
	vl := (price * (float64(per) * 0.01))
	return price + vl
}

func GetBidPrice(price float64) float64 {
	P10 := 10.0
	P100 := P10 * 10.0
	P1_000 := P100 * 10.0
	P10_000 := P1_000 * 10.0
	P100_000 := P10_000 * 10.0
	P1_000_000 := P100_000 * 10.0

	if price >= P1_000_000*2 {
		return P1_000
	} else if price >= P1_000_000 {
		return P100 * 5
	} else if price >= P100_000*5 {
		return P100
	} else if price >= P100_000 {
		return P10 * 5
	} else if price >= P10_000 {
		return P10
	} else if price >= P1_000 {
		return 5.0
	} else if price > P100 {
		return 1.0
	} else if price > P10 {
		return 0.1
	} else {
		return 0.01
	}
}

func BinarySearch(items []int, v int) int {
	min := 0
	max := len(items)
	mid := 0

	if v < items[min] || v > items[max-1] {
		return -1
	}

	for min <= max {
		mid = (min + max) / 2
		// fmt.Println(min, mid, max, items[mid], v)
		if items[mid] == v {
			break
		} else if v > items[mid] {
			min = mid + 1
		} else {
			max = mid - 1
		}
		//fmt.Println(min, max)
	}

	// fmt.Println("=============================")
	if min >= len(items) {
		return -1
	} else {
		return mid
	}
}

func GetPrintFloat64(v float64) string {
	return fmt.Sprintf("%.4f", v)
}

func IsExistsFileOrDir(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func GetTodayYYYYMMDD() string {
	return time.Now().Format("yyyymmdd")
}
