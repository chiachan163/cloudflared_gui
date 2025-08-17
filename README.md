# CLOUDFLARED_GUI


## mac环境编辑

安装 mingw 工具链（交叉编译 Windows GUI 必须要有）：

```
brew install mingw-w64
```
然后编译时加上 CC 环境变量，让 Go 用 Windows 的 gcc：

```
# 64位
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o releases/cloudflared_gui.exe 
# 32位
GOOS=windows GOARCH=386 CC=x86_64-w64-mingw32-gcc go build -o cloudflared_gui_32.exe cloudflared_gui.go

```
这样会正确选中 Windows 的 OpenGL 代码。

