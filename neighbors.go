package shufflereworked

import "math/rand"

func shuffleNeighbors(arr []string, antiNeighbors [][]string) []string {
	for _, anti := range antiNeighbors {
		a, b := anti[0], anti[1]
		for i := 0; i < len(arr)-1; i++ {
			if (arr[i] == a && arr[i+1] == b) || (arr[i] == b && arr[i+1] == a) {
				perm := rand.Perm(len(arr))
				for j := range perm {
					prev := arr[i]
					arr[i] = arr[j]
					arr[j] = prev
				}
			}
		}
		for i := 1; i < len(arr)-1; i++ {
			if arr[i] == a && (arr[i-1] == b || arr[i+1] == b) {
				perm := rand.Perm(len(arr))
				for j := range perm {
					prev := arr[i]
					arr[i] = arr[j]
					arr[j] = prev
				}
			}
			if arr[i] == b && (arr[i-1] == a || arr[i+1] == a) {
				perm := rand.Perm(len(arr))
				for j := range perm {
					prev := arr[i]
					arr[i] = arr[j]
					arr[j] = prev
				}
			}
		}
	}
	return arr
}

func antiNeighborsSatisfied(arr []string, antiNeighbors [][]string) bool {
	for _, anti := range antiNeighbors {
		a, b := anti[0], anti[1]
		for i := 0; i < len(arr)-1; i++ {
			if (arr[i] == a && arr[i+1] == b) || (arr[i] == b && arr[i+1] == a) {
				return false
			}
		}
		for i := 1; i < len(arr)-1; i++ {
			if arr[i] == a && (arr[i-1] == b || arr[i+1] == b) {
				return false
			}
			if arr[i] == b && (arr[i-1] == a || arr[i+1] == a) {
				return false
			}
		}
	}

	return true
}
