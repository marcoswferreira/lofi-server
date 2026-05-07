import { useEffect, useState } from 'react'
import { fetchStations } from './api/client'
import type { Station, User } from './api/client'
import Player from './components/Player/Player'
import Auth from './components/Layout/Auth'
import Chat from './components/Chat/Chat'
import Pomodoro from './components/Pomodoro/Pomodoro'
import CosmicBackground from './components/Layout/CosmicBackground'
import Admin from './components/Layout/Admin'
import Invitations from './components/Layout/Invitations'
import { Rocket, Radio, ShieldCheck, Menu, X, Settings } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'

function App() {
  const [stations, setStations] = useState<Station[]>([])
  const [currentStationId, setCurrentStationId] = useState<string | null>(null)
  const [user, setUser] = useState<User | null>(null)
  const [, setToken] = useState<string | null>(localStorage.getItem('token'))
  const [isSidebarOpen, setIsSidebarOpen] = useState(false)
  const [isAdminMode, setIsAdminMode] = useState(false)

  const loadStations = () => {
    fetchStations().then(data => {
      setStations(data)
      if (data.length > 0 && (!currentStationId || !data.find(s => s.id === currentStationId))) {
        setCurrentStationId(data[0].id)
      }
    })
  }

  useEffect(() => {
    loadStations()
  }, [user])

  const handleAuth = (user: User, token: string) => {
    setUser(user)
    setToken(token)
    localStorage.setItem('token', token)
  }

  return (
    <div className="relative min-h-screen text-slate-200 overflow-hidden font-satoshi selection:bg-accent-blue/30">
      <CosmicBackground />
      
      {/* Premium Header */}
      <header className="fixed top-0 w-full z-50 px-4 md:px-8 h-16 flex items-center justify-between pointer-events-none">
        <motion.div 
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          className="flex items-center gap-4 pointer-events-auto"
        >
          <motion.div 
            whileHover={{ scale: 1.1 }}
            className="w-8 h-8 md:w-10 md:h-10 bg-white flex items-center justify-center rounded-lg shadow-xl cursor-pointer"
            onClick={() => setIsSidebarOpen(!isSidebarOpen)}
          >
            {isSidebarOpen ? <X size={16} className="text-space-black" /> : <Menu size={16} className="text-space-black lg:hidden" />}
            <Rocket size={18} className="text-space-black hidden lg:block" />
          </motion.div>
          <div>
            <h1 className="text-base md:text-xl font-black tracking-tighter text-white font-orbitron text-glow">LOFI SPACE</h1>
            <p className="text-[7px] font-black tracking-[0.4em] uppercase text-accent-cyan opacity-60 hidden md:block">Orbital Station</p>
          </div>
        </motion.div>

        <div className="flex items-center gap-4 pointer-events-auto">
          {user && <Invitations onStatusChange={loadStations} />}

          <button 
            onClick={() => setIsAdminMode(!isAdminMode)}
            className={`p-2 rounded-full transition-all ${isAdminMode ? 'bg-accent-blue text-white shadow-lg shadow-accent-blue/40' : 'bg-white/5 text-white/40 hover:bg-white/10'}`}
          >
            <Settings size={14} />
          </button>

          <div className="hidden sm:flex items-center gap-2 bg-white/5 border border-white/10 px-3 py-1 rounded-full backdrop-blur-md">
            <ShieldCheck size={10} className="text-accent-cyan" />
            <span className="text-[8px] font-black uppercase tracking-widest text-slate-400">Sync: Active</span>
          </div>
          
          {user ? (
            <button className="glass px-3 py-1.5 rounded-full flex items-center gap-2 border-accent-blue/30">
               <div className="w-1 h-1 bg-accent-blue rounded-full animate-ping" />
               <span className="text-[9px] font-black uppercase tracking-widest">{user.username}</span>
            </button>
          ) : (
            <button className="bg-accent-blue/20 hover:bg-accent-blue text-white px-5 py-1.5 rounded-full text-[9px] font-black uppercase tracking-[0.2em] border border-accent-blue/30 transition-all">
              Portal
            </button>
          )}
        </div>
      </header>

      {/* Main Layout */}
      <main className="relative z-10 pt-20 pb-6 px-4 md:px-6 h-screen overflow-y-auto">
        <AnimatePresence mode="wait">
          {isAdminMode ? (
            <motion.div
              key="admin"
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 1.05 }}
              className="w-full h-full"
            >
              <Admin />
            </motion.div>
          ) : (
            <motion.div
              key="app"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="flex flex-col lg:flex-row h-full gap-6"
            >
              {/* Sidebar */}
              <AnimatePresence>
                {(isSidebarOpen || (typeof window !== 'undefined' && window.innerWidth >= 1024)) && (
                  <motion.aside 
                    initial={{ x: -300, opacity: 0 }}
                    animate={{ x: 0, opacity: 1 }}
                    exit={{ x: -300, opacity: 0 }}
                    className={`
                      fixed inset-y-0 left-0 z-40 w-64 lg:relative lg:w-72 flex flex-col gap-4 p-4 lg:p-0 
                      bg-space-black/95 backdrop-blur-2xl lg:bg-transparent lg:backdrop-blur-none
                      ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
                      transition-transform duration-300 ease-in-out lg:flex
                    `}
                  >
                    <div className="glass rounded-[24px] p-5 flex-1 flex flex-col gap-5 overflow-hidden">
                      <div className="flex items-center justify-between opacity-30 px-1">
                        <span className="text-[8px] font-black uppercase tracking-[0.4em]">Frequencies</span>
                        <Radio size={10} />
                      </div>
                      
                      <div className="space-y-1.5 overflow-y-auto pr-2 custom-scrollbar">
                        {stations.map(s => (
                          <div 
                            key={s.id}
                            className={`
                              px-4 py-3 rounded-xl cursor-pointer transition-all duration-300 border text-[10px] font-black uppercase tracking-widest
                              ${currentStationId === s.id 
                                ? 'bg-white border-white text-space-black shadow-lg' 
                                : 'bg-white/[0.01] hover:bg-white/5 border-white/5 text-slate-500 hover:text-slate-300'}
                            `}
                            onClick={() => {
                              setCurrentStationId(s.id);
                              if (window.innerWidth < 1024) setIsSidebarOpen(false);
                            }}
                          >
                            {s.name}
                          </div>
                        ))}
                      </div>
                    </div>
                    <div className="scale-90 origin-bottom">
                      <Pomodoro />
                    </div>
                  </motion.aside>
                )}
              </AnimatePresence>

              {/* Player Area */}
              <section className="flex-1 flex flex-col items-center justify-center min-h-[350px]">
                <AnimatePresence mode="wait">
                  {currentStationId && (
                    <motion.div
                      key={currentStationId}
                      initial={{ opacity: 0, scale: 0.98 }}
                      animate={{ opacity: 1, scale: 1 }}
                      className="w-full h-full max-w-2xl flex items-center justify-center"
                    >
                      <Player stationId={currentStationId} />
                    </motion.div>
                  )}
                </AnimatePresence>
              </section>

              {/* Chat Area */}
              <aside className="w-full lg:w-80 flex flex-col h-[400px] lg:h-full gap-6">
                {currentStationId && <Chat stationId={currentStationId} username={user?.username || 'Anonymous'} />}
                {!user && (
                  <div className="glass rounded-[24px] p-5 hidden lg:block scale-95 origin-top">
                    <Auth onAuth={handleAuth} />
                  </div>
                )}
              </aside>
            </motion.div>
          )}
        </AnimatePresence>
      </main>

      {/* Simplified Footer */}
      <footer className="fixed bottom-4 w-full px-8 flex justify-between items-center pointer-events-none opacity-20 hidden md:flex">
        <span className="text-[7px] font-black uppercase tracking-[1em]">v4.0.2-OPT</span>
        <span className="text-[7px] font-black uppercase tracking-[1em]">432Hz</span>
      </footer>

      {/* Grainy Texture & Scanlines */}
      <div 
        className="fixed inset-0 pointer-events-none z-[100] opacity-[0.02] mix-blend-overlay"
        style={{ backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E")` }}
      />
    </div>
  )
}

export default App
