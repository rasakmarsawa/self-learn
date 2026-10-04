type LoginResponse = {
  message: string
  AccessToken: string
}

export async function login(
  username: string,
  password: string
): Promise<LoginResponse> {
  const response = await fetch('http://localhost:8081/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({
      username,
      password,
    }),
  })

  const data = await response.json()

  if (!response.ok) {
    throw new Error(data.error || 'Login failed')
  }

  return data
}

export async function authenticatedHello(
  accessToken: string,
): Promise<{
  data: unknown
  accessToken: string
}> {
  let currentAccessToken = accessToken

  let response = await callAuthenticatedHello(currentAccessToken)

  if (response.status === 401) {
    currentAccessToken = await refreshAccessToken()

    response = await callAuthenticatedHello(currentAccessToken)
  }

  const data = await response.json()

  if (!response.ok) {
    throw new Error(data.error || 'Request failed')
  }

  return {
    data,
    accessToken: currentAccessToken,
  }
}

async function callAuthenticatedHello(
  accessToken: string
) {
  return fetch(
    'http://localhost:8081/authenticated-hello',
    {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    }
  )
}

async function refreshAccessToken(): Promise<string> {
  const response = await fetch(
    'http://localhost:8081/refresh',
    {
      method: 'POST',
      credentials: 'include',
    }
  )

  const data = await response.json()

  if (!response.ok) {
    throw new Error(
      data.error || 'Failed to refresh access token'
    )
  }

  return data.AccessToken
}