# golang构建数据服务

使用 xorm 实现[博客](http://blog.csdn.net/pmlpml/article/details/78602290#四构建数据服务)中的程序，从**编程效率**、**程序结构**、**服务性能**等角度对比 database/sql 与 orm 实现的异同



### 编程效率

使用xorm，编程效率得到了极大的提高。xorm提供了大量方便使用的API，例如：

- 创建引擎

```go
engine, err := xorm.NewEngine(driverName, dataSourceName)
```

- 自动同步结构体到数据库

```go
err := engine.Sync2(new(StructName))
```

- 插入一条或多条记录

```go
affected, err := engine.Insert(&user)
affected, err := engine.Insert(&user1, &user2)
```

- 查询单条记录（可利用非空数据user来控制条件）

```go
has, err := engine.Get(&user)
```

- 查询多条记录（同样可加入各种条件）

```go
everyone := make([]Userinfo, 0)
err := engine.Find(&everyone)
```

以上只是很小一部分常使用的API以及它们最简单基础的使用方式。利用xorm提高的丰富的API以及它们便于扩展的使用方式，能够非常方便地完成与数据库的交互，编程效率相比原生的database/sql有了极大提高。



### 程序结构

data/sql编程使用的是java经典的“entity - dao - service”层次结构模型。而在xorm中，实现了dao的自动化。即，在service层中，就可以完成与数据库的交互，省去了dao（持久层）结构。

例如，想要在表中插入一个数据，只需要在service层中实现：

```go
func (*UserInfoAtomicService) Save(u *UserInfo) error {
	_, err := engine.Insert(&u)
	checkErr(err)
	return err
}
```

不仅程序结构上得到了很大的简化，编程效率也有很大的提升。



### 功能测试

post数据到网站

```shell
$ curl -d "username=Bob&departname=Software" http://localhost:8080/service/userinfo
```

```
{
  "UID": 5,
  "UserName": "Bob",
  "DepartName": "Software",
  "CreateAt": "2017-11-28T13:22:23.550331819+08:00"
}
```

查询上传的数据

```shell
curl http://localhost:8080/service/userinfo?userid=
```

```shell
[
  ...
  
  {
    "UID": 5,
    "UserName": "Bob",
    "DepartName": "Software",
    "CreateAt": "2017-11-28T00:00:00Z"
  }
]
```



### 服务性能

利用ApacheBench进行压力测试，比较data/sql与xorm：



**data/sql**

```shell
$ ab -n 1000 -c 100 http://localhost:8080/service/userinfo?userid=
```

```shell
This is ApacheBench, Version 2.3 <$Revision: 1757674 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /service/userinfo?userid=
Document Length:        559 bytes

Concurrency Level:      100
Time taken for tests:   0.353 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      683000 bytes
HTML transferred:       559000 bytes
Requests per second:    2831.95 [#/sec] (mean)
Time per request:       35.311 [ms] (mean)
Time per request:       0.353 [ms] (mean, across all concurrent requests)
Transfer rate:          1888.89 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   0.9      0       5
Processing:     1   34  29.5     27     137
Waiting:        1   33  29.5     27     137
Total:          1   34  29.7     27     139
WARNING: The median and mean for the initial connection time are not within a normal deviation
        These results are probably not that reliable.

Percentage of the requests served within a certain time (ms)
  50%     27
  66%     32
  75%     36
  80%     38
  90%     81
  95%    120
  98%    129
  99%    131
 100%    139 (longest request)
```

- 吞吐率(Request per second)：2831.95 [requests/sec]
- 用户平均等待时间(Time per request)：35.311 [ms]
- 服务器平均等待时间(Time per request:across all concurrent requests)：0.353 [ms]



**xorm**

```shell
$ ab -n 1000 -c 100 http://localhost:8080/service/userinfo?userid=
```

```shell
This is ApacheBench, Version 2.3 <$Revision: 1757674 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /service/userinfo?userid=
Document Length:        444 bytes

Concurrency Level:      100
Time taken for tests:   0.286 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      568000 bytes
HTML transferred:       444000 bytes
Requests per second:    3502.31 [#/sec] (mean)
Time per request:       28.553 [ms] (mean)
Time per request:       0.286 [ms] (mean, across all concurrent requests)
Transfer rate:          1942.69 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   0.9      0       8
Processing:     1   27  12.9     27      76
Waiting:        1   27  12.9     26      75
Total:          1   28  13.0     28      76
WARNING: The median and mean for the initial connection time are not within a normal deviation
        These results are probably not that reliable.

Percentage of the requests served within a certain time (ms)
  50%     28
  66%     33
  75%     36
  80%     38
  90%     43
  95%     50
  98%     57
  99%     62
 100%     76 (longest request)
```

- 吞吐率(Request per second)：3502.31 [requests/sec]
- 用户平均等待时间(Time per request)：28.553 [ms]
- 服务器平均等待时间(Time per request:across all concurrent requests)：0.286 [ms]



可以看到，使用xorm的测试中，吞吐量更大，用户平均等待时间和服务器平均等待时间都更短，因此orm对比database/sql，性能更好。
