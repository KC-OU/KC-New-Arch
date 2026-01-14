package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Color codes for CLI
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorBold    = "\033[1m"
)

// Admin credentials
const (
	DefaultAdminPassword = "admin123"
)

// User represents a quiz user with their progress
type User struct {
	ID        string                      `json:"id"`
	Name      string                      `json:"name"`
	CreatedAt time.Time                   `json:"created_at"`
	Scores    map[string]map[string]Score `json:"scores"` // category -> module -> score
}

// Score tracks user performance in a module
type Score struct {
	Correct   int       `json:"correct"`
	Total     int       `json:"total"`
	LastTaken time.Time `json:"last_taken"`
}

// Question represents a quiz question
type Question struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"` // index of correct answer
	Category string   `json:"category"`
	Module   string   `json:"module"`
}

// QuizData holds all quiz questions
type QuizData struct {
	Questions []Question `json:"questions"`
}

// AdminConfig stores admin password
type AdminConfig struct {
	Password string `json:"password"`
}

var (
	currentUser   *User
	quizData      QuizData
	adminConfig   AdminConfig
	reader        *bufio.Reader
	cacheDir      string
	usersFile     string
	questionsFile string
	adminFile     string
)

func main() {
	reader = bufio.NewReader(os.Stdin)

	// Setup cache directory
	setupCacheDirectory()

	// Load or create data files
	loadData()

	// User login/registration
	userLogin()

	// Main menu loop
	for {
		showMainMenu()
	}
}

func setupCacheDirectory() {
	// Get user cache directory
	homeDir, err := os.UserCacheDir()
	if err != nil {
		homeDir, _ = os.UserHomeDir()
		cacheDir = filepath.Join(homeDir, ".cyber-quiz")
	} else {
		cacheDir = filepath.Join(homeDir, "cyber-quiz")
	}

	// Create cache directory if it doesn't exist
	os.MkdirAll(cacheDir, 0755)

	// Set file paths
	usersFile = filepath.Join(cacheDir, "users.json")
	questionsFile = filepath.Join(cacheDir, "questions.json")
	adminFile = filepath.Join(cacheDir, "admin.json")

	printColor(ColorCyan, fmt.Sprintf("📁 Data stored in: %s\n", cacheDir))
	time.Sleep(1 * time.Second)
}

func loadData() {
	// Load admin config
	if _, err := os.Stat(adminFile); err == nil {
		data, _ := os.ReadFile(adminFile)
		json.Unmarshal(data, &adminConfig)
	} else {
		// Create default admin config
		adminConfig.Password = DefaultAdminPassword
		saveAdminConfig()
	}

	// Load questions
	if _, err := os.Stat(questionsFile); err == nil {
		data, _ := os.ReadFile(questionsFile)
		json.Unmarshal(data, &quizData)
	} else {
		// Create default questions
		createDefaultQuestions()
		saveQuestions()
	}
}

func createDefaultQuestions() {
	quizData.Questions = []Question{
		// CompTIA PenTest+ Questions
		{
			ID:       "pt1",
			Question: "What is the primary purpose of a penetration test?",
			Options: []string{
				"To fix all vulnerabilities",
				"To identify and exploit vulnerabilities in a controlled manner",
				"To install security software",
				"To train employees",
			},
			Answer:   1,
			Category: "CompTIA",
			Module:   "PenTest+",
		},
		{
			ID:       "pt2",
			Question: "Which phase comes first in a penetration testing methodology?",
			Options: []string{
				"Exploitation",
				"Reporting",
				"Planning and Reconnaissance",
				"Post-Exploitation",
			},
			Answer:   2,
			Category: "CompTIA",
			Module:   "PenTest+",
		},
		{
			ID:       "pt3",
			Question: "What tool is commonly used for network scanning?",
			Options: []string{
				"Wireshark",
				"Nmap",
				"Metasploit",
				"John the Ripper",
			},
			Answer:   1,
			Category: "CompTIA",
			Module:   "PenTest+",
		},
		{
			ID:       "pt4",
			Question: "What does OSINT stand for?",
			Options: []string{
				"Operating System Intelligence",
				"Open Source Intelligence",
				"Online Security Internet",
				"Organized Security Interface",
			},
			Answer:   1,
			Category: "CompTIA",
			Module:   "PenTest+",
		},
		{
			ID:       "pt5",
			Question: "Which of the following is a social engineering attack?",
			Options: []string{
				"SQL Injection",
				"Buffer Overflow",
				"Phishing",
				"XSS Attack",
			},
			Answer:   2,
			Category: "CompTIA",
			Module:   "PenTest+",
		},

		// Cisco CCNA Questions
		{
			ID:       "ccna1",
			Question: "What is the default administrative distance of OSPF?",
			Options: []string{
				"90",
				"100",
				"110",
				"120",
			},
			Answer:   2,
			Category: "Cisco",
			Module:   "CCNA",
		},
		{
			ID:       "ccna2",
			Question: "Which layer of the OSI model does a switch operate at?",
			Options: []string{
				"Layer 1 - Physical",
				"Layer 2 - Data Link",
				"Layer 3 - Network",
				"Layer 4 - Transport",
			},
			Answer:   1,
			Category: "Cisco",
			Module:   "CCNA",
		},
		{
			ID:       "ccna3",
			Question: "What is the maximum number of usable hosts in a /26 subnet?",
			Options: []string{
				"30",
				"62",
				"126",
				"254",
			},
			Answer:   1,
			Category: "Cisco",
			Module:   "CCNA",
		},
		{
			ID:       "ccna4",
			Question: "Which protocol is used by ping?",
			Options: []string{
				"TCP",
				"UDP",
				"ICMP",
				"ARP",
			},
			Answer:   2,
			Category: "Cisco",
			Module:   "CCNA",
		},
		{
			ID:       "ccna5",
			Question: "What does STP stand for?",
			Options: []string{
				"Simple Transfer Protocol",
				"Spanning Tree Protocol",
				"Secure Transmission Protocol",
				"Switch Transport Protocol",
			},
			Answer:   1,
			Category: "Cisco",
			Module:   "CCNA",
		},
	}
}

func userLogin() {
	clearScreen()
	printBoxHeader("Cyber Learning Quiz Application", ColorCyan)
	fmt.Println()

	printColor(ColorYellow, "Are you a:\n")
	fmt.Println("1. New User")
	fmt.Println("2. Returning User")
	printColor(ColorYellow, "\nEnter choice (1-2): ")

	choice := readInput()

	switch choice {
	case "1":
		createNewUser()
	case "2":
		loginExistingUser()
	default:
		printColor(ColorRed, "Invalid choice. Creating new user...\n")
		time.Sleep(1 * time.Second)
		createNewUser()
	}
}

func createNewUser() {
	clearScreen()
	printBoxHeader("New User Registration", ColorGreen)
	fmt.Println()

	printColor(ColorYellow, "Enter your name: ")
	name := readInput()

	// Generate user-friendly ID
	userID := generateUserID(name)

	currentUser = &User{
		ID:        userID,
		Name:      name,
		CreatedAt: time.Now(),
		Scores:    make(map[string]map[string]Score),
	}

	saveUser()

	printColor(ColorGreen, fmt.Sprintf("\n✓ Welcome, %s! Your User ID is: ", name))
	printColor(ColorBold+ColorCyan, userID+"\n")
	printColor(ColorYellow, "\nPress Enter to continue...")
	readInput()
}

func generateUserID(name string) string {
	// Convert name to lowercase and remove spaces
	baseName := strings.ToLower(strings.ReplaceAll(name, " ", ""))

	// Load existing users to find next number
	users := loadAllUsers()
	count := 1

	for {
		userID := fmt.Sprintf("%s%d", baseName, count)
		exists := false

		for _, u := range users {
			if u.ID == userID {
				exists = true
				break
			}
		}

		if !exists {
			return userID
		}
		count++
	}
}

func loginExistingUser() {
	users := loadAllUsers()

	if len(users) == 0 {
		printColor(ColorRed, "\nNo existing users found. Creating new user...\n")
		time.Sleep(1 * time.Second)
		createNewUser()
		return
	}

	clearScreen()
	printBoxHeader("Returning Users", ColorBlue)
	fmt.Println()

	for i, user := range users {
		printColor(ColorCyan, fmt.Sprintf("%d. ", i+1))
		printColor(ColorWhite, fmt.Sprintf("%s ", user.Name))
		printColor(ColorYellow, fmt.Sprintf("(ID: %s)\n", user.ID))
	}

	printColor(ColorYellow, "\nEnter user number: ")
	choice := readInput()

	var userIndex int
	fmt.Sscanf(choice, "%d", &userIndex)
	userIndex--

	if userIndex >= 0 && userIndex < len(users) {
		currentUser = &users[userIndex]
		printColor(ColorGreen, fmt.Sprintf("\n✓ Welcome back, %s!\n", currentUser.Name))
	} else {
		printColor(ColorRed, "Invalid selection. Creating new user...\n")
		time.Sleep(1 * time.Second)
		createNewUser()
		return
	}

	printColor(ColorYellow, "Press Enter to continue...")
	readInput()
}

func showMainMenu() {
	clearScreen()
	printBoxHeader(fmt.Sprintf("User: %s (%s)", currentUser.Name, currentUser.ID), ColorCyan)
	printColor(ColorBold+ColorMagenta, "           MAIN MENU\n")
	fmt.Println()

	fmt.Println("1. 📝 Take Quiz")
	fmt.Println("2. 📊 View Scores")
	fmt.Println("3. 👤 Switch User")
	fmt.Println("4. 🔧 Admin Panel")
	fmt.Println("5. ❌ Exit")
	printColor(ColorYellow, "\nEnter choice (1-5): ")

	choice := readInput()

	switch choice {
	case "1":
		selectQuizModule()
	case "2":
		viewScores()
	case "3":
		userLogin()
	case "4":
		adminPanel()
	case "5":
		printColor(ColorGreen, "\nThank you for using Cyber Learning Quiz!\n")
		printColor(ColorCyan, "Your progress has been saved.\n")
		os.Exit(0)
	default:
		printColor(ColorRed, "Invalid choice. Press Enter to continue...")
		readInput()
	}
}

func adminPanel() {
	clearScreen()
	printBoxHeader("Admin Authentication", ColorRed)
	fmt.Println()

	printColor(ColorYellow, "Enter admin password: ")
	password := readInput()

	if password != adminConfig.Password {
		printColor(ColorRed, "\n✗ Access Denied! Incorrect password.\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	printColor(ColorGreen, "\n✓ Access Granted!\n")
	time.Sleep(1 * time.Second)

	for {
		clearScreen()
		printBoxHeader("Admin Panel", ColorRed)
		printColor(ColorBold+ColorYellow, "⚠ Administrator Mode Active ⚠\n")
		fmt.Println()

		fmt.Println("1. ➕ Add New Question")
		fmt.Println("2. ➖ Remove Question")
		fmt.Println("3. 📁 Add New Module")
		fmt.Println("4. 🗑️  Remove Module")
		fmt.Println("5. 👥 Manage Users")
		fmt.Println("6. 📋 List All Questions")
		fmt.Println("7. 🔑 Change Admin Password")
		fmt.Println("8. ⬅️  Back to Main Menu")
		printColor(ColorYellow, "\nEnter choice (1-8): ")

		choice := readInput()

		switch choice {
		case "1":
			addNewQuestion()
		case "2":
			removeQuestion()
		case "3":
			addNewModule()
		case "4":
			removeModule()
		case "5":
			manageUsers()
		case "6":
			listAllQuestions()
		case "7":
			changeAdminPassword()
		case "8":
			return
		default:
			printColor(ColorRed, "Invalid choice. Press Enter to continue...")
			readInput()
		}
	}
}

func selectQuizModule() {
	clearScreen()
	printBoxHeader("Select Quiz Module", ColorBlue)
	fmt.Println()

	modules := getAvailableModules()

	if len(modules) == 0 {
		printColor(ColorRed, "No quiz modules available.\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	idx := 1
	moduleList := make(map[int]struct{ category, module string })

	for category, mods := range modules {
		printColor(ColorCyan+ColorBold, fmt.Sprintf("\n%s:\n", category))
		for _, mod := range mods {
			questionCount := countQuestions(category, mod)
			printColor(ColorWhite, fmt.Sprintf("  %d. ", idx))
			printColor(ColorGreen, fmt.Sprintf("%s ", mod))
			printColor(ColorYellow, fmt.Sprintf("(%d questions)\n", questionCount))
			moduleList[idx] = struct{ category, module string }{category, mod}
			idx++
		}
	}

	printColor(ColorRed, fmt.Sprintf("\n%d. Back to Main Menu\n", idx))
	printColor(ColorYellow, "\nEnter choice: ")

	var choice int
	fmt.Sscanf(readInput(), "%d", &choice)

	if choice == idx {
		return
	}

	if mod, exists := moduleList[choice]; exists {
		takeQuiz(mod.category, mod.module)
	} else {
		printColor(ColorRed, "Invalid choice.\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
	}
}

func takeQuiz(category, module string) {
	questions := getQuestionsByModule(category, module)

	if len(questions) == 0 {
		printColor(ColorRed, "No questions available for this module.\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	correct := 0
	total := len(questions)

	for i, q := range questions {
		clearScreen()
		printColor(ColorCyan+ColorBold, fmt.Sprintf("╔════════════════════════════════════════╗\n"))
		printColor(ColorCyan, fmt.Sprintf("║ %s - %s\n", category, module))
		printColor(ColorYellow, fmt.Sprintf("║ Question %d of %d\n", i+1, total))
		printColor(ColorCyan+ColorBold, fmt.Sprintf("╚════════════════════════════════════════╝\n\n"))

		printColor(ColorWhite+ColorBold, q.Question+"\n\n")

		for j, opt := range q.Options {
			printColor(ColorCyan, fmt.Sprintf("%d. ", j+1))
			fmt.Println(opt)
		}

		printColor(ColorYellow, "\nYour answer (1-4): ")
		var answer int
		fmt.Sscanf(readInput(), "%d", &answer)
		answer--

		if answer == q.Answer {
			printColor(ColorGreen+ColorBold, "\n✓ Correct!\n")
			correct++
		} else {
			printColor(ColorRed+ColorBold, "\n✗ Incorrect. ")
			printColor(ColorGreen, fmt.Sprintf("The correct answer was: %s\n", q.Options[q.Answer]))
		}

		printColor(ColorYellow, "\nPress Enter to continue...")
		readInput()
	}

	// Save score
	saveScore(category, module, correct, total)

	// Show results
	clearScreen()
	printBoxHeader("Quiz Completed", ColorGreen)
	fmt.Println()

	percentage := float64(correct) / float64(total) * 100

	printColor(ColorCyan, fmt.Sprintf("Module: %s - %s\n", category, module))
	printColor(ColorWhite, fmt.Sprintf("Score: %d/%d ", correct, total))

	if percentage >= 80 {
		printColor(ColorGreen+ColorBold, fmt.Sprintf("(%.1f%%) 🎉\n", percentage))
	} else if percentage >= 60 {
		printColor(ColorYellow+ColorBold, fmt.Sprintf("(%.1f%%) 👍\n", percentage))
	} else {
		printColor(ColorRed+ColorBold, fmt.Sprintf("(%.1f%%) 📚\n", percentage))
	}

	printColor(ColorYellow, "\nPress Enter to continue...")
	readInput()
}

func viewScores() {
	clearScreen()
	printBoxHeader("Your Scores", ColorMagenta)
	fmt.Println()

	if len(currentUser.Scores) == 0 {
		printColor(ColorYellow, "No scores recorded yet. Take a quiz to get started!\n")
	} else {
		for category, modules := range currentUser.Scores {
			printColor(ColorCyan+ColorBold, fmt.Sprintf("\n%s:\n", category))
			for module, score := range modules {
				percentage := float64(score.Correct) / float64(score.Total) * 100

				printColor(ColorWhite, fmt.Sprintf("  %s: ", module))
				printColor(ColorYellow, fmt.Sprintf("%d/%d ", score.Correct, score.Total))

				if percentage >= 80 {
					printColor(ColorGreen, fmt.Sprintf("(%.1f%%) ", percentage))
				} else if percentage >= 60 {
					printColor(ColorYellow, fmt.Sprintf("(%.1f%%) ", percentage))
				} else {
					printColor(ColorRed, fmt.Sprintf("(%.1f%%) ", percentage))
				}

				printColor(ColorCyan, fmt.Sprintf("- Last taken: %s\n",
					score.LastTaken.Format("2006-01-02 15:04")))
			}
		}
	}

	printColor(ColorYellow, "\nPress Enter to continue...")
	readInput()
}

func addNewQuestion() {
	clearScreen()
	printBoxHeader("Add New Question", ColorGreen)
	fmt.Println()

	printColor(ColorYellow, "Enter Category (e.g., CompTIA, Cisco): ")
	category := readInput()

	printColor(ColorYellow, "Enter Module (e.g., PenTest+, CCNA): ")
	module := readInput()

	printColor(ColorYellow, "\nEnter Question: ")
	question := readInput()

	options := make([]string, 4)
	for i := 0; i < 4; i++ {
		printColor(ColorCyan, fmt.Sprintf("Enter Option %d: ", i+1))
		options[i] = readInput()
	}

	printColor(ColorYellow, "\nEnter correct answer number (1-4): ")
	var answer int
	fmt.Sscanf(readInput(), "%d", &answer)
	answer--

	if answer < 0 || answer > 3 {
		printColor(ColorRed, "Invalid answer number.\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	newQuestion := Question{
		ID:       fmt.Sprintf("q%d", time.Now().Unix()),
		Question: question,
		Options:  options,
		Answer:   answer,
		Category: category,
		Module:   module,
	}

	quizData.Questions = append(quizData.Questions, newQuestion)
	saveQuestions()

	printColor(ColorGreen+ColorBold, "\n✓ Question added successfully!\n")
	printColor(ColorYellow, "Press Enter to continue...")
	readInput()
}

func removeQuestion() {
	clearScreen()
	printBoxHeader("Remove Question", ColorRed)
	fmt.Println()

	if len(quizData.Questions) == 0 {
		printColor(ColorRed, "No questions available to remove.\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	for i, q := range quizData.Questions {
		printColor(ColorCyan, fmt.Sprintf("%d. ", i+1))
		printColor(ColorYellow, fmt.Sprintf("[%s - %s] ", q.Category, q.Module))
		printColor(ColorWhite, fmt.Sprintf("%s\n", q.Question))
	}

	printColor(ColorYellow, fmt.Sprintf("\nEnter question number to remove (1-%d) or 0 to cancel: ", len(quizData.Questions)))
	var choice int
	fmt.Sscanf(readInput(), "%d", &choice)

	if choice == 0 {
		return
	}

	if choice < 1 || choice > len(quizData.Questions) {
		printColor(ColorRed, "Invalid choice.\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	// Remove question
	quizData.Questions = append(quizData.Questions[:choice-1], quizData.Questions[choice:]...)
	saveQuestions()

	printColor(ColorGreen+ColorBold, "\n✓ Question removed successfully!\n")
	printColor(ColorYellow, "Press Enter to continue...")
	readInput()
}

func addNewModule() {
	clearScreen()
	printBoxHeader("Add New Module", ColorGreen)
	fmt.Println()

	printColor(ColorCyan, "To add a new module, simply add questions with the new category/module name.\n")
	printColor(ColorCyan, "Modules are created automatically when you add questions.\n\n")

	printColor(ColorYellow, "Enter new Category name: ")
	category := readInput()

	printColor(ColorYellow, "Enter new Module name: ")
	module := readInput()

	printColor(ColorGreen, fmt.Sprintf("\nNew module '%s - %s' will be created when you add questions to it.\n", category, module))
	printColor(ColorYellow, "Would you like to add a question now? (y/n): ")

	if strings.ToLower(readInput()) == "y" {
		addNewQuestion()
	} else {
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
	}
}

func removeModule() {
	clearScreen()
	printBoxHeader("Remove Module", ColorRed)
	fmt.Println()

	modules := getAvailableModules()

	if len(modules) == 0 {
		printColor(ColorRed, "No modules available to remove.\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	idx := 1
	moduleList := make(map[int]struct{ category, module string })

	for category, mods := range modules {
		printColor(ColorCyan+ColorBold, fmt.Sprintf("\n%s:\n", category))
		for _, mod := range mods {
			questionCount := countQuestions(category, mod)
			printColor(ColorWhite, fmt.Sprintf("  %d. %s (%d questions)\n", idx, mod, questionCount))
			moduleList[idx] = struct{ category, module string }{category, mod}
			idx++
		}
	}

	printColor(ColorYellow, fmt.Sprintf("\nEnter module number to remove (1-%d) or 0 to cancel: ", idx-1))
	var choice int
	fmt.Sscanf(readInput(), "%d", &choice)

	if choice == 0 {
		return
	}

	if mod, exists := moduleList[choice]; exists {
		printColor(ColorRed+ColorBold, fmt.Sprintf("\n⚠ WARNING: This will delete all questions in %s - %s!\n", mod.category, mod.module))
		printColor(ColorYellow, "Are you sure? (yes/no): ")

		if strings.ToLower(readInput()) == "yes" {
			// Remove all questions from this module
			newQuestions := []Question{}
			for _, q := range quizData.Questions {
				if !(q.Category == mod.category && q.Module == mod.module) {
					newQuestions = append(newQuestions, q)
				}
			}
			quizData.Questions = newQuestions
			saveQuestions()

			printColor(ColorGreen+ColorBold, "\n✓ Module removed successfully!\n")
		} else {
			printColor(ColorYellow, "\nCancelled.\n")
		}
	} else {
		printColor(ColorRed, "Invalid choice.\n")
	}

	printColor(ColorYellow, "Press Enter to continue...")
	readInput()
}

func manageUsers() {
	clearScreen()
	printBoxHeader("User Management", ColorMagenta)
	fmt.Println()

	users := loadAllUsers()

	if len(users) == 0 {
		printColor(ColorYellow, "No users found.\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	for i, user := range users {
		printColor(ColorCyan, fmt.Sprintf("%d. ", i+1))
		printColor(ColorWhite, fmt.Sprintf("%s ", user.Name))
		printColor(ColorYellow, fmt.Sprintf("(ID: %s) ", user.ID))
		printColor(ColorGreen, fmt.Sprintf("- Created: %s\n", user.CreatedAt.Format("2006-01-02")))
	}

	fmt.Println("\n1. Delete User")
	fmt.Println("2. Back")
	printColor(ColorYellow, "\nEnter choice: ")

	choice := readInput()

	if choice == "1" {
		printColor(ColorYellow, "Enter user number to delete: ")
		var userNum int
		fmt.Sscanf(readInput(), "%d", &userNum)
		userNum--

		if userNum >= 0 && userNum < len(users) {
			printColor(ColorRed+ColorBold, fmt.Sprintf("\n⚠ WARNING: Delete user %s?\n", users[userNum].Name))
			printColor(ColorYellow, "Type 'DELETE' to confirm: ")

			if readInput() == "DELETE" {
				users = append(users[:userNum], users[userNum+1:]...)
				data, _ := json.MarshalIndent(users, "", "  ")
				os.WriteFile(usersFile, data, 0644)

				printColor(ColorGreen+ColorBold, "\n✓ User deleted successfully!\n")
			} else {
				printColor(ColorYellow, "\nCancelled.\n")
			}
		}
	}

	printColor(ColorYellow, "Press Enter to continue...")
	readInput()
}

func listAllQuestions() {
	clearScreen()
	printBoxHeader("All Questions", ColorBlue)
	fmt.Println()

	modules := getAvailableModules()

	for category, mods := range modules {
		printColor(ColorCyan+ColorBold, fmt.Sprintf("\n%s:\n", category))
		for _, mod := range mods {
			questions := getQuestionsByModule(category, mod)
			printColor(ColorGreen, fmt.Sprintf("\n  %s (%d questions):\n", mod, len(questions)))
			for i, q := range questions {
				printColor(ColorYellow, fmt.Sprintf("    %d. ", i+1))
				printColor(ColorWhite, fmt.Sprintf("%s\n", q.Question))
			}
		}
	}

	printColor(ColorYellow, "\nPress Enter to continue...")
	readInput()
}

func changeAdminPassword() {
	clearScreen()
	printBoxHeader("Change Admin Password", ColorRed)
	fmt.Println()

	printColor(ColorYellow, "Enter current password: ")
	current := readInput()

	if current != adminConfig.Password {
		printColor(ColorRed, "\n✗ Incorrect password!\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	printColor(ColorYellow, "Enter new password: ")
	newPass := readInput()

	printColor(ColorYellow, "Confirm new password: ")
	confirmPass := readInput()

	if newPass != confirmPass {
		printColor(ColorRed, "\n✗ Passwords don't match!\n")
		printColor(ColorYellow, "Press Enter to continue...")
		readInput()
		return
	}

	adminConfig.Password = newPass
	saveAdminConfig()

	printColor(ColorGreen+ColorBold, "\n✓ Admin password changed successfully!\n")
	printColor(ColorYellow, "Press Enter to continue...")
	readInput()
}

// Helper functions

func getAvailableModules() map[string][]string {
	modules := make(map[string]map[string]bool)

	for _, q := range quizData.Questions {
		if modules[q.Category] == nil {
			modules[q.Category] = make(map[string]bool)
		}
		modules[q.Category][q.Module] = true
	}

	result := make(map[string][]string)
	for category, mods := range modules {
		for mod := range mods {
			result[category] = append(result[category], mod)
		}
	}

	return result
}

func getQuestionsByModule(category, module string) []Question {
	var questions []Question
	for _, q := range quizData.Questions {
		if q.Category == category && q.Module == module {
			questions = append(questions, q)
		}
	}
	return questions
}

func countQuestions(category, module string) int {
	count := 0
	for _, q := range quizData.Questions {
		if q.Category == category && q.Module == module {
			count++
		}
	}
	return count
}

func saveScore(category, module string, correct, total int) {
	if currentUser.Scores == nil {
		currentUser.Scores = make(map[string]map[string]Score)
	}

	if currentUser.Scores[category] == nil {
		currentUser.Scores[category] = make(map[string]Score)
	}

	currentUser.Scores[category][module] = Score{
		Correct:   correct,
		Total:     total,
		LastTaken: time.Now(),
	}

	saveUser()
}

func saveUser() {
	users := loadAllUsers()

	found := false
	for i, u := range users {
		if u.ID == currentUser.ID {
			users[i] = *currentUser
			found = true
			break
		}
	}

	if !found {
		users = append(users, *currentUser)
	}

	data, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile(usersFile, data, 0644)
}

func loadAllUsers() []User {
	var users []User

	if _, err := os.Stat(usersFile); err == nil {
		data, _ := os.ReadFile(usersFile)
		json.Unmarshal(data, &users)
	}

	return users
}

func saveQuestions() {
	data, _ := json.MarshalIndent(quizData, "", "  ")
	os.WriteFile(questionsFile, data, 0644)
}

func saveAdminConfig() {
	data, _ := json.MarshalIndent(adminConfig, "", "  ")
	os.WriteFile(adminFile, data, 0644)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func readInput() string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func printColor(color, text string) {
	fmt.Print(color + text + ColorReset)
}

func printBoxHeader(title, color string) {
	fmt.Println(color + "╔════════════════════════════════════════╗" + ColorReset)
	fmt.Printf(color+"║ "+ColorReset+ColorBold+"%-38s"+ColorReset+color+" ║"+ColorReset+"\n", title)
	fmt.Println(color + "╚════════════════════════════════════════╝" + ColorReset)
}
