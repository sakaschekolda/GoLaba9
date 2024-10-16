package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8000"

var token string

// User структура для хранения информации о пользователе
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// AuthRequest структура для хранения данных авторизации
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Функция для отправки запроса на сервер
func sendRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := baseURL + endpoint
	var jsonData []byte
	var err error

	if body != nil {
		jsonData, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

// Функция для авторизации пользователя
func login(username, password string) error {
	authRequest := AuthRequest{Username: username, Password: password}
	resp, err := sendRequest("POST", "/login", authRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("login failed: %s", body)
	}

	var tokenResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}

	token = tokenResponse["token"]
	return nil
}

// Функция для вывода пользователей
func listUsers() {
	resp, err := sendRequest("GET", "/users", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Error:", string(body))
		return
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Age: %d\n", user.ID, user.Name, user.Email, user.Age)
	}
}

// Функция для добавления нового пользователя
func createUser(name, email string, age int) {
	user := User{Name: name, Email: email, Age: age}
	resp, err := sendRequest("POST", "/users", user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Error:", string(body))
		return
	}

	var createdUser User
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("User created: ID: %d, Name: %s, Email: %s, Age: %d\n", createdUser.ID, createdUser.Name, createdUser.Email, createdUser.Age)
}

// Функция для обновления пользователя
func updateUser(id int, name, email string, age int) {
	user := User{ID: id, Name: name, Email: email, Age: age}
	resp, err := sendRequest("PUT", fmt.Sprintf("/users/%d", id), user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Error:", string(body))
		return
	}

	fmt.Printf("User updated: ID: %d, Name: %s, Email: %s, Age: %d\n", user.ID, user.Name, user.Email, user.Age)
}

// Функция для удаления пользователя
func deleteUser(id int) {
	resp, err := sendRequest("DELETE", fmt.Sprintf("/users/%d", id), nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Error:", string(body))
		return
	}

	fmt.Printf("User with ID %d deleted\n", id)
}

// Главное меню
func mainMenu() {
	for {
		fmt.Println("\nМеню:")
		fmt.Println("1. Авторизация")
		fmt.Println("2. Список пользователей")
		fmt.Println("3. Создать пользователя")
		fmt.Println("4. Обновить пользователя")
		fmt.Println("5. Удалить пользователя")
		fmt.Println("6. Выход")
		fmt.Print("Выбор: ")

		var option int
		fmt.Scan(&option)

		switch option {
		case 1:
			var username, password string
			fmt.Print("Введите имя: ")
			fmt.Scan(&username)
			fmt.Print("Введите пароль: ")
			fmt.Scan(&password)

			if err := login(username, password); err != nil {
				fmt.Println("Ошибка авторизации:", err)
			} else {
				fmt.Println("Успешно")
			}
		case 2:
			listUsers()
		case 3:
			var name, email string
			var age int
			fmt.Print("Введите имя: ")
			fmt.Scan(&name)
			fmt.Print("Введите email: ")
			fmt.Scan(&email)
			fmt.Print("Введите возраст: ")
			fmt.Scan(&age)
			createUser(name, email, age)
		case 4:
			var id, age int
			var name, email string
			fmt.Print("Выберите ID пользователя: ")
			fmt.Scan(&id)
			fmt.Print("Введите новое имя: ")
			fmt.Scan(&name)
			fmt.Print("Введите новый email: ")
			fmt.Scan(&email)
			fmt.Print("Введите новый возраст: ")
			fmt.Scan(&age)
			updateUser(id, name, email, age)
		case 5:
			var id int
			fmt.Print("Введите ID, который хотите удалить: ")
			fmt.Scan(&id)
			deleteUser(id)
		case 6:
			fmt.Println("Выход...")
			os.Exit(0)
		default:
			fmt.Println("Ошибка, попробуйте еще раз.")
		}
	}
}

func main() {
	mainMenu()
}
