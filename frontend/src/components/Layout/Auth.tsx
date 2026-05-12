import React, { useState } from 'react';
import { login, register } from '../../api/client';
import type { User } from '../../api/client';
import { Mail, Lock, User as UserIcon, ArrowRight } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

interface AuthProps {
  onAuth: (user: User, token: string) => void;
}

const Auth: React.FC<AuthProps> = ({ onAuth }) => {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({ username: '', email: '', password: '' });
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const res = isLogin 
        ? await login({ email: formData.email, password: formData.password })
        : await register(formData);
      onAuth(res.user, res.token);
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Authentication failed';
      setError(errorMessage);
    }
  };

  return (
    <div className="space-y-10">
      <div className="flex flex-col gap-2">
        <div className="flex items-center gap-3 opacity-60">
          <UserIcon size={14} className="text-accent-blue" />
          <h3 className="text-[10px] font-black uppercase tracking-[0.5em] text-white">
            {isLogin ? 'Access Identity' : 'New Subject'}
          </h3>
        </div>
        <p className="text-[9px] text-slate-600 font-bold uppercase tracking-widest">Biometric override enabled</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        {!isLogin && (
          <div className="relative group">
            <UserIcon size={14} className="absolute left-6 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-accent-blue transition-colors" />
            <input 
              type="text" 
              placeholder="Designation" 
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              className="w-full bg-white/[0.02] border border-white/5 rounded-[20px] pl-14 pr-6 py-4 text-sm focus:outline-none focus:border-accent-blue/30 focus:bg-white/[0.04] transition-all"
              required
            />
          </div>
        )}
        <div className="relative group">
          <Mail size={14} className="absolute left-6 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-accent-blue transition-colors" />
          <input 
            type="email" 
            placeholder="Neural Link (Email)" 
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            className="w-full bg-white/[0.02] border border-white/5 rounded-[20px] pl-14 pr-6 py-4 text-sm focus:outline-none focus:border-accent-blue/30 focus:bg-white/[0.04] transition-all"
            required
          />
        </div>
        <div className="relative group">
          <Lock size={14} className="absolute left-6 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-accent-blue transition-colors" />
          <input 
            type="password" 
            placeholder="Access Key" 
            value={formData.password}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            className="w-full bg-white/[0.02] border border-white/5 rounded-[20px] pl-14 pr-6 py-4 text-sm focus:outline-none focus:border-accent-blue/30 focus:bg-white/[0.04] transition-all"
            required
          />
        </div>

        <AnimatePresence>
          {error && (
            <motion.p 
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              className="text-[9px] text-accent-magenta font-black uppercase tracking-widest text-center"
            >
              ⚠ Error: {error}
            </motion.p>
          )}
        </AnimatePresence>

        <motion.button 
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          type="submit"
          className="w-full bg-white text-space-black py-4 rounded-[20px] font-black text-[10px] uppercase tracking-[0.4em] hover:bg-accent-blue hover:text-white transition-all shadow-xl shadow-white/5 flex items-center justify-center gap-3"
        >
          {isLogin ? 'Initiate Login' : 'Create Profile'}
          <ArrowRight size={14} />
        </motion.button>
      </form>
      
      <p 
        onClick={() => setIsLogin(!isLogin)} 
        className="text-[9px] text-center font-black uppercase tracking-[0.3em] text-slate-600 hover:text-accent-blue transition-colors cursor-pointer"
      >
        {isLogin ? "// Access Request Required" : "// Protocol Restoration"}
      </p>
    </div>
  );
};

export default Auth;
