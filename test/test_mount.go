package main

import (
    "fmt"
    "./filehelper"
)

func main(){
    //volume,err := filehelper.Mount("/dev/sda1", "/tmp/p", "ext2", 0)
    deviceString := map [string]int {"/tmp/root":0,"/tmp/home":1}
    volume , err := filehelper.Unionfs(deviceString, "/tmp/real-root")
    fmt.Println(volume)
    fmt.Println(err)
}
