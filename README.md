# pacsym

## How to use
Only has one feature currently, which is to automatically create symlinks from the /usr/pkg directory. The package manager assumes that all software installed is installed with ``/usr/pkg/<packagename>/<packageversion>``. The same software may not have two versions of the same software installed at the same time, it would cause a slot conflict which the software doesn't know yet how to handle.
```
pacsym sync
```

## Note
To follow true with the nature of LFS, I recommend to look at the source code of the package manager so you truly understand what it's doing to your system. 
