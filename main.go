package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func clearConsole() {
	var cmdName string
	if runtime.GOOS == "windows" {
		cmdName = "cls"
	} else {
		cmdName = "clear"
	}

	cmd := exec.Command(cmdName)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

type Node struct {
	data  int
	left  *Node
	right *Node
}

func CreateTree(root *Node, data int) *Node {
	if root == nil {
		return &Node{data: data}
	}

	if data == root.data { // Excluding the addition of identical symbols
		return root
	} else if data > root.data {
		root.left = CreateTree(root.left, data)
	} else {
		root.right = CreateTree(root.right, data)
	}

	return root
}

func printTree(r *Node, l int) {
	if r == nil {
		return
	}

	printTree(r.right, l+1)

	for i := 0; i < l; i++ {
		fmt.Print(" ")
	}

	fmt.Println(r.data)

	printTree(r.left, l+1)
}

func search(r *Node, data int) *Node {
	if r == nil {
		return nil
	}

	if r.data == data {
		return r
	}

	if data > r.data {
		return search(r.left, data)
	} else {
		return search(r.right, data)
	}
}

func countOccurrences(r *Node, data int) int {
	if r == nil {
		return 0
	}

	if r.data == data {
		return 1 + countOccurrences(r.left, data) + countOccurrences(r.right, data)
	}

	if data > r.data {
		return countOccurrences(r.left, data)
	} else {
		return countOccurrences(r.right, data)
	}
}

func main() {
	var root *Node = nil
	var choice int
	var value int

	for {
		fmt.Println("\n========== МЕНЮ ==========")
		fmt.Println("1. Добавить элемент")
		fmt.Println("2. Найти элемент")
		fmt.Println("3. Вывести дерево")
		fmt.Println("4. Подсчитать количество вхождений")
		fmt.Println("0. Выход")
		fmt.Println("==========================")
		fmt.Print("Выберите действие: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			clearConsole()
			fmt.Print("Введите число для добавления: ")
			fmt.Scan(&value)

			if search(root, value) != nil {
				fmt.Printf("Элемент %d уже существует.\n", value)
			} else {
				root = CreateTree(root, value)
				fmt.Printf("Элемент %d добавлен.\n", value)
			}

		case 2:
			clearConsole()
			fmt.Print("Введите число для поиска: ")
			fmt.Scan(&value)

			result := search(root, value)
			if result != nil {
				fmt.Printf("Элемент %d найден в дереве.\n", value)
			} else {
				fmt.Printf("Элемент %d не найден в дереве.\n", value)
			}
			fmt.Scan()

		case 3:
			clearConsole()
			if root == nil {
				fmt.Println("Дерево пустое.")
			} else {
				fmt.Println("\nДерево:")
				printTree(root, 0)
			}
			fmt.Scan()

		case 4:
			clearConsole()
			fmt.Print("Введите число: ")
			fmt.Scan(&value)

			count := countOccurrences(root, value)
			fmt.Printf(
				"Количество вхождений элемента %d: %d\n",
				value, count,
			)
			fmt.Scan()

		case 0:
			clearConsole()
			fmt.Println("Программа завершена.")
			return

		default:
			fmt.Println("Ошибка: такого пункта меню нет.")
		}
	}
}
