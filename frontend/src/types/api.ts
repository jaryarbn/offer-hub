/** 后端统一响应结构；业务失败时 data 为 null。 */
export interface ApiResponse<T> {
  code: number
  msg: string
  data: T | null
}
