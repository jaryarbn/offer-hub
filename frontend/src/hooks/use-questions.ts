import { useQuery } from "@tanstack/react-query"

import { listQuestions } from "@/services/question"
import type { ListQuestionParams } from "@/types/question"

export function useQuestions(params: ListQuestionParams) {
  return useQuery({
    queryKey: ["questions", params],
    queryFn: () => listQuestions(params),
  })
}
