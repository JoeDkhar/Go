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
		fmt.Println("\n==================================")
		fmt.Println("|      EMPLOYEE MANAGEMENT        |")
		fmt.Println("==================================")
		fmt.Println("| 1. Add Employee                |")
		fmt.Println("| 2. View Employee               |")
		fmt.Println("| 3. Update Employee             |")
		fmt.Println("| 4. Delete Employee             |")
		fmt.Println("| 5. Exit                        |")
		fmt.Println("==================================")
		fmt.Print("\nEnter your choice (1-5): ")

		fmt.Scan(&choice)

		switch choice {
		case 1:
			if !isEmployeeSet {
				fmt.Println("\n==================================")
				fmt.Println("|         ADD EMPLOYEE           |")
				fmt.Println("==================================")
				fmt.Print("Enter Employee Name     : ")
				fmt.Scan(&employeeName)
				fmt.Print("Enter Employee ID       : ")
				fmt.Scan(&employeeID)
				fmt.Print("Enter Employee Salary   : ")
				fmt.Scan(&salary)
				fmt.Print("Enter Years of Service : ")
				fmt.Scan(&yearsOfService)
				fmt.Print("Enter Department       : ")
				fmt.Scan(&department)
				isActive = true
				isEmployeeSet = true
				fmt.Println("\n✅ Employee added successfully!")
				fmt.Println("==================================")
			} else {
				fmt.Println("\n❌ Employee already exists! Use update option.")
			}

		case 2:
			if isEmployeeSet {
				fmt.Println("\n==================================")
				fmt.Println("|        EMPLOYEE DETAILS        |")
				fmt.Println("==================================")
				fmt.Printf("| Name          : %-14s|\n", employeeName)
				fmt.Printf("| Employee ID   : %-14d|\n", employeeID)
				fmt.Printf("| Salary        : ₹%-13.2f|\n", salary)
				fmt.Printf("| Department    : %-14s|\n", department)
				fmt.Printf("| Service Years : %-14d|\n", yearsOfService)
				fmt.Printf("| Active Status : %-14t|\n", isActive)
				fmt.Println("==================================")
			} else {
				fmt.Println("\n❌ No employee data exists!")
			}

		case 3:
			if isEmployeeSet {
				fmt.Println("\n==================================")
				fmt.Println("|       UPDATE EMPLOYEE          |")
				fmt.Println("==================================")
				fmt.Print("Enter new Name       : ")
				fmt.Scan(&employeeName)
				fmt.Print("Enter new Salary       : ")
				fmt.Scan(&salary)
				fmt.Print("Enter new Department   : ")
				fmt.Scan(&department)
				fmt.Print("Enter new Service Years: ")
				fmt.Scan(&yearsOfService)
				fmt.Println("\n✅ Employee updated successfully!")
				fmt.Println("==================================")
			} else {
				fmt.Println("\n❌ No employee data exists!")
			}

		case 4:
			if isEmployeeSet {
				fmt.Println("\n==================================")
				fmt.Println("|       DELETE EMPLOYEE          |")
				fmt.Println("==================================")
				isEmployeeSet = false
				isActive = false
				fmt.Println("\n✅ Employee deleted successfully!")
				fmt.Println("==================================")
			} else {
				fmt.Println("\n❌ No employee data exists!")
			}

		case 5:
			fmt.Println("\n==================================")
			fmt.Println("|          GOODBYE!              |")
			fmt.Println("==================================")
			return

		default:
			fmt.Println("\n❌ Invalid choice! Please enter 1-5")
		}
	}
}
