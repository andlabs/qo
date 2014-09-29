qo is a simple program that builds C and C++ programs (though I can add more supported languages later).

Usage is simple:

```
$ cd project
$ ls
file.c file2.c file3.c file4.c project.h
$ qo
...
$ ./project
```

There's no makefile or configure or cmake or whatever needed; it grabs what it needs directly from the source and header files. In fact, here's a list of all the files and their extensions (case insensitive) that qo knows about:

* C files: `.c`
* C++ files: `.cpp`, `.cxx`, `.c++`, `.cc`
	* note the case-insensitive part; `.C` is recognized as C, not C++
* C header files: `.h`, `.hpp`, `.hxx`, `.h++`, `.hh`
* Objective-C files: `.m`
* Objective-C++ files: `.mm`
* Windows Resource files: `.rc`
* TODO also need to add:
	* gresource files: `.xml` with root tag `<gresources>`
	* Qt moc files (will need some way to distinguish; same as C headers)
	* Qt Designer files: `.ui` as XML with root tag `<ui>`
	* anything else (send ideas!)
* TODO can these be embedded?
	* gettext files: `.po`
	* Qt Linguist files: `.ts` as XML with root t ag `<TS>`
	* anything else (send ideas!)

qo also automatically builds with as many reasonable compiler diagnostics as possible enabled.

There are ways to customize the build: the $CFLAGS, $CXXFLAGS, and $LDFLAGS environment variables, some command-line options, and special directives in the source and header files. These directives are of the form

```
// #qo thing: arguments...
```

where the whitespace up to and including the first space after `#qo` is significant.

The directives are

```
CFLAGS
CXXFLAGS
LDFLAGS
	just like the environment variables
pkg-config
	passes named packages to pkg-config and adds the results to CFLAGS, CXXFLAGS, and LDFLAGS
LIBS
	adds named libraries to LDFLAGS. Intended to make cross-compiling with MinGW and MSVC easier; for instance, LIBS: user32 will do -luser32 on MinGW and user32.lib on MSVC
```

Debug builds are simple: just pass `-g` to `qo`.

Cross-compiling is also simple: there's `-os`, `-arch`, and `-tc` commands for specifying target OS, architecture, and toolchain. (`-os` may change.)

To make cross-compiling easy, there are ways to mark a certain source file or directory as being only for a certain target OS, architecture, or both:

- For files, the filename must have `_OS`, `_arch`, or `_OS_arch` before the extension.
- For directories, the directory name must consist entirely of `OS`, `arch`, or `OS_arch`.

The list of supported OSs and architectures can be gathered with `-os list` and `-arch list`, respectively.

For MinGW, use the default `gcc` on Windows and `mingwcc` on other OSs if you have the correct cross-compiler toolchain set up.

`-x` shows a verbose build.

All C files are assumed C99; all C++ files C++11.

For MSVC builds, large address awareness is implied.

For FreeBSD builds, clang is **required** due to GNU triplet madness and gcc versioning conflicts. This means that FreeBSD 8 users will need to install clang manually.

### A note on optional features
qo does not support the notion of optional features: everything in the recursive directory tree of the current directory is compiled. I personally don't like features being optional; if something really needs to be conditional, it should be a plugin, and there's no reason to ship a gimped or feature-incomplete version of a program. I don't like how graphviz packages in at least Ubuntu don't sihp with pic support (even though I'm probably the only person int he world that still uses troff).

### Notes on MSVC
The version of MSVC used defines how much C99 or C++11 can be used.

The following seem to be undocumented as being MinGW extensions to the rc format:
- arithmetic expressions
