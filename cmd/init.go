package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/thzoid/issue-mafia/util"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config file",
	Long:  "Execute the config file creation wizard that sets up your .issue-mafia file, containing the issue-mafia hooks repository.",
	Run: func(cmd *cobra.Command, args []string) {
		dirIsRepo, dirHasConfig := util.IsRepo("."), util.HasConfig(".")
		// Check if remote repository has hook files
		if !dirIsRepo {
			fmt.Println("Hang on! This does not look like a Git repository (which means issue-mafia won't be able to synchronize hooks).\nDo you want to proceed anyway? \u001b[90m(\u001b[1mY\u001b[0m\u001b[90m/\u001b[1mn\u001b[0m\u001b[90m)\u001b[0m: \u001b[1m")
			var answer string
			fmt.Scanf("%s", &answer)
			fmt.Print("\u001b[0m")
			fmt.Println()
			// Validate answer
			if answer != "Y" && answer != "y" {
				util.WarningLogger.Fatalln("no files generated.")
			}
		}
		if dirHasConfig {
			fmt.Println("It appears this directory already contains an issue-mafia configuration file.\nContinuing the process will overwrite settings. Do you want to proceed anyway? \u001b[90m(\u001b[1mY\u001b[0m\u001b[90m/\u001b[1mn\u001b[0m\u001b[90m)\u001b[0m: \u001b[1m")
			var answer string
			fmt.Scanf("%s", &answer)
			fmt.Print("\u001b[0m")
			fmt.Println()
			// Validate answer
			if answer != "Y" && answer != "y" {
				util.WarningLogger.Fatalln("no files generated.")
			}
		}
		fmt.Println("Welcome to issue-mafia!\nPlease, type the repository with which you would like to synchronize Git hooks:")
		fmt.Print("\u001b[90mgithub.com/\u001b[0m\u001b[1m")

		// Get repository
		var repo string
		fmt.Scanf("%s", &repo)
		fmt.Print("\u001b[0m")
		fmt.Println()

		// Validate repository
		re := regexp.MustCompile(`^[-a-zA-Z0-9_]+\/[-a-zA-Z0-9_]+$`)
		if !re.Match([]byte(repo)) {
			util.ErrorLogger.Fatalln("invalid repository address")
		}

		// Check repository existence
		util.InfoLogger.Println("checking if repository is accessible...")
		if status := util.FetchRepository(repo); status != 200 {
			util.ErrorLogger.Fatalln("could not access the specified repository. received status", fmt.Sprintf("%d", status)+".")
		}
		fmt.Println()

		// Get branch
		fmt.Println("Please, specify the branch that issue-mafia should look for hooks \u001b[90m(default is \u001b[1mmain\u001b[0m\u001b[90m)\u001b[0m:\u001b[1m")
		var branch string
		fmt.Scanf("%s", &branch)
		fmt.Print("\u001b[0m")

		if branch == "" {
			branch = "main"
		} else {
			// Validate branch
			re = regexp.MustCompile(`^[-a-zA-Z0-9_]+$`)
			if !re.Match([]byte(branch)) {
				util.ErrorLogger.Fatalln("invalid branch name")
			}
		}

		// Check branch existence
		util.InfoLogger.Println("checking repository integrity...")
		files, status := util.FetchIntersectingFiles(repo, branch)
		if files == nil {
			util.ErrorLogger.Fatalln("could not access the specified branch. received status", fmt.Sprintf("%d", status)+".")
		}

		// Check if remote repository has hook files
		if len(files) == 0 {
			fmt.Println()
			fmt.Println("This does not look like an issue-mafia repository. Do you want to add it anyway? \u001b[90m(\u001b[1mY\u001b[0m\u001b[90m/\u001b[1mn\u001b[0m\u001b[90m)\u001b[0m: \u001b[1m")
			var answer string
			fmt.Scanf("%s", &answer)
			fmt.Print("\u001b[0m")
			fmt.Println()
			// Validate answer
			if answer != "Y" && answer != "y" {
				util.WarningLogger.Fatalln("no files generated.")
			}
		}

		// Generate config file
		util.InfoLogger.Println("generating file...")
		f, _ := os.OpenFile(".issue-mafia", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		fmt.Fprintf(f, "%s %s", repo, branch)
		f.Close()
		fmt.Println()
		fmt.Println("Configuration file created successfully! Run \u001b[90missue-mafia\u001b[0m to synchronize hooks.")
		fmt.Println()
		fmt.Println("\u001b[1mAlways make sure you trust the repository that you are executing scripts from!\u001b[0m")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
