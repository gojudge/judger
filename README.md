# judger 

<table>
	<tr>
		<td>linux-drone</td>
		<td>windows</td>
		<td>linux-travis</td>
	</tr>
	<tr>
		<td>
			<a title="drone" href="https://drone.io/github.com/duguying/judger/latest">
				<img src="https://drone.io/github.com/duguying/judger/status.png" />
			</a>
		</td>
		<td>
			<a title="appveyor" href="https://ci.appveyor.com/project/duguying/judger">
				<img src="https://ci.appveyor.com/api/projects/status/s20r7g9jfgxngiik?svg=true" />
			</a>
		</td>
		<td>
			<a title="travis" href="https://travis-ci.org/duguying/judger">
				<img src="https://api.travis-ci.org/duguying/judger.png" />
			</a>
		</td>
	</tr>
</table>

the judger server of online judge system

## Build in Linux

```shell
go get
go build
cd sandbox/c/build
cmake ..
make
```

## Install from Docker

```shell
docker pull duguying/judger
```

```shell
mkdir /var/goj/judger
docker run -d -p 1004:1004 -p 1005:1005 -v /var/goj/judger:/data duguying/judger
```

## License #

MIT License
