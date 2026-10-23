package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Graph struct {
	N   int
	Mat [][]int
}

func NewGraph(n int, p float64) Graph {
	g := Graph{N: n, Mat: make([][]int, n)}
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

func readStartVertex(in *bufio.Reader, n int) int {
	for {
		fmt.Printf("Введите начальную вершину (1-%d): ", n)
		s, _ := in.ReadString('\n')
		v, err := strconv.Atoi(strings.TrimSpace(s))
		if err == nil && v >= 1 && v <= n {
			return v - 1
		}
		fmt.Println("Ошибка: введите целое число из указанного диапазона.")
	}
}

func (g Graph) DFS(v int, visited []bool, order []int, index *int) {
	visited[v] = true
	order[*index] = v
	*index++
	for to := 0; to < g.N; to++ {
		if v != to && g.Mat[v][to] == 1 && !visited[to] {
			g.DFS(to, visited, order, index)
		}
	}
}

func (g Graph) DFSRecursive(start int) []int {
	visited := make([]bool, g.N)
	order := make([]int, g.N)
	index := start
	g.DFS(start, visited, order, &index)
	for v := 0; v < g.N; v++ {
		if !visited[v] {
			g.DFS(v, visited, order, &index)
		}
	}
	return order
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

func (g Graph) DFSAdjacencyList(start int) []int {
	list := g.AdjacencyList()
	visited := make([]bool, g.N)
	order := make([]int, 0, g.N)
	var dfs func(int)
	dfs = func(v int) {
		visited[v] = true
		order = append(order, v)
		for _, to := range list[v] {
			if !visited[to] {
				dfs(to)
			}
		}
	}
	dfs(start)
	for v := 0; v < g.N; v++ {
		if !visited[v] {
			dfs(v)
		}
	}
	return order
}

func (g Graph) DFSIterative(start int) []int {
	visited := make([]bool, g.N)
	order := make([]int, 0, g.N)

	for v := 0; v < g.N; v++ {
		if !visited[v] {
			stack := []int{v}
			for len(stack) > 0 {
				last := len(stack) - 1
				v := stack[last]
				stack = stack[:last]
				if visited[v] {
					continue
				}
				visited[v] = true
				order = append(order, v)
				for to := g.N - 1; to >= 0; to-- {
					if g.Mat[v][to] == 1 && !visited[to] {
						stack = append(stack, to)
					}
				}
			}
		}
	}
	return order
}

func (g Graph) DFSIterativeList(start int) []int {
	visited := make([]bool, g.N)
	order := make([]int, 0, g.N)
	list := g.AdjacencyList()

	for v := 0; v < g.N; v++ {
		if !visited[v] {
			stack := []int{v}
			for len(stack) > 0 {
				last := len(stack) - 1
				v := stack[last]
				stack = stack[:last]
				if visited[v] {
					continue
				}
				visited[v] = true
				order = append(order, v)
				slices.Reverse(list[v])
				for _, to := range list[v] {
					if !visited[to] {
						stack = append(stack, to)
					}
				}
			}
		}
	}
	return order
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

func readInt(in *bufio.Reader, prompt string, min, max int) int {
	for {
		fmt.Print(prompt)
		s, _ := in.ReadString('\n')
		v, err := strconv.Atoi(strings.TrimSpace(s))
		if err == nil && v >= min && v <= max {
			return v
		}
		fmt.Printf("Введите целое число от %d до %d.\n", min, max)
	}
}

func readProbability(in *bufio.Reader) float64 {
	for {
		fmt.Print("Вероятность появления ребра (0-100%): ")
		s, _ := in.ReadString('\n')
		s = strings.TrimSpace(strings.ReplaceAll(s, ",", "."))
		p, err := strconv.ParseFloat(s, 64)
		if err == nil && p >= 0 && p <= 100 {
			return p
		}
		fmt.Println("Введите число от 0 до 100, например 30 или 75.5.")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	in := bufio.NewReader(os.Stdin)
	var g Graph
	var probability float64 = 30

	for {
		fmt.Println("\n========== ЛАБОРАТОРНАЯ РАБОТА №7 ==========")
		fmt.Println("Обход графа в глубину")
		fmt.Println("1. Создать случайный граф")
		fmt.Println("2. Вывести матрицу смежности")
		fmt.Println("3. Рекурсивный обход по матрице")
		fmt.Println("4. Обход по спискам смежности")
		fmt.Println("5. Нерекурсивный обход")
		fmt.Println("6. Вывести списки смежности")
		fmt.Println("7. Указать вероятность ребра")
		fmt.Println("8. Нерекурсивный обход списка")
		fmt.Println("0. Выход")
		choice := readInt(in, "Выберите пункт меню: ", 0, 8)

		switch choice {
		case 0:
			fmt.Println("Работа программы завершена.")
			return
		case 1:
			n := readInt(in, "Количество вершин (1-100): ", 1, 100)
			probability = readProbability(in)
			g = NewGraph(n, probability)
			fmt.Printf("Граф создан: вершин — %d, вероятность каждого ребра — %.2f%%.\n", n, probability)
			g.PrintMatrix()
		case 2:
			g.PrintMatrix()
		case 3:
			if g.N == 0 {
				fmt.Println("Сначала создайте граф (пункт 1).")
				continue
			}
			start := readStartVertex(in, g.N)
			startTime := time.Now()
			printOrder("Рекурсивный обход в глубину по матрице смежности:", g.DFSRecursive(start))
			elapsed := time.Since(startTime)
			fmt.Printf("Время выполнения: %v\n", elapsed)
		case 4:
			if g.N == 0 {
				fmt.Println("Сначала создайте граф (пункт 1).")
				continue
			}
			start := readStartVertex(in, g.N)
			startTime := time.Now()
			printOrder("Обход в глубину по спискам смежности:", g.DFSAdjacencyList(start))
			elapsed := time.Since(startTime)
			fmt.Printf("Время выполнения: %v\n", elapsed)
		case 5:
			if g.N == 0 {
				fmt.Println("Сначала создайте граф (пункт 1).")
				continue
			}
			start := readStartVertex(in, g.N)
			startTime := time.Now()
			printOrder("Нерекурсивный обход в глубину с помощью стека:", g.DFSIterative(start))
			elapsed := time.Since(startTime)
			fmt.Printf("Время выполнения: %v\n", elapsed)
		case 6:
			g.PrintAdjacencyList()
		case 7:
			probability = readProbability(in)
			fmt.Printf("Установлена вероятность %.2f%%. Она будет использована при следующей генерации графа.\n", probability)
		case 8:
			if g.N == 0 {
				fmt.Println("Сначала создайте граф (пункт 1).")
				continue
			}
			start := readStartVertex(in, g.N)
			startTime := time.Now()
			printOrder("Нерекурсивный обход в глубину с помощью стека:", g.DFSIterativeList(start))
			elapsed := time.Since(startTime)
			fmt.Printf("Время выполнения: %v\n", elapsed)
		}
	}
}
