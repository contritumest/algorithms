package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Graph struct {
	N   int
	Mat [][]int
}

type Node struct {
	Value int
	Next  *Node
}

type LinkedQueue struct {
	Head *Node
	Last *Node
}

func (q *LinkedQueue) Push(value int) {
	node := &Node{
		Value: value,
		Next:  nil,
	}

	if q.Last == nil {
		q.Head = node
		q.Last = node
		return
	}

	q.Last.Next = node
	q.Last = node
}

func (q *LinkedQueue) Pop() (int, bool) {
	if q.Head == nil {
		return 0, false
	}

	value := q.Head.Value
	q.Head = q.Head.Next

	if q.Head == nil {
		q.Last = nil
	}

	return value, true
}

func (q *LinkedQueue) IsEmpty() bool {
	return q.Head == nil
}

func NewGraph(n int, p float64) Graph {
	g := Graph{
		N:   n,
		Mat: make([][]int, n),
	}

	for i := range g.Mat {
		g.Mat[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Float64()*100 < p {
				g.Mat[i][j] = 1
				g.Mat[j][i] = 1
			}
		}
	}

	for i := 0; i < n; i++ {
		if rand.Float64()*100 < p {
			g.Mat[i][i] = 1
		}
	}

	return g
}

func (g Graph) PrintMatrix() {
	if g.N == 0 {
		fmt.Println("Граф ещё не создан.")
		return
	}

	fmt.Println("\nМатрица смежности:")
	fmt.Print("    ")

	for j := 0; j < g.N; j++ {
		fmt.Printf("%3d", j+1)
	}
	fmt.Println()

	for i := 0; i < g.N; i++ {
		fmt.Printf("%3d ", i+1)

		for j := 0; j < g.N; j++ {
			fmt.Printf("%3d", g.Mat[i][j])
		}

		fmt.Println()
	}
}

func (g Graph) AdjacencyList() [][]int {
	list := make([][]int, g.N)

	for i := 0; i < g.N; i++ {
		for j := 0; j < g.N; j++ {
			if g.Mat[i][j] == 1 {
				list[i] = append(list[i], j)
			}
		}
	}

	return list
}

func (g Graph) PrintAdjacencyList() {
	if g.N == 0 {
		fmt.Println("Граф ещё не создан.")
		return
	}

	fmt.Println("\nСписки смежности:")
	list := g.AdjacencyList()

	for i, neighbors := range list {
		fmt.Printf("%d: ", i+1)

		if len(neighbors) == 0 {
			fmt.Println("нет смежных вершин")
			continue
		}

		for k, v := range neighbors {
			if k > 0 {
				fmt.Print(", ")
			}
			fmt.Print(v + 1)
		}

		fmt.Println()
	}
}

func (g Graph) BFSMatrix(start int, order *[]int) {
	visited := make([]bool, g.N)
	*order = make([]int, 0, g.N)

	for startVertex := 0; startVertex < g.N; startVertex++ {
		v := start
		if startVertex > 0 {
			v = startVertex - 1
			if v >= start {
				v++
			}
		}

		if visited[v] {
			continue
		}

		queue := make([]int, 0, g.N)

		queue = append(queue, v)
		visited[v] = true

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			*order = append(*order, current)

			for to := 0; to < g.N; to++ {
				if g.Mat[current][to] == 1 && !visited[to] {
					queue = append(queue, to)
					visited[to] = true
				}
			}
		}
	}
}

func (g Graph) BFSLinkedQueue(start int, order *[]int) {
	visited := make([]bool, g.N)
	*order = make([]int, 0, g.N)

	for startVertex := 0; startVertex < g.N; startVertex++ {
		v := start
		if startVertex > 0 {
			v = startVertex - 1
			if v >= start {
				v++
			}
		}

		if visited[v] {
			continue
		}

		var queue LinkedQueue

		queue.Push(v)
		visited[v] = true

		for !queue.IsEmpty() {
			current, _ := queue.Pop()
			*order = append(*order, current)

			for to := 0; to < g.N; to++ {
				if g.Mat[current][to] == 1 && !visited[to] {
					queue.Push(to)
					visited[to] = true
				}
			}
		}
	}
}

func (g Graph) BFSAdjacencyList(start int, order *[]int) {
	list := g.AdjacencyList()

	visited := make([]bool, g.N)
	*order = make([]int, 0, g.N)

	for startVertex := 0; startVertex < g.N; startVertex++ {
		v := start
		if startVertex > 0 {
			v = startVertex - 1
			if v >= start {
				v++
			}
		}

		if visited[v] {
			continue
		}

		queue := make([]int, 0, g.N)

		queue = append(queue, v)
		visited[v] = true

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			*order = append(*order, current)

			for _, to := range list[current] {
				if !visited[to] {
					queue = append(queue, to)
					visited[to] = true
				}
			}
		}
	}
}

func printOrder(title string, order []int) {
	fmt.Println("\n" + title)
	fmt.Print("Порядок посещения вершин: ")

	for i, v := range order {
		if i > 0 {
			fmt.Print(" -> ")
		}

		fmt.Print(v + 1)
	}

	fmt.Println()
}

func readInt(in *bufio.Reader, prompt string, min, max int) (int, error) {
	fmt.Print(prompt)

	var v int
	_, err := fmt.Fscan(in, &v)
	if err != nil {
		return 0, fmt.Errorf("Ошибка ввода: %v", err)
	}

	if v >= min && v <= max {
		return v, nil
	}

	return 0, fmt.Errorf("Ошибка диапазона")
}

func readProbability(in *bufio.Reader) (float64, error) {
	fmt.Print("Вероятность появления ребра (0-100%): ")

	var p float64
	_, err := fmt.Fscan(in, &p)
	if err != nil {
		return 0, fmt.Errorf("Ошибка ввода вероятности: %v", err)
	}

	if p >= 0 && p <= 100 {
		return p, nil
	}
	return 0, fmt.Errorf("Ошибка диапазона вероятности")
}

func checkTime(g Graph, start int, f func(int, *[]int), wg *sync.WaitGroup, tm *time.Duration, order1 *[]int) {
	defer wg.Done()
	startTime := time.Now()
	f(start, order1)
	*tm = time.Since(startTime)
}

func compareTime(g Graph, start int, writeResult bool) {
	fmt.Println("\nСравнение времени выполнения BFS:")

	var order1, order2 []int
	var time1, time2 time.Duration
	var wg sync.WaitGroup

	wg.Add(1)
	go checkTime(g, start, g.BFSMatrix, &wg, &time1, &order1)
	wg.Add(1)
	go checkTime(g, start, g.BFSLinkedQueue, &wg, &time2, &order2)

	wg.Wait()

	if writeResult {
		printOrder("BFS со слайсом:", order1)
		printOrder("BFS с самостоятельно реализованной очередью:", order2)
	}

	fmt.Printf("\nВремя BFS со слайсом: %v\n", time1)
	fmt.Printf("Время BFS с односвязным списком: %v\n", time2)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	in := bufio.NewReader(os.Stdin)

	var g Graph
	var probability float64 = 30

	var order []int

	for {
		fmt.Println("\n========== ЛАБОРАТОРНАЯ РАБОТА №8 ==========")
		fmt.Println("Обход графа в ширину (BFS)")
		fmt.Println("1. Создать случайный граф")
		fmt.Println("2. Вывести матрицу смежности")
		fmt.Println("3. BFS по матрице смежности (слайс)")
		fmt.Println("4. BFS по матрице (самостоятельная очередь)")
		fmt.Println("5. Сравнить время двух реализаций BFS")
		fmt.Println("6. Вывести списки смежности")
		fmt.Println("7. BFS по спискам смежности")
		fmt.Println("8. Указать вероятность появления ребра")
		fmt.Println("0. Выход")

		choice, err := readInt(in, "Выберите пункт меню: ", 0, 8)
		for err != nil {
			fmt.Println("Ошибка:", err)
			choice, err = readInt(in, "Выберите пункт меню: ", 0, 8)
		}

		switch choice {
		case 0:
			fmt.Println("Работа программы завершена.")
			return

		case 1:
			n, err := readInt(in, "Количество вершин (1-100): ", 1, 100)
			for err != nil {
				fmt.Println("Ошибка:", err)
				n, err = readInt(in, "Количество вершин (1-100): ", 1, 100)
			}

			g = NewGraph(n, probability)

			fmt.Printf(
				"Граф создан: вершин — %d, вероятность каждого ребра — %.2f%%.\n",
				n,
				probability,
			)

			g.PrintMatrix()

		case 2:
			g.PrintMatrix()

		case 3:
			if g.N == 0 {
				fmt.Println("Сначала создайте граф (пункт 1).")
				continue
			}

			start, err := readInt(in, "Введите начальную вершину: ", 1, g.N)
			for err != nil {
				fmt.Println("Ошибка:", err)
				start, err = readInt(in, "Введите начальную вершину: ", 1, g.N)
			}
			start--

			startTime := time.Now()
			g.BFSMatrix(start, &order)
			elapsed := time.Since(startTime)

			printOrder("Обход в ширину по матрице смежности:", order)
			fmt.Printf("Время выполнения: %v\n", elapsed)

		case 4:
			if g.N == 0 {
				fmt.Println("Сначала создайте граф (пункт 1).")
				continue
			}

			start, err := readInt(in, "Введите начальную вершину: ", 1, g.N)
			for err != nil {
				fmt.Println("Ошибка:", err)
				start, err = readInt(in, "Введите начальную вершину: ", 1, g.N)
			}
			start--

			startTime := time.Now()
			g.BFSLinkedQueue(start, &order)
			elapsed := time.Since(startTime)

			printOrder("Обход в ширину с самостоятельно реализованной очередью:", order)
			fmt.Printf("Время выполнения: %v\n", elapsed)

		case 5:
			if g.N == 0 {
				fmt.Println("Сначала создайте граф (пункт 1).")
				continue
			}

			start, err := readInt(in, "Введите начальную вершину: ", 1, g.N)
			for err != nil {
				fmt.Println("Ошибка:", err)
				start, err = readInt(in, "Введите начальную вершину: ", 1, g.N)
			}

			compareTime(g, start, true)

		case 6:
			g.PrintAdjacencyList()

		case 7:
			if g.N == 0 {
				fmt.Println("Сначала создайте граф (пункт 1).")
				continue
			}

			start, err := readInt(in, "Введите начальную вершину: ", 1, g.N)
			for err != nil {
				fmt.Println("Ошибка:", err)
				start, err = readInt(in, "Введите начальную вершину: ", 1, g.N)
			}
			start--

			startTime := time.Now()
			g.BFSAdjacencyList(start, &order)
			elapsed := time.Since(startTime)

			printOrder("Обход в ширину по спискам смежности:", order)
			fmt.Printf("Время выполнения: %v\n", elapsed)

		case 8:
			probability, err := readProbability(in)
			for err != nil {
				fmt.Println(err)
				probability, err = readProbability(in)
			}

			fmt.Printf(
				"Установлена вероятность %.2f%%. Она будет использована при следующей генерации графа.\n",
				probability,
			)
		}
	}
}
