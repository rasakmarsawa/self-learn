import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'

import { useAuth } from './context/AuthContext'

import './App.css'

function App() {
  const { accessToken } = useAuth()

  if (!accessToken) {
    return <LoginPage />
  }

  return <DashboardPage />
}

export default App