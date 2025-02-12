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
	ErrInvalidID        = errors.New("invalid employee ID")
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
	if emp.ID <= 0 {
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
	es.learningChan <- emp
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

			// Update performance thresholds
			performanceThresholds[emp.Position] = avgPerformance

			// Log insights
			fmt.Printf("Learning Insight - Position: %s, Average Performance: %.2f\n",
				emp.Position, avgPerformance)
		}

		// Simulate processing time
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// Initialize the system
	system := NewEmployeeSystem()

	// Example usage with error handling
	err := system.AddEmployee(Employee{
		ID:       1,
		Name:     "John Doe",
		Position: "Developer",
		Salary:   75000,
	})
	if err != nil {
		fmt.Printf("Error adding employee: %v\n", err)
		return
	}

	// Add more employees
	employees := []Employee{
		{ID: 2, Name: "Jane Smith", Position: "Developer", Salary: 80000},
		{ID: 3, Name: "Bob Johnson", Position: "Manager", Salary: 95000},
	}

	for _, emp := range employees {
		if err := system.AddEmployee(emp); err != nil {
			fmt.Printf("Error adding employee %s: %v\n", emp.Name, err)
		}
	}

	// Update performance ratings
	if err := system.UpdatePerformance(1, 4.5); err != nil {
		fmt.Printf("Error updating performance: %v\n", err)
	}

	if err := system.UpdatePerformance(2, 4.8); err != nil {
		fmt.Printf("Error updating performance: %v\n", err)
	}

	// Get and display all employees
	allEmployees := system.GetAllEmployees()
	fmt.Println("\nAll Employees:")
	for _, emp := range allEmployees {
		fmt.Printf("ID: %d, Name: %s, Position: %s, Performance: %.2f\n",
			emp.ID, emp.Name, emp.Position, emp.Performance)
	}

	// Keep the program running to allow self-learning to process
	time.Sleep(2 * time.Second)
}
