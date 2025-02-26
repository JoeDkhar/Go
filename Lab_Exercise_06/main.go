package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Department constants using iota
const (
	HR = iota
	Engineering
	Finance
	Marketing
	Operations
)

// DepartmentToString converts department constant to string
func DepartmentToString(dept int) string {
	switch dept {
	case HR:
		return "HR"
	case Engineering:
		return "Engineering"
	case Finance:
		return "Finance"
	case Marketing:
		return "Marketing"
	case Operations:
		return "Operations"
	default:
		return "Unknown"
	}
}

// StringToDepartment converts string to department constant
func StringToDepartment(dept string) (int, error) {
	switch strings.ToLower(dept) {
	case "hr":
		return HR, nil
	case "engineering":
		return Engineering, nil
	case "finance":
		return Finance, nil
	case "marketing":
		return Marketing, nil
	case "operations":
		return Operations, nil
	default:
		return -1, errors.New("invalid department")
	}
}

// Custom error types
var (
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrInvalidID        = errors.New("invalid employee ID")
	ErrDuplicateID      = errors.New("employee ID already exists")
	ErrInvalidInput     = errors.New("invalid input")
)

// Employee struct to store employee information
type Employee struct {
	ID         int
	Name       string
	Position   string
	Salary     float64
	Department int
	JoinDate   time.Time
}

// CalculateExperience calculates years of experience
func (e *Employee) CalculateExperience() float64 {
	duration := time.Since(e.JoinDate)
	return duration.Hours() / 24 / 365
}

// String returns a formatted string representation of the employee
func (e *Employee) String() string {
	return fmt.Sprintf(
		"ID: %d\nName: %s\nPosition: %s\nSalary: $%.2f\nDepartment: %s\nJoin Date: %s\nExperience: %.1f years",
		e.ID, e.Name, e.Position, e.Salary, DepartmentToString(e.Department),
		e.JoinDate.Format("2006-01-02"), e.CalculateExperience(),
	)
}

// EmployeeManager interface defines operations for managing employees
type EmployeeManager interface {
	AddEmployee(e *Employee) error
	RemoveEmployee(id int) error
	UpdateEmployee(e *Employee) error
	GetEmployee(id int) (*Employee, error)
	ListEmployees() ([]*Employee, error)
	FilterEmployees(filter func(*Employee) bool) []*Employee
}

// InMemoryEmployeeManager implements EmployeeManager interface using in-memory storage
type InMemoryEmployeeManager struct {
	employees map[int]*Employee
	nextID    int
}

// NewInMemoryEmployeeManager creates a new InMemoryEmployeeManager
func NewInMemoryEmployeeManager() *InMemoryEmployeeManager {
	return &InMemoryEmployeeManager{
		employees: make(map[int]*Employee),
		nextID:    1,
	}
}

// AddEmployee adds a new employee to the manager
func (m *InMemoryEmployeeManager) AddEmployee(e *Employee) error {
	if e == nil {
		return ErrInvalidInput
	}

	if e.ID == 0 {
		// Auto-assign ID if not provided
		e.ID = m.nextID
		m.nextID++
	} else if _, exists := m.employees[e.ID]; exists {
		return ErrDuplicateID
	}

	// Store a copy of the employee
	employeeCopy := *e
	m.employees[e.ID] = &employeeCopy
	return nil
}

// RemoveEmployee removes an employee by ID
func (m *InMemoryEmployeeManager) RemoveEmployee(id int) error {
	if _, exists := m.employees[id]; !exists {
		return ErrEmployeeNotFound
	}
	delete(m.employees, id)
	return nil
}

// UpdateEmployee updates an existing employee
func (m *InMemoryEmployeeManager) UpdateEmployee(e *Employee) error {
	if e == nil || e.ID == 0 {
		return ErrInvalidInput
	}

	if _, exists := m.employees[e.ID]; !exists {
		return ErrEmployeeNotFound
	}

	// Store a copy of the updated employee
	employeeCopy := *e
	m.employees[e.ID] = &employeeCopy
	return nil
}

// GetEmployee retrieves an employee by ID
func (m *InMemoryEmployeeManager) GetEmployee(id int) (*Employee, error) {
	employee, exists := m.employees[id]
	if !exists {
		return nil, ErrEmployeeNotFound
	}

	// Return a copy to prevent modification of the original
	employeeCopy := *employee
	return &employeeCopy, nil
}

// ListEmployees returns a list of all employees
func (m *InMemoryEmployeeManager) ListEmployees() ([]*Employee, error) {
	employees := make([]*Employee, 0, len(m.employees))
	for _, emp := range m.employees {
		// Create a copy to prevent modification of the original
		employeeCopy := *emp
		employees = append(employees, &employeeCopy)
	}
	return employees, nil
}

// FilterEmployees returns employees that match the filter criteria
func (m *InMemoryEmployeeManager) FilterEmployees(filter func(*Employee) bool) []*Employee {
	result := make([]*Employee, 0)
	for _, emp := range m.employees {
		if filter(emp) {
			// Create a copy to prevent modification of the original
			employeeCopy := *emp
			result = append(result, &employeeCopy)
		}
	}
	return result
}

// AddMultipleEmployees demonstrates a variadic function to add multiple employees
func AddMultipleEmployees(manager EmployeeManager, employees ...*Employee) []error {
	errors := make([]error, 0)
	for _, emp := range employees {
		if err := manager.AddEmployee(emp); err != nil {
			errors = append(errors, fmt.Errorf("error adding employee ID %d: %w", emp.ID, err))
		}
	}
	return errors
}

// Helper functions for user interaction

// readString reads a string from the user
func readString(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// readInt reads an integer from the user
func readInt(reader *bufio.Reader, prompt string) (int, error) {
	input, err := readString(reader, prompt)
	if err != nil {
		return 0, err
	}

	if input == "" {
		return 0, nil // Allow empty input for optional fields
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("%w: please enter a valid number", ErrInvalidInput)
	}
	return value, nil
}

// readFloat reads a float from the user
func readFloat(reader *bufio.Reader, prompt string) (float64, error) {
	input, err := readString(reader, prompt)
	if err != nil {
		return 0, err
	}

	if input == "" {
		return 0, nil // Allow empty input for optional fields
	}

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: please enter a valid number", ErrInvalidInput)
	}
	return value, nil
}

// readDate reads a date from the user
func readDate(reader *bufio.Reader, prompt string) (time.Time, error) {
	input, err := readString(reader, prompt+" (YYYY-MM-DD): ")
	if err != nil {
		return time.Time{}, err
	}

	if input == "" {
		return time.Now(), nil // Default to current date if empty
	}

	date, err := time.Parse("2006-01-02", input)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: please enter a valid date in YYYY-MM-DD format", ErrInvalidInput)
	}
	return date, nil
}

// readDepartment reads a department from the user
func readDepartment(reader *bufio.Reader) (int, error) {
	fmt.Println("\nAvailable departments:")
	fmt.Println("1. HR")
	fmt.Println("2. Engineering")
	fmt.Println("3. Finance")
	fmt.Println("4. Marketing")
	fmt.Println("5. Operations")

	choice, err := readInt(reader, "Select department (1-5): ")
	if err != nil {
		return -1, err
	}

	switch choice {
	case 1:
		return HR, nil
	case 2:
		return Engineering, nil
	case 3:
		return Finance, nil
	case 4:
		return Marketing, nil
	case 5:
		return Operations, nil
	default:
		return -1, fmt.Errorf("%w: please select a valid department", ErrInvalidInput)
	}
}

// Interactive console functions

// addEmployeeInteractive adds an employee through user interaction
func addEmployeeInteractive(manager EmployeeManager, reader *bufio.Reader) error {
	fmt.Println("\n=== Add New Employee ===")

	name, err := readString(reader, "Name: ")
	if err != nil {
		return err
	}

	position, err := readString(reader, "Position: ")
	if err != nil {
		return err
	}

	salary, err := readFloat(reader, "Salary: ")
	if err != nil {
		return err
	}

	department, err := readDepartment(reader)
	if err != nil {
		return err
	}

	joinDate, err := readDate(reader, "Join Date")
	if err != nil {
		return err
	}

	employee := &Employee{
		Name:       name,
		Position:   position,
		Salary:     salary,
		Department: department,
		JoinDate:   joinDate,
	}

	err = manager.AddEmployee(employee)
	if err != nil {
		return err
	}

	fmt.Printf("\nEmployee added successfully with ID: %d\n", employee.ID)
	return nil
}

// updateEmployeeInteractive updates an employee through user interaction
func updateEmployeeInteractive(manager EmployeeManager, reader *bufio.Reader) error {
	fmt.Println("\n=== Update Employee ===")

	id, err := readInt(reader, "Enter employee ID to update: ")
	if err != nil {
		return err
	}

	employee, err := manager.GetEmployee(id)
	if err != nil {
		return err
	}

	fmt.Println("\nCurrent employee information:")
	fmt.Println(employee)
	fmt.Println("\nEnter new information (leave blank to keep current value):")

	name, err := readString(reader, fmt.Sprintf("Name [%s]: ", employee.Name))
	if err != nil {
		return err
	}
	if name != "" {
		employee.Name = name
	}

	position, err := readString(reader, fmt.Sprintf("Position [%s]: ", employee.Position))
	if err != nil {
		return err
	}
	if position != "" {
		employee.Position = position
	}

	salaryStr, err := readString(reader, fmt.Sprintf("Salary [%.2f]: ", employee.Salary))
	if err != nil {
		return err
	}
	if salaryStr != "" {
		salary, err := strconv.ParseFloat(salaryStr, 64)
		if err != nil {
			return fmt.Errorf("%w: please enter a valid number", ErrInvalidInput)
		}
		employee.Salary = salary
	}

	fmt.Println("\nUpdate department? (y/n)")
	updateDept, err := readString(reader, "Choice: ")
	if err != nil {
		return err
	}

	if strings.ToLower(updateDept) == "y" {
		department, err := readDepartment(reader)
		if err != nil {
			return err
		}
		employee.Department = department
	}

	fmt.Println("\nUpdate join date? (y/n)")
	updateDate, err := readString(reader, "Choice: ")
	if err != nil {
		return err
	}

	if strings.ToLower(updateDate) == "y" {
		joinDate, err := readDate(reader, "Join Date")
		if err != nil {
			return err
		}
		employee.JoinDate = joinDate
	}

	err = manager.UpdateEmployee(employee)
	if err != nil {
		return err
	}

	fmt.Println("\nEmployee updated successfully!")
	return nil
}

// removeEmployeeInteractive removes an employee through user interaction
func removeEmployeeInteractive(manager EmployeeManager, reader *bufio.Reader) error {
	fmt.Println("\n=== Remove Employee ===")

	id, err := readInt(reader, "Enter employee ID to remove: ")
	if err != nil {
		return err
	}

	employee, err := manager.GetEmployee(id)
	if err != nil {
		return err
	}

	fmt.Println("\nEmployee to remove:")
	fmt.Println(employee)

	confirm, err := readString(reader, "\nAre you sure you want to remove this employee? (y/n): ")
	if err != nil {
		return err
	}

	if strings.ToLower(confirm) != "y" {
		fmt.Println("\nOperation cancelled.")
		return nil
	}

	err = manager.RemoveEmployee(id)
	if err != nil {
		return err
	}

	fmt.Println("\nEmployee removed successfully!")
	return nil
}

// searchEmployeesInteractive searches for employees through user interaction
func searchEmployeesInteractive(manager EmployeeManager, reader *bufio.Reader) error {
	fmt.Println("\n=== Search Employees ===")
	fmt.Println("1. Search by name")
	fmt.Println("2. Search by department")
	fmt.Println("3. Search by salary range")
	fmt.Println("4. Search by experience")

	option, err := readInt(reader, "\nSelect search option: ")
	if err != nil {
		return err
	}

	var employees []*Employee

	switch option {
	case 1:
		name, err := readString(reader, "Enter name (partial match allowed): ")
		if err != nil {
			return err
		}

		employees = manager.FilterEmployees(func(e *Employee) bool {
			return strings.Contains(strings.ToLower(e.Name), strings.ToLower(name))
		})

	case 2:
		department, err := readDepartment(reader)
		if err != nil {
			return err
		}

		employees = manager.FilterEmployees(func(e *Employee) bool {
			return e.Department == department
		})

	case 3:
		minSalary, err := readFloat(reader, "Enter minimum salary: ")
		if err != nil {
			return err
		}

		maxSalary, err := readFloat(reader, "Enter maximum salary: ")
		if err != nil {
			return err
		}

		employees = manager.FilterEmployees(func(e *Employee) bool {
			return e.Salary >= minSalary && e.Salary <= maxSalary
		})

	case 4:
		minExp, err := readFloat(reader, "Enter minimum years of experience: ")
		if err != nil {
			return err
		}

		employees = manager.FilterEmployees(func(e *Employee) bool {
			return e.CalculateExperience() >= minExp
		})

	default:
		return fmt.Errorf("%w: please select a valid option", ErrInvalidInput)
	}

	if len(employees) == 0 {
		fmt.Println("\nNo employees found matching the criteria.")
		return nil
	}

	fmt.Printf("\nFound %d employee(s):\n\n", len(employees))
	for i, emp := range employees {
		fmt.Printf("=== Employee %d ===\n", i+1)
		fmt.Println(emp)
		fmt.Println()
	}

	return nil
}

// displayAllEmployees displays all employees
func displayAllEmployees(manager EmployeeManager) error {
	employees, err := manager.ListEmployees()
	if err != nil {
		return err
	}

	if len(employees) == 0 {
		fmt.Println("\nNo employees found.")
		return nil
	}

	fmt.Printf("\n=== All Employees (%d) ===\n\n", len(employees))
	for i, emp := range employees {
		fmt.Printf("=== Employee %d ===\n", i+1)
		fmt.Println(emp)
		fmt.Println()
	}

	return nil
}

// addSampleData adds sample data to the manager
func addSampleData(manager EmployeeManager) {
	// Create sample employees
	employees := []*Employee{
		{
			Name:       "John Doe",
			Position:   "Software Engineer",
			Salary:     85000,
			Department: Engineering,
			JoinDate:   time.Date(2020, 5, 15, 0, 0, 0, 0, time.Local),
		},
		{
			Name:       "Jane Smith",
			Position:   "HR Manager",
			Salary:     75000,
			Department: HR,
			JoinDate:   time.Date(2019, 3, 10, 0, 0, 0, 0, time.Local),
		},
		{
			Name:       "Michael Johnson",
			Position:   "Finance Director",
			Salary:     110000,
			Department: Finance,
			JoinDate:   time.Date(2018, 1, 5, 0, 0, 0, 0, time.Local),
		},
		{
			Name:       "Emily Williams",
			Position:   "Marketing Specialist",
			Salary:     65000,
			Department: Marketing,
			JoinDate:   time.Date(2021, 8, 22, 0, 0, 0, 0, time.Local),
		},
		{
			Name:       "Robert Brown",
			Position:   "Operations Manager",
			Salary:     90000,
			Department: Operations,
			JoinDate:   time.Date(2019, 11, 7, 0, 0, 0, 0, time.Local),
		},
	}

	// Add employees using variadic function
	errors := AddMultipleEmployees(manager, employees...)
	if len(errors) > 0 {
		fmt.Println("Errors adding sample data:")
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}

// displayMenu displays the main menu
func displayMenu() {
	fmt.Println("\n======= Employee Management System =======")
	fmt.Println("1. Add Employee")
	fmt.Println("2. View All Employees")
	fmt.Println("3. Update Employee")
	fmt.Println("4. Remove Employee")
	fmt.Println("5. Search Employees")
	fmt.Println("6. Add Sample Data")
	fmt.Println("0. Exit")
	fmt.Println("=========================================")
}

// main function - entry point of the application
func main() {
	// Create employee manager
	manager := NewInMemoryEmployeeManager()

	// Create reader for user input
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Employee Management System!")

	for {
		displayMenu()

		choice, err := readInt(reader, "Enter your choice: ")
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		switch choice {
		case 1:
			err = addEmployeeInteractive(manager, reader)
		case 2:
			err = displayAllEmployees(manager)
		case 3:
			err = updateEmployeeInteractive(manager, reader)
		case 4:
			err = removeEmployeeInteractive(manager, reader)
		case 5:
			err = searchEmployeesInteractive(manager, reader)
		case 6:
			addSampleData(manager)
			fmt.Println("\nSample data added successfully!")
			err = nil
		case 0:
			fmt.Println("\nThank you for using the Employee Management System. Goodbye!")
			return
		default:
			err = fmt.Errorf("%w: please select a valid option", ErrInvalidInput)
		}

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
