import os
import config
import information

#程序运行
def steam():
    if len(os.popen('tasklist | findstr steam.exe').readlines()) == 1:
        os.system('taskkill /f /t /im steam.exe')
        os.startfile("d:/software/steam/steam.exe")
    else:
        os.startfile("d:/software/steam/steam.exe")


# 询问要修改的注册表
def reviseRegistry ():
    while True:
        print(f'\n当前登录账号：{config.currentAccountInformation["AutoLoginUser"]} 输入0返回')
        print(f'可选账号：{config.allUserAccount}')
        tmp = input("请输入数字：")
        try:
            x = eval(tmp)
            if type(x)==int:
                if x == 0:
                    return
                else:
                    if x <= len(config.allUserAccount):
                        information.modifRregistry(x-1)
                        print(f'修改成功，当前账号{config.allUserAccount[x-1]}')
                        steam()
                    else:
                        print("ERROE：超出范围")
        except:
           print("ERROE：请输入数字")

# showAllUserSteamID 显示所有用户SteamID
def showAllUserSteamID():
    
    for i in config.allUserAccountInformation:
        print(f'{i}的SteamID：{config.allUserAccountInformation[i]["SteamID"]}')
        
# 查询单个用户SteamID
def queryUserSteamID():
    while True:
        print(f'\n可选账号：{config.allUserAccount},输入0返回')
        tmp = input("请输入数字：")
        try:
            x = eval(tmp)
            if type(x)==int:
                if x == 0:
                    return
                else:
                    if x <= len(config.allUserAccount):
                        print(f'{config.allUserAccount[x-1]}的 SteamID 为：{config.allUserAccountInformation[config.allUserAccount[x-1]]["SteamID"]}')
                    else:
                        print("ERROE：超出范围")
        except:
           print("ERROE：请输入数字")

def function(num):   
    config.switch.get(num)()