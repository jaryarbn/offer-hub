import apiClient from '@/lib/axios'
import type { ApiResponse } from '@/types/api'

/** 登录接口返回的 camelCase 用户 DTO。 */
export interface PasswordAuthUserDTO {
  userId: string
  username: string
  nickName: string
  avatar: string
  sex: number
  vip: boolean
  phone: string
  email: string
  userStatus: number
  userType: number
}

/** 登录接口成功响应中的 data。 */
export interface PasswordLoginDataDTO {
  token: string
  userInfo: PasswordAuthUserDTO
}

/** 登出接口成功响应中的 data。 */
export interface PasswordLogoutDataDTO {
  message: string
}

/** 用户信息接口返回的 snake_case DTO。 */
export interface UserInfoDTO {
  user_id: string
  username: string
  nick_name: string
  avatar: string
  vip: boolean
  sex: number
  phone: string
  email: string
  introduction: string
  avatar_url: string
  user_status: number
  user_type: number
  create_time: string
  update_time: string
}

/**
 * Auth Context 使用的统一 snake_case 用户结构。
 * 登录接口不返回个人简介、头像地址和时间字段，因此这些字段保持可选。
 */
export type UserInfo = Omit<
  UserInfoDTO,
  'introduction' | 'avatar_url' | 'create_time' | 'update_time'
> &
  Partial<Pick<UserInfoDTO, 'introduction' | 'avatar_url' | 'create_time' | 'update_time'>>

function assertSuccess<T>(response: ApiResponse<T>): void {
  if (response.code !== 0) {
    throw new Error(response.msg)
  }
}

function requireResponseData<T>(response: ApiResponse<T>, fallbackMessage: string): T {
  assertSuccess(response)

  if (response.data === null || response.data === undefined) {
    throw new Error(fallbackMessage)
  }

  return response.data
}

/** 将登录接口的 camelCase 用户 DTO 转换为 Auth Context 使用的 snake_case 结构。 */
export function mapPasswordAuthUserToUserInfo(user: PasswordAuthUserDTO): UserInfo {
  return {
    user_id: user.userId,
    username: user.username,
    nick_name: user.nickName,
    avatar: user.avatar,
    sex: user.sex,
    vip: user.vip,
    phone: user.phone,
    email: user.email,
    user_status: user.userStatus,
    user_type: user.userType,
  }
}

/** 使用用户名和密码注册账号。 */
export async function postAuthRegister(username: string, password: string): Promise<void> {
  const response = await apiClient.post<ApiResponse<null>>('/auth/register', {
    username,
    password,
  })

  assertSuccess(response.data)
}

/** 使用用户名和密码登录，成功时返回 token 和 camelCase userInfo。 */
export async function postAuthLogin(
  username: string,
  password: string
): Promise<PasswordLoginDataDTO> {
  const response = await apiClient.post<ApiResponse<PasswordLoginDataDTO>>('/auth/login', {
    username,
    password,
  })

  return requireResponseData(response.data, '登录响应数据为空')
}

/**
 * 登出当前账号。
 * Authorization Header 由 apiClient 的请求拦截器自动添加。
 */
export async function postAuthLogout(): Promise<PasswordLogoutDataDTO> {
  const response = await apiClient.post<ApiResponse<PasswordLogoutDataDTO>>('/auth/logout')

  return requireResponseData(response.data, '登出响应数据为空')
}

/** 获取当前用户的完整 snake_case 用户信息。 */
export async function getUserInfo(): Promise<UserInfoDTO> {
  const response = await apiClient.get<ApiResponse<UserInfoDTO>>('/api/v1/user_info/get_user_info')

  return requireResponseData(response.data, '用户信息响应数据为空')
}
