# Light-App

使用HTML，CSS，JS开发轻量级跨平台桌面APP开发模式（基于[webview](https://github.com/zserge/webview)）。


### 环境配置

由于使用了Go语言的包，所以需要安装Go环境。[Go安装](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/01.1.md)以及[GOPATH设置](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/01.2.md)

项目需要在GOPATH对应的`src`路径下(例如：GOPATH为:`/Users/username/Documents/go`, 则项目对应路径应为:`/Users/username/Documents/go/src/light-app`)


安装GO依赖： `go get`

跨平台打包需要安装Docker。[docker安装](https://www.docker.com/products/docker-desktop)

### 起步

项目目录及文件：

```
- main.go // 主文件
- assets
    - config.json  // 项目配置文件
    - index.html   // 入口HTML
    - js/css/img   // 相关静态资源
- asset
    - asset.go     // go-bindata之后生成的静态资源依赖    
```

**`config.json`配置**

```
{
  "name": "Light",          // App名称
  "width": 1080,            // 窗口宽度
  "height": 740,            // 窗口高度
  "title": "Test App",      // 窗口标题
  "resizeable": true,       // 是否可缩放
  "debug": true,            // debug模式
  "icon": "Light.icns",     // App图标路径
  "output_path": "/Users/neil/Desktop/"  // 输出路径
}
```

**Bridge接口**

- 属性
    
    `os`： 系统平台，Mac：`darwin`， Windows： `windows`， Linux: `linux`
    
    `arch`: 系统架构
    
    `hostname`： 系统用户
    
    `tempPath`： 临时路径
    
    `currentPath`： 程序当前执行路径

- 方法

    *没有返回值*
    
    `Bridge.init()`： 初始化Bridge，调用这个方法之后，Bridge的属性生效。
    
    `Bridge.exit()`:  退出程序
    
    `Bridge.message(title, content)`： 消息弹窗
    
    `Bridge.info(title, content)`： 信息弹窗
    
    `Bridge.warn(title, content)`： 警告弹窗
    
    `Bridge.error(title, content)`： 错误弹窗
    
    `Bridge.setEnv(key, value)`: 设置环境变量
    
    `Bridge.setColor(r, g, b, a)`: 设置标题栏颜色
    
    `Bridge.setFullScreen(bool)`: 设置是否全屏
    
    *有返回值*
    
    `Bridge.getEnv(key)`: 获取环境变量值
        接收返回值： `eventEmitter.on('getEnv', (value)=>{})`
        
    `Bridge.makeDir(path)`: 创建文件夹
        接收返回值： `eventEmitter.on('makeDir', (err)=>{})`
        
    `Bridge.remove(path)`: 删除文件
        接收返回值： `eventEmitter.on('remove', (err)=>{})`
        
    `Bridge.removeAll(path)`: 删除文件夹
        接收返回值： `eventEmitter.on('removeAll', (err)=>{})`
        
    `Bridge.renameFile(oldpath, newpath)`: 重命名文件
        接收返回值： `eventEmitter.on('renameFile', (err)=>{})`
        
    `Bridge.openFile(path)`: 打开文件
        接收返回值： `eventEmitter.on('openFile', (jsonString)=>{})`

    `Bridge.openDir(path)`: 打开文件夹
        接收返回值： `eventEmitter.on('openDir', (jsonString)=>{})` 
        
    `Bridge.readFile(path)`： 读取文件
        接收返回值： `eventEmitter.on('readFile', (fileContent)=>{})` 
        
    `Bridge.writeFile(content)`： 创建文件
        接收返回值： `eventEmitter.on('writeFile', (err)=>{})` 
    
JS中，通过`eventEmitter.on(方法名, (返回值)=>{})`接受Bridge有返回值方法的返回值。  

### 打包

1、 `go-bindata -o=asset/asset.go -pkg=asset assets/...`

2、 `xgo --targets=darwin-10.10/.,windows/. .`