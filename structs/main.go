package main

import (
	"errors"
	"fmt"

	"temp/types"
)

type Address struct {
	Name   string
	Street string
}

type Person struct {
	Address
	Age int
}

type Company struct {
	Name     string
	Location string
}

type Employee struct {
	Company
	Position string
}

type Combined struct {
	Person
	Employee
}

func main() {
	comb := Combined{
		Person: Person{
			Address: Address{
				Name:   "John Doe",
				Street: "123 Main St",
			},
			Age: 30,
		},
		Employee: Employee{
			Company: Company{
				Name:     "ACME Inc",
				Location: "New York",
			},
			Position: "Engineer",
		},
	}

	fmt.Println("Person Name:", comb.Name)
}

//func main() {
//	var err error
//	println("address err outside:", &err)
//
//	if true {
//		err := errors.New("ошибка в блоке if") // Затенение переменной err
//		println("address err inside:", &err)
//		fmt.Println("Внутри блока if:", err)
//	}
//
//	fmt.Println("Вне блока if:", err) // Ожидается, что err будет nil, но это не так
//	println("address err outside:", &err)
//}

// 106eb8
// func main() {
// 	var counter atomic.Int32
// 	var wg sync.WaitGroup

// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)

// 		go func() {
// 			defer wg.Done()

// 			counter.Add(1)
// 		}()
// 	}

// 	wg.Wait()

//		println(counter.Load())
//	}
func createChan(n int) chan int {
	ch := make(chan int, 1) // создаем канал
	ch <- n                 // отправляем данные в канал
	return ch               // возвращаем канал
}

func someFunc() error {
	return errors.New("some error")
}

type Flight [2]string

type Node struct {
	Value int
	Next  *Node
}

func createList() *Node {
	head := &Node{Value: 1}
	head.Next = &Node{Value: 2}
	head.Next.Next = &Node{Value: 3}
	head.Next.Next.Next = &Node{Value: 4}
	head.Next.Next.Next.Next = &Node{Value: 5}
	return head
}

//type Reader interface {
//	Read(p []byte) (n int, err error)
//}

//type ReadWriter interface {
//	Reader
//	Writer
//}

//{
//{1,0,0,1,0},
//{0,0,0,1,0},
//{1,1,0,1,0},
//{0,0,0,1,0},
//}

// 3

// func main() {
// 	mp := Field{
// 		mp: [][]int{
// 			{0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
// 			{0, 0, 0, 1, 1, 0, 0, 1, 0, 0},
// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
// 			{1, 0, 0, 0, 0, 0, 1, 1, 1, 0},
// 			{0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
// 			{0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
// 			{0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
// 			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
// 		},
// 	}

// 	fmt.Println(mp.findShipCount())
// }

type Field struct {
	mp [][]int
}

// func (f Field) findShipCount() int {
// 	var count int

// 	for i := 0; i < len(f.mp); i++ {
// 		for j := 0; j < len(f.mp); j++ {
// 			if f.mp[i][j] == 1 {
// 				f.mp[i][j] = 0

// 				rightNeiborhood := f.bfs(i, j-1) // true
// 				leftNeiborhood := f.bfs(i, j+1)  // false
// 				upNeiborhood := f.bfs(i-1, j)    // true
// 				downNeiborhood := f.bfs(i+1, j)  // true

// 				if upNeiborhood && downNeiborhood && leftNeiborhood && rightNeiborhood {
// 					count++
// 				}
// 			}
// 		}
// 	}

// 	return count
// }

// func (f Field) bfs(i, j int) bool {
// 	if i < 0 || j < 0 || j > len(f.mp[0])-1 || i > len(f.mp)-1 {
// 		return true
// 	}

// 	if f.mp[i][j] == 0 {
// 		return true
// 	}

// 	f.mp[i][j] = 0

// 	return f.bfs(i, j-1) && f.bfs(i, j+1) && f.bfs(i-1, j) && f.bfs(i+1, j)
// }

//
//func worker() chan int {
//	ch := make(chan int, 1)
//
//	go func() {
//		time.Sleep(3 * time.Second)
//
//		ch <- 1
//	}()
//
//	return ch
//}

//func main() {
//	//paddedNumber := fmt.Sprintf("%010d", 43)
//	//fmt.Println(paddedNumber)
//	fmt.Println((0.2 + 0.1) == 0.3)
//}

//
//list := createList()
//
//list = reverse(list)
//
//tmp := list
//
//for tmp != nil {
//fmt.Println(tmp.Value)
//tmp = tmp.Next
//}

func reverse(head *Node) *Node {
	if head == nil {
		return head
	}

	var tmp, prev, next *Node

	tmp = head

	for tmp != nil {
		next = tmp.Next
		tmp.Next = prev
		prev = tmp
		tmp = next
	}

	return prev
}

//plusOne([]int{7, 2, 8, 5, 0, 9, 1, 2, 9, 5, 3, 6, 6, 7, 3, 2, 8, 4, 3, 7, 9, 5, 7, 7, 4, 7, 4, 9, 4, 7, 0, 1, 1, 1, 7, 4, 0, 0, 9})
// plusOne([]int{9, 9, 9})

// flights := []Flight{
// 	{"Москва", "Белград"},
// 	{"Ереван", "Москва"},
// }
// fmt.Println(findPath(flights)) // Output: [A B C D]

//deferArgsValue()

func deferArgsValue() {
	nums := 1 << 5 // 32

	defer fmt.Println(nums)

	nums = nums >> 1 //16

	fmt.Println("done")
}

func plusOne(digits []int) []int {
	n := len(digits)

	// Start from the least significant digit (rightmost)
	for i := n - 1; i >= 0; i-- {
		if digits[i] < 9 {
			// If the digit is less than 9, simply increment it and return
			digits[i]++
			return digits
		}
		// If the digit is 9, set it to 0 and continue to the next digit
		digits[i] = 0
	}

	// If all digits were 9, we need to add an extra digit at the beginning
	newDigits := make([]int, n+1)
	newDigits[0] = 1
	return newDigits
}

func findPath(flights []Flight) []string {
	// Создаем карту для хранения городов и их соседей
	adjacencyMap := make(map[string]string)
	incomingMap := make(map[string]int)

	// Заполняем карту соседями и считаем входящие ребра
	for _, flight := range flights {
		from, to := flight[0], flight[1]
		adjacencyMap[from] = to

		incomingMap[to]++
		incomingMap[from] += 0
	}

	var start string
	for city, count := range incomingMap {
		if count == 0 {
			start = city
			break
		}
	}

	path := []string{start}
	for {
		next, exists := adjacencyMap[start]
		if !exists {
			break
		}
		path = append(path, next)
		start = next
	}
	return path
}
