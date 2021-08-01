# Golang Replacer

Golang Replacer is find a string in directory name, in files name or inside the files and replace with New one in base directory that give.

You just give a base directory, Old String and New string that want to replace.

Example:

```golang
    var replace replace.Replace
	replace.DIR = "/"
	replace.Old = "replacer word"
	replace.New = "new word"
	replace.Find()
```