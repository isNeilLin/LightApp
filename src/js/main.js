window.onload = function () {
    Bridge.init();
    registerClickEvents()
    registerEventEmitters();
}

function registerClickEvents() {
    os.addEventListener('click', function () {
        Bridge.message('os', Bridge.data.os)
    })
    arch.addEventListener('click', function () {
        Bridge.message('arch', Bridge.data.arch)
    })
    hostname.addEventListener('click', function () {
        Bridge.message('hostname', Bridge.data.hostname)
    })
    tempPath.addEventListener('click', function () {
        Bridge.message('tempPath', Bridge.data.tempPath)
    })
    currentPath.addEventListener('click', function () {
        Bridge.message('currentPath', Bridge.data.currentPath)
    })
    message.addEventListener('click', function () {
        Bridge.message('message', 'This is message content')
    })
    info.addEventListener('click', function () {
        Bridge.info('info', 'This is info content')
    })
    warn.addEventListener('click', function () {
        Bridge.info('warn', 'This is warn content')
    })
    error.addEventListener('click', function () {
        Bridge.info('error', 'This is error content')
    })
    setEnv.addEventListener('click', function () {
        let key = envKey.value;
        let value = envValue.value;
        Bridge.setEnv(key, value)
    })
    getEnv.addEventListener('click', function () {
        let key = getEnvKey.value;
        Bridge.getEnv(key)
    })
    setColor.addEventListener('click', function () {
        Bridge.setColor(244,123,54,23)
    })
    fullScreen.addEventListener('click', function () {
        Bridge.setFullScreen(true)
    })
    unFullScreen.addEventListener('click', function () {
        Bridge.setFullScreen(false)
    })
    makeDir.addEventListener('click', function () {
        let path = dirPath.value;
        Bridge.makeDir(path)
    })
    remove.addEventListener('click', function () {
        let path = filePath.value;
        Bridge.remove(path)
    })
    removeAll.addEventListener('click', function () {
        let path = rmDirPath.value;
        Bridge.removeAll(path)
    })
    renameFile.addEventListener('click', function () {
        var oldpath = oldpath.value;
        var newpath = newpath.value;
        Bridge.renameFile(oldpath, newpath)
    })
    openFile.addEventListener('click', function () {
        Bridge.openFile('Choose File')
    })
    openDir.addEventListener('click', function () {
        Bridge.openDir('Choose Dir')
    })
    readFile.addEventListener('click', function () {
        let path = readpath.value;
        Bridge.readFile(path)
    })
    writeFile.addEventListener('click', function () {
        Bridge.writeFile('write from lightApp')
    })
}

function registerEventEmitters() {
    eventListener.on('getEnv', function(value){
        Bridge.message('from getEnv', value)
    })
    eventListener.on('openFile', function(content){
        Bridge.message('from openFile', content)
    })
    eventListener.on('openDir', function(content){
        Bridge.message('from openDir', content)
    })
    eventListener.on('readFile', function(content){
        fileContent.innerText = content;
    })
}