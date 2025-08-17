package main

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	cmd    *exec.Cmd
	cancel context.CancelFunc
)

func startCloudflared(hostname, url string, status *widget.Label, btn *widget.Button) {
	if cmd != nil {
		status.SetText("已在运行")
		return
	}

	ctx, c := context.WithCancel(context.Background())
	cancel = c

	cmd = exec.CommandContext(ctx, "cloudflared", "access", "rdp", "--hostname", hostname, "--url", url)

	// 输出日志到文件，方便调试
	logFile, _ := os.Create("cloudflared.log")
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	err := cmd.Start()
	if err != nil {
		status.SetText("启动失败: " + err.Error())
		cmd = nil
		return
	}

	status.SetText("运行中...")
	btn.SetText("停止")

	// 异步等待进程退出
	go func() {
		cmd.Wait()
		cmd = nil
		status.SetText("未运行")
		btn.SetText("启动")
	}()
}

func stopCloudflared(status *widget.Label, btn *widget.Button) {
	if cmd == nil {
		status.SetText("未运行")
		return
	}
	cancel()
	time.Sleep(500 * time.Millisecond)
	cmd = nil
	status.SetText("已停止")
	btn.SetText("启动")
}

func launchRDPConnection(address string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("mstsc", "/v:"+address)
	case "darwin": // macOS
		cmd = exec.Command("open", "rdp://"+address)
	default: // Linux或其他系统
		// 尝试使用xfreerdp或其他RDP客户端
		cmd = exec.Command("xfreerdp", "/v:"+address)
	}

	return cmd.Start()
}

func main() {
	a := app.New()
	w := a.NewWindow("Cloudflared RDP 控制")
	w.Resize(fyne.NewSize(400, 250))

	hostnameEntry := widget.NewEntry()
	hostnameEntry.SetText("rdp.3yfist.com")

	urlEntry := widget.NewEntry()
	urlEntry.SetText("rdp://localhost:3388")

	status := widget.NewLabel("未运行")

	var btn *widget.Button
	btn = widget.NewButton("启动", func() {
		if cmd == nil {
			startCloudflared(hostnameEntry.Text, urlEntry.Text, status, btn)
		} else {
			stopCloudflared(status, btn)
		}
	})

	rdpBtn := widget.NewButton("连接RDP", func() {
		address := "localhost:3388"
		// 从urlEntry中提取地址，如果格式正确的话
		if urlEntry.Text != "" {
			// 简单提取地址部分 (例如从 "rdp://localhost:3388" 提取 "localhost:3388")
			if len(urlEntry.Text) > 6 && urlEntry.Text[:6] == "rdp://" {
				address = urlEntry.Text[6:]
			} else {
				address = urlEntry.Text
			}
		}
		err := launchRDPConnection(address)
		if err != nil {
			status.SetText("RDP启动失败: " + err.Error())
		} else {
			status.SetText("RDP连接已启动")
		}
	})

	form := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Hostname", hostnameEntry),
			widget.NewFormItem("URL", urlEntry),
		),
		btn,
		rdpBtn,
		status,
	)

	w.SetContent(form)

	// 捕获 Ctrl+C 优雅退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if cmd != nil {
			stopCloudflared(status, btn)
		}
		os.Exit(0)
	}()

	w.ShowAndRun()
}
