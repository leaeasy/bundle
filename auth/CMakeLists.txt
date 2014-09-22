cmake_minimum_required(VERSION 2.8)


macro(buildgo NAME OUTPUT)
    message("Build golang module: " ${NAME})
    execute_process(
        COMMAND  go build -o ${OUTPUT}
        WORKING_DIRECTORY ${PROJECT_SOURCE_DIR}/auth/${NAME}
        )
endmacro()

buildgo(filesystem bundle_filesystem_helper) 

install(FILES ${PROJECT_SOURCE_DIR}/auth/filesystem/bundle_filesystem_helper
    DESTINATION /usr/lib/libexec/
    PERMISSIONS OWNER_READ OWNER_EXECUTE
                GROUP_READ GROUP_EXECUTE
                WORLD_READ WORLD_EXECUTE
    )

install(FILES ${PROJECT_SOURCE_DIR}/auth/conf/com.linuxdeepin.bundle.filesystem.conf
    DESTINATION /etc/dbus-1/system.d/
    PERMISSIONS OWNER_READ OWNER_WRITE
                GROUP_READ 
                WORLD_READ 
    )

install(FILES ${PROJECT_SOURCE_DIR}/auth/conf/com.linuxdeepin.bundle.filesystem.service
    DESTINATION /usr/share/dbus-1/system-services/
    PERMISSIONS OWNER_READ OWNER_WRITE
                GROUP_READ 
                WORLD_READ 
    )
