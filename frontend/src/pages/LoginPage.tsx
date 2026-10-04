import { useState } from 'react'

import { login } from '../api/auth'
import { useAuth } from '../context/AuthContext'

import '../App.css'

function LoginPage() {
  const { login: loginUser } = useAuth()

  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    setError('')
    setLoading(true)

    try {
      const data = await login(username, password)

      console.log('Message:', data.message)
      console.log('Access Token:', data.AccessToken)

      loginUser(data.AccessToken)
    } catch (error) {
      if (error instanceof Error) {
        setError(error.message)
      } else {
        setError('Login failed')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="login">
      <form className="login-card" onSubmit={handleSubmit}>
        <h1>Login</h1>

        <p className="login-subtitle">
          Enter your credentials to continue
        </p>

        <div className="field">
          <label htmlFor="username">Username</label>

          <input
            id="username"
            name="username"
            type="text"
            autoComplete="username"
            value={username}
            onChange={(event) => setUsername(event.target.value)}
          />
        </div>

        <div className="field">
          <label htmlFor="password">Password</label>

          <input
            id="password"
            name="password"
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
        </div>

        {error && <p>{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </main>
  )
}

export default LoginPage