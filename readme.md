# Secret

## 安装

```
$ git clone https://github.com/koyeo/secret
$ cd secret
$ ./build.sh
$ mv secret /local/bin/
$ secret -h
```

## 用法

1. 加密

```
$ secret e -key 12345678123456781234567812345678 123
```

2. 解密

```
$ secret d -key 12345678123456781234567812345678 Toz3lwpOYfVm84VgdsgbYQ==
```
