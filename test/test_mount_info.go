package main

import (
    "fmt"
    "./filehelper"
)

func main(){
    volumes,_ := filehelper.GetMountedVolumes()
    for _, volume := range volumes {
        fmt.Println(volume)
    }
    t,e := filehelper.MountInfoForPath("/proc")
    fmt.Println(t)
    fmt.Println(e)
}
