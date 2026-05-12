import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';

const Pomodoro: React.FC = () => {
  const [minutes, setMinutes] = useState(25);
  const [seconds, setSeconds] = useState(0);
  const [isActive, setIsActive] = useState(false);

  useEffect(() => {
    let interval: ReturnType<typeof setInterval> | null = null;
    if (isActive) {
      interval = setInterval(() => {
        if (seconds > 0) setSeconds(seconds - 1);
        else if (minutes > 0) { setMinutes(minutes - 1); setSeconds(59); }
        else { setIsActive(false); }
      }, 1000);
    } else if (interval) {
      clearInterval(interval);
    }
    return () => {
      if (interval) clearInterval(interval);
    };
  }, [isActive, minutes, seconds]);

  const progress = ((25 * 60 - (minutes * 60 + seconds)) / (25 * 60)) * 100;

  return (
    <div className="glass rounded-[32px] p-8 flex flex-col items-center gap-6 group hover:border-accent-purple/30 transition-colors">
      <div className="flex items-center gap-2 opacity-40 group-hover:opacity-100 transition-opacity">
        <span className="text-[10px] font-black uppercase tracking-[0.4em]">Orbital Timer</span>
      </div>

      <div className="relative w-40 h-40 flex items-center justify-center">
        {/* Orbital Ring */}
        <svg className="absolute inset-0 w-full h-full -rotate-90">
          <circle cx="80" cy="80" r="70" fill="none" stroke="rgba(255,255,255,0.05)" strokeWidth="4" />
          <motion.circle 
            cx="80" cy="80" r="70" fill="none" stroke="#bb9af7" strokeWidth="4" 
            strokeDasharray="440" 
            animate={{ strokeDashoffset: 440 - (440 * progress) / 100 }}
            strokeLinecap="round"
          />
        </svg>

        <div className="text-center z-10">
          <span className="text-4xl font-black text-white text-glow">
            {String(minutes).padStart(2, '0')}:{String(seconds).padStart(2, '0')}
          </span>
        </div>
      </div>

      <div className="flex gap-4 w-full">
        <button 
          onClick={() => setIsActive(!isActive)}
          className="flex-1 bg-white text-space-black py-3 rounded-2xl font-black text-[10px] uppercase tracking-widest hover:bg-accent-purple hover:text-white transition-all shadow-xl shadow-white/5 active:scale-95"
        >
          {isActive ? 'Pause' : 'Ignite'}
        </button>
      </div>
    </div>
  );
};

export default Pomodoro;
