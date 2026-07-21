/**
 * 插入题库系列和题库测试数据。
 *
 * 使用方法：
 * docker exec -i mongodb-dev mongosh < testdata/insert_bank_and_series.js
 */

use("offer_hub");

// 清空旧数据，保证脚本可以安全重复执行。
db.question_bank_series.drop();
db.question_bank.drop();
print("已清空旧数据");

const now = new Date();
const jobName = "后端开发";
const int64 = (value) => NumberLong(String(value));

const series = [
  {
    series_id: "series-computer-fundamentals",
    series_name: "计算机基础",
    job_name: jobName,
    order: int64(1),
    create_time: now,
    update_time: now,
  },
  {
    series_id: "series-database-cache",
    series_name: "数据库与缓存",
    job_name: jobName,
    order: int64(2),
    create_time: now,
    update_time: now,
  },
  {
    series_id: "series-programming-language",
    series_name: "编程语言",
    job_name: jobName,
    order: int64(3),
    create_time: now,
    update_time: now,
  },
  {
    series_id: "series-java-ecosystem",
    series_name: "Java 生态",
    job_name: jobName,
    order: int64(4),
    create_time: now,
    update_time: now,
  },
  {
    series_id: "series-go-ecosystem",
    series_name: "Golang 生态",
    job_name: jobName,
    order: int64(5),
    create_time: now,
    update_time: now,
  },
  {
    series_id: "series-distributed-architecture",
    series_name: "分布式与微服务",
    job_name: jobName,
    order: int64(6),
    create_time: now,
    update_time: now,
  },
  {
    series_id: "series-middleware",
    series_name: "中间件",
    job_name: jobName,
    order: int64(7),
    create_time: now,
    update_time: now,
  },
  {
    series_id: "series-engineering",
    series_name: "工程化与运维",
    job_name: jobName,
    order: int64(8),
    create_time: now,
    update_time: now,
  },
  {
    series_id: "series-algorithm",
    series_name: "算法与数据结构",
    job_name: jobName,
    order: int64(9),
    create_time: now,
    update_time: now,
  },
];

const banks = [
  {
    bank_id: "bank-computer-network",
    series_id: "series-computer-fundamentals",
    bank_name: "计算机网络题库",
    bank_logo: "/assets/question-bank/computer-network.png",
    desc: "TCP/IP、HTTP、DNS 等计算机网络基础",
    job_name: jobName,
    order: int64(1),
  },
  {
    bank_id: "bank-operating-system",
    series_id: "series-computer-fundamentals",
    bank_name: "操作系统题库",
    bank_logo: "/assets/question-bank/operating-system.png",
    desc: "进程、线程、内存与文件系统",
    job_name: jobName,
    order: int64(2),
  },
  {
    bank_id: "bank-mysql",
    series_id: "series-database-cache",
    bank_name: "MySQL 题库",
    bank_logo: "/assets/question-bank/mysql.png",
    desc: "索引、事务、锁与查询优化",
    job_name: jobName,
    order: int64(1),
  },
  {
    bank_id: "bank-redis",
    series_id: "series-database-cache",
    bank_name: "Redis 题库",
    bank_logo: "/assets/question-bank/redis.png",
    desc: "数据结构、持久化、缓存与集群",
    job_name: jobName,
    order: int64(2),
  },
  {
    bank_id: "bank-mongodb",
    series_id: "series-database-cache",
    bank_name: "MongoDB 题库",
    bank_logo: "/assets/question-bank/mongodb.png",
    desc: "文档模型、索引、聚合与副本集",
    job_name: jobName,
    order: int64(3),
  },
  {
    bank_id: "bank-java",
    series_id: "series-programming-language",
    bank_name: "Java 基础题库",
    bank_logo: "/assets/question-bank/java.png",
    desc: "Java 语法、集合与并发基础",
    job_name: jobName,
    order: int64(1),
  },
  {
    bank_id: "bank-go",
    series_id: "series-programming-language",
    bank_name: "Golang 基础题库",
    bank_logo: "/assets/question-bank/go.png",
    desc: "Go 语言、接口、并发与运行时",
    job_name: jobName,
    order: int64(2),
  },
  {
    bank_id: "bank-jvm",
    series_id: "series-java-ecosystem",
    bank_name: "JVM 题库",
    bank_logo: "/assets/question-bank/jvm.png",
    desc: "JVM 内存、类加载与垃圾回收",
    job_name: jobName,
    order: int64(1),
  },
  {
    bank_id: "bank-spring",
    series_id: "series-java-ecosystem",
    bank_name: "Spring 题库",
    bank_logo: "/assets/question-bank/spring.png",
    desc: "IoC、AOP、事务与 Spring Boot",
    job_name: jobName,
    order: int64(2),
  },
  {
    bank_id: "bank-go-concurrency",
    series_id: "series-go-ecosystem",
    bank_name: "Golang 并发题库",
    bank_logo: "/assets/question-bank/go-concurrency.png",
    desc: "Goroutine、Channel 与 G-M-P 调度",
    job_name: jobName,
    order: int64(1),
  },
  {
    bank_id: "bank-distributed-system",
    series_id: "series-distributed-architecture",
    bank_name: "分布式系统题库",
    bank_logo: "/assets/question-bank/distributed-system.png",
    desc: "一致性、选举、幂等与分布式事务",
    job_name: jobName,
    order: int64(1),
  },
  {
    bank_id: "bank-microservice",
    series_id: "series-distributed-architecture",
    bank_name: "微服务题库",
    bank_logo: "/assets/question-bank/microservice.png",
    desc: "服务治理、限流、熔断与链路追踪",
    job_name: jobName,
    order: int64(2),
  },
  {
    bank_id: "bank-message-queue",
    series_id: "series-middleware",
    bank_name: "消息队列题库",
    bank_logo: "/assets/question-bank/message-queue.png",
    desc: "Kafka、RabbitMQ 与消息可靠性",
    job_name: jobName,
    order: int64(1),
  },
  {
    bank_id: "bank-docker-kubernetes",
    series_id: "series-engineering",
    bank_name: "Docker 与 Kubernetes 题库",
    bank_logo: "/assets/question-bank/docker-kubernetes.png",
    desc: "容器、镜像、Pod 与服务编排",
    job_name: jobName,
    order: int64(1),
  },
  {
    bank_id: "bank-linux",
    series_id: "series-engineering",
    bank_name: "Linux 题库",
    bank_logo: "/assets/question-bank/linux.png",
    desc: "Linux 命令、进程、网络与性能排查",
    job_name: jobName,
    order: int64(2),
  },
  {
    bank_id: "bank-algorithm-data-structure",
    series_id: "series-algorithm",
    bank_name: "算法与数据结构题库",
    bank_logo: "/assets/question-bank/algorithm.png",
    desc: "常见数据结构、排序、搜索与动态规划",
    job_name: jobName,
    order: int64(1),
  },
].map((bank) => ({
  ...bank,
  create_time: now,
  update_time: now,
}));

const seriesResult = db.question_bank_series.insertMany(series);
const bankResult = db.question_bank.insertMany(banks);

db.question_bank_series.createIndex({ series_id: 1 }, { unique: true });
db.question_bank.createIndex({ bank_id: 1 }, { unique: true });
db.question_bank.createIndex({ series_id: 1, order: 1 });

print(`Inserted series: ${Object.keys(seriesResult.insertedIds).length}`);
print(`Inserted banks: ${Object.keys(bankResult.insertedIds).length}`);
print("数据插入完成！");
