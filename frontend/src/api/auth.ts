const BASE_URL = import.meta.env.VITE_API_URL;

export async function login(pin: string) {
  const res = await fetch(`${BASE_URL}/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ pin }),
  })

  if (!res.ok) {
    throw new Error("Login failed")
  }
}

export async function logout() {
  await fetch(`${BASE_URL}/auth/logout`, {
    method: "POST",
    credentials: "include",
  })
}

export async function getMe() {
  const res = await fetch(`${BASE_URL}/auth/me`, {
    credentials: "include",
  })

  if (!res.ok) return null
  return res.json()
}

export async function resetDemo() {
  const res = await fetch(`${BASE_URL}/demo/reset`, {
    method: "POST",
    credentials: "include",
  })

  if (!res.ok) {
    throw new Error("Failed to reset demo")
  }
}
