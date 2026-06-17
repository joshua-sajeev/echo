const BASE_URL = import.meta.env.VITE_API_URL;

async function apiFetch<T>(path: string): Promise<T> {
  const res = await fetch(
    `${BASE_URL}${path}`,
    {
      credentials: "include",
    }
  );

  if (!res.ok) {
    throw new Error(`HTTP ${res.status}`);
  }

  return res.json();
}

export function getDashboard() {
  return apiFetch("/dashboard");
}
