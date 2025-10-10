package main

import (
	"fmt"
)

// 泛型接口：Greet 方法的输入参数为泛型 T
type Greeter[T any] interface {
	Greet(name T) string
}

// 泛型函数类型：输入为 T，返回 string
type MyFunc[T any] func(T) string

// 为 MyFunc[T] 添加方法，实现 Greeter[T] 接口
func (f MyFunc[T]) Greet(name T) string {
	return f(name)
}

func main() {
	// 示例1: 使用 string 类型（与原示例兼容）
	greetString := MyFunc[string](func(name string) string {
		return "Hello, " + name + "!"
	})
	var gString Greeter[string] = greetString
	fmt.Println(gString.Greet("World")) // 输出: Hello, World!

	// 示例2: 使用 int 类型（展示泛型灵活性）
	greetInt := MyFunc[int](func(num int) string {
		return fmt.Sprintf("Number: %d", num)
	})
	var gInt Greeter[int] = greetInt
	fmt.Println(gInt.Greet(42)) // 输出: Number: 42

	// 示例3: 使用自定义类型
	type Person struct {
		Name string
	}
	greetPerson := MyFunc[Person](func(p Person) string {
		return "Hello, " + p.Name + "!"
	})
	var gPerson Greeter[Person] = greetPerson
	fmt.Println(gPerson.Greet(Person{Name: "Alice"})) // 输出: Hello, Alice!
}
