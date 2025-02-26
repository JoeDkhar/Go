package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	MinSalary = 20000
	MaxSalary = 2000000
)

type PositionStats struct {
	AvgPerformance float64
	EmployeeCount  int
	TotalSalary    float64
	LastUpdated    time.Time
}

type Employee struct {
	ID          int
	Name        string
	Position    string
	Salary      float64
	Performance float64
	LastUpdated time.Time
}

type EmployeeSystem struct {
	employees     map[int]Employee
	performance   map[int][]float64
	positionStats map[string]PositionStats
	mutex         sync.RWMutex
	learningChan  chan Employee
	done          chan struct{} // Add this channel for cleanup
	ctx           context.Context
	cancel        context.CancelFunc
}

var (
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrInvalidID        = errors.New("ID must be 100 or greater")
	ErrDuplicateID      = errors.New("employee ID already exists")
	ErrInvalidName      = errors.New("name must be 2-50 characters and contain only letters")
	ErrInvalidPosition  = errors.New("position must be 2-50 characters")
	ErrInvalidSalary    = errors.New("salary must be between 30000 and 500000")
	ErrInvalidRating    = errors.New("performance rating must be between 0 and 5")
)

// Input handling functions
func readString(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readInt(prompt string) (int, error) {
	input := readString(prompt)
	return strconv.Atoi(input)
}

func readFloat(prompt string) (float64, error) {
	input := readString(prompt)
	return strconv.ParseFloat(input, 64)
}

// Validation functions
func validateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 50 {
		return ErrInvalidName
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return ErrInvalidName
		}
	}
	return nil
}

func validateSalary(salary float64) error {
	if salary < MinSalary || salary > MaxSalary {
		return fmt.Errorf("salary must be between %.2f and %.2f", MinSalary, MaxSalary)
	}
	return nil
}

func validateRating(rating float64) error {
	if rating < 0 || rating > 5 {
		return ErrInvalidRating
	}
	return nil
}

func NewEmployeeSystem() *EmployeeSystem {
	ctx, cancel := context.WithCancel(context.Background())
	system := &EmployeeSystem{
		employees:     make(map[int]Employee),
		performance:   make(map[int][]float64),
		positionStats: make(map[string]PositionStats),
		learningChan:  make(chan Employee, 100),
		done:          make(chan struct{}), // Initialize done channel
		ctx:           ctx,
		cancel:        cancel,
	}
	go system.selfLearning()
	return system
}

func (es *EmployeeSystem) AddEmployee(emp Employee) error {
	if emp.ID < 100 {
		return ErrInvalidID
	}
	if err := validateName(emp.Name); err != nil {
		return err
	}
	if err := validateSalary(emp.Salary); err != nil {
		return err
	}

	es.mutex.Lock()
	defer es.mutex.Unlock()

	if _, exists := es.employees[emp.ID]; exists {
		return ErrDuplicateID
	}

	emp.LastUpdated = time.Now()
	es.employees[emp.ID] = emp
	es.performance[emp.ID] = []float64{}

	select {
	case es.learningChan <- emp:
	case <-time.After(100 * time.Millisecond):
		fmt.Printf("Warning: Learning system busy, skipped analysis for %s\n", emp.Name)
	}
	return nil
}

func (es *EmployeeSystem) UpdateEmployee(emp Employee) error {
	if emp.ID < 100 {
		return ErrInvalidID
	}
	if err := validateName(emp.Name); err != nil {
		return err
	}
	if err := validateSalary(emp.Salary); err != nil {
		return err
	}

	es.mutex.Lock()
	defer es.mutex.Unlock()

	if _, exists := es.employees[emp.ID]; !exists {
		return ErrEmployeeNotFound
	}

	emp.LastUpdated = time.Now()
	es.employees[emp.ID] = emp
	return nil
}

func (es *EmployeeSystem) GetEmployee(id int) (Employee, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	emp, exists := es.employees[id]
	if !exists {
		return Employee{}, ErrEmployeeNotFound
	}
	return emp, nil
}

func (es *EmployeeSystem) UpdatePerformance(id int, rating float64) error {
	if err := validateRating(rating); err != nil {
		return err
	}

	es.mutex.Lock()
	defer es.mutex.Unlock()

	emp, exists := es.employees[id]
	if !exists {
		return ErrEmployeeNotFound
	}

	es.performance[id] = append(es.performance[id], rating)

	total := 0.0
	for _, r := range es.performance[id] {
		total += r
	}
	emp.Performance = total / float64(len(es.performance[id]))
	emp.LastUpdated = time.Now()
	es.employees[id] = emp

	select {
	case es.learningChan <- emp:
	default:
		// Non-blocking send to learning channel
	}
	return nil
}

func (es *EmployeeSystem) GetAllEmployees() []Employee {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	employees := make([]Employee, 0, len(es.employees))
	for _, emp := range es.employees {
		employees = append(employees, emp)
	}
	return employees
}

func (es *EmployeeSystem) Shutdown() {
	close(es.done) // Signal the goroutine to stop
}

func (es *EmployeeSystem) selfLearning() {
	for {
		select {
		case emp := <-es.learningChan:
			es.mutex.Lock()
			stats := PositionStats{
				LastUpdated: time.Now(),
			}

			var totalPerf float64
			var count int
			var totalSalary float64

			for _, e := range es.employees {
				if e.Position == emp.Position {
					totalPerf += e.Performance
					totalSalary += e.Salary
					count++
				}
			}

			if count > 0 {
				stats.AvgPerformance = totalPerf / float64(count)
				stats.EmployeeCount = count
				stats.TotalSalary = totalSalary
				es.positionStats[emp.Position] = stats
			}
			es.mutex.Unlock()

			fmt.Printf("\n🤖 Learning System Update:\n")
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("Position: %s\n", emp.Position)
			fmt.Printf("Employees in Position: %d\n", count)
			fmt.Printf("Average Performance: %.2f\n", stats.AvgPerformance)
			if count > 0 {
				fmt.Printf("Average Salary: %.2f\n", totalSalary/float64(count))
			}
			fmt.Printf("Last Updated: %s\n", stats.LastUpdated.Format("15:04:05"))
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		case <-es.ctx.Done():
			return // Exit goroutine cleanly
		}
	}
}

func getEmployeeInput() (Employee, error) {
	id, err := readInt("Enter Employee ID (must be 100 or greater): ")
	if err != nil {
		return Employee{}, fmt.Errorf("invalid ID format: %v", err)
	}

	name := readString("Enter Name: ")
	if err := validateName(name); err != nil {
		return Employee{}, err
	}

	position := readString("Enter Position: ")
	if len(position) < 2 {
		return Employee{}, ErrInvalidPosition
	}

	salary, err := readFloat("Enter Salary: ")
	if err != nil {
		return Employee{}, fmt.Errorf("invalid salary format: %v", err)
	}
	if err := validateSalary(salary); err != nil {
		return Employee{}, err
	}

	return Employee{
		ID:       id,
		Name:     name,
		Position: position,
		Salary:   salary,
	}, nil
}

func main() {
	system := NewEmployeeSystem()
	defer system.Shutdown() // Ensure cleanup happens

	fmt.Printf("\nWelcome to Employee Management System\n")
	fmt.Printf("Valid salary range: %.2f - %.2f\n", MinSalary, MaxSalary)

	for {
		fmt.Println("\n=== Employee Management System ===")
		fmt.Println("1. Add Employee")
		fmt.Println("2. Update Employee")
		fmt.Println("3. View Employee")
		fmt.Println("4. Update Performance")
		fmt.Println("5. View All Employees")
		fmt.Println("6. Exit")

		choice, err := readInt("Enter your choice (1-6): ")
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			emp, err := getEmployeeInput()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			if err := system.AddEmployee(emp); err != nil {
				fmt.Printf("Error adding employee: %v\n", err)
			} else {
				fmt.Println("Employee added successfully!")
			}

		case 2:
			emp, err := getEmployeeInput()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			if err := system.UpdateEmployee(emp); err != nil {
				fmt.Printf("Error updating employee: %v\n", err)
			} else {
				fmt.Println("Employee updated successfully!")
			}

		case 3:
			id, err := readInt("Enter Employee ID: ")
			if err != nil {
				fmt.Println("Invalid ID format")
				continue
			}
			emp, err := system.GetEmployee(id)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("\nEmployee Details:\n")
				fmt.Printf("ID: %d\n", emp.ID)
				fmt.Printf("Name: %s\n", emp.Name)
				fmt.Printf("Position: %s\n", emp.Position)
				fmt.Printf("Salary: %.2f\n", emp.Salary)
				fmt.Printf("Performance: %.2f\n", emp.Performance)
				fmt.Printf("Last Updated: %s\n", emp.LastUpdated.Format("2006-01-02 15:04:05"))
			}

		case 4:
			id, err := readInt("Enter Employee ID: ")
			if err != nil {
				fmt.Println("Invalid ID format")
				continue
			}
			rating, err := readFloat("Enter Performance Rating (0-5): ")
			if err != nil {
				fmt.Println("Invalid rating format")
				continue
			}
			if err := system.UpdatePerformance(id, rating); err != nil {
				fmt.Printf("Error updating performance: %v\n", err)
			} else {
				fmt.Println("Performance updated successfully!")
			}

		case 5:
			employees := system.GetAllEmployees()
			if len(employees) == 0 {
				fmt.Println("No employees found!")
				continue
			}
			fmt.Println("\nAll Employees:")
			fmt.Println("----------------------------------------")
			for _, emp := range employees {
				fmt.Printf("ID: %d\n", emp.ID)
				fmt.Printf("Name: %s\n", emp.Name)
				fmt.Printf("Position: %s\n", emp.Position)
				fmt.Printf("Salary: %.2f\n", emp.Salary)
				fmt.Printf("Performance: %.2f\n", emp.Performance)
				fmt.Printf("Last Updated: %s\n", emp.LastUpdated.Format("2006-01-02 15:04:05"))
				fmt.Println("----------------------------------------")
			}

		case 6:
			fmt.Println("Thank you for using the Employee Management System!")
			system.Shutdown()
			return

		default:
			fmt.Println("Invalid choice! Please enter a number between 1 and 6.")
		}
	}
}
