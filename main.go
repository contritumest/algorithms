package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	P = 0.4
)

type Graph struct {
	Names []string // метки вершин
	M     [][]int  // матрица смежности
}

type ListGraph struct {
	Names []string
	L     [][]int
}

var reader = bufio.NewReader(os.Stdin)

func newMatrix(n int) [][]int {
	m := make([][]int, n)
	for i := range m {
		m[i] = make([]int, n)
	}
	return m
}

func randomGraph(n int) *Graph {
	g := &Graph{Names: make([]string, n), M: newMatrix(n)}
	for i := 0; i < n; i++ {
		g.Names[i] = fmt.Sprintf("v%d", i+1)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Float64() < 0.4 {
				g.M[i][j] = 1
				g.M[j][i] = 1
			}
		}
	}
	return g
}

func toList(g *Graph) *ListGraph {
	n := len(g.M)
	lg := &ListGraph{Names: append([]string{}, g.Names...), L: make([][]int, n)}
	for i := 0; i < n; i++ {
		lg.L[i] = []int{}
		for j := 0; j < n; j++ {
			if g.M[i][j] == 1 {
				lg.L[i] = append(lg.L[i], j)
			}
		}
	}
	return lg
}

func printMatrix(g *Graph) {
	n := len(g.M)
	if n == 0 {
		fmt.Println("  (пустой граф)")
		return
	}
	w := 0
	for _, s := range g.Names {
		if l := len([]rune(s)); l > w {
			w = l
		}
	}
	fmt.Printf("%*s |", w, "")
	for _, s := range g.Names {
		fmt.Printf(" %*s", w, s)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", (w+1)*(n+1)+2))
	for i := 0; i < n; i++ {
		fmt.Printf("%*s |", w, g.Names[i])
		for j := 0; j < n; j++ {
			fmt.Printf(" %*d", w, g.M[i][j])
		}
		fmt.Println()
	}
}

func printList(lg *ListGraph) {
	if len(lg.L) == 0 {
		fmt.Println("  (пустой граф)")
		return
	}
	for i, nb := range lg.L {
		parts := make([]string, len(nb))
		for k, j := range nb {
			parts[k] = lg.Names[j]
		}
		fmt.Printf("  %s: %s\n", lg.Names[i], strings.Join(parts, ", "))
	}
}

func identifyMatrix(g *Graph, u, v int) *Graph {
	n := len(g.M)
	keep := []int{}
	for i := 0; i < n; i++ {
		if i != u && i != v {
			keep = append(keep, i)
		}
	}
	k := len(keep)
	res := &Graph{Names: make([]string, k+1), M: newMatrix(k + 1)}
	for a, i := range keep {
		res.Names[a] = g.Names[i]
		for b, j := range keep {
			res.M[a][b] = g.M[i][j]
		}
		if g.M[u][i] == 1 || g.M[v][i] == 1 {
			res.M[a][k] = 1
			res.M[k][a] = 1
		}
	}
	res.Names[k] = "z"
	return res
}

func contractMatrix(g *Graph, u, v int) (*Graph, error) {
	if g.M[u][v] != 1 {
		return nil, fmt.Errorf("вершины %s и %s не смежны - ребра нет", g.Names[u], g.Names[v])
	}
	return identifyMatrix(g, u, v), nil
}

func splitMatrix(g *Graph, v int) *Graph {
	n := len(g.M)
	res := &Graph{Names: make([]string, n+1), M: newMatrix(n + 1)}
	for i := 0; i < n; i++ {
		res.Names[i] = g.Names[i]
		for j := 0; j < n; j++ {
			res.M[i][j] = g.M[i][j]
		}
	}
	res.Names[v] = g.Names[v] + "'"
	res.Names[n] = g.Names[v] + "''"
	cnt := 0
	for w := 0; w < n; w++ {
		if g.M[v][w] == 1 {
			if cnt%2 == 1 {
				res.M[v][w], res.M[w][v] = 0, 0
				res.M[n][w], res.M[w][n] = 1, 1
			}
			cnt++
		}
	}
	res.M[v][n], res.M[n][v] = 1, 1
	return res
}

func sortedSet(set map[int]bool) []int {
	r := make([]int, 0, len(set))
	for x := range set {
		r = append(r, x)
	}
	sort.Ints(r)
	return r
}

func identifyList(lg *ListGraph, u, v int) *ListGraph {
	n := len(lg.L)
	newIdx := make(map[int]int)
	keep := []int{}
	for i := 0; i < n; i++ {
		if i != u && i != v {
			newIdx[i] = len(keep)
			keep = append(keep, i)
		}
	}
	k := len(keep)
	q := make(map[int]bool)
	for _, w := range lg.L[u] {
		if w != u && w != v {
			q[w] = true
		}
	}
	for _, w := range lg.L[v] {
		if w != u && w != v {
			q[w] = true
		}
	}
	res := &ListGraph{Names: make([]string, k+1), L: make([][]int, k+1)}
	for a, i := range keep {
		res.Names[a] = lg.Names[i]
		res.L[a] = []int{}
		for _, w := range lg.L[i] {
			if w != u && w != v {
				res.L[a] = append(res.L[a], newIdx[w])
			}
		}
		if q[i] {
			res.L[a] = append(res.L[a], k)
		}
	}
	res.Names[k] = "z"
	z := make(map[int]bool)
	for w := range q {
		z[newIdx[w]] = true
	}
	res.L[k] = sortedSet(z)
	return res
}

func contractList(lg *ListGraph, u, v int) (*ListGraph, error) {
	for _, w := range lg.L[u] {
		if w == v {
			return identifyList(lg, u, v), nil
		}
	}
	return nil, fmt.Errorf("вершины %s и %s не смежны - ребра нет", lg.Names[u], lg.Names[v])
}

func splitList(lg *ListGraph, v int) *ListGraph {
	n := len(lg.L)
	res := &ListGraph{Names: make([]string, n+1), L: make([][]int, n+1)}
	for i := 0; i < n; i++ {
		res.Names[i] = lg.Names[i]
		res.L[i] = append([]int{}, lg.L[i]...)
	}
	res.Names[v] = lg.Names[v] + "'"
	res.Names[n] = lg.Names[v] + "''"
	res.L[n] = []int{}

	stay := []int{}
	for idx, w := range lg.L[v] {
		if idx%2 == 1 {
			for t, x := range res.L[w] {
				if x == v {
					res.L[w][t] = n
				}
			}
			sort.Ints(res.L[w])
			res.L[n] = append(res.L[n], w)
		} else {
			stay = append(stay, w)
		}
	}
	res.L[v] = append(stay, n)
	res.L[n] = append(res.L[n], v)
	sort.Ints(res.L[v])
	sort.Ints(res.L[n])
	return res
}

func unionGraph(a, b *Graph) *Graph {
	n := max(len(a.M), len(b.M))
	res := &Graph{Names: make([]string, n), M: newMatrix(n)}
	for i := 0; i < n; i++ {
		res.Names[i] = fmt.Sprintf("v%d", i+1)
		for j := 0; j < n; j++ {
			if (i < len(a.M) && j < len(a.M) && a.M[i][j] == 1) ||
				(i < len(b.M) && j < len(b.M) && b.M[i][j] == 1) {
				res.M[i][j] = 1
			}
		}
	}
	return res
}
func intersectGraph(a, b *Graph) *Graph {
	n := min(len(a.M), len(b.M))
	res := &Graph{Names: make([]string, n), M: newMatrix(n)}
	for i := 0; i < n; i++ {
		res.Names[i] = fmt.Sprintf("v%d", i+1)
		for j := 0; j < n; j++ {
			if a.M[i][j] == 1 && b.M[i][j] == 1 {
				res.M[i][j] = 1
			}
		}
	}
	return res
}

func ringSum(a, b *Graph) *Graph {
	n := max(len(a.M), len(b.M))
	full := newMatrix(n)
	get := func(g *Graph, i, j int) int {
		if i < len(g.M) && j < len(g.M) {
			return g.M[i][j]
		}
		return 0
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			full[i][j] = get(a, i, j) ^ get(b, i, j)
		}
	}
	keep := []int{}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if full[i][j] == 1 {
				keep = append(keep, i)
				break
			}
		}
	}
	res := &Graph{Names: make([]string, len(keep)), M: newMatrix(len(keep))}
	for x, i := range keep {
		res.Names[x] = fmt.Sprintf("v%d", i+1)
		for y, j := range keep {
			res.M[x][y] = full[i][j]
		}
	}
	return res
}

func cartesianProduct(a, b *Graph) *Graph {
	n1, n2 := len(a.M), len(b.M)
	n := n1 * n2
	res := &Graph{Names: make([]string, n), M: newMatrix(n)}
	for i := 0; i < n1; i++ {
		for k := 0; k < n2; k++ {
			res.Names[i*n2+k] = fmt.Sprintf("(%d,%d)", i+1, k+1)
		}
	}
	for i := 0; i < n1; i++ {
		for k := 0; k < n2; k++ {
			for j := 0; j < n1; j++ {
				for l := 0; l < n2; l++ {
					if (i == j && b.M[k][l] == 1) || (k == l && a.M[i][j] == 1) {
						res.M[i*n2+k][j*n2+l] = 1
					}
				}
			}
		}
	}
	return res
}

func pickGraph(g1, g2 *Graph) (*Graph, string) {
	num := 0
	fmt.Scanf("%d", &num)
	if num == 1 {
		return g1, "G1"
	}
	return g2, "G2"
}

func requireGraphs(g1, g2 *Graph) bool {
	if g1 == nil || g2 == nil {
		fmt.Println("Сначала сгенерируйте графы")
		return false
	}
	return true
}

func showGraphs(g1, g2 *Graph) {
	fmt.Println("\nМатрица смежности G1:")
	printMatrix(g1)
	fmt.Println("\nМатрица смежности G2:")
	printMatrix(g2)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var g1, g2 *Graph

	for {
		fmt.Println("1. Сгенерировать графы")
		fmt.Println("2. Показать матрицы смежности")
		fmt.Println("3. Показать списки смежности")
		fmt.Println("4. Унарные операции (матрицы): отождествление / стягивание / расщепление")
		fmt.Println("5. Унарные операции (списки): отождествление / стягивание / расщепление")
		fmt.Println("6. Объединение, пересечение, кольцевая сумма")
		fmt.Println("7. Декартово произведение")
		fmt.Println("0. Выход")
		var choice int
		fmt.Printf("Ваш выбор: ")
		fmt.Scanf("%d", &choice)

		if choice == 0 {
			fmt.Println("Выход.")
			return
		}
		if choice != 1 && !requireGraphs(g1, g2) {
			continue
		}

		switch choice {
		case 1:
			var n1, n2 int
			fmt.Printf("Число вершин G1 (1..15): ")
			fmt.Scanf("%d", &n1)
			fmt.Printf("Число вершин G2 (1..15): ")
			fmt.Scanf("%d", &n2)
			g1, g2 = randomGraph(n1), randomGraph(n2)
			showGraphs(g1, g2)

		case 2:
			showGraphs(g1, g2)

		case 3:
			fmt.Println("\nСписки смежности G1:")
			printList(toList(g1))
			fmt.Println("\nСписки смежности G2:")
			printList(toList(g2))

		case 4:
			g, name := pickGraph(g1, g2)
			fmt.Printf("\nГраф %s:\n", name)
			printMatrix(g)
			n := len(g.M)
			fmt.Println("\n  1. Отождествление вершин\n  2. Стягивание ребра\n  3. Расщепление вершины")
			var op int
			fmt.Printf("Операция: ")
			fmt.Scanf("%d", &op)
			switch op {
			case 1, 2:
				if n < 2 {
					fmt.Println("Нужно минимум 2 вершины.")
					continue
				}
				var u, v int
				fmt.Printf("Вершина u (1..%d): ", n)
				fmt.Scanf("%d", &u)
				u--
				fmt.Printf("Вершина v (1..%d): ", n)
				fmt.Scanf("%d", &v)
				v--
				if u == v {
					fmt.Println("Нужно выбрать две разные вершины.")
					continue
				}
				if op == 1 {
					fmt.Println("\nРезультат отождествления:")
					printMatrix(identifyMatrix(g, u, v))
				} else {
					res, err := contractMatrix(g, u, v)
					if err != nil {
						fmt.Println("Ошибка:", err)
						continue
					}
					fmt.Println("\nРезультат стягивания ребра:")
					printMatrix(res)
				}
			case 3:
				var v int
				fmt.Printf("Вершина для расщепления (1..%d): ", n)
				fmt.Scanf("%d", &v)
				v--
				fmt.Println("\nРезультат расщепления:")
				printMatrix(splitMatrix(g, v))
			}

		case 5:
			g, name := pickGraph(g1, g2)
			lg := toList(g)
			fmt.Printf("\nГраф %s (списки смежности):\n", name)
			printList(lg)
			n := len(lg.L)
			fmt.Println("\n  1. Отождествление вершин\n  2. Стягивание ребра\n  3. Расщепление вершины")
			var op int
			fmt.Printf("Операция: ")
			fmt.Scanf("%d", &op)
			switch op {
			case 1, 2:
				if n < 2 {
					fmt.Println("Нужно минимум 2 вершины.")
					continue
				}
				var u, v int
				fmt.Printf("Вершина u (1..%d): ", n)
				fmt.Scanf("%d", &u)
				u--
				fmt.Printf("Вершина v (1..%d): ", n)
				fmt.Scanf("%d", &v)
				v--
				if u == v {
					fmt.Println("Нужно выбрать две разные вершины.")
					continue
				}
				if op == 1 {
					fmt.Println("\nРезультат отождествления:")
					printList(identifyList(lg, u, v))
				} else {
					res, err := contractList(lg, u, v)
					if err != nil {
						fmt.Println("Ошибка:", err)
						continue
					}
					fmt.Println("\nРезультат стягивания ребра:")
					printList(res)
				}
			case 3:
				var v int
				fmt.Printf("Вершина для расщепления (1..%d): ", n)
				fmt.Scanf("%d", &v)
				v--
				fmt.Println("\nРезультат расщепления:")
				printList(splitList(lg, v))
			}

		case 6:
			showGraphs(g1, g2)
			fmt.Println("\n  1. Объединение G1 ∪ G2\n  2. Пересечение G1 ∩ G2\n  3. Кольцевая сумма G1 ⊕ G2\n  4. Все три операции")
			var op int
			fmt.Printf("Операция: ")
			fmt.Scanf("%d", &op)
			if op == 1 || op == 4 {
				fmt.Println("\nG1 ∪ G2:")
				printMatrix(unionGraph(g1, g2))
			}
			if op == 2 || op == 4 {
				fmt.Println("\nG1 ∩ G2:")
				printMatrix(intersectGraph(g1, g2))
			}
			if op == 3 || op == 4 {
				fmt.Println("\nG1 ⊕ G2 (без изолированных вершин):")
				printMatrix(ringSum(g1, g2))
			}

		case 7:
			if len(g1.M)*len(g2.M) > 100 {
				fmt.Println("Результат слишком велик для вывода (более 100 вершин). Сгенерируйте графы меньшего размера.")
				continue
			}
			showGraphs(g1, g2)
			fmt.Println("\nG1 x G2 (вершина (i,k) = (zi,vk)):")
			printMatrix(cartesianProduct(g1, g2))
		}
	}
}
