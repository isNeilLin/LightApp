# Light-App

使用HTML，CSS，JS开发轻量级跨平台桌面APP开发模式（基于[webview](https://github.com/zserge/webview)）。

谈到使用HTML/CSS/JS开发跨平台的桌面应用，不得不提起[Electron](https://electronjs.org), 但是由于Electron内置了Node.js和Chromium，所以Electron开发的程序即使只有一个很简单的页面，体积也非常大，普遍100M以上。

前段时间偶然发现了一个跨平台的[webview库](https://github.com/zserge/webview)，觉得可以做些文章。就使用这个[webview库](https://github.com/zserge/webview)的golang的[Binding API](https://godoc.org/github.com/zserge/webview)，在JS中注入一个Bridge,提供了一些方法。当然和Electron的完整性不能相提并论，不过开发一些内部使用的简单桌面APP是可以胜任的。不过打包之后大小基本和静态资源大小持平，以示例Demo为例，打包之后只有十几M的大小。


### 环境配置

由于使用了Go语言的包，所以需要安装Go环境。[Go安装](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/01.1.md)以及[GOPATH设置](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/01.2.md)

项目需要在GOPATH对应的`src`路径下(例如：GOPATH为:`/Users/username/Documents/go`, 则项目对应路径应为:`/Users/username/Documents/go/src/light-app`)


安装GO依赖： `go get`

跨平台打包需要安装Docker。[docker安装](https://www.docker.com/products/docker-desktop)

### 起步

项目目录及文件：

```
- main.go // 主文件
- src
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
  "output_path": "./build"  // 输出路径
  "buildTarget": [
    {
      "os": "darwin",       // 打包平台
      "arch": "amd64"       // 平台架构
    },
    {
      "os": "windows",      // 打包平台
      "arch": "amd64"       // 平台架构
    }
  ]
}
```

> JS中可通过`window.Config.data`访问config.json中的属性


**Bridge接口**

- 属性
    
    `data.os`： 系统平台，Mac：`darwin`， Windows： `windows`， Linux: `linux`
    
    `data.arch`: 系统架构
    
    `data.username`： 用户登录名
    
    `data.storagePath`：系统分配的程序存储路径
    
    `data.homePath`: home路径
    
    `data.tempPath`： 临时路径
    
    `data.currentPath`： 程序当前执行路径

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
        接收返回值： `eventListener.on('getEnv', (value)=>{})`
        
    `Bridge.makeDir(path)`: 创建文件夹
        接收返回值： `eventListener.on('makeDir', (err)=>{})`
        
    `Bridge.remove(path)`: 删除文件
        接收返回值： `eventListener.on('remove', (err)=>{})`
        
    `Bridge.removeAll(path)`: 删除文件夹
        接收返回值： `eventListener.on('removeAll', (err)=>{})`
        
    `Bridge.renameFile(oldpath, newpath)`: 重命名文件
        接收返回值： `eventListener.on('renameFile', (err)=>{})`
        
    `Bridge.openFile(dialogTitle)`: 打开文件
        接收返回值： `eventListener.on('openFile', (jsonString)=>{})`

    `Bridge.openDir(dialogTitle)`: 打开文件夹
        接收返回值： `eventListener.on('openDir', (jsonString)=>{})` 
        
    `Bridge.readFile(path)`： 读取文件
        接收返回值： `eventListener.on('readFile', (fileContent)=>{})` 
        
    `Bridge.writeFile(content)`： 创建文件
        接收返回值： `eventListener.on('writeFile', (err)=>{})` 
    
JS中，通过`eventListener.on(方法名, (返回值)=>{})`接受Bridge有返回值方法的返回值。  

### 打包

1、 `go-bindata -o=asset/asset.go -pkg=asset src/...`

2、 `node build.js`
