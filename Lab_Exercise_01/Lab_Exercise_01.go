//to demonstrate the declaration and intialization of variables of various types including integers, float, strings and boolean you are to demonstrate the use of control flow such as if, for and switch statestements

package main

import "fmt"

func main() {
	fmt.Print("Testing")

	var (
		employeeName   string
		employeeID     int32
		salary         float64
		isActive       bool
		department     string
		yearsOfService int
		choice         int
		isEmployeeSet  bool = false
	)

	for {
		fmt.Println("\n=== Employee Management System ===")
		fmt.Println("1. Add Employee")
		fmt.Println("2. View Employee")
		fmt.Println("3. Update Employee")
		fmt.Println("4. Delete Employee")
		fmt.Println("5. Exit")
		fmt.Print("\nEnter your choice (1-5): ")

		fmt.Scan(&choice)

		switch choice {
		case 1:
			if !isEmployeeSet {
				fmt.Println("---Add Employee---")
				fmt.Print("Enter the Employee's name::")
				fmt.Scan(&employeeName)
				fmt.Print("Enter the Employee's ID::")
				fmt.Scan(&employeeID)
				fmt.Print("Enter the Employee's salary::")
				fmt.Scan(&salary)
				fmt.Print("Enter the Employee's years of service::")
				fmt.Scan(&yearsOfService)
				fmt.Print("Enter the Employee's department::")
				fmt.Scan(&department)
				isActive = true
				isEmployeeSet = true
				fmt.Println("Employee details added successfully ")
			} else {
				fmt.Println("Employee already exists| Use update option")
			}
		case 2:
			if isEmployeeSet {
				fmt.Printf("---View Employee---\nEmployee's Name: %s\nEmployee's ID: %d\nEmployee's Salary: %.2f\nEmployee's Department: %s\nEmployee's Years of Service: %d\nEmployee's State: %t\n",
					employeeName, employeeID, salary, department, yearsOfService, isActive)

			} else {
				fmt.Println("No employee data exists!")
			}

		case 3:
			if isEmployeeSet {
				fmt.Println("---Employee Updation---")
				fmt.Print("Enter new Name: ")
				fmt.Scan(&employeeName)
				fmt.Print("Enter new Salary: ")
				fmt.Scan(&salary)
				fmt.Print("Enter new Department: ")
				fmt.Scan(&department)
				fmt.Println("Employee Updated Successfully!")
			} else {
				fmt.Println("No employee data exists!")
			}
		case 4:
			if isEmployeeSet {
				fmt.Println("---Deleting Employee Data---")
				employeeName = ""
				employeeID = 0
				salary = 0.0
				isActive = false
				department = ""
				yearsOfService = 0
				isEmployeeSet = false
			} else {
				fmt.Println("No employee data exists!")
			}
		case 5:
			fmt.Println("\nExiting Program...")
			return
		default:
			fmt.Println("\nInvalid Choice! Please try again.")
		}

	}
}
