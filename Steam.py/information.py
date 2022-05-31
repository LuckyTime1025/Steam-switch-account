import re
import winreg
import config

# registryQuery 查询注册表
def registryQuery () :
    # 打开注册表
    for i in config.registryKey:
        # 查询打开的注册表中的值
        config.currentAccountInformation[i] = winreg.QueryValueEx(config.k, i)[0]
    return
    
# modifRregistry修改注册表
def modifRregistry(i):
    winreg.SetValueEx(config.k,config.registryKey[0],0,winreg.REG_SZ,config.allUserAccount[i])

# getValue 获取所有已登录的账号信息
def getValue(information,lineNumber):
    # x, y 定位
    x = lineNumber
    y = 2
    # 循环获取账号信息
    steamId = dict()
    while True:
        # 判断是否获取完毕
        if information[x+y].lstrip().replace("\"", "").replace("\n", "") == "}":
            return
        # 格式化获取 SteamID 的值
        steamId['SteamID']=int(information[x+y+2].replace("\"SteamID\"", "").lstrip().replace("\"", ""))
        # 把 SteamID 的值存储到对应账号中
        userName = information[x+y].lstrip().replace("\"", "").replace("\n", "")
        config.allUserAccountInformation[userName] = steamId
        config.allUserAccount.append(userName)
        x += 2
        y += 2

# getlineNumber 获取行号
def getlineNumber ():
    # 打开文件对象
    file = open("D:\software\Steam\config\config.vdf", mode="r")
    # 获取文件内容
    information = file.read()
    file.close()
    # 获取行号
    # match.start()是 information 中匹配开始的位置
    match = re.search("Accounts", information)
    lineNumber = information.count('\n', 0, match.start())
    return lineNumber

# readFile 读取文件    
def readFile() :
    path = f"{config.currentAccountInformation['SteamPath']}/config/config.vdf"
    # 打开文件对象
    file = open(path, mode="r")
    # 按行读取文件内容
    information = file.readlines()
    # 关闭文件
    file.close()
    # 获取行号
    lineNumber = getlineNumber()
    # 格式化内容
    accouts = information[lineNumber].lstrip().replace("\"", "").replace("\n", "")
    # 判断内容是否正确
    if accouts == "Accounts":
        getValue(information,lineNumber)