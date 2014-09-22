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
#include <Environment.h>
#include <Types.h>
#include <sys/types.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

static
BOOL drop_privileges()
{
    /*
     * This will set both the real and effective GIDs back to the caller user's one,
     * thus dropping the effective GID of 0 which belongs to the root user.
     */
    if (setregid(getgid(), getgid()) != 0) {
        perror("setregid");
        return FALSE;
    }

    /*
     * This will set both the real and effective UIDs back to the caller user's one,
     * thus dropping the effective UID of 0 which belongs to the root user.
     */
    if (setreuid(getuid(), getuid()) != 0) {
        perror("setreuid");
        return FALSE;
    }

    return TRUE;
}

static
void parse_arguments(int argc, char* argv[], const char** root, char* const** cmdline)
{
    if (argc < 2) {
        fprintf(stderr, " [**] Not enough arguments.\n");
        exit(EXIT_FAILURE);
    }

    (*root) = NULL;
    (*cmdline) = &argv[1];

    if (strcmp(argv[1], "--chroot") == 0) {
        if (argc < 3) {
            fprintf(stderr, " [**] Not enough arguments.\n");
            exit(EXIT_FAILURE);
        }

        (*root) = argv[2];
        (*cmdline) = &argv[3];
    }
}

/**
 * For the bundle launching to work properly, the required special filesystems
 * (/dev, /sys, /proc, ...) must be mounted before invoking this helper launcher.
 *
 * The usage of this helper is:
 *
 *      cbundle-launch-helper [ --chroot <new-root> ] <executable> [ <arg1> ... <argN> ]
 *
 * where "<new-root>" is the directory where to chroot, if needed, "<executable>"
 * is the path to the executable file to launch (relative, if applicable, to "<new-root>")
 * and the following, "<arg1>" through "<argN>", are optional arguments to pass on the
 * newly-launched executable command line.
 */
int main(int argc, char* argv[])
{
    const char* root = NULL;
    char* const* cmdline = NULL;

    parse_arguments(argc, argv, &root, &cmdline);

    if (root != NULL) {
        printf(" [::] Chrooting into \"%s\"...\n", root);

        if (!prepare_chroot(root)) {
            return EXIT_FAILURE;
        }
    }

    if (!drop_privileges()) {
        fprintf(stderr, " [**] Unable to drop root privileges!\n");
        return EXIT_FAILURE;
    }

    printf(" [::] Dropped privileges to those of UID %u and GID %u.\n", geteuid(), getegid());

    if (!process_environment()) {
        fprintf(stderr, " [**] Unable to process the environmental variables!\n");
        return EXIT_FAILURE;
    }

    printf(" [::] Launching executable \"%s\"...\n", cmdline[0]);
    return execv(cmdline[0], &cmdline[0]);
}

