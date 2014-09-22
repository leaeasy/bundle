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

#ifndef CBUNDLE_TYPES_H
#define CBUNDLE_TYPES_H

/**
 * Define a custom boolean type.
 */
enum boolean_values {
    INVALID_BOOL = -1,
    FALSE = 0,
    TRUE = 1
};

typedef enum boolean_values BOOL;

/**
 * Add a custom macro to check if a certain file is executable or not.
 * Has to be used on the st_mode field of a struct stat.
 */
#define S_ISEXE(m) \
    ( S_ISREG((m)) && (((m) & S_IXUSR) | ((m) & S_IXGRP) | ((m) & S_IXOTH)) )

#endif /* CBUNDLE_TYPES_H */