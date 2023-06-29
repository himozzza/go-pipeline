package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

func main() {
	path := "zodiak_tg"
	url := "https://github.com/himozzza/zodiak_tg"

	var start bool
	dirs, _ := os.ReadDir("./")
	for _, dir := range dirs {
		if dir.Name() == path {
			start = true
			break
		}
	}

	if start {
		os.Chdir(path)
		git_init(path)
	} else {
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
		})
		if err != nil {
			os.WriteFile("clone_repo_error", []byte(err.Error()), 0777)
		}
		os.Chdir(path)
		build_release()
		git_init(path)
	}
}

func git_init(path string) {

	exec.Command("./target/release/zodiak_tg").Start() // Заменить для универсальности!!!

	for {
		storage, _ := git.PlainOpen("./")
		w, err := storage.Worktree()
		if err != nil {
			os.WriteFile("../open_dir_error", []byte(err.Error()), 0777)
		}
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			fmt.Println(err)
		} else {
			build_release()
			pid_id, _ := exec.Command("pgrep", "-f", path).Output()
			exec.Command("kill", "-9", strings.Split(string(pid_id), "\n")[0]).Run()
			exec.Command("./target/release/zodiak_tg").Start() // Заменить для универсальности!!!
		}
		time.Sleep(5 * time.Hour)
	}
}

func build_release() {
	exec.Command("cargo", "build", "--release").Output()
}
