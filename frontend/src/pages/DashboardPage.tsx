import { useState } from 'react'

import Navbar from '../components/Navbar'
import { authenticatedHello } from '../api/auth'
import { useAuth } from '../context/AuthContext'

import '../App.css'

const stats = [
  { label: 'Queued messages', value: '1,284' },
  { label: 'Active workers', value: '6' },
  { label: 'Published today', value: '9,431' },
  { label: 'Failed deliveries', value: '12' },
]

function DashboardPage() {
  const { accessToken, login, logout } = useAuth()

  const [response, setResponse] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleAuthenticatedHello() {
    setLoading(true)
    setResponse('')

    try {
      const result = await authenticatedHello(accessToken)

      setResponse(JSON.stringify(result.data))

      if (result.accessToken !== accessToken) {
        login(result.accessToken)
      }
    } catch (error) {
      if (error instanceof Error) {
        setResponse(error.message)
      } else {
        setResponse('Request failed')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="dashboard">
      <Navbar username="Admin" onLogout={logout} />

      <main className="dashboard-body">
        <h1>Dashboard</h1>

        <p className="dashboard-subtitle">
          Overview of your message queue
        </p>

        <div className="stat-grid">
          {stats.map((stat) => (
            <div key={stat.label} className="stat-card">
              <span className="stat-value">{stat.value}</span>
              <span className="stat-label">{stat.label}</span>
            </div>
          ))}
        </div>

        <div className="panel">
          <h2>Authenticated API</h2>

          <button
            type="button"
            onClick={handleAuthenticatedHello}
            disabled={loading}
          >
            {loading ? 'Calling API...' : 'Call Authenticated Hello'}
          </button>

          {response && <p>{response}</p>}
        </div>

        <div className="panel">
          <h2>Recent messages</h2>

          <table className="table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Queue</th>
                <th>Status</th>
                <th>Received</th>
              </tr>
            </thead>

            <tbody>
              <tr>
                <td>msg_10241</td>
                <td>my_queue</td>
                <td>
                  <span className="badge badge-ok">Delivered</span>
                </td>
                <td>2 minutes ago</td>
              </tr>

              <tr>
                <td>msg_10240</td>
                <td>my_queue</td>
                <td>
                  <span className="badge badge-ok">Delivered</span>
                </td>
                <td>5 minutes ago</td>
              </tr>

              <tr>
                <td>msg_10239</td>
                <td>my_queue</td>
                <td>
                  <span className="badge badge-warn">Retrying</span>
                </td>
                <td>11 minutes ago</td>
              </tr>
            </tbody>
          </table>
        </div>
      </main>
    </div>
  )
}

export default DashboardPage