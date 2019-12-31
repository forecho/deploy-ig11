# GitHub WebHooks 部署工具



## 使用

使用配置文件

```sh
cp config_example.json config.json
```

修改配置文件内容


## 配置服务器

本地或者此工具所放置的服务器生成 `ssh` 密钥：

```sh
ssh-keygen -t rsa -b 4096 -C "for_deploy" -f ~/.ssh/for_deploy
```

然后把 `~/.ssh/for_deplo.pub` 文件内容依次拷贝到 `config` 配置的服务器的 `~/.ssh/authorized_keys` 文件中
