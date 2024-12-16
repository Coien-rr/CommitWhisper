package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/Coien-rr/CommitWhisper/pkg/utils"
)

var (
	gitRepoPath string
	once        sync.Once
)

func discoveryRepoPath() {
	gitRepoPath = ""
	curPath, err := os.Getwd()
	if err != nil {
		return
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	for {
		if curPath == userHomeDir || curPath == "/" {
			return
		}

		gitDir := filepath.Join(curPath, ".git")

		if _, err = os.Stat(gitDir); err == nil {
			gitRepoPath = curPath
			return
		}

		curPath = filepath.Dir(curPath)
	}
}

func getRepoPath() string {
	once.Do(discoveryRepoPath)
	return gitRepoPath
}

func IsGitRepo() bool {
	if repoPath := getRepoPath(); repoPath == "" {
		return false
	}
	return true
}

func GetGitDiff() (string, error) {
	unstagedFiles, err := getUnstagedFiles(getRepoPath())
	if err != nil {
		return "", fmt.Errorf("GetGitDiffError: %w", err)
	}

	stagedFiles, err := getStagedFiles(getRepoPath())
	if err != nil {
		return "", fmt.Errorf("GetGitDiffError: %w", err)
	}

	stagedFiles, excludedFiles := excludeFiles(stagedFiles)

	notifyFiles(unstagedFiles, stagedFiles, excludedFiles)

	return getStagedDiffDetails(stagedFiles, getRepoPath())
}

func getUnstagedFiles(repoPath string) ([]string, error) {
	modifiedCmd := exec.Command("git", "ls-files", "--modified")
	modifiedCmd.Dir = repoPath
	modifiedOut, err := modifiedCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("getUnstagedFilesError(modified): %w", err)
	}

	otherCmd := exec.Command("git", "ls-files", "--other", "--exclude-standard")
	otherCmd.Dir = repoPath
	otherOut, err := otherCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("getUnstagedFilesError(other): %w", err)
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

func getStagedFiles(repoPath string) ([]string, error) {
	getStagedFilesCmd := exec.Command("git", "diff", "--staged", "--name-only")
	getStagedFilesCmd.Dir = repoPath
	gitDiffFilesOut, err := getStagedFilesCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("getStagedFilesError: %w", err)
	}

	stagedFiles := splitFilesFromOutput(gitDiffFilesOut)

	return stagedFiles, nil
}

func excludeFiles(stagedFiles []string) (includedFiles, excludedFiles []string) {
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

	for _, file := range stagedFiles {
		isExcluded := false
		for _, excludeExt := range excludeExtensions {
			if strings.Contains(file, excludeExt) {
				isExcluded = true
				break
			}
		}
		if isExcluded {
			excludedFiles = append(excludedFiles, file)
		} else {
			includedFiles = append(includedFiles, file)
		}
	}

	return
}

// func getDiff(diffFiles []string) (string, error) {
// 	filesWithoutLocks := make([]string, 0)
// 	for _, file := range diffFiles {
// 		if !strings.Contains(file, ".lock") && !strings.Contains(file, "-lock.") {
// 			filesWithoutLocks = append(filesWithoutLocks, file)
// 		}
// 	}
// 	return "", nil
// }

func getStagedDiffDetails(stagedFiles []string, repoPath string) (string, error) {
	getStagedDiffCmd := exec.Command("git", "diff", "--staged")
	getStagedDiffCmd.Dir = repoPath
	getStagedDiffCmd.Args = append(getStagedDiffCmd.Args, stagedFiles...)
	diffDetails, err := getStagedDiffCmd.Output()
	if err != nil {
		return "", fmt.Errorf("getDiffDetailsError: %w", err)
	}

	return string(diffDetails), nil
}

func notifyFiles(unstagedFiles, stagedFiles, excludedFiles []string) {
	notifyUnStagedFiles(unstagedFiles)
	notifyExcludedFiles(excludedFiles)
	notifyStagedFiles(stagedFiles)
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
			"Commit messages are generated for these Staged files:",
			stagedFiles,
		)
	}
}

func notifyUnStagedFiles(unstagedFiles []string) {
	if len(unstagedFiles) != 0 {
		utils.WhisperPrinter.WarningDisplayLists(
			"Unstaged changes of these files will not be used for Commit messages generating:",
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
