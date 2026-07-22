import { Route, Routes } from 'react-router-dom'

import { AppProvider } from '@/components/provider/app-provider'
import { HomePage } from '@/pages/home-page'
import { NotFoundPage } from '@/pages/not-found-page'
import { QuestionCollections } from '@/pages/QuestionCollections'
import { QuestionDetail } from '@/pages/QuestionDetail'
import { Questions } from '@/pages/Questions'

export default function App() {
  return (
    <AppProvider>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/questions-collection" element={<QuestionCollections />} />
        <Route path="/questions" element={<Questions />} />
        <Route path="/questions/:question_id" element={<QuestionDetail />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </AppProvider>
  )
}
