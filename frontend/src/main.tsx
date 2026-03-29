import React from 'react'
import ReactDOM from 'react-dom/client'

function App() {
  return (
    <main style={{ fontFamily: 'ui-sans-serif, system-ui, sans-serif', padding: '2rem', lineHeight: 1.5 }}>
      <h1 style={{ marginTop: 0 }}>Synthify</h1>
      <p>Frontend dev server is running.</p>
      <ul>
        <li>API base URL: {import.meta.env.VITE_API_BASE_URL ?? 'not set'}</li>
        <li>Next step: wire Connect RPC clients and graph features.</li>
      </ul>
    </main>
  )
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
