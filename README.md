# FireWalker Project

Firewalker Project is about gathering typical _FileSystem_ tasks,
build a simple and robust library to handle typically cases that
involve reading, writing, moving, translating, compressing files and
so on. 

Firewalker is a really fast and flexible _Concurrent_ FileSystem
walker loaded with quite a few maintanance utilities.  Here are some
of _FireWalkers_ existing features:

### Features

- It is really fast..
- It counts files and file sizes
- It can _Glob_ files really fast, e.g *.go

### Coming Soon

- Plugins: will allow a huge array of ways to manipulate files.
- **Glob(pattern)** find files that match the globbed pattern
- **Stats** Gather important stats about all of the files from a dir
- **Compress** Compress the given set of files
- **Copy** Copy the files from src dir to dst dir
- **Cat** Dump the content of every matching file
- **Tar** Create a tar file out of a matching set of files
- **Flatten** Flatten the path of all files to a common prefix
- **JSON/YAML/XML/HTML** Read/Write
- **SCSS->CSS**
- **Thumbnail**

```
% fw glob '*.txt' -> count, cat -> JSON -> gzip -> write all.txt.gz
count -> 
```
