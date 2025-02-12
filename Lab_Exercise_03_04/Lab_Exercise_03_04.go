package main

import (
	"fmt"
	"strings"
)

type MenuOption int

const (
	AddEmployee MenuOption = iota + 1
	DisplayEmployees
	UpdateEmployeeSalary
	ExitProgram
)

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
	Position   string
}

var (
	departments      = [4]string{"IT", "HR", "Finance", "Marketing"}
	employeesList    []*Employee               // Store pointers to reflect updates
	employees        = make(map[int]*Employee) // Map of Employee pointers
	deptEmployees    = make(map[string][]*Employee)
	salaryThresholds = map[string]float64{
		"Junior":   30000,
		"Senior":   50000,
		"Lead":     80000,
		"Manager":  100000,
		"Director": 150000,
	}
)

// Validate input data
func validate(field string, value interface{}, isUpdate bool) error {
	switch field {
	case "id":
		id, ok := value.(int)
		if !ok || id <= 0 {
			return fmt.Errorf("invalid ID: must be positive")
		}
		if !isUpdate {
			if _, exists := employees[id]; exists {
				return fmt.Errorf("employee ID %d already exists", id)
			}
		} else {
			if _, exists := employees[id]; !exists {
				return fmt.Errorf("employee ID %d not found", id)
			}
		}

	case "name":
		name, ok := value.(string)
		if !ok || strings.TrimSpace(name) == "" {
			return fmt.Errorf("name cannot be empty")
		}

	case "salary":
		salary, ok := value.(float64)
		if !ok || salary < 0 {
			return fmt.Errorf("invalid salary: must be non-negative")
		}

	case "department":
		dept, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid department")
		}
		valid := false
		for _, d := range departments {
			if d == dept {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid department: must be one of %v", departments)
		}
	}
	return nil
}

// Check if employee is eligible for promotion (based only on salary)
func (e *Employee) checkPromotion() bool {
	currentPosition := e.Position
	switch {
	case currentPosition == "Junior" && e.Salary >= salaryThresholds["Senior"]:
		e.Position = "Senior"
	case currentPosition == "Senior" && e.Salary >= salaryThresholds["Lead"]:
		e.Position = "Lead"
	case currentPosition == "Lead" && e.Salary >= salaryThresholds["Manager"]:
		e.Position = "Manager"
	case currentPosition == "Manager" && e.Salary >= salaryThresholds["Director"]:
		e.Position = "Director"
	default:
		return false
	}
	return true
}

// Display all employees
func displayAllEmployees() {
	if len(employeesList) == 0 {
		fmt.Println("No employees found!")
		return
	}

	deptCounts := make(map[string]int)
	for _, emp := range employeesList {
		deptCounts[emp.Department]++
	}

	fmt.Println("\n+-----+------------------+---------------+------------+-------------+")
	fmt.Printf("| %-3s | %-16s | %-13s | %-10s | %-11s |\n",
		"ID", "Name", "Department", "Position", "Salary")
	fmt.Println("+-----+------------------+---------------+------------+-------------+")

	for _, emp := range employeesList {
		fmt.Printf("| %-3d | %-16s | %-13s | %-10s | %-11.2f |\n",
			emp.ID,
			emp.Name,
			emp.Department,
			emp.Position,
			emp.Salary)
	}

	fmt.Println("+-----+------------------+---------------+------------+-------------+")
	fmt.Printf("Total Employees: %d\n", len(employeesList))

	fmt.Println("\nDepartment Breakdown:")
	for dept, count := range deptCounts {
		fmt.Printf("%s: %d employees\n", dept, count)
	}
}

func checkPosition(salary float64) string {
	oldPosition := ""
	newPosition := ""

	// Determine position based on salary
	switch {
	case salary >= salaryThresholds["Director"]:
		newPosition = "Director"
	case salary >= salaryThresholds["Manager"]:
		newPosition = "Manager"
	case salary >= salaryThresholds["Lead"]:
		newPosition = "Lead"
	case salary >= salaryThresholds["Senior"]:
		newPosition = "Senior"
	default:
		newPosition = "Junior"
	}

	// Return new position
	if oldPosition != "" && oldPosition != newPosition {
		fmt.Printf("Position changed from %s to %s\n", oldPosition, newPosition)
	}
	return newPosition
}

func addEmployee(id int, name string, department string, salary float64) error {
	position := checkPosition(salary)
	emp := &Employee{
		ID:         id,
		Name:       name,
		Department: department,
		Salary:     salary,
		Position:   position,
	}
	employeesList = append(employeesList, emp)
	employees[id] = emp
	deptEmployees[department] = append(deptEmployees[department], emp)
	return nil
}

// Update employee salary and check for promotion
func updateEmployee(id int, salary float64) bool {
	if emp, exists := employees[id]; exists {
		emp.Salary = salary
		if emp.checkPromotion() {
			fmt.Printf("🎉 Congratulations! %s has been promoted to %s\n", emp.Name, emp.Position)
		}
		return true
	}
	return false
}

func updateEmployeeSalary(id int, newSalary float64) error {
	emp, exists := employees[id]
	if !exists {
		return fmt.Errorf("employee ID %d not found", id)
	}

	oldPosition := emp.Position
	newPosition := checkPosition(newSalary)

	// Update employee details
	emp.Salary = newSalary
	emp.Position = newPosition

	// Update maps
	employees[id] = emp

	// Update department map
	for i, e := range deptEmployees[emp.Department] {
		if e.ID == id {
			deptEmployees[emp.Department][i].Salary = newSalary
			deptEmployees[emp.Department][i].Position = newPosition
			break
		}
	}

	if oldPosition != newPosition {
		fmt.Printf("Employee %d position updated: %s -> %s\n", id, oldPosition, newPosition)
	}

	return nil
}

func displayEmployees() {
	if len(employees) == 0 {
		fmt.Println("No employees to display")
		return
	}

	fmt.Println("\n=== Employee List ===")
	fmt.Printf("%-5s %-20s %-15s %-12s %-10s\n", "ID", "Name", "Department", "Salary", "Position")
	fmt.Println(strings.Repeat("-", 65))

	for _, emp := range employees {
		// Ensure position is up to date with current salary
		currentPosition := checkPosition(emp.Salary)
		if emp.Position != currentPosition {
			emp.Position = currentPosition
			// Update position in deptEmployees
			for i, deptEmp := range deptEmployees[emp.Department] {
				if deptEmp.ID == emp.ID {
					deptEmployees[emp.Department][i].Position = currentPosition
				}
			}
		}
		fmt.Printf("%-5d %-20s %-15s %-12.2f %-10s\n",
			emp.ID, emp.Name, emp.Department, emp.Salary, emp.Position)
	}
	fmt.Println(strings.Repeat("-", 65))
}

// Display menu
func displayMenu() {
	fmt.Println("\n==================================")
	fmt.Println("|      EMPLOYEE MANAGEMENT       |")
	fmt.Println("==================================")
	fmt.Println("| 1. Add Employee                |")
	fmt.Println("| 2. Display All Employees       |")
	fmt.Println("| 3. Update Employee Salary      |")
	fmt.Println("| 4. Exit                        |")
	fmt.Println("==================================")
	fmt.Printf("\nAvailable Departments: %v\n", departments)
	fmt.Print("\nEnter your choice (1-4): ")
}

// Main function
func main() {
	for {
		displayMenu()
		var choice MenuOption
		fmt.Scan(&choice)

		switch choice {
		case AddEmployee:
			fmt.Println("\n==================================")
			fmt.Println("|         ADD EMPLOYEE           |")
			fmt.Println("==================================")

			var id int
			var name, dept string
			var salary float64

			for {
				fmt.Print("Enter Employee ID: ")
				fmt.Scan(&id)
				if err := validate("id", id, false); err != nil {
					fmt.Printf("\n❌ Error: %v\n", err)
					continue
				}
				break
			}

			for {
				fmt.Print("Enter Employee Name: ")
				fmt.Scan(&name)
				if err := validate("name", name, false); err != nil {
					fmt.Printf("\n❌ Error: %v\n", err)
					continue
				}
				break
			}

			for {
				fmt.Print("Enter Department: ")
				fmt.Scan(&dept)
				if err := validate("department", dept, false); err != nil {
					fmt.Printf("\n❌ Error: %v\n", err)
					continue
				}
				break
			}

			for {
				fmt.Print("Enter Salary: ")
				fmt.Scan(&salary)
				if err := validate("salary", salary, false); err != nil {
					fmt.Printf("\n❌ Error: %v\n", err)
					continue
				}
				break
			}

			addEmployee(id, name, dept, salary)
			fmt.Println("\n✅ Employee added successfully!")

		case DisplayEmployees:
			fmt.Println("\n==================================")
			fmt.Println("|      EMPLOYEE DETAILS          |")
			fmt.Println("==================================")
			displayAllEmployees()

		case UpdateEmployeeSalary:
			var id int
			var salary float64
			fmt.Print("Enter Employee ID: ")
			fmt.Scan(&id)
			fmt.Print("Enter New Salary: ")
			fmt.Scan(&salary)

			if updateEmployee(id, salary) {
				fmt.Println("\n✅ Salary updated successfully!")
			} else {
				fmt.Println("\n❌ Employee not found!")
			}

		case ExitProgram:
			fmt.Println("\n👋 Goodbye!")
			return

		default:
			fmt.Println("\n❌ Invalid choice! Please enter 1-4")
		}
	}
}
