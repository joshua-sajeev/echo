const BASE_URL = import.meta.env.VITE_API_URL; 

export async function getDashboard() {
  const res = await fetch(
    `${BASE_URL}/dashboard`,
    {
      credentials: "include",
    }
  );

  if (!res.ok) {
    throw new Error("Failed to fetch dashboard");
  }

  return res.json();
}
