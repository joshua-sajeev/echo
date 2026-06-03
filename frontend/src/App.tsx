import { useEffect, useState } from "react"
import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
} from "react-router-dom"

import Login from "./pages/Login"
import Dashboard from "./pages/Dashboard"
import { getMe } from "./api/auth"
import CreateTransactionsPage from "./pages/CreateTransactions";
import AllTransactionsPage from "./pages/AllTransactions";
import EditTransactions from "./pages/EditTransactions";
export default function App() {
  const [user, setUser] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  console.log("VITE_API_URL =", import.meta.env.VITE_API_URL);
  useEffect(() => {
    async function init() {
      try {
        const me = await getMe()

        if (me?.authenticated) {
          setUser(me)
        } else {
          setUser(null)
        }
      } catch {
        setUser(null)
      } finally {
        setLoading(false)
      }
    }

    init()
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen text-gray-400">
        Loading...
      </div>
    )
  }

  return (
    <BrowserRouter>
      <Routes>

        <Route
          path="/"
          element={
            user?.authenticated ? (
              <Navigate to="/dashboard" replace />
            ) : (
              <Navigate to="/login" replace />
            )
          }
        />

        <Route
          path="/login"
          element={
            user?.authenticated ? (
              <Navigate to="/dashboard" replace />
            ) : (
              <Login setUser={setUser} />
            )
          }
        />

        <Route
          path="/dashboard"
          element={
            user?.authenticated ? (
              <Dashboard
                user={user}
                setUser={setUser}
              />
            ) : (
              <Navigate to="/login" replace />
            )
          }
        />

        <Route path="/transactions/new" element={<CreateTransactionsPage />} />
        <Route path="/transactions/" element={<AllTransactionsPage />} />
        <Route
          path="/transactions/:id/edit"
          element={
            <EditTransactions
              setUser={setUser}
            />
          }
        />
      </Routes>
    </BrowserRouter>
  )
}
