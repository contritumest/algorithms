package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	N = 6   // количество вершин
	P = 0.4 // вероятность наличия ребра
)

// Генерация матрицы смежности для неориентированного графа
func generateAdjacencyMatrix(n int, p float64, rng *rand.Rand) [][]int {
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Float64() <= p {
				matrix[i][j] = 1
				matrix[j][i] = 1
			}
		}
	}
	return matrix
}

func printMatrix(matrix [][]int, rowLabels []string, colLabels []string, title string) {
	fmt.Println()
	fmt.Println(title)

	// Заголовок столбцов
	fmt.Print("     ")
	for _, c := range colLabels {
		fmt.Printf("%4s", c)
	}
	fmt.Println()

	// Строки
	for i, row := range matrix {
		fmt.Printf("%4s ", rowLabels[i])
		for _, v := range row {
			fmt.Printf("%4d", v)
		}
		fmt.Println()
	}
}

// Степени вершин по матрице смежности
func degreesFromAdjacency(matrix [][]int) []int {
	n := len(matrix)
	deg := make([]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			deg[i] += matrix[i][j]
		}
	}
	return deg
}

// Классификация вершин
type VertexClassification struct {
	Isolated   []int
	Terminal   []int
	Dominating []int
}

func analyzeVertices(deg []int, n int) VertexClassification {
	var res VertexClassification
	for i, d := range deg {
		switch {
		case d == 0:
			res.Isolated = append(res.Isolated, i+1)
		case d == 1:
			res.Terminal = append(res.Terminal, i+1)
		}
		if d == n-1 {
			res.Dominating = append(res.Dominating, i+1)
		}
	}
	return res
}

type Edge struct {
	U int
	V int
}

func buildIncidenceMatrix(adj [][]int) ([][]int, []Edge) {
	n := len(adj)
	var edges []Edge
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if adj[i][j] == 1 {
				edges = append(edges, Edge{i, j})
			}
		}
	}

	m := len(edges)
	inc := make([][]int, n)
	for i := range inc {
		inc[i] = make([]int, m)
	}
	for k, e := range edges {
		inc[e.U][k] = 1
		inc[e.V][k] = 1
	}
	return inc, edges
}

func printIntSlice(label string, s []int) {
	if len(s) == 0 {
		fmt.Printf("%s нет\n", label)
		return
	}
	fmt.Printf("%s %v\n", label, s)
}

func sizeGraph(adj [][]int) int {
	sum := 0
	for _, row := range adj {
		for _, v := range row {
			sum += v
		}
	}
	return sum
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Метки вершин
	vertexLabels := make([]string, N)
	for i := range vertexLabels {
		vertexLabels[i] = fmt.Sprintf("v%d", i+1)
	}

	// Матрица смежности
	fmt.Println("============================================================")
	fmt.Println("ЗАДАНИЕ 1. Матрица смежности")
	fmt.Println("============================================================")

	adj := generateAdjacencyMatrix(N, P, rng)
	printMatrix(adj, vertexLabels, vertexLabels, "Матрица смежности графа G:")

	// Размер графа
	sizeAdj := sizeGraph(adj) / 2
	fmt.Printf("\nРазмер графа = %d\n", sizeAdj)

	// Степени вершин
	degAdj := degreesFromAdjacency(adj)
	fmt.Println("\nСтепени вершин:")
	for i, d := range degAdj {
		fmt.Printf("  deg(v%d) = %d\n", i+1, d)
	}

	cls1 := analyzeVertices(degAdj, N)
	fmt.Println()
	printIntSlice("Изолированные вершины:  ", cls1.Isolated)
	printIntSlice("Концевые вершины:       ", cls1.Terminal)
	printIntSlice("Доминирующие вершины:   ", cls1.Dominating)

	// Матрица инцидентности
	fmt.Println("\n============================================================")
	fmt.Println("ЗАДАНИЕ 2. Матрица инцидентности")
	fmt.Println("============================================================")

	inc, edges := buildIncidenceMatrix(adj)

	// Метки столбцов
	edgeLabels := make([]string, len(edges))
	for i := range edgeLabels {
		edgeLabels[i] = fmt.Sprintf("e%d", i+1)
	}

	printMatrix(inc, vertexLabels, edgeLabels, "Матрица инцидентности:")

	fmt.Println("\nСписок рёбер:")
	for k, e := range edges {
		fmt.Printf("  e%d = (v%d, v%d)\n", k+1, e.U+1, e.V+1)
	}

	// Размер графа
	sizeInc := len(edges)
	fmt.Printf("\nРазмер графа = %d \n", sizeInc)

	// Степени вершин
	degInc := make([]int, N)
	for i, row := range inc {
		for _, v := range row {
			degInc[i] += v
		}
	}
	fmt.Println("\nСтепени вершин:")
	for i, d := range degInc {
		fmt.Printf("  deg(v%d) = %d\n", i+1, d)
	}

	cls2 := analyzeVertices(degInc, N)
	fmt.Println()
	printIntSlice("Изолированные вершины:  ", cls2.Isolated)
	printIntSlice("Концевые вершины:       ", cls2.Terminal)
	printIntSlice("Доминирующие вершины:   ", cls2.Dominating)
}
