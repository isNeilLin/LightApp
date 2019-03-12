const config = require("./src/config.json");
const path = require("path");
const fs = require("fs");
const { spawn } = require("child_process");

function getTemplate(name){
    return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
<key>CFBundleExecutable</key>
<string>${name}</string>
<key>CFBundleIconFile</key>
<string>${config.icon}</string>
<key>CFBundleIdentifier</key>
<string>com.example.yours</string>
<key>NSHighResolutionCapable</key>
<true/>
<key>LSUIElement</key>
<true/>
</dict>
</plist>`
};

function main() {
    let targets = config.buildTarget.map(target=>{
        let os = target.os === 'darwin' ? 'darwin-10.10' : target.os;
        return `${os}/${target.arch}`
    })
    targets = `--targets=${targets.join(',')}`
    compile = spawn("xgo", ["-out", `${config.output_path}/${config.name}`, targets, "."])
    compile.on('error', function (err) {
        console.log(err)
    })
    compile.on('close',function () {
        createMacApp()
    })
}

function createMacApp() {
    let outputPath = path.join(__dirname, config.output_path);
    fs.readdir(outputPath, (err, files)=>{
        if (err) {
            throw Error(err)
        }
        files.filter(file=>file.indexOf('darwin') !== -1).map(file=>{
            fs.mkdirSync(path.join(outputPath, `${file}.app`), {recursive: true})
            fs.mkdirSync(path.join(outputPath, `${file}.app`, "Contents"), {recursive: true})
            fs.mkdirSync(path.join(outputPath, `${file}.app`, "Contents", "MacOS"), {recursive: true})
            fs.mkdirSync(path.join(outputPath, `${file}.app`, "Contents", "Resources"), {recursive: true})
            fs.renameSync(path.join(__dirname, config.icon), path.join(outputPath, `${file}.app`, "Contents", "Resources", config.icon))
            fs.renameSync(path.join(outputPath, file), path.join(outputPath, `${file}.app`, "Contents", "MacOS", file))
            fs.writeFileSync(path.join(outputPath, `${file}.app`, "Contents", "Info.plist"), getTemplate(file), "utf8")
        })
    });
}

main()