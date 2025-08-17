# CLOUDFLARED_GUI

本工具是基于cloudflared，实现无公网IP远程访问局域网的方式。
工具依赖于cloudflared，自行安装
 https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.msi

## cloudflared操作（待补充）

## windows下编译

1. 安装mingw64
https://github.com/niXman/mingw-builds-binaries/releases

解压后，把目录/bin添加到PATH环境变量

2. 安装fyne
go install fyne.io/tools/cmd/fyne@latest

3. 检查环境（这一步可选）
https://geoffrey-artefacts.fynelabs.com/github/andydotxyz/fyne-io/setup/latest/

下载这个工具检查

4. 编译
```
mkdir releases
fyne package -os windows -icon Icon.png
```
