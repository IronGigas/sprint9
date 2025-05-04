package main

import (
	"fmt"
	//"math/rand/v2"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

var wg sync.WaitGroup

// generateRandomElements generates random elements.
func generateRandomElements(size int) []int {
	var slice []int

	if size <= 0 {
		return []int{} // empty slice
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	randRange := rand.New(randSource)

	for range size {
		//randomNumber := (rand.IntN(900_000_000) + 1)  //через /rand/2 было бы проще, но в задании хотят именно /rand
		randomNumber := randRange.Intn(900_000_000)
		slice = append(slice, randomNumber)
	}

	return slice
}

// maximum returns the maximum number of elements.
func maximum(data []int) int {

	max := 0

	if len(data) == 0 {
		return max
	}

	if len(data) == 1 {
		return data[0]
	}

	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

// maxChunks returns the maximum number of elements in a chunks.
func maxChunks(data []int) int {
	partOfSlice := len(data) / CHUNKS
	max := 0

	if len(data) == 0 {
		return max
	}

	if len(data) == 1 {
		return data[0]
	}

	resultChan := make(chan int, CHUNKS) // буферизированный канал для сбора результатов
	wg.Add(CHUNKS)
	for i := 0; i < CHUNKS; i++ {
		firstIndex := i * partOfSlice
		lastIndex := firstIndex + partOfSlice
		preparedSlice := data[firstIndex:lastIndex]

		go func(slice []int) {
			defer wg.Done()
			maxInPreparedSlice := maximum(slice)
			resultChan <- maxInPreparedSlice // отправляем результат в канал
		}(preparedSlice)

	}

	wg.Wait()
	close(resultChan)

	sliceOfMax := make([]int, CHUNKS)
	for maxVal := range resultChan {
		sliceOfMax = append(sliceOfMax, maxVal)
	}

	max = maximum(sliceOfMax)
	return max

}

func main() {
	fmt.Printf("Генерируем %d целых чисел\n\n", SIZE) // \n\n тут и далее для красоты, просто чтобы отделить в консоли логические элементы
	slice := (generateRandomElements(SIZE))

	fmt.Println("Ищем максимальное значение в один поток")
	start := time.Now()
	max := maximum(slice)
	elapsed := time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n\n", max, elapsed)

	fmt.Printf("Ищем максимальное значение в %d потоков\n", CHUNKS)
	start = time.Now()
	max = maxChunks(slice)
	elapsed = time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms", max, elapsed)

}
