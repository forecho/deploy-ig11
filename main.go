package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/koding/multiconfig"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// Config 环境以及服务器配置信息
type Config struct {
	GitURL              string `json:"git_url"`               // git 链接
	GithubWebhookSecret string `json:"github_webhook_secret"` // GitHub webhooks 密钥
	CloneAfterCmd       string `json:"clone_after_cmd"`       // git clone 之后执行的命令
	Server              []struct {
		RsyncRemotePath string `json:"rsync_remote_path"` // 远程路径
		SyncAfterCmd    string `json:"sync_after_cmd"`    // 同步文件之后执行的命令
	} `json:"server"`
}

func main() {
	m := multiconfig.NewWithPath("config.json")
	conf := new(Config)
	m.MustLoad(conf) // Panic's if there is any error

	hook, _ := github.New(github.Options.Secret(conf.GithubWebhookSecret))

	http.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			updateCode(conf)
			fmt.Printf("%+v", release)
		}
	})
	http.ListenAndServe(":8090", nil)
}

func updateCode(conf *Config) {
	// 默认为 web 项目
	webName := "web"
	gitClone(conf.GitURL, webName)
	fmt.Println(conf.GitURL)

	// cmd(conf.CloneAfterCmd)

	for _, value := range conf.Server {
		fmt.Println(value.RsyncRemotePath)
		rsync(value.RsyncRemotePath, webName)
		cmd(value.SyncAfterCmd)
		fmt.Println("====================================")
		fmt.Println("================成功=================")
		fmt.Println("====================================")
	}
}

func gitClone(gitURL string, webName string) {
	cmd("rm -rf " + webName)
	cmd(fmt.Sprintf("git clone --depth=1 %s %s", gitURL, webName))
}

func rsync(rsyncRemotePath string, webName string) {
	cmd(fmt.Sprintf("rsync -av ./%s/* %s", webName, rsyncRemotePath))
}

func cmd(cmd string) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s", cmd)
	}
	fmt.Println(string(out))
}
