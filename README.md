# GitHub WebHooks 部署工具

## 使用

使用配置文件

```sh
cp config_example.json config.json
```

修改配置文件内容，参数说明参考代码 [mian.go](https://github.com/forecho/ig11/blob/master/main.go#L12) 的代码。


## 部署

### 配置服务器

本地或者此工具所放置的服务器生成 `ssh` 密钥：

```sh
ssh-keygen -t rsa -b 4096 -C "deploy" -f ~/.ssh/deploy
```

然后把 `~/.ssh/deploy.pub` 文件内容依次拷贝到 `config` 配置的服务器的 `~/.ssh/authorized_keys` 文件中


### 其他配置（参考）

如果部署的是私有项目，可参考下面的代码：


- 给 GitHub 项目配置 `Deploy keys`，找到具体项目的 Settings 页面，然后找到 Deploy keys，然后可以把上面生成的公钥添加上去。（建议不要勾选 `Allow write access`，不勾选只具备拉取权限)
- 然后还要执行以下命令：

```
eval `ssh-agent -s`
ssh-add ~/.ssh/deploy
ssh-add -l
```

### 部署

去 [Releases](https://github.com/forecho/ig11/releases) 页面下载已经打包好的 Linux 版本，然后在服务器上执行下面命令即可部署。

```shell
mkdir -p /var/log/app/ig11/ && touch /var/log/app/ig11/script.log
/usr/bin/nohup $PWD/ig11 > /var/log/app/ig11/script.log 2>&1 &
```

检查确认命令：

```shell
lsof -i:8090
ps -ax | grep ig11
```

注：

- 默认是 8090 端口，目前只能通过修改代码来修改
- 每次修改配置文件之后，要 kill 掉进程，然后再开启


## 使用

部署完了之后就去 GitHub 对应的项目配置 WebHooks，找到具体项目的 Settings 页面，然后找到 Webhooks，然后添加就可以了.

- `Payload URL`：如果你没有使用 nginx 配置反向代理的话，这里写 `你服务器的ip:8090`
- `Content type`：默认就可以
- `Secret`：自己随机生成，与你第一步配置文件一致就可以
- `events`：选 `Send me everything.` 或者自定义事件中的 `Releases` (本项目只有收到 `Releases` 事件才触发部署)

