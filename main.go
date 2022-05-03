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

type AutoLoginUserMaterial struct {
	SteamPath string
	value     string
	name      string
}
type ActiveUserMaterial struct {
	SteamPath string
	value     int
	name      string
}

// SteamAllUser 存放用户数据
var SteamAllUser = make(map[string]map[string]int)
var Steam = map[string]string{"AutoLoginUser": "AutoLoginUser", "SteamExe": "SteamExe", "SteamPath": "SteamPath"}

var AutoLoginUser AutoLoginUserMaterial
var ActiveUser ActiveUserMaterial
var username []string

// RegistryQuery 查询注册表
func RegistryQuery(Steam *map[string]string, path string) {
	k, err := registry.OpenKey(registry.CURRENT_USER, path, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer func(k registry.Key) {
		_ = k.Close()
	}(k)
	switch path {
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
		for y := range SteamAllUser {
			SteamAllUser[y]["ActiveUser"], _ = strconv.Atoi(keys[x])
			x++
		}
	}
}

// Document 读取文件，获取数据
func Document(path string, information *map[string]map[string]int) {
	content, _ := ioutil.ReadFile(path)
	regex := regexp.MustCompile(`[^\d\w]`)
	if regex == nil { //解释失败，返回nil
		fmt.Println("regexp err")
	}
	documentContent := regex.ReplaceAllString(string(content), "")
	subscript := strings.Index(documentContent, "Accounts") + 8
	index := 0
	for i := subscript; i < len(documentContent); i++ {
		if documentContent[i] == 'S' {
			if documentContent[i:i+7] == "SteamID" {
				index = i + 24
			}
		}
	}
	data := strings.ReplaceAll(documentContent[subscript:index], "SteamID", ":")
	for i := 0; i < len(data); i++ {
		SteamUser := make(map[string]int)
		if data == "" {
			break
		}
		if data[i] == ':' {
			SteamUser["SteamID"], _ = strconv.Atoi(data[i+1 : i+18])
			username = append(username, data[0:i])
			(*information)[data[0:i]] = SteamUser
			data = data[i+18:]
			i = 0
		}
	}
}

//获取当前选择账号的信息
func information(informationLabel *widget.Label, value *string) {
	informationLabel.SetText(*value + "\nSteamID:\n" + strconv.Itoa(SteamAllUser[*value]["SteamID"]) + "\nActiveUse:\n" + strconv.Itoa(SteamAllUser[*value]["ActiveUser"]))
}

// SteamTasks 关闭Steam
func SteamTasks(strkey string, strExeName string) bool {

	buf := bytes.Buffer{}
	cmd := exec.Command("wmic", "process", "get", "name,executablepath")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return false
	}

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

// ReviseRegistry 修改注册表
func (material AutoLoginUserMaterial) ReviseRegistry() {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, material.SteamPath, registry.ALL_ACCESS)
	err = key.SetStringValue(material.name, material.value)

	if err != nil {
		log.Fatal(err)
	}
}
func (material ActiveUserMaterial) ReviseRegistry() {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, material.SteamPath, registry.ALL_ACCESS)
	err = key.SetQWordValue(material.name, uint64(material.value))

	if err != nil {
		log.Fatal(err)
	}
}

// RunSteam 运行Steam
func RunSteam(commandLine string) {
	s := strings.Split(commandLine, ",")
	var err error
	var sI syscall.StartupInfo
	var pI syscall.ProcessInformation
	var pS syscall.SecurityAttributes
	if strings.EqualFold(s[len(s)-1], "msi") {
		commandLine = "msiexec /a" + commandLine
	}
	//argv := syscall.StringToUTF16Ptr(commandLine)
	argv, _ := syscall.UTF16PtrFromString(commandLine)
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

	AutoLoginUser = AutoLoginUserMaterial{
		name:      "AutoLoginUser",
		value:     userlist.Selected,
		SteamPath: "Software\\Valve\\Steam",
	}
	ActiveUser = ActiveUserMaterial{
		name:      "ActiveUser",
		value:     SteamAllUser[userlist.Selected]["ActiveUser"],
		SteamPath: "Software\\Valve\\Steam\\ActiveProcess",
	}
	go AutoLoginUser.ReviseRegistry()
	go ActiveUser.ReviseRegistry()
	time.Sleep(10000)
	RunSteam(Steam["SteamExe"])
}

func main() {
	path := []string{"Software\\Valve\\Steam", "Software\\Valve\\Steam\\Users"}
	//查询注册表
	RegistryQuery(&Steam, path[0])
	go RegistryQuery(&Steam, path[1])

	MainWindow := app.New().NewWindow("Steam")
	MainWindow.Resize(fyne.Size{Width: 250, Height: 300})
	MainWindow.SetFixedSize(true)

	informationLabel := widget.NewLabel(Steam["AutoLoginUser"] + "\nSteamID:\n" + strconv.Itoa(SteamAllUser[Steam["AutoLoginUser"]]["SteamID"]) + "\nActiveUse:\n" + strconv.Itoa(SteamAllUser[Steam["AutoLoginUser"]]["ActiveUser"]))
	informationLabel.Resize(fyne.Size{Width: 200, Height: 50})

	userlist := widget.NewSelect(username, func(value string) {
		information(informationLabel, &value)
	})

	button := widget.NewButton("Run Steam", func() {
		if SteamTasks("steam.exe", "steam.exe") {
			c := exec.Command("taskkill.exe", "/f", "/im", "steam.exe")
			_ = c.Start()
		} else {
			implement(userlist)
		}
	})
	content := container.New(layout.NewBorderLayout(userlist, button, informationLabel, nil),
		userlist, button, informationLabel)

	MainWindow.SetContent(content)
	MainWindow.ShowAndRun()
}
