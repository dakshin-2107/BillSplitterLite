import './App.css'
import { Routes, Route } from 'react-router-dom'
import Home from './components/Home'
import JoinSession from './components/JoinSession'
import Ping from './components/debug/ping'
import Header from './components/common/header'
import Footer from './components/common/footer'

function App() {
  return (
    <div className="app-layout">
      <Header />
      <main className="main-content">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/join/:billId" element={<JoinSession />} />
          <Route path="/ping" element={<Ping />} />
        </Routes>
      </main>
      <Footer />
    </div>
  )
}

export default App
