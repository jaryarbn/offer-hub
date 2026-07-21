import axios from "axios"

const TOKEN_KEY = "token"

const apiClient = axios.create({
  baseURL: "",
  timeout: 10_000,
})

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY)

  if (token) {
    config.headers.set("Authorization", `Bearer ${token}`)
  }

  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  (error: unknown) => {
    if (axios.isAxiosError(error)) {
      if (error.response?.status === 401) {
        localStorage.removeItem(TOKEN_KEY)
      }

      return Promise.reject(error.response?.data ?? error)
    }

    return Promise.reject(error)
  },
)

export { apiClient }
export default apiClient
