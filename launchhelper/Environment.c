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

#include <Environment.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>

/**
 * Declare the following external char**, which is an array of char* that will
 * be passed to our executable from the outside.
 */
extern char** environ;

/**
 * Define a list of strings, which will contain the environmental variable names
 * that should be added or removed at the end of process_environment().
 *
 * The "name" field is required, while the "value" field is optional (in case of
 * the "to be removed" list, the variable value is useless).
 */
struct node {
    const char*  name;
    const char*  value;
    struct node* next;
};

typedef struct node* list;

/**
 * Define the environmental variable name prefix and its length (in characters).
 */
#define ENVVAR_PFX     "CBUNDLE_"
#define ENVVAR_PFX_LEN 8

static
struct node* node_new(const char* name, const char* value)
{
    struct node* ptr = (struct node*) malloc(sizeof(struct node));

    if (!ptr) {
        return NULL;
    }

    memset(ptr, 0, sizeof(struct node));
    ptr->name = strdup(name);

    if (!ptr->name) {
        free(ptr);
        return NULL;
    }

    if (value) {
        ptr->value = strdup(value);

        if (!ptr->value) {
            free((void*) ptr->name);
            free(ptr);
            return NULL;
        }
    }

    return ptr;
}

static
void node_delete(struct node* ptr)
{
    free((void*) ptr->value);
    free((void*) ptr->name);
    free(ptr);
}

static
BOOL list_push_front(list* ptr, const char* name, const char* value)
{
    struct node* head = node_new(name, value);

    if (!head) {
        return FALSE;
    }

    head->next = (*ptr);
    (*ptr) = head;

    return TRUE;
}

static
BOOL rename_cbundle_variable(const char* envvar, list* to_add, list* to_remove)
{
    char* copy = strdup(envvar);
    const char* name;
    const char* value;

    if (!copy) {
        return FALSE;
    }

    name = strtok(copy, "=");
    value = strtok(NULL, "\0");

    /**
     * Skip the variable name prefix when adding the variable to the "to_add" list.
     */
    list_push_front(to_add, &name[ENVVAR_PFX_LEN], value);
    list_push_front(to_remove, name, NULL);
    free(copy);

    return TRUE;
}

BOOL process_environment()
{
    list to_add = NULL;
    list to_remove = NULL;
    char** it;

    /**
     * First, scan all the environmental variables and take note of the
     * variables that are to be added or removed from the environment.
     */
    for (it = environ; (*it) != NULL; ++it) {
        if (strncmp(*it, ENVVAR_PFX, ENVVAR_PFX_LEN) != 0) {
            continue;
        }

        if (!rename_cbundle_variable(*it, &to_add, &to_remove)) {
            return FALSE;
        }
    }

    /**
     * Then, scan the built lists of variables, and actually add them to the
     * environment or remove them the environment. The list is freed while
     * iterating on it.
     */
    while (to_add != NULL && to_remove != NULL) {
        struct node* next_add = to_add->next;
        struct node* next_remove = to_remove->next;

        if (unsetenv(to_remove->name) != 0) {
            fprintf(stderr, " [**] Cannot unset variable \"%s\": %s\n", to_remove->name, strerror(errno));
        }

        if (setenv(to_add->name, to_add->value, TRUE) != 0) {
            fprintf(stderr, " [**] Cannot set variable \"%s\": %s\n", to_add->name, strerror(errno));
        }

        node_delete(to_add);
        node_delete(to_remove);

        to_add = next_add;
        to_remove = next_remove;
    }

    return TRUE;
}
