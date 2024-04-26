package shufflereworked

import (
	"errors"
	"fmt"
	"time"
)

/*
	Алгоритм сортировки массива с учетом соседничества и не повторения элементов последовательно

	Пример: [1 1 1 1 1 1 1 1 1 1 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3]
		Вес каждого элемента:
			| 1 - 10 |
			| 3 - 20 |

	Результат  : [1 3 1 3 1 3 1 3 1 3 3 3 3 3 3 1 3 3 3 3 3 3 1 3 1 3 1 3 1 3]

	Принцип работы алгоритма:
		1. Подсчитываем вес каждого элемента
		2. Находим элемент с минимальным весом и ставим его в середину массива
		3. Далее массив заполняется слева и справа относительно своего центра
		4. Для каждого следующего элемента проверяем предыдущий элемент. Если они одинаковые либо срабатывает проверка на соседство,
			то пропускаем его и ставим следующий по весу (если таковой имеется)
		5. В конце заполняем массив оставшимися элементами. Здесь есть риск повторения элементов, т.к. остатки поставить некуда.


	Соседство описывается так:

	neighbors = [][]string{
		{"element1", "element2"},
		{"element3", "element4"},
	}
*/

var (
	ErrUniqElements   = errors.New("Minimum number of unique elements must be 2")
	ErrShuffleTimeout = errors.New("Shufflle timeout")
	DefaultTimeout    = 10 * time.Second //default timeout
)

func trackTime(showLog bool) func() {
	if showLog == false {
		return func() {}
	}

	start := time.Now()
	return func() {
		fmt.Printf("Shuffle execution time %s\n", time.Since(start))
	}
}

func filterArr(ss [][]string, filter func([]string) bool) (ret [][]string) {
	for _, s := range ss {
		if filter(s) {
			ret = append(ret, s)
		}
	}
	return
}

func ShuffleReworked(initArr []string, neighbors [][]string, showLog bool, videoDuration, slotCount, timeoutSec int) ([]string, error) {
	defer trackTime(showLog)()
	var res = make([]string, slotCount)

	timeout := DefaultTimeout

	if timeoutSec > 0 {
		timeout = time.Duration(timeoutSec) * time.Second
	}

	//Формируем Map с информацией о весе элементов массива
	weightMap := calculateArrayWeight(initArr)
	weightMap.Log(showLog)

	if len(weightMap) < 2 {
		return initArr, ErrUniqElements
	}

	/*---1. Нужно найти середину массива и поставить туда значение с наименьшим весом*/
	minKey, _ := weightMap.GetMinWeightValue()
	middle := len(res) / 2
	res[middle] = minKey

	weightMap[minKey] = (weightMap[minKey]) - 1

	if weightMap[minKey] == 0 {
		delete(weightMap, minKey)
	}

	step := 1
	prev := minKey //храним предыдущий элемент массива

	for {
		leftStep := middle - step
		rightStep := middle + step

		maxKey, maxValue := weightMap.GetMaxWeightValueByNeighbors(prev, neighbors, true)

		if maxValue == 1 {
			// 1 - означает, что в мапе остались элементы с весом 1. Из них выбирать не кого. Поэтому выходим из цикла
			break
		}

		if leftStep >= 0 {
			res[leftStep] = maxKey
		}
		if rightStep < len(res) {
			res[rightStep] = maxKey
		}

		if leftStep < 0 && rightStep >= len(res) {
			break
		}

		weightMap[maxKey] = (weightMap[maxKey]) - 2

		if weightMap[maxKey] == 0 {
			delete(weightMap, maxKey)
		}

		step++
		prev = maxKey
	}

	//Бежим по левому краю массива и заполняем его оставшимися элементами
	prev = ""
	for i := middle; i >= 0; i-- {
		if res[i] == "" {

			prev = res[i+1] //предыдущий элемент стоящий справа

			maxKey, _ := weightMap.GetMaxWeightValueByNeighbors(prev, neighbors, true)
			res[i] = maxKey

			if weightMap[maxKey] <= 1 {
				delete(weightMap, maxKey)
			} else {
				weightMap[maxKey] = (weightMap[maxKey]) - 1
			}
		}
	}

	//Бежим по правому краю массива и заполняем его оставшимися элементами
	prev = ""
	for i := middle; i < len(res)-1; i++ {
		if res[i] == "" {

			prev = res[i-1] //предыдущий элемент стоящий слева

			maxKey, _ := weightMap.GetMaxWeightValueByNeighbors(prev, neighbors, false)
			res[i] = maxKey

			if weightMap[maxKey] <= 1 {
				delete(weightMap, maxKey)
			} else {
				weightMap[maxKey] = (weightMap[maxKey]) - 1
			}
		}
	}

	/*--- В этом цикле остаток элементов распределяется по оставшимся пустым ячейкам массива ---*/
	for i := range res {
		if res[i] == "" {
			maxKey, _ := weightMap.GetMaxWeightValue()
			res[i] = maxKey

			if weightMap[maxKey] <= 1 {
				delete(weightMap, maxKey)
			} else {
				weightMap[maxKey] = (weightMap[maxKey]) - 1
			}
		}
	}

	var shuffledArr []string
	i := 0
	for end := time.Now().Add(timeout); ; {
		shuffledArr = shuffleNeighbors(res, neighbors)
		if antiNeighborsSatisfied(shuffledArr, neighbors) {
			break
		}

		if i&0x0f == 0 { // Check in every 16th iteration
			if time.Now().After(end) {
				fmt.Println("ShuffleReworked completed by timeout")
				return nil, ErrShuffleTimeout
			}
		}
		i++
	}

	return shuffledArr, nil
}
