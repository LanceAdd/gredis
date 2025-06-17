# gredis simple版本
## 前提概要
1. 由于goframe的gredis过度封装导致无法使用底层go-redis的很多原生方法，所以自定义一个返回go-redis原生UniversalClient客户端
2. 保留了原版从gcfg中读取配置文件的功能，但是增加clear方法用于使用类似nacos作为远程配置中心时，配置文件发生变化时，清除配置缓存和旧的客户端实例
3. 需要在nacos配置中的OnChange方法中调用Clear，当配置文件发生变化时会回调清除现有配置缓存和客户端实例，不会干扰gf全局缓存，只会修改gredis局部缓存
4. goframe v3版本计划废除gredis，可以期待一下

## 使用注意
没啥要注意的，只要记得Clear()就行，测试用例里只包含一个incr用例和pipeline用例，因为完全就是go-redis的client的，直接看这玩意怎么用就行