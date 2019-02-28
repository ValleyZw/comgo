# Comgo

Comgo is a [GO](https://en.wikipedia.org/wiki/Go_(programming_language)) based project for parsing [COMTRADE](https://en.wikipedia.org/wiki/Comtrade) files and plotting transient waveforms.


### Description

##### [COMTRADE](https://standards.ieee.org/findstds/standard/C37.111-2013.html) - Common format for Transient Data Exchange for power systems

##### COMTRADE files

>Each **COMTRADE** record has a set of up to four files associated with it.

|Type|Name|Description|Usage|
|:---:|:---:|:---:|:---:|
|xxxxxxxx.HDR|Header file|(Optional) ASCII text file|(Desired format) Up to user|
|xxxxxxxx.CFG|Configuration file|(Essential) ASCII text file|(Specific format) Interprets .DAT file|
|xxxxxxxx.DAT|Data file|(Essential) ASCII or binary format|(Specific format) Store value for channels|
|xxxxxxxx.INF|Information file|(Optional) ASCII or binary format|(Desired format) Contains extra information|

>Useful sites [powergridapp](http://www.powergridapp.com/), [pycomtrade](https://github.com/miguelmoreto/pycomtrade)

### Usage

a. Download and install it

```sh
    $ go get github.com/ValleyZW/comgo
```

b. Import it in your code

```go
    import "github.com/ValleyZW/comgo"
```

c. Init private variable
```go
    var cfg comgo.CFG
```

d. Open and read cfg
```go
    file, err := os.Open(cfgFile)
    err := cfg.ReadCFG(file)
```

e. Open and read dat
```go
    file, err := os.Open(datFile)
    err := cfg.ReadDAT(file)
```

f. Get value of specific channel
```go
points, err := cfg.GetAnalogChannelData(channelNum)
```

### Examples
See [examples](./examples/README.md)
