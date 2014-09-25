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

There's no makefile or configure or cmake or whatever needed; it grabs what it needs directly from the source and header files. That being said, there are ways to customize the build: the $CFLAGS, $CXXFLAGS, and $LDFLAGS environment variables, some command-line options, and special directives in the source and header files. These directives are of the form

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
```

Debug builds are simple: just pass `-g` to `qo`.

Cross-compiling is also simple: there's `-os`, `-arch`, and `-tc` commands for specifying target OS, architecture, and toolchain. (This may change.)

TODO conditional compilation via the filename, just like `go build`

`-x` shows a verbose build.

All C files are assumed C99; all C++ files C++11.

Still early in development, still rather unpolished, but suggestions welcome!

### A note on optional features
qo does not support the notion of optional features: everything in the recursive directory tree of the current directory is compiled. I personally don't like features being optional; if something really needs to be conditional, it should be a plugin, and there's no reason to ship a gimped or feature-incomplete version of a program. I don't like how graphviz packages in at least Ubuntu don't sihp with pic support (even though I'm probably the only person int he world that still uses troff).
