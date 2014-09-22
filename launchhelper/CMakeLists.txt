include_directories(${CMAKE_CURRENT_SOURCE_DIR})

add_definitions(-D_GNU_SOURCE)

add_executable(cbundle-launch-helper
    LaunchHelper.c
    Environment.c
    FilesystemUtilities.c)

install(FILES ${CMAKE_CURRENT_BINARY_DIR}/cbundle-launch-helper
        DESTINATION lib/libexec/
        PERMISSIONS OWNER_READ OWNER_EXECUTE
                    GROUP_READ GROUP_EXECUTE
                    WORLD_READ WORLD_EXECUTE
                    SETUID)
