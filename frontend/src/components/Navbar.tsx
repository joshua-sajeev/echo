import { useNavigate } from "react-router-dom"
import { logout } from "../api/auth"

export default function Navbar({ setUser }: any) {
  const navigate = useNavigate()

  async function handleLogout() {
    try {
      await logout()
    } finally {
      setUser(null)
      navigate("/login", { replace: true })
    }
  }

  return (
    <nav
      style={{
        background: "var(--bg)",
        borderBottom: "1px solid var(--border)",
        color: "var(--text)",
      }}
      className="w-full flex items-center justify-between px-6 py-4 backdrop-blur-md sticky top-0"
    >
      <div className="flex flex-col leading-tight">
        <span
          style={{ color: "var(--text-h)" }}
          className="text-xl font-semibold tracking-tight"
        >
          Echo
        </span>

        <span
          className="text-xs"
          style={{ color: "var(--text)" }}
        >
          Personal Finance Dashboard
        </span>
      </div>

      <button
        onClick={handleLogout}
        style={{
          background: "var(--accent-bg)",
          color: "var(--accent)",
          border: "1px solid var(--accent-border)",
        }}
        className="px-4 py-2 text-sm rounded-lg transition-all hover:opacity-90 active:scale-95"
      >
        Logout
      </button>
    </nav>
  )
}
