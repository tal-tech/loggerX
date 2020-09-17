## loggerX 日志组件库

[[开发文档]](https://www.yuque.com/tal-tech/loggerx)

​	日志作为整个代码行为的记录，是程序执行逻辑和异常最直接的反馈 ，xesLogger日志组件，插件化支持，支持标准输出和高性能磁盘写入,多样化的配置方式，使用起来简单方便，7个日志级别满足项目中各种需求。

#### 1.下载安装  

```shell
go get -u github.com/tal-tech/loggerX 
```


#### 2.main函数中初始化logger：  
* ##### 使用xml配置文件：  

   ```golang
   logger.InitLogger("conf/log.xml")  //xml配置地址  
   defer logger.Close() 
   ```

   - log.xml配置文件示例：

     ```xml
      xxxxxxxxxx <logging>    <filter enabled="false">         <tag>stdout</tag>           <!-- 控制台输出日志 -->        <type>console</type>             <!-- level is (:?FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) -->        <level>DEBUG</level>        <!-- 日志级别 -->    </filter>    <filter enabled="true">        <!-- 文件输出日志 -->        <tag>goentry</tag>                <type>file_trace</type>          <level>INFO</level>           <!-- 日志存储路径 -->         <property name="filename">/home/logs/xeslog/teacherpanel/teacherpanel.log</property>         <!--        %T - Time (15:04:05 MST) 时间格式        %t - Time (15:04)        %D - Date (2006/01/02)   日期格式        %d - Date (01/02/06)        %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)  打印级别        %S - Source          %M - Message        It ignores unknown format strings (and removes them)        Recommended: "[%D %T] [%L] (%S) %M"        -->        <!-- 格式化输出日志 -->         <property name="format">%G %L %S %M</property>        <property name="rotate">true</property> <!-- true enables log rotation, otherwise append --> <property name="maxsize">0M</property> <!-- \d+[KMG]? Suffixes are in terms of 2**10 -->        <property name="maxlines">0K</property> <!-- \d+[KMG]? Suffixes are in terms of thousands -->        <property name="daily">true</property> <!-- Automatically rotates when a log message is written after midnight -->    </filter></logging>
     ```

     

* ##### 使用自定义配置：  

   ```xml
      config := logger.NewLogConfig()
      config.Level="WARNING"  //更新其他配置     
      logger.InitLogWithConfig(config)  
      defer logger.Close() 
   ```
   * 配置说明：

     ```golang
        config格式及默认值 
        //存储路径  默认值 /home/logs/xeslog/default/default.log
        LogPath string  
        //日志级别  默认值 INFO   
        Level string  
        //日志标签 多日志时使用  默认值 default  
        Tag string    
        //日志格式  默认值  "%G %L %S %M"  G=>isoTime L=>level  S=>source  M=>msg    
        Format string  
        //最大行数切割  默认值  0K  支持K\M\G\k\m\g 单位1000  
        RotateLines string  
        //最大容量切割  默认值  0M  支持K\M\G\k\m\g 单位1024  
        RotateSize string  
        //按日期切割   默认值 true  
        RotateDaily bool   
     ```

   * ##### 使用ini配置：   

     ```golang
       logMap := confutil.GetConfStringMap("log")  //通过配置文件转为map[string]string     
        config := logger.NewLogConfig()
        config.SetConfigMap(logMap)    
        logger.InitLogWithConfig(config)   
        defer logger.Close() 
     ```

     * ini文件配置内容:

       ```golang
          //指定字段更新，不指定使用默认值    
          [log]      
          LogPath=/home/logs/xeslog/logpath/log.log   
          Level=DEBUG   
          RotateSize=2G   
          RotateDaily=false 
       ```

#### 3.打印日志方法

* 支持不同级别打印日志的方法：  

```golang
logger.I(tag string, args interface{}, v ...interface{})  
logger.T(tag string, args interface{}, v ...interface{})  
logger.D(tag string, args interface{}, v ...interface{})  
logger.W(tag string, args interface{}, v ...interface{})  
logger.E(tag string, args interface{}, v ...interface{})  
logger.C(tag string, args interface{}, v ...interface{})  
logger.F(tag string, args interface{}, v ...interface{})  //级别同CRITICAL，但触发panic  
logger.Ix(ctx context.Context, tag string, args interface{}, v ...interface{})  
logger.Tx(ctx context.Context, tag string, args interface{}, v ...interface{})  
logger.Dx(ctx context.Context, tag string, args interface{}, v ...interface{})  
logger.Wx(ctx context.Context, tag string, args interface{}, v ...interface{})   
logger.Ex(ctx context.Context, tag string, args interface{}, v ...interface{})  
```
* 使用用例:   

```golang
logger.I("TESTFUNCTAG","data save to mysql, uid:%d ,name:%s",100,"学生1")   
logger.Ix(ctx,"TESTFUNCTAG","data save to mysql, uid:%d ,name:%s",100,"学生1")  
logger.E("TESTFUNCTAG","get redis error：%v, uid:%d ,name:%s",err,100,"学生1")   
logger.Ex(ctx,"TESTFUNCTAG","get redis error:%v, uid:%d ,name:%s",err,100,"学生1")
```
* 支持携带ctx的打印方法:   

```golang
每次调用携带全局变量，支持log特殊需求，如链路追踪等     
默认支持两个特性：  //一次请求开始时写入   
ctx = context.WithValue(ctx, "logid", id)   //id类型string，通过id区分一次完整请求的所有日志      
ctx = context.WithValue(ctx, "start", time.Now())   //每条日志会计算出相对接口开始时间耗时     
```
* 其他   

```
日志支持接入网校链路追踪系统,使用参考example目录
```
####  4.生成error 
```golang  
调用：    
logger.NewError("error",SYSTEM_DEFAULT)
func NewError(err interface{}, ext ...XesError)       
//err传参支持string，error类型(会自动解析rpc server错误)，表示错误根本原因    
//ext XesError，Xes错误码，对外输出错误信息，不传默认系统异常    
type XesError struct {        
    Code int    
    Msg  string    
} 
```
* 其他参考    

```
error_test.go
```
#### 5.日志插件

​     在我们其它的日志组件中，内部打印日志的方式为 `logger.X` 的方式，默认引入xesLogger日志库实例进行打印，若您的项目中采用其它的日志库，可以使用以下方式引入您的日志库。

* 支持引入logrus实例

  ```golang
  //支持logrus日志库使用LoggerX接口格式打印日志 log为logrus全局实例
  	logger.AccessLogLib("logrus", log)
  ```

* 支持引入zap实例

  ```
  //支持Zap日志库使用LoggerX接口格式打印日志 log为zap全局实例
  	logger.AccessLogLib("zap", log)
  ```

* 若现有日志格式不符合预期

  ```
  参考builders目录下builder文件，实现打印日志接口中LoggerX与Build方法即可。
  ```

* 使用参考

  ```
  日志支持接入其他日志库实例,请参考example目录 accessLogLibDemo.go中的demo
  ```

#### 6.注意事项：

 * logger库是并发不安全的，所以全局只能有一个实例。在写单元测时，有可能会多次初始化，此时一定要在包测试完之后进行logger.Close()操作。否则可能会出现如下错误
 ```shell
    FileLogTraceWriter("/xxx/xxx.log"): Rotate: rename /xxx/xxx.log.1: no such file or dircotry\n
 ```
* logger库初始化的时候会先进行一次关闭操作，如果在init方法中使用的logger日志打印，数据写入channel，main函数初始化时进行关闭channel操作时会造成panic。


