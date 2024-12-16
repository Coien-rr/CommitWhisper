package git

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGetChangedFiles(t *testing.T) {
	dir, removeDirFunc := createTempGitDir(t)
	defer removeDirFunc()
	gotFiles, err := getUnstagedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	wantFiles := getExpectFiles(t, dir)

	assertFiles(t, wantFiles, gotFiles)
}

func showGitStatus(t *testing.T, dirPath string) {
	gitStatusCmd := exec.Command("git", "status")
	gitStatusCmd.Dir = dirPath
	if out, err := gitStatusCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to exec git status in repo: %v\n%s", err, out)
	}
}

func createTempGitDir(t *testing.T) (string, func()) {
	t.Helper()

	tempGitDir, err := os.MkdirTemp("", "tempGit")
	if err != nil {
		t.Fatalf("could not create temp dir %v", err)
	}

	removeDir := func() {
		os.RemoveAll(tempGitDir)
	}

	gitInitCmd := exec.Command("git", "init")
	gitInitCmd.Dir = tempGitDir
	if out, err := gitInitCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to init git repo: %v\n%s", err, out)
	}

	initGitIgnore(t, tempGitDir)
	execGitAdd(t, tempGitDir)

	createTempFile(t, tempGitDir)
	execGitAdd(t, tempGitDir)

	showGitStatus(t, tempGitDir)

	createTempFile(t, tempGitDir)

	return tempGitDir, removeDir
}

func execGitAdd(t *testing.T, dirPath string) {
	gitAddCmd := exec.Command("git", "add", ".")
	gitAddCmd.Dir = dirPath
	if out, err := gitAddCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to exec git add in repo: %v\n%s", err, out)
	}
}

func initGitIgnore(t *testing.T, dirPath string) {
	t.Helper()
	path := dirPath + "/.gitignore"
	gitIgnoreFile, err := os.Create(path)
	if err != nil {
		t.Fatal("initGitIgnoreError: %w", err)
	}

	gitIgnoreFile.Write([]byte("*.png"))

	defer gitIgnoreFile.Close()
}

func createTempFile(t *testing.T, dirPath string) {
	t.Helper()

	tempFileExts := [3]string{"txt", "go", "png"}
	tempFileContent := map[string]string{
		"txt": "test",
		"go": `package main
import "fmt"
func main() {
	fmt.Println("Hello, World!")
}`,
	}

	for _, ext := range tempFileExts {
		pattern := "test_*." + ext
		fileWriter, err := os.CreateTemp(dirPath, pattern)
		if err != nil {
			t.Fatal("createTempFileError: %w", err)
		}
		if content, exists := tempFileContent[ext]; exists {
			fileWriter.Write([]byte(content))
		}
	}
}

func getExpectFiles(t *testing.T, dirPath string) []string {
	t.Helper()

	allFiles := make([]string, 0)
	filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			allFiles = append(allFiles, filepath.Base(path))
		}

		if d.IsDir() && path != dirPath {
			return filepath.SkipDir
		}

		return nil
	})

	files := make([]string, 0)

	for _, file := range allFiles {
		if file != "" && !strings.Contains(file, ".png") {
			files = append(files, file)
		}
	}
	return files
}

func assertFiles(t *testing.T, wantFiles, gotFiles []string) {
	t.Helper()

	if !reflect.DeepEqual(gotFiles, wantFiles) {
		t.Errorf("want: %v, but got: %v", wantFiles, gotFiles)
	}
}
