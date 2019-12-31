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
	SSHKeyPath          string `json:"ssh_key_path"`          // SSH 登录服务器的密钥文件
	Server              []struct {
		SSH             string `json:"ssh"`               // ssh 登录信息
		RsyncRemotePath string `json:"rsync_remote_path"` // 远程服务器项目路径
		SyncAfterCmd    string `json:"sync_after_cmd"`    // 同步文件之后执行的命令
	} `json:"server"`
}

func main() {
	m := multiconfig.NewWithPath("config.json")
	conf := new(Config)
	m.MustLoad(conf) // Panic's if there is any error

	hook, _ := github.New(github.Options.Secret(conf.GithubWebhookSecret))

	http.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		// 方便调试
		// updateCode(conf)

		payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		fmt.Println("收到请求")

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

	if len(conf.CloneAfterCmd) > 0 {
		command(fmt.Sprintf("cd %s  && %s", webName, conf.CloneAfterCmd))
	}

	for _, value := range conf.Server {
		fmt.Println(value.RsyncRemotePath)
		sshRemotePath := fmt.Sprintf("%s:%s", value.SSH, value.RsyncRemotePath)
		rsync(conf, sshRemotePath, webName)

		sshCmd := fmt.Sprintf("ssh -i %s %s \"%s\"", conf.SSHKeyPath, value.SSH, value.SyncAfterCmd)
		command(sshCmd)
		fmt.Println("====================================")
		fmt.Println("================成功=================")
		fmt.Println("====================================")
	}
}

func gitClone(gitURL string, webName string) {
	command("rm -rf " + webName)
	command(fmt.Sprintf("git clone --depth=1 %s %s", gitURL, webName))
}

func rsync(conf *Config, sshRemotePath string, webName string) {
	command(fmt.Sprintf("rsync -avuz -e 'ssh -i %s' ./%s/* %s", conf.SSHKeyPath, webName, sshRemotePath))
}

func command(cmd string) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s", cmd)
	}
	fmt.Println(string(out))
}
