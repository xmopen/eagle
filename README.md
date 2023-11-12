# <img src="https://blogs-1303903194.cos.ap-beijing.myqcloud.com/blogs/169979015623b03ab5d5d64ee7a85dc2a1166268ee_1014961312 (1).png" style="width:48px;height:48px"> Eagle 高性能Golang协程池
高性能golang协程池

- goroutine  复用
- 非阻塞


Java中的线程池: Java中线程池的作用是为了避免重复开线程带来的损耗，因为每次打开一个线程都需要进行上下文切换，损耗太大，所以需要常驻线程以及idle线
程。
Java中需要固定开启N个线程而不销毁掉，同时不断添加任务来让运行的线程执行

Golang协程：无需进行上下文切换，轻量，开启一个协程才2K，1M内存可以开启512个协程，无需担心上下文切换或者内存开销，那为什么还需要协程池呢？ 
Golang的协程主要是为了避免开启大量的goroutine才出现的，虽然Golang的协程很轻量，但是在海量协程情况下也是需要进行限制的，也就是说Golang协程池
需要限制的是不能海量开启goroutine

尽可能、最大化的复用处于idle状态下的Worker