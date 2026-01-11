import './App.css'
import { Routes, Route } from 'react-router-dom'
import Home from './components/Home'
import JoinSession from './components/JoinSession'
import Ping from './components/debug/ping'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/join/:billId" element={<JoinSession />} />
      <Route path="/ping" element={<Ping />} />
    </Routes>
  )
}

export default App
