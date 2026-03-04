import ReactDOM from 'react-dom/client'
import { HashRouter } from 'react-router-dom'
import { Toaster } from 'sonner'
import App from '@/app/App'
import '@/index.css'

const root = document.getElementById('app')
if (!root) {
  throw new Error('Root element #app not found')
}

ReactDOM.createRoot(root).render(
  <HashRouter>
    <App />
    <Toaster richColors position="top-right" />
  </HashRouter>
)
