# comgo - cli
This is a command line demo of comgo.

## Tutorials

a. `cd` to current folder.
```sh
    $ cd yourpath/cg
```

b. `build` the demo:
```sh
    $ go build
```
c. `run` to check current version:

```sh
   $ cg -v
    comgo version cg[0.0.1]
```

d. `enjoy` the simple demo

### Usage

a. looking for help:

```sh
   $ cg -h [or] cg --help
    comgo is a comtrade file parser.
        
    Usage:
        detail:  cg [-f] filepath [-d]
        parse:   cg [-f] filepath [-c] channel No.
    options:
        -f	--file		 cfg file path
        -h	--help		 information about the commands
        -c	--channel	 channel No. to save
        -d	--detail	 provide analog channel names
        -v	--version	 print netgo version
```

b. print available analog channels:

```sh
   $ cg -f ..\data\test1.cfg -d [or] cg --file ..\data\test1.cfg --detail 
```

e. save a named channel data to .csv (the same folder of .cfg file):

```sh
   $ cg -f ..\data\test1.cfg -c 10
   success!
```
f. (Optional) [just for fun](http://patorjk.com/software/taag/#p=display&f=Isometric3&t=comgo) - you can test cmd demo
  
```sh
    $ cg
       ____   U  ___ u  __  __     ____    U  ___ u 
    U /"___|   \/"_ \/U|' \/ '|uU /"___|u   \/"_ \/ 
    \| | u     | | | |\| |\/| |/\| |  _ /   | | | | 
     | |/__.-,_| |_| | | |  | |  | |_| |.-,_| |_| | 
      \____|\_)-\___/  |_|  |_|   \____| \_)-\___/  
     _// \\      \\   <<,-,,-.    _)(|_       \\    
    (__)(__)    (__)   (./  \.)  (__)__)     (__) 
    Cmd line:
```