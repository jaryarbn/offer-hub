/**
 * 插入题目测试数据。
 *
 * 使用方法：
 * docker exec -i mongodb-dev mongosh < testdata/insert_question.js
 */

use("offer_hub");

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

db.question.drop();
print("已清空旧题目数据");

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
      hot_degree: NumberInt(100 - groupIndex * 3 - index),
      view_count: NumberInt(3200 - questionOrder * 37),
      thumbs_up_count: NumberInt(520 - questionOrder * 6),
      dislike_count: NumberInt(questionOrder % 9),
      order: int64(index + 1),
      create_time: now,
      update_time: now,
    });
  });
});

questions.push(
  {
    question_id: "question-049",
    bank_list: ["bank-redis", "bank-distributed-system"],
    job_name: jobName,
    title: "分布式锁应该如何正确实现？",
    content: "应使用唯一持有者标识、原子加锁和带校验的原子释放，并设置合理过期时间。对强一致性要求高的场景还要考虑租约续期、故障切换和 fencing token。",
    difficulty: NumberInt(3),
    tags: ["Redis", "分布式锁", "幂等"],
    status: NumberInt(1),
    vip: true,
    hot_degree: NumberInt(99),
    view_count: NumberInt(2860),
    thumbs_up_count: NumberInt(486),
    dislike_count: NumberInt(7),
    order: int64(4),
    create_time: now,
    update_time: now,
  },
  {
    question_id: "question-050",
    bank_list: ["bank-mysql", "bank-distributed-system"],
    job_name: jobName,
    title: "本地消息表如何保证最终一致性？",
    content: "业务数据和消息记录在同一本地事务中提交，后台任务可靠投递消息并更新状态，消费端通过幂等处理应对重复投递。",
    difficulty: NumberInt(3),
    tags: ["MySQL", "分布式事务", "最终一致性"],
    status: NumberInt(1),
    vip: true,
    hot_degree: NumberInt(97),
    view_count: NumberInt(2540),
    thumbs_up_count: NumberInt(442),
    dislike_count: NumberInt(5),
    order: int64(4),
    create_time: now,
    update_time: now,
  },
);

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
