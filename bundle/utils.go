package bundle

// #include <stdlib.h>
import "C"
import "unsafe"

import (
    "os"
    "strings"
    "io"
    "os/exec"
)

const helperPath = "/usr/libexec/cbundle-launch-helper"

func IsEnvExists(envName string) (ok bool) {
    for _, e := range os.Environ() {
        if strings.HasPrefix(e, envName+"=") {
            ok = true
            break
        }
    }
    return
}

func UnsetEnv(envName string) (err error) {
    doUnsetEnvC(envName)
    envs := os.Environ()
    newEnvsData := make(map[string] string)
    for _, e := range envs {
        a := strings.SplitN(e, "=", 2)
        var name, value string
        if len(a) == 2 {
            name = a[0]
            value = a[1]
        } else {
            name = a[0]
            value = ""
        }
        if name != envName {
            newEnvsData[name] = value
        }
    }
    os.Clearenv()
    for e, v := range newEnvsData {
        err = os.Setenv(e, v)
        if err != nil {
            return
        }
    }
    return

}

func doUnsetEnvC(envName string) {
    cname := C.CString(envName)
    defer C.free(unsafe.Pointer(cname))
    C.unsetenv(cname)
}

//copy a file from srcName to dstName
func CopyFile( srcName string, dstName string) (written int64, err error){
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()
    dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return
    }
    defer dst.Close()
    return io.Copy(dst, src)
}

//Depends xhost, maybe should write by C
func StartChrootedProcess(process string, unionpath string, arguments []string) {
    unionRoot := unionpath
    command :=exec.Command("xhost","local:")
    err:= command.Run()
    if err != nil {
        panic(err)
    }
    arguments_pre := []string{"--chroot", unionRoot, process}
    arguments = append( arguments_pre, arguments...)
    Exec(helperPath, arguments...)
}

func StartNormalProcess(process string, arguments []string) {
    execute := []string{process}
    arguments = append(execute, arguments...)
    Exec(helperPath, arguments...)
}

func Exec(name string, args ...string){
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()
}
