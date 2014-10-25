### 传输协议

#### 客户端

"string\03"

每次发送一个字符串，字符串的结尾以ASCII码0x03结束，由于0x03是不可见字符，故串的正文部分不应该包含0x03不可见字符。

字符串的正文格式为json格式如下

##### 密码验证

```json
{
	"action":	"login",			//action name
	"ip":		"127.0.0.1",		//client ip
	"password": "password",			//password string
	"time":		123456789			//timestamp
}
```

##### 提交任务

```json
{
	"action":	"task_add",			//action name
	"id":		12,					//task id, task id must be unique
	"time":		123456789,			//timestamp
	"language":	"C",				//language
	"code":		"here is code"		//code
}
```

##### 查询任务状态

```json
{
	"action":	"task_status",		//action name
	"id":		12,					//task id
	"time":		123456789
}
```

#### 服务端

服务端返回数据也是json格式

##### 密码验证

登录成功

```json
{
	"result":	true,				//bool, login result
	"os":		"linux x86",
	"language":	{					//language:compiler
					"C":	"gcc",
					"C++":	"g++",
					"Java":	"javac version 1.7"
				},
	"time":		123456789			//server time stamp
}
```

登录失败

```json
{
	"result":	false				//bool, login result
}
```

##### 任务提交响应

提交成功

```json
{
	"result":	true,
	"message":	"task in queue",
	"time":		123456789
}
```

提交失败

```json
{
	"result":	false,
	"message":	"task commit failed, the reason is blablabla...",
	"time":		123456789
}
```

##### 任务状态查询

查询失败(通常指任务不存在)

```json
{
	"result":	false,
	"message":	"task is not exist",
	"time":		123456789
}
```

队列中

```json
{
	"result":	true,
	"status":	"waitting",			//只存在waitting,success,failed三种类型
	"message":	"task in queue",
	"time":		123456789
}
```

通过

```json
{
	"result":	true,
	"status":	"success",
	"log":		"here is compiler log and run log or any other logs",
	"time":		123456789
}
```

失败

```json
{
	"result":	true,
	"status":	"failed",
	"log":		"here is compiler log and run log or any other logs",
	"time":		123456789
}
```

##### 系统状态不可用

```json
{
	"result":	false,
	"message":	"system failure",
	"time":		123456789
}
```

#### 注意

注意，在发送数据的时候数据中不能够包含注释，否则会抛出异常，此问题不准备解决，为了节省流量发送的数据中不应该包含注释。
