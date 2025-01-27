//to demonstrate the declaration and intialization of variables of various types including integers, float, strings and boolean you are to demonstrate the use of control flow such as if, for and switch statestements

package main

import "fmt"

func main() {

	const (
		AddEmployee = iota + 1
		ViewEmployee
		UpdateEmployee
		DeleteEmployee
		Exit
	)

	const (
		Minimum_Salary = 20000.00
		BaseEmployeeID = 100
	)

	// Department types
	const (
		IT         = "IT"
		HR         = "HR"
		Finance    = "Finance"
		Operations = "Operations"
		Marketing  = "Marketing"
	)

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
	fmt.Println("\n==================================")
	fmt.Println("|      EMPLOYEE MANAGEMENT       |")
	fmt.Println("==================================")
	fmt.Println("|      NOTE                      |")
	fmt.Println("==================================")
	fmt.Println("| Minimum_Salary = 20000.00      |")
	fmt.Println("| Base Employee ID : 100         |")
	fmt.Println("==================================")
	fmt.Println("|    AVAIABLE DEPARTMENTS        |")
	fmt.Println("==================================")
	fmt.Println("|IT                              |")
	fmt.Println("|HR                              |")
	fmt.Println("|Finance                         |")
	fmt.Println("|Operations                      |")
	fmt.Println("|Marketing                       |")
	fmt.Println("==================================")

	for {
		fmt.Println("\n==================================")
		fmt.Println("|      EMPLOYEE MANAGEMENT       |")
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
		case AddEmployee:
			if !isEmployeeSet {
				fmt.Println("\n==================================")
				fmt.Println("|         ADD EMPLOYEE           |")
				fmt.Println("==================================")
				fmt.Print("Enter Employee Name     : ")
				fmt.Scan(&employeeName)

				//self learning component: Employee Name Validation
				if len(employeeName) == 0 {
					fmt.Println("\n❌ Name cannot be empty")
					continue
				}

				fmt.Print("Enter Employee ID       : ")
				fmt.Scan(&employeeID)

				//self learning component: Employee ID Validation
				if employeeID <= BaseEmployeeID {
					fmt.Println("\n❌Invalid EmployeeID\nNot Adhering to the Base Employee ID")
					continue
				}

				fmt.Print("Enter Employee Salary   : ")
				fmt.Scan(&salary)

				//self learning component: Salary Validation
				if salary < Minimum_Salary {
					fmt.Println("\n❌Salary must adhere to the minimum salary")
					continue
				}

				fmt.Print("Enter Years of Service : ")
				fmt.Scan(&yearsOfService)

				fmt.Print("Enter Department       : ")
				fmt.Scan(&department)

				// self learning component: Department validation
				switch department {
				case IT, HR, Finance, Operations, Marketing:
					isActive = true
					isEmployeeSet = true
				default:
					fmt.Println("\n❌ Invalid department! Please choose from available departments")
					continue
				}

				isActive = true
				isEmployeeSet = true
				fmt.Println("\n✅ Employee added successfully!")
				fmt.Println("==================================")
			} else {
				fmt.Println("\n❌ Employee already exists! Use update option.")
			}

		case ViewEmployee:
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

		case UpdateEmployee:
			if isEmployeeSet {
				fmt.Println("\n==================================")
				fmt.Println("|       UPDATE EMPLOYEE          |")
				fmt.Println("==================================")
				fmt.Print("Enter new Name       : ")
				fmt.Scan(&employeeName)
				fmt.Print("Enter new Salary       : ")
				fmt.Scan(&salary)

				//self learning component: Salary Validation
				if salary < Minimum_Salary {
					fmt.Println("\n❌Salary must adhere to the minimum salary")
					continue
				}
				fmt.Print("Enter new Department   : ")
				fmt.Scan(&department)

				// self learning component: Department validation
				switch department {
				case IT, HR, Finance, Operations, Marketing:
					isActive = true
					isEmployeeSet = true
				default:
					fmt.Println("\n❌ Invalid department! Please choose from available departments")
					continue
				}

				fmt.Print("Enter new Service Years: ")
				fmt.Scan(&yearsOfService)
				fmt.Println("\n✅ Employee updated successfully!")
				fmt.Println("==================================")
			} else {
				fmt.Println("\n❌ No employee data exists!")
			}

		case DeleteEmployee:
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

		case Exit:
			fmt.Println("\n==================================")
			fmt.Println("|          GOODBYE!              |")
			fmt.Println("==================================")
			return

		default:
			fmt.Println("\n❌ Invalid choice! Please enter 1-5")
		}
	}
}
