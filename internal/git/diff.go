package git

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/Coien-rr/CommitWhisper/pkg/utils"
)

func IsGitRepo() bool {
	statusCmd := exec.Command("git", "status")
	out, err := statusCmd.CombinedOutput()
	if err != nil &&
		strings.Contains(
			string(out),
			"fatal: not a git repository (or any of the parent directories): .git",
		) {
		return false
	}

	return true
}

func GetGitDiff() (string, error) {
	changedFiles, err := getUnstagedChangedFiles("./")
	if err != nil {
		return "", fmt.Errorf("GetGitDiffError: %w", err)
	}

	diffInfo, err := getDiff(changedFiles)
	if err != nil {
		return "", fmt.Errorf("GetGitDiffError: %w", err)
	}

	return diffInfo, nil
}

func getUnstagedChangedFiles(dirPath string) ([]string, error) {
	modifiedCmd := exec.Command("git", "ls-files", "--modified")
	modifiedCmd.Dir = dirPath
	modifiedOut, err := modifiedCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("GetChangedFilesError(modified): %w", err)
	}

	otherCmd := exec.Command("git", "ls-files", "--other", "--exclude-standard")
	otherCmd.Dir = dirPath
	otherOut, err := otherCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("GetChangedFilesError(other): %w", err)
	}

	files := make([]string, 0)
	files = append(files, strings.Split(string(modifiedOut), "\n")...)
	for _, file := range strings.Split(string(otherOut), "\n") {
		if file == "" {
			continue
		}
		files = append(files, file+"(󱙝 Untracked)")

	}

	filteredFiles := make([]string, 0)
	for _, file := range files {
		if file != "" {
			filteredFiles = append(filteredFiles, file)
		}
	}

	sort.Strings(filteredFiles)

	return filteredFiles, nil
}

func getDiff(diffFiles []string) (string, error) {
	excludeExtensions := []string{
		".lock",
		"-lock.",
		".svg",
		".png",
		".jpg",
		".jpeg",
		".webp",
		".gif",
	}

	excludedFiles := make([]string, 0)

	for _, file := range diffFiles {
		for _, excludeExt := range excludeExtensions {
			if strings.Contains(file, excludeExt) {
				excludedFiles = append(excludedFiles, file)
			}
		}
	}

	notifyExcludedFiles(excludedFiles)

	filesWithoutLocks := make([]string, 0)
	for _, file := range diffFiles {
		if !strings.Contains(file, ".lock") && !strings.Contains(file, "-lock.") {
			filesWithoutLocks = append(filesWithoutLocks, file)
		}
	}

	diffStagedFiles, err := getDiffFiles()
	if err != nil {
		return "", fmt.Errorf("RuntimeError: %w", err)
	}

	// TODO: refactor
	notifyUnStagedFiles(diffFiles)
	notifyStagedFiles(diffStagedFiles)

	return getDiffDetails()
}

func getDiffDetails() (string, error) {
	gitDiffCmd := exec.Command("git", "diff", "--staged", "--")
	// gitDiffCmd.Args = append(gitDiffCmd.Args, filesWithoutLocks...)
	diffDetails, err := gitDiffCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run git diff: %w", err)
	}

	return string(diffDetails), nil
}

func getDiffFiles() ([]string, error) {
	gitDiffFilesCmd := exec.Command("git", "diff", "--staged", "--name-only")
	gitDiffFilesOut, err := gitDiffFilesCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run [git diff --staged --name-only] : %w", err)
	}

	stagedFiles := splitFilesFromOutput(gitDiffFilesOut)

	return stagedFiles, nil
}

func notifyExcludedFiles(excludedFiles []string) {
	if len(excludedFiles) != 0 {
		utils.WhisperPrinter.WarningDisplayLists(
			"Some files are excluded by default from 'git diff'. No commit messages are generated for this files:",
			excludedFiles,
		)
	}
}

func notifyStagedFiles(stagedFiles []string) {
	if len(stagedFiles) != 0 {
		utils.WhisperPrinter.InfoDisplayLists(
			"Commit messages are generated for these staged files:",
			stagedFiles,
		)
	}
}

func notifyUnStagedFiles(unstagedFiles []string) {
	if len(unstagedFiles) != 0 {
		utils.WhisperPrinter.WarningDisplayLists(
			"Some unstaged changes of these files will not be used for Commit messages generating:",
			unstagedFiles,
		)
	}
}

func splitFilesFromOutput(buf []byte) []string {
	files := strings.Split(string(buf), "\n")
	filterFiles := make([]string, 0)
	for _, file := range files {
		if file != "" {
			filterFiles = append(filterFiles, file)
		}
	}

	return filterFiles
}
