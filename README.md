# pacsym

A package manager powered by symlinks.

## How to use

The package manager assumes that all software installed is installed with `/usr/pkg/<packagename>/<packageversion>`. The same software may not have two versions of the same software installed at the same time, it would cause a slot conflict which the software doesn't know yet how to handle.

To build software you pass the `build` command, where ``<URL>`` is replaced with a url leading to a tarball with the source code to the software.

```
pacsym build <URL> [OPTIONS] [MAKEFLAGS]
```

If the package is stored locally you may pass the `--local` flag and write

```
pacsym build </PATH/TO/TARBALL> --local ...OPTIONS] [MAKEFLAGS]
```

When the software is done compiling you have to put it into the `/usr/pkg/<PACKAGENAME>/<PACKAGEVERSION` hierarchy by passing the `install` command.

```
pacsym <PACKAGENAME> <PACKAGEVERSION> [MAKEFLAGS]
```

You need to put in a `<PACKAGENAME>` and `<PACKAGEVERSION>` to be used with the previously built package.

Now that you've compiled and installed the packages into the `/usr/pkg/` hierarchy they have to be symlinked, which will make symlinks in the place where they would usually be installed in accordance to the Makefile of the package you installed.

```
pacsym sync
```

## Dependencies

The only dependecy for `pacsym` is `go` and `wget`, where `wget` is optional and only used for URL powered builds.

## How to install

All you have to do to install `pacsym` is to download the repository, execute

```
make
make install
```

and it will install, compile and symlink in compliance to `pacsym`.

After you've installed the package, use
```
pacsym clean
```
to clear the ``/usr/pkgsrc/`` directory, which is important as ``pacsync`` doesn't work correctly when the build directory is still in the ``/usr/pkgsrc/``  directory. The reason it's not automatically cleared after installation is so you still can still interact with the code after install.

## Note

To follow true with the nature of LFS, I recommend to look at the source code of the package manager so you truly understand what it's doing to your system.
