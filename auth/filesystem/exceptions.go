package main

const (
    // Process run time error code
    FailedtoStart = iota
    Crashed
    Timedout
    WriteError
    ReadError
    UnknowError
)

type ProcessExecutionException struct{
    code int
    msg string
}

const (
    //Bundle runtime error code
    ModuleLoadException = iota
    FileAccessException
    MountPointException
    NoSuchElementException
    NullPointerException
    ProcessException
    UnsupportedOperationsException
)

type ExitErrorException struct{
    code int
    msg string
}

var (
    MsgExitErrorException = "Command %s exited with status %s"
    MsgProcessExecutionException = "Command execution failed with error %s"
    MsgFileAccessException = "Unable to open %1"
    MsgFileAccessExceptionWithExitCode = "Unable to open %s: %s"
    MsgModuleLoadException = "Unable to load a needed module: %s"
    MsgMountException = "Unable to mount %s"
    MsgUnmountException = "Unable to unmount %s"
    MsgNoSuchElementException = "No such element %s"
    MsgNullPointerException = "Null Pointer"
    MsgUnsupportedOperationsException = "Unsupported operation"

    KMod_CreateContext = "Couldn't create a new kmod context"
    KMod_LookupError = "Couldn't lookup the requested module alias"
    KMod_NoSuckModule = "Could't find the requested module"
    KMod_Blacklisted = "The required module is blacklisted"
    KMod_NoError = "No Error"
)
