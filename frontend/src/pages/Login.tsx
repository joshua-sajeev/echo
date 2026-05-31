import { useState } from "react"
import { useNavigate } from "react-router-dom"

import { login, getMe } from "../api/auth"

const PIN_LENGTH = 6

export default function Login({ setUser }: any) {
  const navigate = useNavigate()

  const [pin, setPin] = useState("")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  async function loginWithPin(value: string) {
    setLoading(true)
    setError("")

    try {
      await login(value)

      const me = await getMe()

      if (me?.authenticated) {
        setUser(me)
        navigate("/dashboard", { replace: true })
        return
      }

      throw new Error()
    } catch {
      setError("Invalid PIN")
      setPin("")
    } finally {
      setLoading(false)
    }
  }

  function addDigit(digit: string) {
    if (loading) return

    setError("")

    if (pin.length >= PIN_LENGTH) {
      return
    }

    const next = pin + digit

    setPin(next)

    if (next.length === PIN_LENGTH) {
      setTimeout(() => {
        loginWithPin(next)
      }, 100)
    }
  }

  function backspace() {
    if (loading) return

    setError("")
    setPin((p) => p.slice(0, -1))
  }

  const digits = [
    "1", "2", "3",
    "4", "5", "6",
    "7", "8", "9",
  ]

  return (
    <div
      className="min-h-screen flex items-start justify-center px-6 pt-20"
      style={{ background: "var(--bg)" }}
    >
      <div className="w-full max-w-xs">
        <div className="text-center mb-3">
          <div
            className="w-14 h-14 rounded-xl mx-auto mb-2 flex items-center justify-center border"
            style={{
              background: "var(--code-bg)",
              borderColor: "var(--border)",
            }}
          >
            <span
              className="text-xl font-bold"
              style={{ color: "var(--accent)" }}
            >
              E
            </span>
          </div>

          <h1
            className="text-xl font-semibold"
            style={{ color: "var(--text-h)" }}
          >
            Echo Finance
          </h1>

          <p
            className="text-sm mt-1"
            style={{ color: "var(--text)" }}
          >
            Enter your PIN to continue
          </p>
        </div>

        <div className="flex justify-center gap-3 mb-3">
          {Array.from({ length: PIN_LENGTH }).map((_, i) => (
            <div
              key={i}
              className="w-4 h-4 rounded-full border transition-all"
              style={{
                borderColor:
                  i < pin.length
                    ? "var(--accent)"
                    : "var(--border)",

                background:
                  i < pin.length
                    ? "var(--accent)"
                    : "transparent",
              }}
            />
          ))}
        </div>

        <div className="h-4 text-center mb-2">
          {error && (
            <p className="text-sm text-red-400">
              {error}
            </p>
          )}
        </div>

        <div className="grid grid-cols-3 gap-2.5 max-w-[280px] mx-auto">
          {digits.map((digit) => (
            <button
              key={digit}
              onClick={() => addDigit(digit)}
              disabled={loading}
              className="aspect-square rounded-lg border text-base font-medium transition"
              style={{
                background: "var(--code-bg)",
                borderColor: "var(--border)",
                color: "var(--text-h)",
              }}
            >
              {digit}
            </button>
          ))}

          <button
            disabled
            className="aspect-square rounded-lg border text-base"
            style={{
              background: "var(--code-bg)",
              borderColor: "var(--border)",
              color: "var(--text)",
              opacity: 0.4,
            }}
          >
            🔒
          </button>

          <button
            onClick={() => addDigit("0")}
            disabled={loading}
            className="aspect-square rounded-lg border text-base font-medium"
            style={{
              background: "var(--code-bg)",
              borderColor: "var(--border)",
              color: "var(--text-h)",
            }}
          >
            0
          </button>

          <button
            onClick={backspace}
            disabled={loading}
            className="aspect-square rounded-lg border text-base"
            style={{
              background: "var(--code-bg)",
              borderColor: "var(--border)",
              color: "var(--text)",
            }}
          >
            ⌫
          </button>
        </div>

        {loading && (
          <p
            className="text-center text-sm mt-4"
            style={{ color: "var(--text)" }}
          >
            Logging in...
          </p>
        )}
      </div>
    </div>
  )
}
