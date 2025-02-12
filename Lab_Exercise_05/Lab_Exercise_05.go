package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Employee struct to store employee information
type Employee struct {
	ID          int
	Name        string
	Position    string
	Salary      float64
	Performance float64
}

// EmployeeSystem struct to manage the system
type EmployeeSystem struct {
	employees    map[int]Employee
	performance  map[int][]float64
	mutex        sync.RWMutex
	learningChan chan Employee
}

// Error definitions
var (
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrInvalidID        = errors.New("ID must be 100 or greater")
	ErrDuplicateID      = errors.New("employee ID already exists")
)

// NewEmployeeSystem creates a new instance of EmployeeSystem
func NewEmployeeSystem() *EmployeeSystem {
	system := &EmployeeSystem{
		employees:    make(map[int]Employee),
		performance:  make(map[int][]float64),
		learningChan: make(chan Employee, 100),
	}

	// Start the self-learning goroutine
	go system.selfLearning()
	return system
}

// AddEmployee adds a new employee to the system
func (es *EmployeeSystem) AddEmployee(emp Employee) error {
	if emp.ID < 100 {
		return ErrInvalidID
	}

	es.mutex.Lock()
	defer es.mutex.Unlock()

	if _, exists := es.employees[emp.ID]; exists {
		return ErrDuplicateID
	}

	es.employees[emp.ID] = emp
	es.performance[emp.ID] = []float64{}

	// Send to learning channel
	go func() {
		es.learningChan <- emp
	}()
	return nil
}

// UpdateEmployee updates existing employee information
func (es *EmployeeSystem) UpdateEmployee(emp Employee) error {
	if emp.ID <= 0 {
		return ErrInvalidID
	}

	es.mutex.Lock()
	defer es.mutex.Unlock()

	if _, exists := es.employees[emp.ID]; !exists {
		return ErrEmployeeNotFound
	}

	es.employees[emp.ID] = emp
	return nil
}

// GetEmployee retrieves employee information by ID
func (es *EmployeeSystem) GetEmployee(id int) (Employee, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	emp, exists := es.employees[id]
	if !exists {
		return Employee{}, ErrEmployeeNotFound
	}
	return emp, nil
}

// UpdatePerformance adds a new performance rating for an employee
func (es *EmployeeSystem) UpdatePerformance(id int, rating float64) error {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	if _, exists := es.employees[id]; !exists {
		return ErrEmployeeNotFound
	}

	es.performance[id] = append(es.performance[id], rating)

	// Calculate average performance and update employee
	emp := es.employees[id]
	total := 0.0
	for _, r := range es.performance[id] {
		total += r
	}
	emp.Performance = total / float64(len(es.performance[id]))
	es.employees[id] = emp

	// Send to learning channel
	es.learningChan <- emp
	return nil
}

// GetAllEmployees returns a slice of all employees
func (es *EmployeeSystem) GetAllEmployees() []Employee {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	employees := make([]Employee, 0, len(es.employees))
	for _, emp := range es.employees {
		employees = append(employees, emp)
	}
	return employees
}

// selfLearning implements the self-learning component using goroutines
func (es *EmployeeSystem) selfLearning() {
	performanceThresholds := make(map[string]float64)

	for emp := range es.learningChan {
		// Analyze performance patterns
		es.mutex.RLock()
		positionPerformance := make([]float64, 0)
		for _, e := range es.employees {
			if e.Position == emp.Position {
				positionPerformance = append(positionPerformance, e.Performance)
			}
		}
		es.mutex.RUnlock()

		// Calculate average performance for the position
		if len(positionPerformance) > 0 {
			total := 0.0
			for _, p := range positionPerformance {
				total += p
			}
			avgPerformance := total / float64(len(positionPerformance))
			performanceThresholds[emp.Position] = avgPerformance

			// Enhanced output formatting
			fmt.Printf("\n=== Learning Insight ===\n")
			fmt.Printf("Position: %s\n", emp.Position)
			fmt.Printf("Current Employee ID: %d\n", emp.ID)
			fmt.Printf("Position Average Performance: %.2f\n", avgPerformance)
			fmt.Printf("Employees in Position: %d\n", len(positionPerformance))
			fmt.Printf("Timestamp: %s\n", time.Now().Format("15:04:05"))
			fmt.Printf("=====================\n\n")
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// Add these validation functions
func validateEmployee(emp Employee) error {
	if emp.ID < 100 {
		return ErrInvalidID
	}
	if emp.Name == "" {
		return errors.New("name cannot be empty")
	}
	if emp.Position == "" {
		return errors.New("position cannot be empty")
	}
	if emp.Salary <= 0 {
		return errors.New("salary must be positive")
	}
	return nil
}

func validatePerformanceRating(rating float64) error {
	if rating < 0 || rating > 5 {
		return errors.New("performance rating must be between 0 and 5")
	}
	return nil
}

// Add this menu function
func displayMenu() {
	fmt.Println("\n=== Employee Management System ===")
	fmt.Println("1. Add Employee")
	fmt.Println("2. Update Employee")
	fmt.Println("3. View Employee")
	fmt.Println("4. Update Performance")
	fmt.Println("5. View All Employees")
	fmt.Println("6. Exit")
	fmt.Print("Enter your choice: ")
}

// Add this input function
func getEmployeeInput() (Employee, error) {
	var emp Employee

	fmt.Print("Enter Employee ID (must be 100 or greater): ")
	if _, err := fmt.Scan(&emp.ID); err != nil {
		return Employee{}, errors.New("invalid ID format")
	}
	if emp.ID < 100 {
		return Employee{}, ErrInvalidID
	}

	fmt.Print("Enter Name: ")
	var name string
	if _, err := fmt.Scan(&name); err != nil {
		return Employee{}, errors.New("invalid name format")
	}
	emp.Name = name

	fmt.Print("Enter Position: ")
	var position string
	if _, err := fmt.Scan(&position); err != nil {
		return Employee{}, errors.New("invalid position format")
	}
	emp.Position = position

	fmt.Print("Enter Salary: ")
	if _, err := fmt.Scan(&emp.Salary); err != nil {
		return Employee{}, errors.New("invalid salary format")
	}

	return emp, validateEmployee(emp)
}

// Modify the main function
func main() {
	system := NewEmployeeSystem()
	var choice int

	for {
		displayMenu()
		fmt.Scan(&choice)

		switch choice {
		case 1:
			emp, err := getEmployeeInput()
			if err != nil {
				fmt.Printf("Invalid input: %v\n", err)
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
				fmt.Printf("Invalid input: %v\n", err)
				continue
			}

			if err := system.UpdateEmployee(emp); err != nil {
				fmt.Printf("Error updating employee: %v\n", err)
			} else {
				fmt.Println("Employee updated successfully!")
			}

		case 3:
			var id int
			fmt.Print("Enter Employee ID: ")
			fmt.Scan(&id)

			emp, err := system.GetEmployee(id)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("\nEmployee Details:\n")
				fmt.Printf("ID: %d\nName: %s\nPosition: %s\nSalary: %.2f\nPerformance: %.2f\n",
					emp.ID, emp.Name, emp.Position, emp.Salary, emp.Performance)
			}

		case 4:
			var id int
			var rating float64

			fmt.Print("Enter Employee ID: ")
			fmt.Scan(&id)
			fmt.Print("Enter Performance Rating (0-5): ")
			fmt.Scan(&rating)

			if err := validatePerformanceRating(rating); err != nil {
				fmt.Printf("Invalid rating: %v\n", err)
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
				fmt.Printf("ID: %d\nName: %s\nPosition: %s\nSalary: %.2f\nPerformance: %.2f\n",
					emp.ID, emp.Name, emp.Position, emp.Salary, emp.Performance)
				fmt.Println("----------------------------------------")
			}

		case 6:
			fmt.Println("Exiting program...")
			return

		default:
			fmt.Println("Invalid choice! Please try again.")
		}
	}
}
