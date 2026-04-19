import './App.css'
import { Routes, Route } from 'react-router-dom'
import Home from './components/Home'
import JoinSession from './components/JoinSession'
import Ping from './components/debug/ping'
import Example from './components/Example'
import Header from './components/common/header'
import Footer from './components/common/footer'
import { TooltipProvider } from "@/components/ui/tooltip"

function App() {
  return (
    <div className="app-layout">
      <Header />
      <TooltipProvider>
        <main className="main-content">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/join/:billId" element={<JoinSession />} />
            <Route path="/ping" element={<Ping />} />
            <Route path="/example" element={<Example />} />
          </Routes>
        </main>
      </TooltipProvider>
      <Footer />
    </div>
  )
}

export default App
