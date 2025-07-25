package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/minmaxmean/sigma/siq"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sigma",
	Short: "A CLI tool for processing SIQ files",
	Long: `Sigma is a command-line tool for reading and analyzing SIQ (SIG Pack) files.
It provides detailed information about SIG packages including questions, answers, and metadata.`,
}

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

var markdownCmd = &cobra.Command{
	Use:   "markdown [siq-file] [output-file]",
	Short: "Convert a SIQ file to markdown format",
	Long: `Convert a SIQ file to a readable markdown format including:
- Package metadata and statistics
- All rounds, themes, and questions
- Question content and answers
- Authors, sources, and comments`,
	Args: cobra.ExactArgs(2),
	Run:  runMarkdown,
}

func runMarkdown(cmd *cobra.Command, args []string) {
	siqFile := args[0]
	outputFile := args[1]

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

	// Generate markdown content
	markdown := generateMarkdown(pkg, reader.GetVersion())

	// Write to output file
	err = os.WriteFile(outputFile, []byte(markdown), 0644)
	if err != nil {
		log.Fatal("Failed to write markdown file:", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", siqFile, outputFile)
}

func generateMarkdown(pkg *siq.Package, version int) string {
	var sb strings.Builder

	// Process rounds
	rounds := pkg.Rounds
	if version == 4 {
		// Convert v4 rounds to v5 format for processing
		rounds = make([]siq.Round, len(pkg.RoundsV4))
		for i, roundV4 := range pkg.RoundsV4 {
			rounds[i] = convertV4RoundToV5(roundV4)
		}
	}

	for roundIndex, round := range rounds {
		sb.WriteString(fmt.Sprintf("## Round %d: %s\n\n", roundIndex+1, round.Name))

		// Process themes
		for themeIndex, theme := range round.Themes {
			sb.WriteString(fmt.Sprintf("### Theme %d: %s\n\n", themeIndex+1, theme.Name))

			// Process questions
			for questionIndex, question := range theme.Questions {
				sb.WriteString(fmt.Sprintf("#### Question %d\n\n", questionIndex+1))

				// Question content
				content := question.GetQuestionContent()
				if len(content) > 0 {
					sb.WriteString("**Content**:\n\n")
					for _, item := range content {
						sb.WriteString(fmt.Sprintf("- **%s**: %s", item.Type, item.Value))
						if item.Duration > 0 {
							sb.WriteString(fmt.Sprintf(" (duration: %d)", item.Duration))
						}
						if item.Placement != "" && item.Placement != "screen" {
							sb.WriteString(fmt.Sprintf(" (placement: %s)", item.Placement))
						}
						sb.WriteString("\n")
					}
					sb.WriteString("\n")
				}

				// Right answers
				if len(question.Right) > 0 {
					sb.WriteString("**Right Answer")
					if len(question.Right) > 1 {
						sb.WriteString("s")
					}
					sb.WriteString("**:\n\n")

					if len(question.Right) == 1 {
						sb.WriteString(question.Right[0] + "\n\n")
					} else {
						for i, answer := range question.Right {
							sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, answer))
						}
						sb.WriteString("\n")
					}
				}

				sb.WriteString("---\n\n")
			}
		}
	}

	return sb.String()
}

func convertV4RoundToV5(roundV4 siq.RoundV4) siq.Round {
	themes := make([]siq.Theme, len(roundV4.Themes))
	for i, themeV4 := range roundV4.Themes {
		questions := make([]siq.Question, len(themeV4.Questions))
		for j, questionV4 := range themeV4.Questions {
			questions[j] = siq.ConvertV4ToV5Question(questionV4)
		}
		themes[i] = siq.Theme{
			Name:      themeV4.Name,
			Info:      themeV4.Info,
			Questions: questions,
		}
	}
	return siq.Round{
		Name:   roundV4.Name,
		Info:   roundV4.Info,
		Themes: themes,
	}
}

func init() {
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(markdownCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
