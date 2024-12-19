package main

import (
	"easyyaml/internal/templates"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TabItem 表示一个标签页
type TabItem struct {
	path    string
	content *widget.Entry
}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("EasyYaml Editor")

	// 创建标签容器
	tabs := container.NewAppTabs()

	// 创建新标签页
	createNewTab := func(name string, content string) {
		editor := widget.NewMultiLineEntry()
		editor.SetText(content)
		tabs.Append(container.NewTabItem(name, editor))
		tabs.Select(tabs.Items[len(tabs.Items)-1])
	}

	// 获取当前编辑器
	getCurrentEditor := func() *widget.Entry {
		if tabs.Selected() == nil {
			return nil
		}
		return tabs.Selected().Content.(*widget.Entry)
	}

	// 创建工具栏
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			// 新建文件
			createNewTab("未命名.yaml", "")
		}),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			// 打开文件对话框
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, window)
					return
				}
				if reader == nil {
					return
				}
				defer reader.Close()

				data, err := ioutil.ReadAll(reader)
				if err != nil {
					dialog.ShowError(err, window)
					return
				}

				filename := filepath.Base(reader.URI().String())
				createNewTab(filename, string(data))
			}, window)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".yaml", ".yml"}))
			fd.Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			editor := getCurrentEditor()
			if editor == nil {
				return
			}

			// 保存文件对话框
			fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, window)
					return
				}
				if writer == nil {
					return
				}
				defer writer.Close()

				// 处理文件扩展名
				path := writer.URI().String()
				ext := filepath.Ext(path)
				if ext == "" {
					path += ".yaml"
				} else if ext != ".yaml" && ext != ".yml" {
					// 如果扩展名不是.yaml或.yml，询问用户
					dialog.ShowConfirm("确认",
						fmt.Sprintf("文件扩展名不是.yaml或.yml，是否添加.yaml后缀？当前：%s", ext),
						func(add bool) {
							if add {
								path += ".yaml"
							}
							saveFile(path, editor.Text, window)
						}, window)
					return
				}

				saveFile(path, editor.Text, window)

				// 更新标签页标题
				if tabs.Selected() != nil {
					tabs.Selected().Text = filepath.Base(path)
					tabs.Refresh()
				}
			}, window)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".yaml", ".yml"}))
			fd.Show()
		}),
	)

	// 创建模板选择下拉框
	templateSelect := widget.NewSelect([]string{
		"Kubernetes配置",
		"Docker Compose",
		"Go服务配置",
	}, func(selected string) {
		var template string
		switch selected {
		case "Kubernetes配置":
			template = templates.GetK8sDeploymentTemplate()
		case "Docker Compose":
			template = templates.GetDockerComposeTemplate()
		case "Go服务配置":
			template = templates.GetGoConfigTemplate()
		}
		createNewTab(strings.ToLower(selected)+".yaml", template)
	})

	// 创建主布局
	content := container.NewBorder(
		container.NewVBox(toolbar, templateSelect), // 顶部
		nil,  // 底部
		nil,  // 左侧
		nil,  // 右侧
		tabs, // 中间
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 600))

	// 创建初始标签页
	createNewTab("未命名.yaml", "")

	window.ShowAndRun()
}

// 保存文件的辅助函数
func saveFile(path, content string, window fyne.Window) {
	uri := storage.NewFileURI(path)
	writer, err := storage.Writer(uri)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	defer writer.Close()

	_, err = writer.Write([]byte(content))
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
}
