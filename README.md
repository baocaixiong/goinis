goinis
======

## 有关
goinis 是一个对`ini`文件的简单解析器，当然也做了一些方面的增强

## 特性
- 支持返回 bool, float64, int, int64 以及stirng 类型的值直接返回
- 支持子集section, 并且可以直接由父级section获得子集secion的值
- `key[] = value` 样式表示key 是一个slice类型的值
- 支持keyValue 换行处理，只需要在换行前加上 `-` 即可

## 用例

### config.ini
```ini
	[parent]
	name=johnnihaoyaa
	    -   asfasfhahah
	relation=father
	sex[]=maleqweqw
	-  999
	sex[]=zhangming1
	- 888
	age=32
	boolean=true

	[hasChildren]
	name=has
	[hasChildren.child1]
	name=child1name
```


```go
config, _ := NewConfigFile("config.ini")
s1, _ :=config.GetSection("parent")
v1, _ := s1.GetValue("name")
fmt.Println(v1) // johnnihaoyaa   asfasfhahah(type is interface{})
fmt.Println(s1.MustStringValue("name")) // johnnihaoyaa   asfasfhahah(type is string)
s2, _ := config.GetSection("hasChildren")
sub1, _ := s1.GetSubSection("child1")
v2, _ := sub1.GetValue("name") // child1name(type is interface{})

```

## 参考
- [goconfig](https://github.com/Unknwon/goconfig)
