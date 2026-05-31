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
        background: "#0f1117",
        borderBottom: "1px solid #1e2130",
        color: "#d1d5db",
      }}
      className="w-full flex items-center justify-between px-6 py-4 sticky top-0 z-50"
    >
      {/* Brand */}
      <div className="flex flex-col leading-tight">
        <span
          style={{ color: "#f3f4f6" }}
          className="text-xl font-semibold tracking-tight"
        >
          Echo
        </span>

        <span style={{ color: "#6b7280" }} className="text-xs">
          Personal Finance Dashboard
        </span>
      </div>

      {/* Right side actions */}
      <div className="flex items-center gap-3">
        {/* subtle status pill (optional but very fintech-like) */}

        <button
          onClick={handleLogout}
          className="flex items-center gap-2 px-3 py-2 rounded-lg transition-all active:scale-95"
          style={{
            background: "#161922",
            border: "1px solid #2a2d3a",
            color: "#E24B4A",
          }}
        >
          {/* logout icon */}
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
            <polyline points="16 17 21 12 16 7" />
            <line x1="21" y1="12" x2="9" y2="12" />
          </svg>

          {/* <span className="text-sm font-medium">Logout</span> */}
        </button>
      </div>
    </nav>
  )
}
