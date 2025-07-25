package cmd

import (
	"fmt"
	"log"

	"github.com/minmaxmean/sigma/siq"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read [siq-file]",
	Short: "Read and display information from a SIQ file",
	Long: `Read a SIQ file and display comprehensive information including:
- Package metadata (name, ID, version, difficulty, etc.)
- Statistics (rounds, themes, questions)
- File listing
- Question details with answers`,
	Args: cobra.ExactArgs(1),
	Run:  runRead,
}

func runRead(cmd *cobra.Command, args []string) {
	siqFile := args[0]

	// Open and read the SIQ file
	reader, err := siq.NewSIQReader(siqFile)
	if err != nil {
		log.Fatal("Failed to open SIQ file:", err)
	}
	defer reader.Close()

	// Parse the package
	pkg, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to read SIQ file:", err)
	}

	// Display package information
	fmt.Printf("=== SIG Pack Information ===\n")
	fmt.Printf("Name: %s\n", pkg.Name)
	fmt.Printf("ID: %s\n", pkg.ID)
	fmt.Printf("Version: %s\n", pkg.Version)
	fmt.Printf("Detected Format: v%d\n", reader.GetVersion())
	fmt.Printf("Difficulty: %d/10\n", pkg.Difficulty)
	fmt.Printf("Language: %s\n", pkg.Language)
	fmt.Printf("Publisher: %s\n", pkg.Publisher)
	fmt.Printf("Date: %s\n", pkg.Date)

	// Display statistics
	fmt.Printf("\n=== Statistics ===\n")
	fmt.Printf("Rounds: %d\n", pkg.GetRoundCount())
	fmt.Printf("Themes: %d\n", pkg.GetThemeCount())
	fmt.Printf("Questions: %d\n", pkg.GetQuestionCount())

	// Display files in archive
	fmt.Printf("\n=== Files in Archive ===\n")
	files := reader.ListFiles()
	for _, file := range files {
		fmt.Printf("  %s\n", file)
	}

	// Display questions summary
	fmt.Printf("\n=== Questions Summary ===\n")
	questions := pkg.GetAllQuestions()
	for i, question := range questions {
		fmt.Printf("Question %d:\n", i+1)
		fmt.Printf("  Type: %s\n", question.Type)
		fmt.Printf("  Right answers: %d\n", len(question.Right))
		fmt.Printf("  Wrong answers: %d\n", len(question.Wrong))

		// Show right answers
		if len(question.Right) > 0 {
			fmt.Printf("  Right answer(s):\n")
			for j, answer := range question.Right {
				fmt.Printf("    %d. %s\n", j+1, answer)
			}
		}

		// Show wrong answers
		if len(question.Wrong) > 0 {
			fmt.Printf("  Wrong answer(s):\n")
			for j, answer := range question.Wrong {
				fmt.Printf("    %d. %s\n", j+1, answer)
			}
		}

		// Show question content
		content := question.GetQuestionContent()
		if len(content) > 0 {
			fmt.Printf("  Content items: %d\n", len(content))
			for j, item := range content {
				fmt.Printf("    Item %d: %s (%s)\n", j+1, item.Type, item.Value)
			}
		}
		fmt.Println()
	}

	fmt.Println("SIQ file processed successfully!")
}

// GetReadCmd returns the read command
func GetReadCmd() *cobra.Command {
	return readCmd
}
