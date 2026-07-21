export interface ListQuestionParams {
  bank_id?: string
  keyword?: string
  difficulty?: number
  tags?: string[]
  job_name?: string
  sort_by?: string
  sort_order?: "asc" | "desc"
  page?: number
  page_size?: number
}

export interface Question {
  question_id: string
  bank_list: string[]
  title: string
  content: string
  difficulty: number
  tags: string[]
  status: number
  vip: boolean
  hot_degree: number
  view_count: number
  thumbs_up_count: number
  dislike_count: number
  order: number
  user_tag: number
  user_liked: boolean
  create_time: string
  update_time: string
}

export interface ListQuestionResponseData {
  total: number
  list: Question[]
}
