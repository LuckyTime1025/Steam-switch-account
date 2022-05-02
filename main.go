package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type modifyRegistry interface {
	ReviseRegistry(registrykey string, registryvalue string, SteamPath string)
}

type Material struct {
	SteamPath string
	value     string
	name      string
}

//存放用户数据
var SteamAllUser = make(map[string]map[string]int)
var Steam = map[string]string{"AutoLoginUser": "AutoLoginUser", "SteamExe": "SteamExe", "SteamPath": "SteamPath"}

var AutoLoginUser modifyRegistry
var ActiveUser modifyRegistry

//查询注册表
func RegistryQuery(Steam *map[string]string) {
	path := []string{"Software\\Valve\\Steam", "Software\\Valve\\Steam\\Users"}
	for _, i := range path {
		k, err := registry.OpenKey(registry.CURRENT_USER, i, registry.ALL_ACCESS)
		if err != nil {
			log.Fatal(err)
		}
		defer k.Close()
		switch i {
		case "Software\\Valve\\Steam":
			for _, x := range *Steam {
				s, _, err := k.GetStringValue(x)
				if err != nil {
					log.Fatal(err)
				}
				(*Steam)[x] = s
			}
		case "Software\\Valve\\Steam\\Users":
			//读取文件，获取数据
			Document((*Steam)["SteamPath"]+"\\config\\config.vdf", &SteamAllUser)
			keys, _ := k.ReadSubKeyNames(0)
			x := 0
			for y, _ := range SteamAllUser {
				SteamAllUser[y]["ActiveUser"], _ = strconv.Atoi(keys[x])
				x++
			}
		}
	}
}

//读取文件，获取数据
func Document(path string, information *map[string]map[string]int) {
	content, _ := ioutil.ReadFile(path)
	regex := regexp.MustCompile(`[^\d\w]`)
	if regex == nil { //解释失败，返回nil
		fmt.Println("regexp err")
	}
	document_content := regex.ReplaceAllString(string(content), "")
	subscript := strings.Index(document_content, "Accounts") + 8
	index := 0
	for i := subscript; i < len(document_content); i++ {
		if document_content[i] == 'S' {
			if document_content[i:i+7] == "SteamID" {
				index = i + 24
			}
		}
	}
	data := strings.ReplaceAll(document_content[subscript:index], "SteamID", ":")
	for i := 0; i < len(data); i++ {
		SteamUser := make(map[string]int)
		if data == "" {
			break
		}
		if data[i] == ':' {
			SteamUser["SteamID"], _ = strconv.Atoi(data[i+1 : i+18])
			(*information)[data[0:i]] = SteamUser
			data = data[i+18:]
			i = 0
		}
	}
}

//获取当前选择账号的信息
func information(information_label *widget.Label, value *string) {
	//fmt.Println(SteamAllUser[*value]["ActiveUser"])
	information_label.SetText(*value + "\nSteamID:\n" + strconv.Itoa(SteamAllUser[*value]["SteamID"]) + "\nActiveUse:\n" + strconv.Itoa(SteamAllUser[*value]["ActiveUser"]))
}

//关闭Steam
func SteamTasks(strkey string, strExeName string) bool {

	buf := bytes.Buffer{}
	cmd := exec.Command("wmic", "process", "get", "name,executablepath")
	cmd.Stdout = &buf
	cmd.Run()

	cmd2 := exec.Command("findstr", strkey)
	cmd2.Stdin = &buf
	data, err := cmd2.CombinedOutput()
	if err != nil && err.Error() != "exit status 1" {
		return false
	}

	strData := string(data)
	if strings.Contains(strData, strExeName) {
		return true
	} else {
		return false
	}
}

//修改注册表
func (material Material) ReviseRegistry(registrykey string, registryvalue string, SteamPath string) {
	material.name = registrykey
	value, erro := strconv.Atoi(registryvalue)
	material.SteamPath = SteamPath

	key, _, err := registry.CreateKey(registry.CURRENT_USER, material.SteamPath, registry.ALL_ACCESS)
	if erro != nil {
		err = key.SetStringValue(material.name, registryvalue)
	} else {
		err = key.SetQWordValue(material.name, uint64(value))
	}

	if err != nil {
		log.Fatal(err)
	}
}

//运行Steam
func RunSteam(commandLine string) {
	s := strings.Split(commandLine, ",")
	var err error
	var sI syscall.StartupInfo
	var pI syscall.ProcessInformation
	var pS syscall.SecurityAttributes
	if strings.EqualFold(s[len(s)-1], "msi") {
		commandLine = "msiexec /a" + commandLine
	}
	argv := syscall.StringToUTF16Ptr(commandLine)
	err = syscall.CreateProcess(
		nil,
		argv,
		&pS,
		nil,
		true,
		0,
		nil,
		nil,
		&sI,
		&pI)
	if err != nil {
		log.Fatalf("CreateProcess(%s) failed with %s\n", commandLine, err)
	}
}

//执行
func implement(userlist *widget.Select) {

	AutoLoginUser = new(Material)
	ActiveUser = new(Material)

	go AutoLoginUser.ReviseRegistry("AutoLoginUser", userlist.Selected, "Software\\Valve\\Steam")
	go ActiveUser.ReviseRegistry("ActiveUser", strconv.Itoa(SteamAllUser[userlist.Selected]["ActiveUser"]), "Software\\Valve\\Steam\\ActiveProcess")
	time.Sleep(10000)
	RunSteam(Steam["SteamExe"])
}

func main() {
	//查询注册表
	RegistryQuery(&Steam)

	var username []string
	for i, _ := range SteamAllUser {
		username = append(username, i)
	}

	MainWindow := app.New().NewWindow("Steam")
	MainWindow.Resize(fyne.Size{250, 300})
	MainWindow.SetFixedSize(true)

	information_label := widget.NewLabel(Steam["AutoLoginUser"] + "\nSteamID:\n" + strconv.Itoa(SteamAllUser[Steam["AutoLoginUser"]]["SteamID"]) + "\nActiveUse:\n" + strconv.Itoa(SteamAllUser[Steam["AutoLoginUser"]]["ActiveUser"]))
	information_label.Resize(fyne.Size{200, 50})

	userlist := widget.NewSelect(username, func(value string) {
		information(information_label, &value)
	})

	button := widget.NewButton("Run Steam", func() {
		if SteamTasks("steam.exe", "steam.exe") {
			c := exec.Command("taskkill.exe", "/f", "/im", "steam.exe")
			c.Start()
		} else {
			implement(userlist)
		}
	})
	content := container.New(layout.NewBorderLayout(userlist, button, information_label, nil),
		userlist, button, information_label)

	MainWindow.SetContent(content)
	MainWindow.ShowAndRun()
}
