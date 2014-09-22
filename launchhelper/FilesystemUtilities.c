/*
 * cbundle - A shared library to manage Chakra Linux Bundles.
 * Copyright (C) 2010-2011  The Chakra Project team
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */

#include <FilesystemUtilities.h>
#include <sys/types.h>
#include <sys/mount.h>
#include <sys/stat.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>
#include <sched.h>
#include <mntent.h>

static
char* make_target_path(const char* root, const char* source)
{
    size_t n1 = strlen(root);
    size_t n2 = strlen(source);

    char* target = (char*) calloc(n1 + n2 + 1, sizeof(char));
    if (!target) {
        return NULL;
    }

    strcpy(target, root);
    strcpy(target + n1, source);
    return target;
}

static
BOOL check_root(const char* root)
{
    BOOL returncode = path_is_directory(root);

    switch (returncode) {
    case INVALID_BOOL:
        fprintf(stderr, " [**] Unable to get info about the requested path.\n");
        return FALSE;

    case FALSE:
        fprintf(stderr, " [**] The requested path is not a directory.\n");
        return FALSE;

    default:
        break;
    }

    return TRUE;
}

static
BOOL is_mounted(const char* mountpoint)
{
    FILE* mtab = setmntent(_PATH_MOUNTED, "r");
    struct mntent* it = 0;

    if (mtab == 0) {
        return INVALID_BOOL;
    }

    while (it = getmntent(mtab)) {
        if (strcmp(it->mnt_dir, mountpoint) == 0) {
            endmntent(mtab);
            return TRUE;
        }
    }

    endmntent(mtab);
    return FALSE;
}

static
BOOL rebind(const char* path, const char* root)
{
    char* target = make_target_path(root, path);
    int returncode = 0;

    if (!target) {
        fprintf(stderr, " [**] Memory allocation failure.\n");
        return FALSE;
    }

    if (is_mounted(target) == FALSE) {
        returncode = mount(path, target, NULL, MS_BIND, NULL);
        printf(" [::] Bind mount %s -> %s.\n",path, target); 
    } else {
        printf(" [:Ignore:] Bind mounted %s -> %s.\n",path, target); 
    }

    free(target);

    if (returncode != 0) {
        perror("mount");
        return FALSE;
    }

    return TRUE;
}
static
BOOL mtmpfs(const char* path, const char* root)
{
    char* target = make_target_path(root, path);
    int returncode = 0;

    if (!target) {
        fprintf(stderr, " [**] Memory allocation failure.\n");
        return FALSE;
    }

    if (is_mounted(target) == FALSE) {
        returncode = mount(path, target, "tmpfs", 0, NULL);
        printf(" [::] Mount Tmpfs -> %s.\n", target); 
    } else {
        printf(" [:Ignore:] Mounted Tmpfs -> %s.\n", target); 
    }

    free(target);

    if (returncode != 0) {
        perror("mount");
        return FALSE;
    }

    return TRUE;
}

BOOL prepare_chroot(const char* root)
{
    static char* __filesystems[] = {
        "/dev",
        "/dev/shm",
        "/dev/pts",
        "/sys",
        "/run",
        "/var",
        "/proc",
        "/home",
    };

    static char* __tmpfs[] = {
            "/tmp",
    };

    static size_t __filesystems_count = sizeof(__filesystems) / sizeof(__filesystems[0]);

    size_t i;

    if (!check_root(root)) {
        return FALSE;
    }

    /*
     * Separate the execution of the bundled program from the rest of the system.
     */
    if (unshare(CLONE_NEWNS) != 0) {
        perror("unshare");
        return FALSE;
    }

    for (i = 0; i < __filesystems_count; ++i) {
        if (!rebind(__filesystems[i], root)) {
            while (i != 0) {
                char* target = make_target_path(root, __filesystems[i]);
                if (target != 0) {
                    umount(target);
                    free(target);
                }
            }
            return FALSE;
        }
    }

    static size_t __tmpfs_count = sizeof(__tmpfs) / sizeof(__tmpfs[0]);
    for (i = 0; i < __tmpfs_count; ++i) {
        if (!mtmpfs(__tmpfs[i], root)) {
            while (i != 0) {
                char* target = make_target_path(root, __filesystems[i]);
                if (target != 0) {
                    umount(target);
                    free(target);
                }
            }
            return FALSE;
        }
    }

    /*
     * Enter the chroot.
     */
    if (chroot(root) != 0) {
        fprintf(stderr, " [**] Unable to chroot into \"%s\"!\n", root);
        return FALSE;
    }

    /*
     * Move into a safe path.
     */
    if (chdir("/") != 0) {
        fprintf(stderr, " [**] Unable to change working directory!\n");
        return FALSE;
    }

    return TRUE;
}

BOOL path_is_directory(const char* path)
{
    struct stat info;

    if (stat(path, &info) != 0) {
        perror("stat");
        return INVALID_BOOL;
    }

    return S_ISDIR(info.st_mode) ? TRUE : FALSE;
}

BOOL path_is_executable(const char* path)
{
    struct stat info;

    if (stat(path, &info) != 0) {
        perror("stat");
        return INVALID_BOOL;
    }

    return S_ISEXE(info.st_mode) ? TRUE : FALSE;
}

BOOL check_executable(const char* root, const char* executable)
{
    BOOL returncode;
    char* target = make_target_path(root, executable);

    if (!target) {
        fprintf(stderr, " [**] Memory allocation failure.\n");
        return FALSE;
    }

    returncode = path_is_executable(target);
    free(target);

    switch (returncode) {
    case INVALID_BOOL:
        fprintf(stderr, " [**] Unable to get info about the requested path.\n");
        return FALSE;

    case FALSE:
        fprintf(stderr, " [**] The requested path is not an executable.\n");
        return FALSE;

    default:
        break;
    }

    return TRUE;
}
