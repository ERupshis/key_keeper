package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Простое CLI-приложение")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Введите команду (exit для выхода): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			continue
		}

		// Удаляем символ новой строки из ввода
		command := strings.TrimSpace(input)

		if command == "exit" {
			fmt.Println("Выход из приложения.")
			break
		}

		processCommand(command)
	}

}

func processCommand(command string) {
	fmt.Printf("Вы ввели команду: %s\n", command)
	// Здесь можно добавить логику для обработки конкретных команд
}
