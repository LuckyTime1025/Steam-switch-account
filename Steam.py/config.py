import sys
from traceback import print_last
import registry
import winreg
import information

# registryPath 注册表路径
registryPath = "Software\Valve\Steam"
# registryKey 要查询的注册表项
registryKey = ['AutoLoginUser','SteamPath','SteamExe']

# 所有已登录的账户
allUserAccount = []
# 当前登录的用户账户
currentAccountInformation = dict()
# 已登录账户 SteamID
allUserAccountInformation = dict()

k = winreg.OpenKey(winreg.HKEY_CURRENT_USER, registryPath, reserved=0, access=winreg.KEY_READ|winreg.KEY_WRITE)

switch = {
    #查看所有用户SteamID
    0:registry.showAllUserSteamID,
    #查询单个用户SteamID
    1:registry.queryUserSteamID,
    #修改注册表
    2:registry.reviseRegistry
}

def main():
    information.registryQuery()
    information.readFile()
    while True:
        print("\n1.查看所有用户SteamID 2.查询单个用户SteamID 3.快速切换账号 0.退出")
        tmp = input("请输入数字：")
        try:
            x = eval(tmp)
            if type(x)==int:
                if x == 0:
                    print("Bye")
                    break
                else:
                    if x <= len(switch):
                       registry.function(x-1)
                    else:
                        print("ERROE：超出范围")
        except:
           print("ERROE：请输入数字")
    sys.exit()