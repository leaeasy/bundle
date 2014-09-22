#include <stdio.h>
#include <stddef.h>
#include <unistd.h>
#include <stdlib.h>
#include <errno.h>

#include <libkmod.h>

int load(char *name){
        struct kmod_ctx *ctx;
        int err;
        ctx = kmod_new(NULL, NULL);
        if(ctx == NULL)
                return EXIT_FAILURE;

        struct kmod_list *list = 0;

        err = kmod_module_new_from_lookup(ctx, name, &list); 
        if (err != 0){
                printf("Error creating module from path %d\n", err);
                return EXIT_FAILURE;
        }
        if (!list) {
                printf("No such module %s\n", name);
                return EXIT_FAILURE;
        }

        struct kmod_list *it = 0;
        kmod_list_foreach(it, list) {
                struct kmod_module *mod = kmod_module_get_module(it);

                int rc = kmod_module_probe_insert_module (
                                mod,
                                KMOD_PROBE_APPLY_BLACKLIST,
                                0, 0, 0, 0
                                );
                const char* module = kmod_module_get_name(mod);
                if (rc == 0) {
                        printf("Module %s inserted!\n", module);
                }
                kmod_module_unref(mod);
        }

        kmod_unref(ctx);
        return EXIT_SUCCESS;
}
