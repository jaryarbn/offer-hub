/**
 * 插入题目测试数据。
 *
 * 使用方法：
 * docker exec -i mongodb-dev mongosh < testdata/insert_question.js
 */

use("offer_hub");

try {
const jobName = "后端开发";
const now = new Date();
const int64 = (value) => NumberLong(String(value));

const questionGroups = [
  {
    bank_id: "bank-computer-network",
    questions: [
      ["TCP 三次握手的过程是什么？", "客户端发送 SYN，服务端返回 SYN+ACK，客户端再发送 ACK。三次握手用于确认双方的收发能力并同步初始序列号。", 1, ["TCP", "网络"]],
      ["TCP 和 UDP 的主要区别是什么？", "TCP 面向连接、可靠、有序并提供流量控制；UDP 无连接、不保证可靠性，但报文开销更小，适合实时音视频和 DNS 等场景。", 1, ["TCP", "UDP"]],
      ["HTTP/1.1、HTTP/2 和 HTTP/3 有什么区别？", "HTTP/1.1 支持持久连接，HTTP/2 引入二进制分帧和多路复用，HTTP/3 基于 QUIC，减少了传输层队头阻塞。", 3, ["HTTP", "QUIC"]],
    ],
  },
  {
    bank_id: "bank-operating-system",
    questions: [
      ["进程和线程有什么区别？", "进程是资源分配的基本单位，线程是 CPU 调度的基本单位。同一进程内线程共享地址空间，但各自拥有栈和寄存器上下文。", 1, ["进程", "线程"]],
      ["什么是虚拟内存？", "虚拟内存为进程提供连续的逻辑地址空间，由操作系统和 MMU 映射到物理内存，并可通过换页扩展可用空间。", 2, ["虚拟内存", "分页"]],
      ["死锁产生的四个必要条件是什么？", "互斥、请求并保持、不可剥夺和循环等待同时成立时可能产生死锁，可通过破坏其中一个条件来预防。", 2, ["死锁", "并发"]],
    ],
  },
  {
    bank_id: "bank-mysql",
    questions: [
      ["为什么 InnoDB 使用 B+ 树作为索引结构？", "B+ 树分支因子高、树高低，能减少磁盘 IO；数据集中在叶子节点且叶子节点有序相连，适合等值查询和范围扫描。", 2, ["MySQL", "索引", "B+树"]],
      ["MySQL 的事务隔离级别有哪些？", "SQL 标准定义读未提交、读已提交、可重复读和串行化。InnoDB 默认使用可重复读，并结合 MVCC 和锁处理并发。", 2, ["MySQL", "事务", "MVCC"]],
      ["什么情况下联合索引会失效？", "跳过最左列、对索引列做函数运算、类型隐式转换或使用无法形成有效范围的条件，都可能导致索引不能按预期使用。", 3, ["MySQL", "联合索引", "查询优化"]],
    ],
  },
  {
    bank_id: "bank-redis",
    questions: [
      ["Redis 常用数据结构有哪些？", "Redis 提供 String、Hash、List、Set、Sorted Set、Bitmap、HyperLogLog、Stream 等数据结构。", 1, ["Redis", "数据结构"]],
      ["RDB 和 AOF 持久化有什么区别？", "RDB 保存时间点快照，恢复快但可能丢失最近数据；AOF 记录写命令，数据更完整但文件更大、恢复通常更慢。", 2, ["Redis", "RDB", "AOF"]],
      ["如何处理缓存穿透、击穿和雪崩？", "穿透可用布隆过滤器或缓存空值，击穿可用互斥更新或逻辑过期，雪崩可用随机过期、多级缓存、限流和降级。", 3, ["Redis", "缓存", "高可用"]],
    ],
  },
  {
    bank_id: "bank-mongodb",
    questions: [
      ["MongoDB 适合哪些业务场景？", "MongoDB 适合结构灵活、迭代较快、以文档聚合为主并需要水平扩展的场景，但复杂关联查询通常不如关系数据库自然。", 1, ["MongoDB", "文档数据库"]],
      ["MongoDB 复合索引的最左前缀规则是什么？", "查询从复合索引的第一个字段开始连续匹配时能更有效利用索引，跳过前导字段通常无法完成有效的索引定位。", 2, ["MongoDB", "索引"]],
      ["MongoDB 副本集如何保证高可用？", "副本集由一个 Primary 和多个 Secondary 组成，节点间复制 oplog；Primary 不可用时，满足条件的成员会发起选举产生新 Primary。", 3, ["MongoDB", "副本集", "选举"]],
    ],
  },
  {
    bank_id: "bank-java",
    questions: [
      ["HashMap 的基本实现原理是什么？", "HashMap 使用数组加链表或红黑树存储键值对，通过哈希值定位桶；冲突严重且容量达到条件时链表会树化。", 2, ["Java", "HashMap"]],
      ["Java 中 equals 和 hashCode 有什么约定？", "equals 相等的对象必须具有相同 hashCode；重写 equals 时通常也要重写 hashCode，否则在哈希容器中会出现不一致行为。", 1, ["Java", "对象"]],
      ["synchronized 和 ReentrantLock 有什么区别？", "两者都支持可重入互斥。ReentrantLock 还提供可中断获取、公平锁和多个 Condition，但需要显式在 finally 中释放。", 2, ["Java", "并发", "锁"]],
    ],
  },
  {
    bank_id: "bank-go",
    questions: [
      ["Go 中 slice 的底层结构是什么？", "slice 是对底层数组的描述，包含指针、长度和容量。append 超过容量时会分配新数组并复制已有元素。", 1, ["Go", "Slice"]],
      ["Go 接口的 nil 判断有哪些陷阱？", "接口值包含动态类型和动态值。只有两者都为空时接口才等于 nil；装入一个值为 nil 的具体指针后，接口本身不等于 nil。", 2, ["Go", "Interface", "nil"]],
      ["defer 的执行顺序和参数求值时机是什么？", "defer 按后进先出顺序执行，函数参数在注册 defer 时立即求值，闭包引用的外部变量则在闭包实际执行时读取。", 2, ["Go", "defer"]],
    ],
  },
  {
    bank_id: "bank-jvm",
    questions: [
      ["JVM 运行时数据区包括哪些部分？", "运行时数据区通常包括程序计数器、虚拟机栈、本地方法栈、堆和方法区，其中堆与方法区由线程共享。", 1, ["Java", "JVM", "内存"]],
      ["类加载的双亲委派模型是什么？", "类加载器收到请求后先委派父加载器，父加载器无法完成时才由自己加载，可避免核心类被重复或恶意替换。", 2, ["JVM", "类加载"]],
      ["G1 垃圾收集器的核心思路是什么？", "G1 将堆划分为多个 Region，按收益预测优先回收垃圾较多的区域，并以可预测停顿为目标执行并发标记和混合回收。", 3, ["JVM", "GC", "G1"]],
    ],
  },
  {
    bank_id: "bank-spring",
    questions: [
      ["Spring IoC 解决了什么问题？", "IoC 容器负责对象创建、依赖装配和生命周期管理，使业务代码减少对具体实现和构造过程的直接依赖。", 1, ["Spring", "IoC"]],
      ["Spring AOP 常见的应用场景有哪些？", "AOP 常用于事务、日志、权限、监控等横切关注点，通过代理在目标方法前后织入统一逻辑。", 2, ["Spring", "AOP"]],
      ["为什么同类内部调用可能导致 @Transactional 失效？", "Spring 声明式事务通常依赖代理。同类内部直接调用没有经过代理对象，事务拦截器不会被触发。", 3, ["Spring", "事务", "代理"]],
    ],
  },
  {
    bank_id: "bank-go-concurrency",
    questions: [
      ["Goroutine 和操作系统线程有什么区别？", "Goroutine 由 Go runtime 调度，初始栈较小且可增长，多个 goroutine 会复用较少的操作系统线程。", 1, ["Go", "Goroutine", "GMP"]],
      ["无缓冲 Channel 和有缓冲 Channel 有什么区别？", "无缓冲 Channel 要求发送和接收同时就绪；有缓冲 Channel 在缓冲区未满或非空时允许双方暂时独立执行。", 2, ["Go", "Channel", "并发"]],
      ["如何避免 Goroutine 泄漏？", "需要为阻塞操作提供退出路径，常用 context 传递取消信号，并确保 Channel 的生产者、消费者和关闭责任清晰。", 3, ["Go", "Goroutine", "Context"]],
    ],
  },
  {
    bank_id: "bank-distributed-system",
    questions: [
      ["CAP 定理表达了什么？", "发生网络分区时，分布式系统无法同时保证强一致性和可用性，需要结合业务在一致性与可用性之间选择。", 2, ["分布式", "CAP"]],
      ["如何设计幂等的写接口？", "客户端携带唯一幂等键，服务端使用唯一索引或原子状态记录结果；重复请求返回首次结果，并处理并发与记录过期。", 3, ["分布式", "幂等"]],
      ["Raft 如何完成领导者选举？", "节点超时后转为 Candidate，提高任期并请求投票；获得多数票后成为 Leader，随后通过心跳维持领导关系。", 3, ["分布式", "Raft", "选举"]],
    ],
  },
  {
    bank_id: "bank-microservice",
    questions: [
      ["服务注册与发现的作用是什么？", "服务实例将地址和健康状态注册到注册中心，调用方通过查询或订阅获得可用实例列表，从而支持动态扩缩容。", 1, ["微服务", "服务发现"]],
      ["熔断、限流和降级有什么区别？", "熔断隔离持续失败的依赖，限流控制进入系统的请求量，降级在资源不足或依赖异常时提供简化能力。", 2, ["微服务", "熔断", "限流"]],
      ["分布式链路追踪如何关联一次请求？", "入口生成 Trace ID，各服务为调用生成 Span 并传递上下文，采集器汇总时序、标签和错误信息还原调用链。", 2, ["微服务", "Trace", "可观测性"]],
    ],
  },
  {
    bank_id: "bank-message-queue",
    questions: [
      ["为什么系统会引入消息队列？", "消息队列可用于异步处理、系统解耦和削峰填谷，但会增加消息一致性、重复消费和运维复杂度。", 1, ["消息队列", "异步"]],
      ["如何保证消息至少被消费一次？", "生产端确认写入，Broker 持久化，消费端处理成功后再确认；失败时重试，因此消费者还需要保证业务幂等。", 2, ["消息队列", "可靠性", "幂等"]],
      ["Kafka 分区有什么作用？", "分区是 Kafka 并行读写和扩展吞吐的基本单位，同一分区内消息有序，消费者组内一个分区同一时刻只分配给一个消费者。", 2, ["Kafka", "分区", "消费者组"]],
    ],
  },
  {
    bank_id: "bank-docker-kubernetes",
    questions: [
      ["Docker 镜像和容器有什么区别？", "镜像是只读的分层模板，容器是镜像的运行实例，并在镜像层之上增加可写层。", 1, ["Docker", "容器"]],
      ["Kubernetes Pod 是什么？", "Pod 是 Kubernetes 最小调度单元，其中一个或多个容器共享网络命名空间和部分存储卷，通常共同完成一个紧密协作的任务。", 1, ["Kubernetes", "Pod"]],
      ["Deployment 如何实现滚动更新？", "Deployment 通过调整新旧 ReplicaSet 的副本数逐步替换 Pod，并结合就绪探针和更新策略控制可用性。", 2, ["Kubernetes", "Deployment", "滚动更新"]],
    ],
  },
  {
    bank_id: "bank-linux",
    questions: [
      ["如何查看 Linux 端口被哪个进程占用？", "可以使用 ss -lntp 或 lsof -i :端口 查看监听套接字及对应进程，再结合 ps 获取进程详情。", 1, ["Linux", "网络", "排查"]],
      ["Linux load average 表示什么？", "load average 表示一段时间内处于可运行状态或不可中断睡眠状态的平均任务数，需要结合 CPU 核数和 IO 指标判断。", 2, ["Linux", "负载"]],
      ["如何排查 Linux 服务器 CPU 使用率过高？", "先用 top 或 pidstat 定位进程和线程，再结合调用栈、性能剖析、系统调用及业务日志缩小问题范围。", 2, ["Linux", "CPU", "性能排查"]],
    ],
  },
  {
    bank_id: "bank-algorithm-data-structure",
    questions: [
      ["哈希表的平均时间复杂度是多少？", "在哈希函数和负载因子合理时，插入、查询和删除平均为 O(1)；冲突严重时最坏复杂度会退化。", 1, ["数据结构", "哈希表"]],
      ["快速排序的基本思想是什么？", "选择基准值，将序列划分为小于和大于基准的两部分，再递归排序。平均复杂度 O(n log n)，最坏 O(n²)。", 2, ["算法", "快速排序"]],
      ["动态规划适合解决什么问题？", "动态规划适合具有重叠子问题和最优子结构的问题，通过定义状态、转移方程和初始条件避免重复计算。", 2, ["算法", "动态规划"]],
    ],
  },
];

const supplementalQuestions = {
  "bank-computer-network": [
    ["TCP 四次挥手的过程是什么？", "主动关闭方发送 FIN，对端确认 ACK；对端处理完剩余数据后再发送 FIN，主动关闭方回复 ACK 并进入 TIME_WAIT。", 1, ["TCP", "挥手"]],
    ["为什么 TCP 主动关闭方需要进入 TIME_WAIT？", "TIME_WAIT 让最后一个 ACK 有机会重传，并等待旧连接中的延迟报文失效，避免影响相同四元组的新连接。", 2, ["TCP", "TIME_WAIT"]],
    ["TCP 流量控制和拥塞控制有什么区别？", "流量控制通过接收窗口避免压垮接收方；拥塞控制通过拥塞窗口和慢启动等算法保护网络链路。", 2, ["TCP", "拥塞控制"]],
    ["DNS 域名解析通常经历哪些步骤？", "客户端依次查询浏览器和系统缓存、递归 DNS 服务器；递归服务器再向根、顶级域和权威服务器查询并缓存结果。", 1, ["DNS", "网络"]],
    ["HTTPS 建立连接时 TLS 握手完成了什么？", "握手协商协议和密码套件、验证服务端证书，并通过密钥交换生成会话密钥，后续使用对称加密传输。", 2, ["HTTPS", "TLS"]],
    ["Cookie、Session 和 Token 有什么区别？", "Cookie 是浏览器存储和携带数据的机制；Session 通常在服务端保存状态；Token 常由客户端保存并携带可验证的身份凭证。", 1, ["HTTP", "认证"]],
    ["WebSocket 为什么适合实时双向通信？", "WebSocket 在一次 HTTP Upgrade 后复用长连接，客户端和服务端都能主动发送帧，减少轮询开销。", 2, ["WebSocket", "实时通信"]],
    ["CDN 如何降低用户访问延迟？", "CDN 将静态或可缓存内容分发到靠近用户的边缘节点，并通过 DNS 调度把请求导向合适节点。", 1, ["CDN", "缓存"]],
    ["正向代理和反向代理有什么区别？", "正向代理代表客户端访问外部服务；反向代理代表服务端接收请求，可承担负载均衡、缓存和 TLS 终止。", 1, ["代理", "负载均衡"]],
    ["NAT 解决了什么问题？", "NAT 在私网地址和公网地址之间转换，缓解 IPv4 地址不足并隐藏内部网络，但会增加端到端通信复杂度。", 2, ["NAT", "IPv4"]],
    ["HTTP 缓存中的强缓存和协商缓存有什么区别？", "强缓存根据 Cache-Control 或 Expires 直接复用本地副本；协商缓存携带 ETag 或 Last-Modified 向服务端确认是否变化。", 2, ["HTTP", "缓存"]],
    ["TCP 粘包和拆包产生的原因是什么？", "TCP 是字节流协议，没有消息边界；发送缓冲、MSS 和接收读取方式都可能让多条业务消息合并或拆分，需要应用层协议定界。", 2, ["TCP", "协议设计"]],
  ],
  "bank-operating-system": [
    ["用户态和内核态有什么区别？", "用户态权限受限，不能直接访问硬件和内核内存；系统调用会切换到内核态执行受保护操作。", 1, ["操作系统", "系统调用"]],
    ["进程上下文切换为什么有开销？", "切换需要保存和恢复寄存器、调度状态与地址空间相关信息，还可能破坏 CPU 缓存和 TLB 的局部性。", 2, ["进程", "上下文切换"]],
    ["常见的 CPU 调度算法有哪些？", "常见算法包括先来先服务、短作业优先、时间片轮转、优先级调度和多级反馈队列，各自在吞吐与响应时间间权衡。", 1, ["CPU", "调度"]],
    ["分页和分段内存管理有什么区别？", "分页按固定大小划分地址空间以便物理内存分配；分段按代码、数据等逻辑单元划分，长度可变且更符合程序结构。", 2, ["内存", "分页"]],
    ["什么是缺页异常？", "进程访问的虚拟页不在当前物理内存时触发缺页异常，内核会装入页面、更新页表后重新执行指令。", 2, ["虚拟内存", "缺页"]],
    ["写时复制 Copy-on-Write 的作用是什么？", "父子进程先共享只读物理页，只有某一方写入时才复制页面，可降低 fork 的初始时间和内存消耗。", 2, ["内存", "COW"]],
    ["进程间通信有哪些常见方式？", "常见方式包括管道、消息队列、共享内存、信号、Socket 和信号量，其中共享内存性能高但需要同步。", 1, ["进程", "IPC"]],
    ["select、poll 和 epoll 有什么区别？", "select 有描述符数量限制，poll 使用数组但仍需线性扫描；epoll 在 Linux 中维护关注集合并返回就绪事件，更适合大量连接。", 3, ["Linux", "IO多路复用"]],
    ["文件描述符是什么？", "文件描述符是进程打开文件、Socket 或管道后获得的整数句柄，内核通过它定位对应的打开文件对象。", 1, ["文件系统", "文件描述符"]],
    ["mmap 与普通 read/write 有什么区别？", "mmap 将文件映射到进程虚拟地址空间，可像访问内存一样读写并减少用户态缓冲复制，但仍会发生缺页和回写。", 2, ["mmap", "文件系统"]],
    ["互斥锁和信号量有什么区别？", "互斥锁通常保护临界区且强调持有者释放；信号量维护计数，可控制多个并发许可或用于线程间通知。", 2, ["并发", "同步"]],
    ["僵尸进程和孤儿进程有什么区别？", "僵尸进程已经退出但父进程尚未回收其状态；孤儿进程的父进程先退出，之后会由系统进程接管。", 2, ["进程", "Linux"]],
  ],
  "bank-mysql": [
    ["聚簇索引和二级索引有什么区别？", "InnoDB 聚簇索引叶子节点保存整行数据，二级索引叶子节点保存索引列和主键，查询其他列时可能需要回表。", 2, ["MySQL", "索引"]],
    ["InnoDB 的 MVCC 是如何工作的？", "InnoDB 通过隐藏事务字段、undo log 版本链和 Read View 判断记录可见性，让一致性读在多数场景下避免加锁。", 3, ["MySQL", "MVCC"]],
    ["间隙锁和 Next-Key Lock 的作用是什么？", "间隙锁锁定索引值之间的范围，Next-Key Lock 组合记录锁和间隙锁，用于可重复读下防止范围查询出现幻读。", 3, ["MySQL", "锁"]],
    ["redo log、undo log 和 binlog 分别有什么作用？", "redo log 保证崩溃恢复，undo log 支持回滚和 MVCC，binlog 记录逻辑变更用于复制和数据恢复。", 2, ["MySQL", "日志"]],
    ["MySQL 为什么采用 WAL？", "WAL 先顺序写日志再异步刷新数据页，把随机写转换为顺序写，并能在宕机后根据日志恢复已提交事务。", 2, ["MySQL", "WAL"]],
    ["如何使用 EXPLAIN 分析慢查询？", "重点关注访问类型、命中索引、扫描行数、过滤比例和 Extra 信息，再结合数据分布判断是否需要改写 SQL 或索引。", 2, ["MySQL", "EXPLAIN"]],
    ["深分页为什么慢，如何优化？", "大 OFFSET 仍需扫描并丢弃前面的记录；可使用覆盖索引延迟关联，或基于上一页唯一排序键做游标分页。", 3, ["MySQL", "分页"]],
    ["MySQL 死锁发生后应该如何处理？", "InnoDB 会检测循环等待并回滚代价较小的事务；业务应捕获错误重试，并通过固定加锁顺序和缩短事务降低概率。", 3, ["MySQL", "死锁"]],
    ["主从复制的基本流程是什么？", "主库写入 binlog，从库 IO 线程拉取为 relay log，再由 SQL 线程或并行工作线程重放变更。", 2, ["MySQL", "复制"]],
    ["读写分离可能带来什么一致性问题？", "主从复制存在延迟，刚写入的数据从从库读取可能不可见，可采用读主、等待位点或基于业务的会话一致性策略。", 2, ["MySQL", "读写分离"]],
    ["数据库分库分表后常见的挑战有哪些？", "需要处理路由、全局唯一 ID、跨分片查询和事务、扩容迁移，以及统计分页等操作复杂度。", 3, ["MySQL", "分库分表"]],
    ["为什么不建议使用过长的主键？", "InnoDB 二级索引叶子节点保存主键，主键过长会放大所有二级索引，占用更多内存和磁盘并降低缓存效率。", 2, ["MySQL", "主键"]],
  ],
  "bank-redis": [
    ["Redis 的过期删除策略有哪些？", "Redis 结合惰性删除和定期随机抽查清理过期键，避免全量扫描，同时在访问过期键时立即删除。", 2, ["Redis", "过期策略"]],
    ["Redis 常见的内存淘汰策略有哪些？", "可按所有键或带过期键执行 LRU、LFU、随机淘汰，也可淘汰最短 TTL 或拒绝写入，应根据访问模式选择。", 2, ["Redis", "淘汰策略"]],
    ["Redis 单线程为什么仍然很快？", "核心命令主要在内存中执行，采用高效数据结构和 IO 多路复用，避免了大量线程切换与锁竞争。", 1, ["Redis", "性能"]],
    ["Redis 分布式锁需要满足哪些条件？", "加锁应原子设置唯一持有者和过期时间，解锁时用 Lua 校验持有者后删除，并处理续期、时钟和故障切换风险。", 3, ["Redis", "分布式锁"]],
    ["Redis Cluster 如何分配数据？", "Cluster 将键映射到 16384 个哈希槽，槽分配给不同主节点；客户端根据重定向信息访问负责该槽的节点。", 2, ["Redis", "Cluster"]],
    ["Sentinel 的主要职责是什么？", "Sentinel 监控主从节点，通过主观和客观下线判断故障，选举后执行主从切换并通知客户端新主节点。", 2, ["Redis", "Sentinel"]],
    ["Pipeline 和事务有什么区别？", "Pipeline 批量发送命令以减少网络往返，不保证原子性；MULTI/EXEC 将命令排队后顺序执行，但不提供传统回滚。", 2, ["Redis", "Pipeline"]],
    ["Lua 脚本在 Redis 中有什么用途？", "脚本在服务端原子执行多条命令，适合带条件的扣减、锁释放等逻辑，并减少网络往返。", 2, ["Redis", "Lua"]],
    ["如何发现和治理 Redis 大 Key？", "可通过扫描和内存分析定位大 Key，再拆分数据、渐进删除并限制集合增长，避免阻塞和网络突发。", 2, ["Redis", "大Key"]],
    ["如何处理 Redis 热 Key？", "可以增加本地缓存、拆分 Key、读副本或代理分片，并通过限流和请求合并降低单节点压力。", 3, ["Redis", "热Key"]],
    ["缓存与数据库双写如何降低不一致？", "常用先更新数据库再删除缓存，并通过重试、消息或订阅 binlog 保证删除最终成功；还要设置合理过期时间。", 3, ["Redis", "缓存一致性"]],
    ["Sorted Set 的典型实现与场景是什么？", "Sorted Set 通常结合哈希表和跳表，既能按成员定位又能按分值范围排序，适合排行榜和延时任务。", 2, ["Redis", "SortedSet"]],
  ],
  "bank-mongodb": [
    ["MongoDB 的文档模型有什么特点？", "MongoDB 使用 BSON 文档保存嵌套结构，字段可以灵活演进，适合围绕聚合根一次读取的数据建模。", 1, ["MongoDB", "BSON"]],
    ["MongoDB 中嵌入文档和引用如何选择？", "一起读取且生命周期一致的数据适合嵌入；多对多、独立增长或频繁单独更新的数据更适合引用。", 2, ["MongoDB", "建模"]],
    ["Aggregation Pipeline 的执行方式是什么？", "文档依次经过 $match、$project、$group、$sort 等阶段，每个阶段对输入进行转换并传给下一阶段。", 2, ["MongoDB", "聚合"]],
    ["MongoDB 多文档事务适合什么场景？", "事务用于必须跨多个文档或集合原子提交的操作，但成本高于单文档原子更新，建模时应优先减少跨文档事务。", 3, ["MongoDB", "事务"]],
    ["MongoDB 分片集群由哪些组件组成？", "分片保存数据，Config Server 保存元数据，mongos 负责查询路由；分片键决定数据分布和查询效率。", 2, ["MongoDB", "分片"]],
    ["如何选择合适的 MongoDB 分片键？", "应兼顾高基数、均匀写入和常见查询前缀，避免单调热点，同时减少无法定向到少量分片的广播查询。", 3, ["MongoDB", "分片键"]],
    ["Write Concern 控制什么？", "Write Concern 决定写操作需要多少副本确认以及是否等待日志持久化，在延迟、可用性和持久性之间权衡。", 2, ["MongoDB", "WriteConcern"]],
    ["Read Concern 和 Read Preference 有什么区别？", "Read Concern 控制读取数据的一致性级别；Read Preference 决定从主节点还是符合条件的从节点读取。", 2, ["MongoDB", "一致性"]],
    ["TTL 索引适合什么用途？", "TTL 索引会由后台线程清理超过指定时间的文档，适合会话、临时日志等，但删除并非到期瞬间完成。", 1, ["MongoDB", "TTL"]],
    ["如何使用 explain 分析 MongoDB 查询？", "检查是否使用 IXSCAN、扫描键和文档数量、返回数量及排序阶段，判断索引选择性和查询形状是否合理。", 2, ["MongoDB", "explain"]],
    ["oplog 和 Change Stream 有什么关系？", "副本集使用 oplog 复制变更，Change Stream 在其基础上向应用提供可恢复的变更事件订阅接口。", 3, ["MongoDB", "oplog"]],
    ["MongoDB 大数据量分页如何优化？", "深分页应避免大 skip，可基于稳定排序字段和 _id 记录上一页边界，通过范围条件继续查询。", 2, ["MongoDB", "分页"]],
  ],
  "bank-java": [
    ["String 为什么设计为不可变对象？", "不可变使 String 可安全共享和缓存哈希值，便于作为哈希键，也提高了线程安全和类加载等安全场景的可靠性。", 1, ["Java", "String"]],
    ["ArrayList 和 LinkedList 如何选择？", "ArrayList 随机访问和遍历更高效；LinkedList 已知节点位置时插删方便，但定位仍是线性且对象开销更大。", 1, ["Java", "集合"]],
    ["ConcurrentHashMap 如何支持并发访问？", "Java 8 使用桶数组、CAS 和桶级 synchronized 协调更新，读取大多无锁，冲突严重时链表可转为红黑树。", 3, ["Java", "ConcurrentHashMap"]],
    ["volatile 能保证哪些语义？", "volatile 保证写入对其他线程可见并限制相关指令重排，但复合读改写操作仍不具备原子性。", 2, ["Java", "volatile"]],
    ["Java 线程池的核心参数有哪些？", "核心线程数、最大线程数、空闲时间、任务队列、线程工厂和拒绝策略共同决定并发量、排队与过载行为。", 2, ["Java", "线程池"]],
    ["受检异常和非受检异常有什么区别？", "受检异常要求编译期捕获或声明；RuntimeException 及其子类属于非受检异常，通常表示编程错误或不可恢复条件。", 1, ["Java", "异常"]],
    ["Java 泛型为什么会发生类型擦除？", "大多数泛型类型信息在编译后被擦除为上界并插入必要的转换，以兼容早期 JVM 字节码和类库。", 2, ["Java", "泛型"]],
    ["反射的优点和代价是什么？", "反射能在运行时检查和调用类型，适合框架扩展，但可读性、类型安全和性能通常不如直接调用。", 2, ["Java", "反射"]],
    ["final、finally 和 finalize 有什么区别？", "final 修饰不可重新赋值或继承，finally 是异常处理清理块；finalize 是已弃用的不可靠对象回收回调。", 1, ["Java", "关键字"]],
    ["Java Stream 使用时有哪些常见陷阱？", "Stream 通常只能消费一次，副作用会削弱可读性，并行流还需考虑拆分成本、线程安全和公共线程池干扰。", 2, ["Java", "Stream"]],
    ["Java SPI 的基本机制是什么？", "服务提供者在配置文件中声明实现类，ServiceLoader 按约定发现并实例化，实现接口与具体实现的解耦。", 2, ["Java", "SPI"]],
    ["record 类型适合什么场景？", "record 适合表达不可变数据载体，编译器自动生成构造器、访问器、equals、hashCode 和 toString，但它仍可实现接口和定义方法。", 1, ["Java", "record"]],
  ],
  "bank-go": [
    ["Go 中数组和切片有什么区别？", "数组长度属于类型且按值复制；切片是指向底层数组的描述，包含指针、长度和容量。", 1, ["Go", "数组"]],
    ["Go map 是否支持并发读写？", "普通 map 不支持并发读写，可能触发竞态甚至运行时错误；应使用锁、分片或 sync.Map 等方案。", 1, ["Go", "Map"]],
    ["byte 和 rune 分别表示什么？", "byte 是 uint8 的别名，常表示原始字节；rune 是 int32 的别名，用于表示 Unicode 码点。", 1, ["Go", "字符串"]],
    ["Go 推荐如何处理错误？", "错误作为普通返回值显式传递，调用方应及时判断并补充上下文；可用 errors.Is 和 errors.As 检查包装后的错误链。", 1, ["Go", "Error"]],
    ["panic 和 recover 应该如何使用？", "panic 表示当前流程无法继续，recover 只能在延迟函数中捕获；业务错误通常应返回 error，而不是依赖 panic 控制流程。", 2, ["Go", "panic"]],
    ["值接收者和指针接收者如何选择？", "需要修改接收者、避免复制大对象或保持方法集一致时使用指针接收者；小型不可变值可使用值接收者。", 2, ["Go", "方法"]],
    ["Go 变量为什么会逃逸到堆上？", "当编译器判断变量生命周期可能超过当前栈帧或无法安全放在栈上时会逃逸，具体结果可通过编译器分析查看。", 2, ["Go", "逃逸分析"]],
    ["Go 的 init 函数执行顺序是什么？", "运行时先初始化依赖包，再按文件依赖完成包级变量初始化，随后执行各 init，最后进入 main。", 2, ["Go", "init"]],
    ["Go Module 解决了什么问题？", "Go Module 使用 go.mod 描述模块路径和依赖版本，并通过最小版本选择和校验文件实现可复现依赖管理。", 1, ["Go", "Module"]],
    ["context.Context 应该如何传递？", "Context 通常作为首个参数沿调用链传递，不存入结构体或传 nil，用于取消、截止时间和请求范围元数据。", 2, ["Go", "Context"]],
    ["Go 泛型中的类型约束是什么？", "类型约束是描述允许类型集合的接口，泛型函数可基于约束使用所有候选类型共同支持的操作。", 2, ["Go", "泛型"]],
    ["new 和 make 有什么区别？", "new 为任意类型分配零值并返回指针；make 仅初始化 slice、map 和 channel，并返回可直接使用的值。", 1, ["Go", "内存"]],
  ],
  "bank-jvm": [
    ["Java 对象通常如何在堆上分配？", "线程优先在 TLAB 中通过指针碰撞分配，空间不足时再申请新 TLAB 或走共享分配路径，大对象可能直接进入老年代。", 2, ["JVM", "对象分配"]],
    ["什么对象可以作为 GC Roots？", "活动线程栈中的引用、静态字段、JNI 引用和 JVM 内部引用等可作为根，垃圾回收从这些根做可达性分析。", 2, ["JVM", "GC Roots"]],
    ["Minor GC、Major GC 和 Full GC 有什么区别？", "Minor GC 主要回收新生代；Major GC 常指老年代回收；Full GC 通常回收整个堆并可能处理类元数据。", 2, ["JVM", "GC"]],
    ["什么是 Stop-The-World？", "STW 是 JVM 暂停应用线程以获得稳定运行时状态的阶段，GC、安全点操作等可能触发，目标是尽量缩短停顿。", 1, ["JVM", "STW"]],
    ["JVM 安全点有什么作用？", "线程到达安全点后，JVM 能可靠检查栈和对象引用，用于垃圾回收、反优化等需要全局一致状态的操作。", 2, ["JVM", "Safepoint"]],
    ["新生代为什么常采用复制算法？", "新生代对象死亡率高，复制少量存活对象即可回收整块空间，分配速度快且避免碎片。", 2, ["JVM", "新生代"]],
    ["元空间和永久代有什么区别？", "元空间使用本地内存保存类元数据，取代了受固定堆区域限制的永久代，但仍需设置和监控上限。", 1, ["JVM", "元空间"]],
    ["常见的 OutOfMemoryError 类型有哪些？", "常见有 Java heap space、Metaspace、Direct buffer memory 和 unable to create native thread，应结合错误类型和监控定位资源。", 2, ["JVM", "OOM"]],
    ["如何分析 JVM 内存泄漏？", "先观察堆增长和 GC 后占用，再获取 heap dump，使用支配树、引用链和对象增长对比定位无法释放的持有者。", 3, ["JVM", "内存泄漏"]],
    ["JIT 编译器的作用是什么？", "JIT 根据运行时热点把字节码编译为机器码，并执行内联、逃逸分析等优化，提高长期运行代码性能。", 2, ["JVM", "JIT"]],
    ["逃逸分析可以触发哪些优化？", "编译器确认对象不逃逸时，可能进行标量替换、栈上效果的消除分配和锁消除，但具体取决于优化条件。", 3, ["JVM", "逃逸分析"]],
    ["如何选择 G1、ZGC 或 Parallel GC？", "吞吐优先可考虑 Parallel GC，平衡吞吐和停顿可选 G1，超大堆低停顿场景可评估 ZGC，并用真实负载验证。", 3, ["JVM", "垃圾收集器"]],
  ],
  "bank-spring": [
    ["Spring Bean 的生命周期包括哪些阶段？", "Bean 经实例化、属性填充、Aware 回调、前后置处理、初始化后可用，容器关闭时再执行销毁回调。", 2, ["Spring", "Bean"]],
    ["Spring 常见的 Bean Scope 有哪些？", "常见有 singleton、prototype，以及 Web 环境中的 request、session 等；不同作用域的创建和销毁责任不同。", 1, ["Spring", "Scope"]],
    ["Spring 如何处理单例 Bean 的循环依赖？", "默认情况下可通过提前暴露引用解决部分 setter 注入循环依赖，但构造器循环依赖和某些代理场景仍无法解决。", 3, ["Spring", "循环依赖"]],
    ["Spring AOP 使用 JDK 代理还是 CGLIB？", "目标实现接口时可使用 JDK 动态代理，否则通常使用 CGLIB 子类代理；final 类或方法会限制子类代理。", 2, ["Spring", "AOP"]],
    ["Spring MVC 一次请求的主要流程是什么？", "DispatcherServlet 查找 Handler 与适配器，执行控制器后由返回值处理器、视图解析或消息转换器生成响应。", 2, ["Spring MVC", "请求流程"]],
    ["Spring Boot 自动配置的基本原理是什么？", "自动配置类通过条件注解判断类路径、Bean 和配置属性，在用户未自定义时提供合理默认 Bean。", 2, ["Spring Boot", "自动配置"]],
    ["Spring 事务传播行为解决什么问题？", "传播行为定义已有事务存在时新方法加入、挂起还是创建新事务，如 REQUIRED、REQUIRES_NEW 和 NESTED。", 3, ["Spring", "事务传播"]],
    ["为什么 checked exception 默认可能不回滚事务？", "Spring 声明式事务默认对 RuntimeException 和 Error 回滚，受检异常需通过 rollbackFor 等规则显式配置。", 2, ["Spring", "事务"]],
    ["Spring 事件机制适合什么场景？", "ApplicationEvent 可在进程内解耦发布者和监听器，适合非核心通知；异步监听还需处理线程池、事务边界和失败重试。", 2, ["Spring", "事件"]],
    ["如何统一处理 Spring MVC 参数校验错误？", "使用 Bean Validation 注解声明规则，在控制器启用校验，并通过 ControllerAdvice 统一转换异常响应。", 1, ["Spring MVC", "校验"]],
    ["Spring Profile 和外部化配置如何配合？", "Profile 选择环境相关 Bean，配置文件、环境变量和命令行参数提供外部属性，优先级决定最终值。", 2, ["Spring Boot", "配置"]],
    ["Spring Boot Actuator 能提供哪些能力？", "Actuator 暴露健康、指标、环境和线程等运维端点，生产环境应控制暴露范围并配置鉴权。", 1, ["Spring Boot", "Actuator"]],
  ],
  "bank-go-concurrency": [
    ["Go 调度器中的 G、M、P 分别是什么？", "G 表示 goroutine，M 表示操作系统线程，P 持有执行 Go 代码所需资源并调度本地队列中的 G。", 2, ["Go", "GMP"]],
    ["Channel 应该由谁关闭？", "通常由明确知道不会再发送数据的发送方关闭，接收方不应随意关闭；关闭不是必须操作，也不能重复关闭。", 1, ["Go", "Channel"]],
    ["select 的 default 分支有什么影响？", "存在 default 时，如果其他分支都未就绪，select 会立即执行 default，可能形成忙轮询，需要配合阻塞或退避。", 2, ["Go", "select"]],
    ["sync.Mutex 和 sync.RWMutex 如何选择？", "RWMutex 允许并发读但写互斥，只有读多写少且临界区足够长时才可能受益，否则 Mutex 更简单。", 2, ["Go", "Mutex"]],
    ["WaitGroup 使用时有哪些注意事项？", "在启动 goroutine 前调用 Add，每个任务完成后 Done，等待方调用 Wait；不应在前一轮 Wait 未结束时复用计数。", 1, ["Go", "WaitGroup"]],
    ["Go 原子操作适合什么场景？", "atomic 适合简单计数、标志和指针更新，复杂不变量更适合互斥锁；必须保证所有访问都遵循同一同步协议。", 2, ["Go", "atomic"]],
    ["race detector 如何帮助发现并发问题？", "使用 -race 构建或测试时会插桩监控内存访问，报告没有同步关系的并发读写，但只能发现实际执行路径。", 2, ["Go", "Race"]],
    ["context 取消信号如何在 goroutine 间传播？", "父 Context 取消后 Done Channel 关闭，子 Context 和监听它的 goroutine 可及时退出并释放资源。", 2, ["Go", "Context"]],
    ["什么是 fan-in 和 fan-out 并发模式？", "fan-out 将任务分发给多个工作者并行处理，fan-in 把多个结果流合并到一个 Channel。", 2, ["Go", "并发模式"]],
    ["如何实现有界 Worker Pool？", "固定数量的 worker 从任务 Channel 取任务，通过结果 Channel 返回结果，并结合 Context、关闭顺序和 WaitGroup 管理退出。", 2, ["Go", "Worker Pool"]],
    ["sync.Once 适合解决什么问题？", "sync.Once 保证函数在并发调用下最多执行一次，常用于延迟初始化，但函数 panic 后也被视为已经执行。", 1, ["Go", "sync.Once"]],
    ["sync.Pool 的使用边界是什么？", "sync.Pool 用于临时对象复用以减轻分配压力，池中对象可能在 GC 时被清除，不能保存必须持久存在的状态。", 2, ["Go", "sync.Pool"]],
  ],
  "bank-distributed-system": [
    ["强一致性和最终一致性有什么区别？", "强一致性要求读取立即看到最新成功写入；最终一致性允许短暂差异，但在没有新写入时各副本最终收敛。", 1, ["分布式", "一致性"]],
    ["Quorum 读写如何工作？", "在 N 个副本中要求写入 W 个、读取 R 个，当 R+W>N 时读写集合必有交集，但仍需版本比较处理并发写。", 3, ["分布式", "Quorum"]],
    ["租约 Lease 与普通锁有什么区别？", "租约带明确有效期，持有者需续期；失联后资源可自动回收，但必须处理时钟、暂停和旧持有者继续工作的风险。", 3, ["分布式", "Lease"]],
    ["Fencing Token 为什么能防止过期持有者写入？", "每次获得锁都分配递增令牌，下游只接受比已见令牌更新的请求，从而拒绝旧租约持有者的迟到写入。", 3, ["分布式", "Fencing Token"]],
    ["两阶段提交 2PC 有什么问题？", "协调者先收集准备结果再决定提交，能实现原子性，但存在同步阻塞、协调者故障和长时间占用资源等问题。", 3, ["分布式事务", "2PC"]],
    ["TCC 和 Saga 如何选择？", "TCC 需要业务提供 Try、Confirm、Cancel，控制力强；Saga 把长事务拆为本地事务和补偿，适合长流程但可能暴露中间状态。", 3, ["分布式事务", "Saga"]],
    ["分布式系统为什么需要指数退避和抖动？", "指数退避降低持续失败时的请求压力，随机抖动避免大量客户端同时重试形成惊群和周期性峰值。", 2, ["分布式", "重试"]],
    ["雪花算法生成 ID 的基本结构是什么？", "通常将时间戳、机器标识和同毫秒序列组合为整数，能本地生成趋势递增 ID，但要处理时钟回拨和机器编号。", 2, ["分布式", "唯一ID"]],
    ["一致性哈希解决了什么问题？", "一致性哈希把节点和键映射到环上，节点增删时只迁移邻近区间的数据，虚拟节点可改善分布均匀性。", 2, ["分布式", "一致性哈希"]],
    ["逻辑时钟的作用是什么？", "Lamport 时钟能表达事件的先后因果约束，向量时钟还能识别并发事件，但它们不等同于物理时间。", 3, ["分布式", "逻辑时钟"]],
    ["什么是脑裂，如何降低风险？", "网络分区可能让多个节点都认为自己是主节点，可通过多数派仲裁、租约、Fencing Token 和隔离旧主降低双写风险。", 3, ["分布式", "脑裂"]],
    ["分布式系统中的反熵是什么？", "副本通过摘要、Merkle Tree 或版本信息定期比较并修复差异，使偶发丢失或离线副本最终收敛。", 3, ["分布式", "反熵"]],
  ],
  "bank-microservice": [
    ["API Gateway 通常承担哪些职责？", "网关统一处理路由、认证、限流、协议转换和观测，但应避免堆积过多业务逻辑形成新的单体瓶颈。", 1, ["微服务", "网关"]],
    ["配置中心需要考虑哪些可靠性问题？", "配置应版本化、可审计和回滚，客户端要缓存最后可用值，并处理推送失败、灰度范围和敏感信息保护。", 2, ["微服务", "配置中心"]],
    ["客户端负载均衡和服务端负载均衡有什么区别？", "客户端从服务发现结果中选实例；服务端由独立代理或负载均衡器转发，两者在耦合、观测和流量控制上不同。", 2, ["微服务", "负载均衡"]],
    ["健康检查为什么要区分存活和就绪？", "存活检查判断进程是否需要重启，就绪检查判断实例是否能接收流量，依赖未准备好时不应立即重启进程。", 1, ["微服务", "健康检查"]],
    ["常见限流算法有哪些？", "固定窗口简单但有边界突发，滑动窗口更平滑，漏桶控制输出速率，令牌桶允许一定突发流量。", 2, ["微服务", "限流"]],
    ["熔断器通常有哪些状态？", "关闭状态正常放行，失败达到阈值后进入打开状态快速失败，等待后进入半开状态放少量探测请求决定恢复。", 2, ["微服务", "熔断"]],
    ["舱壁隔离模式解决什么问题？", "为不同依赖或业务分配独立线程池、连接池或并发配额，避免单个故障耗尽全部资源并拖垮系统。", 2, ["微服务", "隔离"]],
    ["调用下游服务时为什么不能无条件重试？", "重试会放大故障流量且可能重复执行非幂等操作，应限定错误类型、次数、总超时并配合退避和幂等。", 2, ["微服务", "重试"]],
    ["微服务 API 如何保持兼容性？", "优先做向后兼容的字段扩展，对破坏性变化进行版本化，并通过契约测试和分阶段迁移保护调用方。", 2, ["微服务", "API"]],
    ["日志、指标和链路追踪如何互补？", "指标发现趋势和告警，追踪定位跨服务请求路径，日志提供具体上下文，三者通过统一标识关联。", 2, ["微服务", "可观测性"]],
    ["灰度发布需要哪些保护措施？", "应按用户或流量比例逐步放量，持续比较错误率和延迟，设置自动停止与快速回滚，并兼容新旧数据格式。", 2, ["微服务", "灰度发布"]],
    ["Service Mesh 解决了什么问题？", "Service Mesh 将服务间流量治理、身份和观测下沉到代理与控制面，减少语言 SDK 重复实现，但增加基础设施复杂度。", 3, ["微服务", "Service Mesh"]],
  ],
  "bank-message-queue": [
    ["消息投递语义有哪些？", "至多一次可能丢失但不重复，至少一次允许重复但尽量不丢失，恰好一次通常需要 Broker 与业务状态共同限定边界。", 2, ["消息队列", "投递语义"]],
    ["消费者如何处理重复消息？", "使用业务唯一键、去重表或状态机把处理设计为幂等，并让重复请求返回已完成结果。", 2, ["消息队列", "幂等"]],
    ["如何保证同一业务键的消息有序？", "将相同业务键路由到同一分区或队列，并由单个消费序列处理；失败重试也不能绕过前序消息。", 3, ["消息队列", "顺序消息"]],
    ["延迟消息有哪些实现方式？", "可使用 Broker 原生延迟能力、分级时间轮或有序集合轮询；需要考虑精度、堆积、取消和重复投递。", 2, ["消息队列", "延迟消息"]],
    ["死信队列有什么用途？", "超过重试次数、过期或无法路由的消息进入死信队列，便于隔离故障、告警和人工或自动补偿。", 1, ["消息队列", "死信队列"]],
    ["消息积压时应该如何处理？", "先定位生产突增或消费变慢原因，再临时扩容消费者、提高分区并行度或降级生产，同时避免无序扩容压垮下游。", 2, ["消息队列", "消息积压"]],
    ["Kafka 的 ISR 是什么？", "ISR 是与 Leader 保持足够同步的副本集合，Leader 故障时优先从 ISR 选新 Leader，以降低数据丢失风险。", 2, ["Kafka", "ISR"]],
    ["Kafka Consumer Offset 应在何时提交？", "至少一次语义下应在业务处理成功后提交；自动提交可能在处理完成前推进位点，导致失败时丢失重试机会。", 2, ["Kafka", "Offset"]],
    ["Consumer Group Rebalance 会带来什么影响？", "分区重新分配期间消费会暂停，未正确提交的进度可能导致重复；应控制成员波动并使用增量协作再均衡等机制。", 3, ["Kafka", "Rebalance"]],
    ["RabbitMQ Exchange 有哪些常见类型？", "direct 按精确路由键，topic 按模式，fanout 广播到绑定队列，headers 根据消息头匹配。", 1, ["RabbitMQ", "Exchange"]],
    ["Transactional Outbox 模式如何工作？", "业务数据和待发事件在同一本地事务写入，独立投递器可靠发布事件并标记状态，消费者仍需幂等。", 3, ["消息队列", "Outbox"]],
    ["Kafka 的幂等生产者解决什么问题？", "生产者为批次携带序列号，Broker 可识别同一会话中的重复写入，降低重试导致的分区内重复消息。", 3, ["Kafka", "幂等生产者"]],
  ],
  "bank-docker-kubernetes": [
    ["容器隔离主要依赖哪些 Linux 能力？", "Namespace 隔离进程、网络和挂载视图，cgroup 限制和统计资源，Capability 等机制进一步收缩权限。", 2, ["Docker", "Namespace"]],
    ["Dockerfile 分层对构建有什么影响？", "每条相关指令形成可缓存层，稳定步骤应放在前面，减少无关文件进入上下文可以提高缓存命中并缩小镜像。", 1, ["Docker", "Dockerfile"]],
    ["多阶段构建有什么好处？", "在构建阶段安装编译工具和依赖，只把产物复制到精简运行镜像，可降低镜像大小和攻击面。", 1, ["Docker", "多阶段构建"]],
    ["Docker Volume 和 Bind Mount 有什么区别？", "Volume 由 Docker 管理生命周期和位置；Bind Mount 直接映射宿主机路径，灵活但更依赖主机环境。", 1, ["Docker", "存储"]],
    ["Kubernetes requests 和 limits 分别有什么作用？", "调度器依据 requests 安排 Pod，运行时通过 limits 限制最大资源；CPU 超限会节流，内存超限可能被 OOMKill。", 2, ["Kubernetes", "资源"]],
    ["Liveness、Readiness 和 Startup Probe 有什么区别？", "存活探针决定是否重启，就绪探针决定是否接流量，启动探针为慢启动应用提供更长初始化窗口。", 2, ["Kubernetes", "Probe"]],
    ["Kubernetes Service 有哪些常见类型？", "ClusterIP 提供集群内访问，NodePort 暴露节点端口，LoadBalancer 请求外部负载均衡，ExternalName 提供 DNS 别名。", 1, ["Kubernetes", "Service"]],
    ["Ingress 和 Service 的职责有什么区别？", "Service 为一组 Pod 提供稳定四层入口，Ingress 通过控制器提供 HTTP 路由、域名和 TLS 等七层能力。", 2, ["Kubernetes", "Ingress"]],
    ["ConfigMap 和 Secret 应该如何使用？", "ConfigMap 保存非敏感配置，Secret 保存敏感字节但默认仅编码不等于加密，仍需 RBAC、静态加密和外部密钥方案。", 2, ["Kubernetes", "配置"]],
    ["HPA 根据什么扩缩容？", "HPA 根据 CPU、内存或自定义指标调整副本数，需要合理 requests、稳定指标和扩缩容行为配置。", 2, ["Kubernetes", "HPA"]],
    ["StatefulSet 适合什么工作负载？", "StatefulSet 为副本提供稳定名称、有序部署和独立持久卷，适合数据库等有状态服务。", 2, ["Kubernetes", "StatefulSet"]],
    ["Pod 一直处于 CrashLoopBackOff 应如何排查？", "查看容器日志、退出码和事件，检查命令、配置、探针、依赖与资源限制，并用前一次容器日志定位重启前错误。", 2, ["Kubernetes", "排障"]],
  ],
  "bank-linux": [
    ["Linux 文件权限中的 755 表示什么？", "所有者拥有读写执行权限，组用户和其他用户拥有读与执行权限；每一位由读 4、写 2、执行 1 相加。", 1, ["Linux", "权限"]],
    ["硬链接和软链接有什么区别？", "硬链接指向同一 inode，通常不能跨文件系统；软链接是保存目标路径的独立文件，可跨文件系统但目标删除后会失效。", 1, ["Linux", "文件系统"]],
    ["Linux 常见进程状态有哪些？", "常见有运行或就绪 R、可中断睡眠 S、不可中断睡眠 D、停止 T 和僵尸 Z，状态有助于判断 CPU 或 IO 问题。", 1, ["Linux", "进程"]],
    ["SIGTERM 和 SIGKILL 有什么区别？", "SIGTERM 可被进程捕获并执行优雅退出；SIGKILL 不能被捕获或忽略，由内核立即终止进程。", 1, ["Linux", "信号"]],
    ["systemd 如何管理服务？", "通过 Unit 文件描述启动命令、依赖、重启策略和环境，systemctl 用于启动、停止、查看状态和设置开机启动。", 2, ["Linux", "systemd"]],
    ["free 命令中的 available 为什么比 free 更有参考价值？", "Linux 会利用空闲内存做页缓存，available 估算无需交换即可提供给新进程的内存，更接近实际余量。", 1, ["Linux", "内存"]],
    ["vmstat 和 iostat 分别适合观察什么？", "vmstat 汇总进程、内存、换页和 CPU，iostat 重点展示块设备吞吐、队列和等待时间。", 2, ["Linux", "性能"]],
    ["磁盘显示有空间却无法创建文件可能是什么原因？", "可能是 inode 用尽、目录或用户配额达到上限、文件系统只读，或进程权限和保留块限制。", 2, ["Linux", "磁盘"]],
    ["Linux Page Cache 的作用是什么？", "内核缓存文件数据以减少磁盘 IO，读写通常先经过页缓存，脏页再由后台或同步操作回写。", 2, ["Linux", "Page Cache"]],
    ["如何查看 Linux TCP 连接状态分布？", "可使用 ss -ant 汇总 ESTABLISHED、TIME-WAIT 等状态，再结合端口、进程和内核指标定位连接异常。", 1, ["Linux", "TCP"]],
    ["journalctl 常用于哪些排查场景？", "journalctl 查询 systemd journal，可按服务、启动批次、时间和优先级过滤系统及服务日志。", 1, ["Linux", "日志"]],
    ["Shell 管道和重定向有什么区别？", "管道把前一命令的标准输出连接到后一命令的标准输入；重定向改变命令输入输出对应的文件或描述符。", 1, ["Linux", "Shell"]],
  ],
  "bank-algorithm-data-structure": [
    ["数组和链表的主要区别是什么？", "数组支持 O(1) 随机访问但中间插删需移动元素；链表定位为 O(n)，已知节点后插删可为 O(1)。", 1, ["数据结构", "数组"]],
    ["栈和队列分别适合哪些场景？", "栈后进先出，适合调用、括号匹配和回溯；队列先进先出，适合任务调度和广度优先搜索。", 1, ["数据结构", "栈"]],
    ["二叉树的前中后序遍历顺序是什么？", "前序为根左右，中序为左根右，后序为左右根；可使用递归或显式栈实现。", 1, ["数据结构", "二叉树"]],
    ["二叉搜索树的查找复杂度是多少？", "平均情况下为 O(log n)，但树退化成链表时为 O(n)；平衡树通过旋转限制高度。", 2, ["数据结构", "BST"]],
    ["堆适合解决哪些问题？", "堆能在 O(1) 查看极值并在 O(log n) 插入删除，适合优先队列、Top K 和堆排序。", 2, ["数据结构", "堆"]],
    ["并查集如何优化？", "路径压缩让查找后的节点直接靠近根，按秩或大小合并避免高树，两者结合后摊还复杂度接近常数。", 2, ["数据结构", "并查集"]],
    ["BFS 和 DFS 如何选择？", "BFS 适合无权图最短路径和分层遍历，DFS 适合回溯、连通性和拓扑相关探索，空间特征也不同。", 1, ["算法", "图"]],
    ["二分查找最常见的边界错误是什么？", "需要统一闭区间或半开区间定义，并相应更新左右边界；查找首尾位置时还要在命中后继续收缩。", 2, ["算法", "二分查找"]],
    ["拓扑排序适用于什么问题？", "拓扑排序为有向无环图生成满足依赖先后的序列，可用入度队列或 DFS 实现，并能检测环。", 2, ["算法", "拓扑排序"]],
    ["LRU 缓存如何做到 O(1) 访问和淘汰？", "哈希表负责 O(1) 定位节点，双向链表维护新旧顺序，访问后移到头部，淘汰尾部节点。", 2, ["数据结构", "LRU"]],
    ["Trie 前缀树适合什么场景？", "Trie 按字符路径共享前缀，适合前缀检索和词典匹配，时间与字符串长度相关但可能占用较多内存。", 2, ["数据结构", "Trie"]],
    ["如何判断一个算法是否适合使用贪心？", "需要证明每一步局部最优选择能扩展为全局最优，常通过交换论证或最优子结构证明，不能仅凭直觉。", 3, ["算法", "贪心"]],
  ],
};

questionGroups.forEach((group) => {
  const additions = supplementalQuestions[group.bank_id];
  if (!additions) {
    throw new Error(`缺少 ${group.bank_id} 的扩充题目`);
  }
  group.questions.push(...additions);
  if (group.questions.length !== 15) {
    throw new Error(`${group.bank_id} 题目数量应为 15，实际为 ${group.questions.length}`);
  }
});

const requiredBankIds = questionGroups.map((group) => group.bank_id);
const existingBankIds = db.question_bank.find(
  { bank_id: { $in: requiredBankIds } },
  { _id: 0, bank_id: 1 },
).toArray().map((bank) => bank.bank_id);
const missingBankIds = requiredBankIds.filter((bankId) => !existingBankIds.includes(bankId));

if (missingBankIds.length > 0) {
  throw new Error(
    `缺少题库数据：${missingBankIds.join(", ")}。请先执行 insert_bank_and_series.js`,
  );
}

let questionOrder = 0;
const questions = [];

questionGroups.forEach((group, groupIndex) => {
  group.questions.forEach(([title, content, difficulty, tags], index) => {
    questionOrder += 1;
    questions.push({
      question_id: `question-${String(questionOrder).padStart(3, "0")}`,
      bank_list: [group.bank_id],
      job_name: jobName,
      title,
      content,
      difficulty: NumberInt(difficulty),
      tags,
      status: NumberInt(1),
      vip: questionOrder % 7 === 0,
      hot_degree: NumberInt(2500 - groupIndex * 100 - index * 4),
      view_count: NumberInt(20000 - questionOrder * 47),
      thumbs_up_count: NumberInt(2400 - questionOrder * 5),
      dislike_count: NumberInt(questionOrder % 19),
      order: int64(index + 1),
      create_time: now,
      update_time: now,
    });
  });
});

questions.push(
  {
    question_id: "question-241",
    bank_list: ["bank-redis", "bank-distributed-system"],
    job_name: jobName,
    title: "分布式锁应该如何正确实现？",
    content: "应使用唯一持有者标识、原子加锁和带校验的原子释放，并设置合理过期时间。对强一致性要求高的场景还要考虑租约续期、故障切换和 fencing token。",
    difficulty: NumberInt(3),
    tags: ["Redis", "分布式锁", "幂等"],
    status: NumberInt(1),
    vip: true,
    hot_degree: NumberInt(3000),
    view_count: NumberInt(23600),
    thumbs_up_count: NumberInt(2680),
    dislike_count: NumberInt(7),
    order: int64(16),
    create_time: now,
    update_time: now,
  },
  {
    question_id: "question-242",
    bank_list: ["bank-mysql", "bank-distributed-system"],
    job_name: jobName,
    title: "本地消息表如何保证最终一致性？",
    content: "业务数据和消息记录在同一本地事务中提交，后台任务可靠投递消息并更新状态，消费端通过幂等处理应对重复投递。",
    difficulty: NumberInt(3),
    tags: ["MySQL", "分布式事务", "最终一致性"],
    status: NumberInt(1),
    vip: true,
    hot_degree: NumberInt(2950),
    view_count: NumberInt(22800),
    thumbs_up_count: NumberInt(2510),
    dislike_count: NumberInt(5),
    order: int64(16),
    create_time: now,
    update_time: now,
  },
);

const questionIds = questions.map((question) => question.question_id);
if (questions.length !== 242 || new Set(questionIds).size !== questions.length) {
  throw new Error(
    `题目数据校验失败：总数 ${questions.length}，唯一 ID 数 ${new Set(questionIds).size}`,
  );
}
if (
  questions.some(
    (question) =>
      !question.question_id ||
      !question.title ||
      !question.content ||
      question.bank_list.length === 0 ||
      ![1, 2, 3].includes(question.difficulty.valueOf()),
  )
) {
  throw new Error("题目数据校验失败：存在必填字段为空或难度值非法的题目");
}

db.question.drop();
print("已清空旧题目数据");

const result = db.question.insertMany(questions);
db.question.createIndex({ question_id: 1 }, { unique: true });
db.question.createIndex({ status: 1, order: 1 });
db.question.createIndex({ bank_list: 1, status: 1 });
db.question.createIndex({ job_name: 1, status: 1, hot_degree: -1 });

questionGroups.forEach((group) => {
  print(`Inserted ${group.bank_id}: ${group.questions.length}`);
});
print(`Inserted questions: ${Object.keys(result.insertedIds).length}`);
print("题目数据插入完成！");
} catch (error) {
  print(`题目数据插入失败：${error.message}`);
  quit(1);
}
