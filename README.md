# qo: a build system for C/C++

qo is a new build system for C and C++ (though I can add other languages later). In contrast to existing build systems, which require the use of not only a Makefile but also an assortment of complex configuration files (or multiple stages thereof), qo doesn't use any. Instead, custom build settings are embedded using simple directives directly into the source code of your program. qo conditionally compiles each source file based on its filename. qo also supports some resource files normally compiled inot the program. Debug builds and cross-compiles are also intended to be done as easily as possible.

Enjoy! Suggestions, fixes, etc. welcome.

## Installing
qo is written in [Go](http://golang.org/). It has no outside dependencies and does not use cgo, so a compiled qo binary is statically linked and ready to run out of the box. Once the project matures more, I may offer prebuilt binaries for download.

## Getting Started
Let's say you have a simple project in a directory:

```
$ cd project
$ ls
file1.c  file2.c  file3.c  file4.c  project.h
```

To build this project as it stands, simply invoke qo with no arguments:

```
$ qo
[  0%] Beginning build
[ 20%] Compiled file1.c
[ 40%] Compiled file3.c
[ 60%] Compiled file4.c
[ 80%] Compiled file2.c
[100%] Linked project
```

You should see the status of the build as it happens (as above), and upon completion, the compiled program will be left as the executable `project` (named after the project directory) in the project directory, ready for running:

```
$ ./project
```

To build a debug version, pass `-g`:

```
$ qo -g
```

To see the individual commands as they happen, pass `-x`.

Note that qo automatically builds with as many reasonable compiler diagnostics as possible enabled.

## What is Built?
qo scans the current directory and all subdirectories for files to build. Files matched have the given (case-insensitive) extensions:

* C files: `.c`
* C++ files: `.cpp`, `.cxx`, `.c++`, `.cc`
	* note the case-insensitive part; `.C` is recognized as C, not C++
* C header files: `.h`, `.hpp`, `.hxx`, `.h++`, `.hh`
* Objective-C files: `.m`
* Objective-C++ files: `.mm`
* Windows Resource files: `.rc`

Files can be excluded from the build if they are meant for a different operating system and/or CPU architecture; this is also done by filename and is described below, under "Cross-Compiling".

C files are assumed to be C99. C++ files are assumed to be C++11.

## Configuring the Build
So how do you specify extra libraries or compiler options for a project? Simple: you include special directives in the source files! Directives take the form

```
// #qo directive: arguments
```

where whitespace up to and including the first space after `#qo` is significant, and where the `//` must be the first thing on the line.

The two most important (and most portable) directives are `pkg-config` and `LIBS`. `pkg-config` passes the package names listed in `arguments` to `pkg-config`, inserting the resultant compiler flags as needed. `LIBS` takes the library names in `arguments` and passes them to the linker, applying the correct argument format for the toolchain in use (see "Cross-Compiling" below). For example:

```
// #qo pkg-config: gtk+-3.0
// #qo LIBS: pwquality sqlite3
```

For more ocntrol over the command lines for compiling each file, the `CFLAGS`, `CXXFLAGS`, and `LDFLAGS` directives pass their `arguments` as extra arguments to the C compiler, C++ compiler, and linker, respectively.

`#qo` directives are assembled from all source files together. That is, do not copy the directives into each source file; bad things will happen.

In addition, the `$CFLAGS`, `$CXXFLAGS`, and `$LDFLAGS` environment variables also change compiler/linker command-line arguments.

## Cross-Compiling
qo tries to make cross-compiling easy. There are three concepts at play:

- the target OS
- the target architecture
- the toolchain, which defines which compilers and linkers to use

By default, qo builds for the system you are presently running (actually the system the qo binary was built for; but this is a limitation of Go). This is called the host. You can change the target OS, target arch, or toolchain with the `-os`, `-arch`, and `-tc` options, respectively. Pass `list` to see a list of supported OSs, architectures, and toolchains.

(qo by default tends toward gcc/clang-based toolchains.)

In addition, qo will omit files and folders from the build if they are intended for a differnet OS and/or architecture than the target. To omit a file, have `_OS`, `_arch`, or `_OS_arch` before the extension. To omit a folder, its name must consist entirely of `OS`, `arch`, or `OS_arch`. For example:

```
file.c                compiled always
file_windows.c        only compiled if targetting Windows
file_386.c            only compiled if targetting architecture 386
file_windows_386.c    only compiled if targetting 386-based Windows
directory/            trasversed always
windows/              only trasversed on Windows
386/                  only trasversed if targetting 386
windows_386/          only trasversed if targetting 386-based Windows
```

## Cross-Compiler Executable Search Order
Under the hood, however, cross-compiling is a very complex and problematic undertaking for historical and practical reasons. qo assumes you have a correctly configured cross-compiler setup for the target OS, architecture, and toolchain (even if it's just the toolchain).

qo makes the following compromise. Given the following terms:

**unqualified binaries** - binaries named `gcc`, `g++`, `clang`, and `clang++`, without any target triplet<br>
**multilib flags** - `-m32` and `-m64`

1. If the `-triplet` option is passed to qo to explicitly specify a triplet to use, that triplet is used, no questions asked. No mulitlib flags will be appended to the command line.
2. Otherwise, if the target is the same as the host, unqualified binaries are run, and multilib flags may or may not be appended.
3. Otherwise, if the target OS is the same as the host OS and host OS is not Windows, if the host arch is either `386` or `amd64` and the target arch is either `386` or `amd64`, a multilib flag is appended, and the unqualified binaries are run.
4. Otherwise, if using clang, a generic target triplet is generated and used.
5. Otherwise, if the target OS is windows, MinGW-w64 binaries are used.
6. Otherwise, an error occurs.

For more information, see [this](http://stackoverflow.com/a/26101710/3408572) and its references.

## Notes
### A note on optional features and cyclic dependencies
qo does not support the notion of optional features: everything in the recursive directory tree of the current directory is compiled. I personally don't like features being optional; if something really needs to be conditional, it should be a plugin, and there's no reason to ship a gimped or feature-incomplete version of a program. I don't like how graphviz packages in at least Ubuntu don't sihp with pic support (even though I'm probably the only person int he world that still uses troff).

In a related vein, cyclic dependencies (usually caused by optional features, as is the case with GLib ↔ GVFS) should also be avoided.

### Notes on MSVC
The version of MSVC used defines how much C99 or C++11 can be used.

The following seem to be undocumented as being MinGW extensions to the rc format:
- arithmetic expressions
