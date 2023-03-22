# git-multi
CLI to clone multiple git repositories in parallel

(under development)

### Installation
```bash
git clone git@github.com:yagikota/git-multi.git
cd git-multi
```

### Usage
```bash
go run main.go multiclone -h
multiclone clones multiple git repositories in parallel

Usage:
  git-multi multiclone [flags]

Flags:
  -h, --help               help for multiclone
      --maxgoroutine int   max number of goroutine (default 10)
```

```bash
go run main.go multiclone --maxgoroutine=10 git@github.com:gin-gonic/gin.git git@github.com:labstack/echo.git git@github.com:beego/beego.git
==> Cloning 3 repositories:
git@github.com:gin-gonic/gin.git ...
git@github.com:labstack/echo.git ...
git@github.com:beego/beego.git ...
https://github.com/gin-gonic/gin (1/3)
https://github.com/labstack/echo (2/3)
https://github.com/beego/beego (3/3)
====================================================================================================
==> (3/3) success
All 3 repositories are successfully cloned
```
maxgoroutine: The number of goroutine should be set accordingly.

### License
MIT
### Author
https://github.com/yagikota
