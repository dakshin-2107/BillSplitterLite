import './App.css'
import { Routes, Route } from 'react-router-dom'
import DashboardPage from './components/dashboard/DashboardPage'
import SplitterFormPage from './components/SplitterFormPage'
import SplitterPage from './components/SplitterPage'
import JoinSession from './components/JoinSession'
import Ping from './components/debug/ping'
import Example from './components/Example'
import Header from './components/common/header'
import Footer from './components/common/footer'
import ErrorPage from './components/common/ErrorPage'
import MobileBlock from './components/common/MobileBlock'
import { TooltipProvider } from "@/components/ui/tooltip"
import { Toaster } from '@ui/sonner'

function App() {
  return (
    <div className="app-layout">
      <MobileBlock />
      <Header />
      <Toaster />
      <TooltipProvider>
        <main className="main-content">
          <Routes>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/splitter" element={<SplitterFormPage />} />
            <Route path="/splitter/:splitId" element={<SplitterPage />} />
            <Route path="/join/:billId" element={<JoinSession />} />
            <Route path="/ping" element={<Ping />} />
            <Route path="/example" element={<Example />} />
            <Route path="*" element={<ErrorPage />} />
          </Routes>
        </main>
      </TooltipProvider>
      <Footer />
    </div>
  )
}

export default App
