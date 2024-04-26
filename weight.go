package shufflereworked

import (
	"fmt"
	"sort"
	"strings"
)

type WeightMap map[string]int

func (this WeightMap) Log(showLog bool) {
	if showLog {
		maxLength := 0
		for e, r := range this {
			mess := fmt.Sprintf("| %s - %d |", e, r)
			if len(mess) > maxLength {
				maxLength = len(mess)
			}
			fmt.Println(mess)
		}
		fmt.Println(strings.Repeat("-", maxLength))
	}
}

func (this WeightMap) GetMinWeightValue() (string, int) {
	keys := make([]string, 0, len(this))
	for key := range this {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return "", 1
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return this[keys[i]] < this[keys[j]]
	})

	return keys[0], this[keys[0]]
}

func (this WeightMap) GetMaxWeightValue() (string, int) {
	keys := make([]string, 0, len(this))
	for key := range this {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return "", 1
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return this[keys[i]] > this[keys[j]]
	})

	return keys[0], this[keys[0]]
}

func (this WeightMap) GetMaxWeightValueByNeighbors(prev string, neighbors [][]string, withSort bool) (string, int) {
	keys := make([]string, 0, len(this))
	for key := range this {

		//если пред и текущий один и тот же рк и если у них нет антисоседов можно пропустить вычесление веса
		if key == prev && len(neighbors) <= 0 {
			continue
		}

		n1 := filterArr(neighbors, func(arr []string) bool {
			for _, elem := range arr {
				if key == elem {
					return true
				}
			}
			return false
		})

		n2 := filterArr(n1, func(arr []string) bool {
			for _, elem := range arr {
				if prev == elem {
					return true
				}
			}
			return false
		})

		if len(n2) == 0 {
			keys = append(keys, key)
		}

	}

	if withSort {
		sort.SliceStable(keys, func(i, j int) bool {
			return this[keys[i]] > this[keys[j]]
		})
	}

	if len(keys) == 0 {
		return "", 1
	}

	return keys[0], this[keys[0]]
}

// Подсчет веса каждого элемента массива
func calculateArrayWeight(arr []string) WeightMap {
	var mp = WeightMap{}

	for _, elem := range arr {
		if _, ok := mp[elem]; ok {
			mp[elem]++
		} else {
			mp[elem] = 1
		}
	}
	return mp
}
