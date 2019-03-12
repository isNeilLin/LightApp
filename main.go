package main

import (
  "net"
  "fmt"
  "io/ioutil"
  "net/http"
  "mime"
  "path/filepath"
  "os"
  "github.com/zserge/webview"
  "io"
  "runtime"
  "bytes"
  "unsafe"
  "encoding/json"
  "light-app/asset"
)

var W webview.WebView

func initServer() string {
  ln, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    fmt.Errorf("ServerError: ", err)
  }
  go func() {
    defer ln.Close()
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      path := r.URL.Path
      if (len(path) == 1 && path == "/") {
        w.Header().Add("Content-Type", mime.TypeByExtension("html"))
        data, _ := asset.Asset("assets/index.html")
        io.Copy(w, bytes.NewBuffer(data))
      } else {
        w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
        newpath := string([]byte(path)[1:])
        data, _ := asset.Asset(newpath)
        io.Copy(w, bytes.NewBuffer(data))
      }
    })
    fmt.Print(http.Serve(ln, nil))
  }()
  return "http://"+ln.Addr().String()
}

type Bridge struct {
  Os string          `json:"os"`
  Arch string        `json:"arch"`
  Hostname string    `json:"hostname"`
  TempPath string    `json:"tempPir"`
  CurrentPath string `json:"currentPath"`
}

type Config struct {
  Width int
  Height int
  Title string
  Resizeable bool
  Debug bool
  Name string
  OutputPath string
  Icon string
}

func (bridge *Bridge)  Init(){
  hostname, _ := os.Hostname()
  wd, _ := os.Getwd()
  bridge.Os = runtime.GOOS
  bridge.Arch = runtime.GOARCH
  bridge.Hostname = hostname
  bridge.CurrentPath = wd
  bridge.TempPath = os.TempDir()
}

func (bridge *Bridge)  Exit(){
  W.Terminate()
}

func (bridge *Bridge)  Message(title string, content string) {
  W.Dialog(webview.DialogTypeAlert, 0, title, content)
}

func (bridge *Bridge)  Info(title string, content string){
  W.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, title, content)
}

func (bridge *Bridge)  Warn(title string, content string){
  W.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, title, content)
}

func (bridge *Bridge)  Error(title string, content string){
  W.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, title, content)
}

func (bridge *Bridge)  SetTitle(title string){
  W.SetTitle(title)
}

func (bridge *Bridge)  SetEnv(key string, value string){
  os.Setenv(key, value)
}

func (bridge *Bridge)  GetEnv(key string){
  value := os.Getenv(key)
  ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('getEnv', '%v')", value))
}

func (bridge *Bridge)  SetColor(r uint8, g uint8, b uint8, a uint8){
  W.SetColor(r,g,b,a)
}

func (bridge *Bridge)  SetFullScreen(fullscreen bool){
  W.SetFullscreen(fullscreen)
}

func (bridge *Bridge)  MakeDir(path string){
  err := os.MkdirAll(path, 0777)
  if err != nil {
    ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('makeDir', '%v')", err))
  }
  ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('makeDir')"))
}

func (bridge *Bridge)  Remove(path string){
  err := os.Remove(path)
  if err != nil {
    ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('remove', '%v')", err))
  }
  ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('remove')"))
}

func (bridge *Bridge)  RemoveAll(path string){
  err := os.RemoveAll(path)
  if err != nil {
    ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('removeAll', '%v')", err))
  }
  ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('removeAll')"))
}

func (bridge *Bridge)  RenameFile(oldpath string, newpath string){
  err := os.Rename(oldpath, newpath)
  if err != nil {
    ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('renameFile', '%v')", err))
  }
  ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('renameFile')"))
}

func (bridge *Bridge)  OpenFile(title string){
  path := W.Dialog(webview.DialogTypeOpen, 0, title, "")
  file, _ := os.Open(path)
  stat, _ := file.Stat()
  var json string
  json = fmt.Sprintf(`{
      name: %v,
      size: %v,
      isDir: %v,
      modTime: %v,
      path: %v
    }`, stat.Name(), stat.Size(), stat.IsDir(), stat.ModTime(), path)
  ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('openFile', `%v`)", json))
}

func (bridge *Bridge)  OpenDir(title string) {
  path := W.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, title, "")
  dir, err := os.Open(path)
  files, err := dir.Readdir(-1)
  if err != nil {
    ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('openDir', '%v')", err))
  }
  var json string
  json = "["
  for _, file := range files {
    jsonStr := fmt.Sprintf(`{
      name: %v,
      size: %v,
      isDir: %v,
      modTime: %v
    }`, file.Name(), file.Size(), file.IsDir(), file.ModTime())
    json += jsonStr + ","
  }
  json += "]"
  ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('openDir', `%v`)", json))
}

func (bridge *Bridge)  ReadFile(path string){
  bs, err := ioutil.ReadFile(path)
  if err != nil {
    ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('readFile', '%v')", err))
  }
  str := (*string)(unsafe.Pointer(&bs))
  ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('readFile', `%v`)", *str))
}

func (bridge *Bridge)  WriteFile(content string){
  path := W.Dialog(webview.DialogTypeSave, 0, "选择路径", "")
  err := ioutil.WriteFile(path, []byte(content), 0644)
  if err != nil {
    ExcuteJS(fmt.Sprintf("window.eventEmitter.emit('writeFile', '%v')", err))
  }
}


func ExcuteJS(js string)  {
  W.Dispatch(func() {
    W.Eval(js)
  })
}

func main()  {
  asset.RestoreAssets("./","./assets")
  data, _ := asset.Asset("assets/config.json")
  config := Config{}
  json.Unmarshal(data, &config)
  url := initServer()
  W = webview.New(webview.Settings{
    Width:          config.Width,
    Height:         config.Height,
    Title:          config.Title,
    URL:            url,
    Resizable:      config.Resizeable,
    Debug:          config.Debug,
  })
  defer W.Exit()
  W.Dispatch(func() {
    W.Eval(`window.eventEmitter = {
    emit: function (eventName, args) {
        var callbackArr = events[eventName]
        if ( callbackArr && callbackArr.length) {
            for (var i = 0; i < callbackArr.length; i++) {
                callbackArr[i](args)
            }
        }
    }
};
window.events = {};
window.eventListener = {
    on: function (eventName, callback) {
        if (!events[eventName]) {
            events[eventName] = []
        }
        events[eventName].push(callback)
    },
    remove: function (eventName, callback) {
        var callbackArr = events[eventName]
        if (callbackArr && callbackArr.length) {
            var index = callbackArr.indexOf(callback);
            callbackArr.splice(index, 1);
        }
    }
}`)
  })
  W.Bind("Bridge", &Bridge{})
  W.Run()
}