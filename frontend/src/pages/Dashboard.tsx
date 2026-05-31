import Navbar from "../components/Navbar"

export default function Dashboard({
  user,
  setUser,
}: any) {
  return (
    <div>
      <Navbar setUser={setUser} />

      <div className="p-6">
        <h1
          style={{ color: "var(--text-h)" }}
          className="text-2xl font-semibold"
        >
          Dashboard
        </h1>

        <p
          style={{ color: "var(--text)" }}
          className="mt-2"
        >
          Welcome 👋
        </p>
      </div>
    </div>
  )
}
