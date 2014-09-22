cmake_minimum_required(VERSION 2.8)

message("Building golang binary")
execute_process(
    COMMAND go build bundler-helper.go
    WORKING_DIRECTORY ${PROJECT_SOURCE_DIR}
)

add_subdirectory(launchhelper)
add_subdirectory(auth)


install(FILES ${PROJECT_SOURCE_DIR}/bundler-helper
    DESTINATION ${CMAKE_INSTALL_PREFIX}/bin/
    PERMISSIONS OWNER_READ OWNER_EXECUTE
                GROUP_READ GROUP_EXECUTE
                WORLD_READ WORLD_EXECUTE
                )

install(FILES ${PROJECT_SOURCE_DIR}/data/bundlelauncher.desktop
    DESTINATION ${CMAKE_INSTALL_PREFIX}/share/applications/
    PERMISSIONS OWNER_READ OWNER_WRITE
                GROUP_READ
                WORLD_READ
    )

install(FILES ${PROJECT_SOURCE_DIR}/data/application-x-cb-bundle.svg
    DESTINATION ${CMAKE_INSTALL_PREFIX}/share/icons/hicolor/scalable/mimetypes/
    PERMISSIONS OWNER_READ OWNER_WRITE
                GROUP_READ
                WORLD_READ
    )

install(FILES ${PROJECT_SOURCE_DIR}/data/x-cb-bundle.xml
    DESTINATION ${CMAKE_INSTALL_PREFIX}/share/mime/packages/
    PERMISSIONS OWNER_READ OWNER_WRITE
                GROUP_READ
                WORLD_READ
    )

install(FILES ${PROJECT_SOURCE_DIR}/data/bundle-loop-devices.conf
    DESTINATION /etc/modprobe.d
    PERMISSIONS OWNER_READ OWNER_WRITE
                GROUP_READ
                WORLD_READ
    )
