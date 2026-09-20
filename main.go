package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func clearConsole() {
	// Определяем команду в зависимости от ОС
	var cmdName string
	if runtime.GOOS == "windows" {
		cmdName = "cls"
	} else {
		cmdName = "clear"
	}

	// Создаём команду и запускаем
	cmd := exec.Command(cmdName)
	cmd.Stdout = os.Stdout // Важно: передаём стандартный вывод в тот же поток, иначе ничего не отобразится
	cmd.Run()
}

type PriorityQueue struct {
	head *Node
	last *Node
	size int
}

func (pq *PriorityQueue) Push(inf string, prior int) {
	newNode := &Node{Inf: inf, Prior: prior}
	if pq.head == nil {
		pq.head = newNode
		pq.last = newNode
	} else {
		current := pq.head
		for current != nil && current.Prior <= prior {
			current = current.Next
		}
		if current == pq.head {
			newNode.Next = pq.head
			pq.head = newNode
		} else {
			prev := pq.head
			for prev.Next != current {
				prev = prev.Next
			}
			prev.Next = newNode
			newNode.Next = current
		}
	}
	pq.size++
}

func (pq *PriorityQueue) Pop() (string, int) {
	if pq.head == nil {
		return "", -1
	}
	inf := pq.head.Inf
	prior := pq.head.Prior
	pq.head = pq.head.Next
	pq.size--

	if pq.head == nil {
		pq.last = nil
	}

	return inf, prior
}

func (pq *PriorityQueue) Size() int {
	return pq.size
}

func (pq *PriorityQueue) Review() {
	current := pq.head
	for current != nil {
		println("Inf:", current.Inf, "Prior:", current.Prior)
		current = current.Next
	}
}

func (pq *PriorityQueue) Find(name string) *Node {
	current := pq.head
	for current != nil {
		if current.Inf == name {
			return current
		}
		current = current.Next
	}
	return nil
}

func (pq *PriorityQueue) Delete(name string) bool {
	if pq.head == nil {
		return false
	}

	if pq.head.Inf == name {
		pq.head = pq.head.Next
		pq.size--
		if pq.head == nil {
			pq.last = nil
		}
		return true
	}

	prev := pq.head
	current := pq.head.Next

	for current != nil {
		if current.Inf == name {
			prev.Next = current.Next
			pq.size--
			if current == pq.last {
				pq.last = prev
			}
			return true
		}
		prev = current
		current = current.Next
	}

	return false
}

var reader = bufio.NewReader(os.Stdin)

func readLine(prompt string) string {
	fmt.Print(prompt)
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}

func readInt(prompt string) int {
	s := readLine(prompt)
	n, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return n
}

func menuPriorityQueue() {
	pq := &PriorityQueue{}

	for {
		fmt.Println("\n[Приоритетная очередь]")
		fmt.Println("1 - Добавить элемент")
		fmt.Println("2 - Извлечь (макс. приоритет)")
		fmt.Println("3 - Просмотр")
		fmt.Println("4 - Поиск")
		fmt.Println("5 - Удаление")
		fmt.Println("0 - Назад")

		choice := readInt("Выбор: ")

		switch choice {
		case 1:
			clearConsole()
			name := readLine("Введите название объекта: ")
			if name == "" {
				fmt.Println("Запись не была произведена")
				continue
			}
			prior := readInt("Введите приоритет: ")
			pq.Push(name, prior)

		case 2:
			clearConsole()
			if p, ok := pq.Pop(); p != "" && ok != -1 {
				fmt.Printf("Извлечён элемент: %s (приоритет %d)\n", p, ok)
			} else {
				fmt.Println("Очередь пуста")
			}

		case 3:
			clearConsole()
			pq.Review()

		case 4:
			clearConsole()
			name := readLine("Введите имя: ")
			if node := pq.Find(name); node != nil {
				fmt.Printf("Найден: %s (приоритет %d)\n", node.Inf, node.Prior)
			} else {
				fmt.Println("Элемент не найден")
			}

		case 5:
			clearConsole()
			name := readLine("Введите имя: ")
			if pq.Delete(name) {
				fmt.Println("Удалено")
			} else {
				fmt.Println("Элемент не найден")
			}

		case 0:
			clearConsole()
			return

		default:
			clearConsole()
			fmt.Println("Неверный выбор")
		}
	}
}

func menuQueue() {
	q := &Queue{}

	for {
		fmt.Println("\n[Очередь]")
		fmt.Println("1 - Enqueue (в конец)")
		fmt.Println("2 - Dequeue (из начала)")
		fmt.Println("3 - Просмотр")
		fmt.Println("0 - Назад")

		choice := readInt("Выбор: ")

		switch choice {
		case 1:
			clearConsole()
			name := readLine("Введите название объекта: ")
			if name == "" {
				fmt.Println("Запись не была произведена")
				continue
			}
			q.Enqueue(name)

		case 2:
			clearConsole()
			if str := q.Dequeue(); str != "" {
				fmt.Printf("Извлечён элемент: %s\n", str)
			} else {
				fmt.Println("Очередь пуста")
			}

		case 3:
			clearConsole()
			q.Review()

		case 0:
			clearConsole()
			return

		default:
			clearConsole()
			fmt.Println("Неверный выбор")
		}
	}
}

func menuStack() {
	s := &Stack{}

	for {
		fmt.Println("\n[Стек]")
		fmt.Println("1 - Push (на вершину)")
		fmt.Println("2 - Pop (с вершины)")
		fmt.Println("3 - Peek (вершина)")
		fmt.Println("4 - Просмотр")
		fmt.Println("0 - Назад")

		choice := readInt("Выбор: ")

		switch choice {
		case 1:
			clearConsole()
			name := readLine("Введите название объекта: ")
			if name == "" {
				fmt.Println("Запись не была произведена")
				continue
			}
			s.Push(name)

		case 2:
			clearConsole()
			if str := s.Pop(); str != "" {
				fmt.Printf("Извлечён элемент: %s\n", str)
			} else {
				fmt.Println("Стек пуст")
			}

		case 3:
			clearConsole()
			if p, ok := s.Peek(); ok {
				fmt.Printf("Вершина стека: %s\n", p.Inf)
			} else {
				fmt.Println("Стек пуст")
			}

		case 4:
			clearConsole()
			s.Review()

		case 0:
			clearConsole()
			return

		default:
			clearConsole()
			fmt.Println("Неверный выбор")
		}
	}
}

func main() {
	for {
		fmt.Println("\n=== Главное меню ===")
		fmt.Println("1 - Приоритетная очередь")
		fmt.Println("2 - Очередь (FIFO)")
		fmt.Println("3 - Стек (LIFO)")
		fmt.Println("0 - Выход")

		choice := readInt("Выбор: ")

		switch choice {
		case 1:
			clearConsole()
			menuPriorityQueue()
		case 2:
			clearConsole()
			menuQueue()
		case 3:
			clearConsole()
			menuStack()
		case 0:
			clearConsole()
			fmt.Println("Выход...")
			return
		default:
			clearConsole()
			fmt.Println("Неверный выбор")
		}
	}
}
