# Mecchi Loader

A tiny app written with [Golang](https://go.dev/) and [Fyne](https://github.com/fyne-io/fyne) for managing [Mecchi](https://github.com/seruva19/mecchi).  
For now it is a simple wrapper over package.json scripts, but it ~~will~~ might evolve.

## Screenshot

![img](/images/mecchi-loader.png)

## Usage

Requires Go 1.20+.  
You can either clone repository and run:

```
go run .
```

Or download distribution package from Releases. If you already have existing 'mecchi' repository folder, put distributable alongside (_not inside_) it. If you don't, just run  it and press 'install' and then 'launch'.

Not tested yet on Linux.

## TODO

- [ ] Automatic installation of Python 
- [ ] Automatic installation of Git
- [ ] Automatic installation of Node.js
- [ ] Troubleshooting wizard
- [ ] Error report generator