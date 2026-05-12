import React, { useState, useEffect, useCallback } from 'react';
import { fetchInvitations, updateInvitation } from '../../api/client';
import type { PlaylistShare } from '../../api/client';
import { Check, X, Bell } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

interface InvitationsProps {
  onStatusChange: () => void;
}

const Invitations: React.FC<InvitationsProps> = ({ onStatusChange }) => {
  const [invitations, setInvitations] = useState<PlaylistShare[]>([]);
  const [isOpen, setIsOpen] = useState(false);

  const loadInvitations = useCallback(() => {
    fetchInvitations()
      .then(data => {
        setInvitations(data);
      })
      .catch(err => {
        console.error('Failed to load invitations', err);
      });
  }, []);

  useEffect(() => {
    loadInvitations();
    const interval = setInterval(loadInvitations, 30000); // Poll every 30s
    return () => clearInterval(interval);
  }, [loadInvitations]);

  const handleAction = async (id: number, status: 'accepted' | 'rejected') => {
    try {
      await updateInvitation(id, status);
      setInvitations(prev => prev.filter(inv => inv.id !== id));
      onStatusChange();
    } catch {
      alert('Error updating invitation');
    }
  };

  return (
    <div className="relative">
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className="relative p-2 bg-white/5 rounded-full text-white/40 hover:bg-white/10 hover:text-white transition-all"
      >
        <Bell size={16} />
        {invitations.length > 0 && (
          <span className="absolute top-0 right-0 w-2 h-2 bg-accent-blue rounded-full border-2 border-space-black" />
        )}
      </button>

      <AnimatePresence>
        {isOpen && (
          <motion.div 
            initial={{ opacity: 0, y: 10, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 10, scale: 0.95 }}
            className="absolute right-0 mt-4 w-64 glass border border-white/10 rounded-2xl shadow-2xl z-[100] p-4 space-y-3"
          >
            <h3 className="text-[10px] font-black text-white/30 uppercase tracking-widest px-1">Invitations</h3>
            
            {invitations.length === 0 ? (
              <p className="text-[9px] text-white/20 text-center py-4 uppercase font-black tracking-widest">No pending shares</p>
            ) : (
              <div className="space-y-2">
                {invitations.map(inv => (
                  <div key={inv.id} className="bg-white/5 p-3 rounded-xl border border-white/5 space-y-2">
                    <p className="text-[10px] text-white/60 leading-tight">
                      Invitation to access <span className="text-white font-bold">"{inv.stationName}"</span>
                    </p>
                    <div className="flex gap-2">
                      <button 
                        onClick={() => handleAction(inv.id, 'accepted')}
                        className="flex-1 py-1 bg-accent-blue/20 text-accent-blue text-[9px] font-black rounded-lg hover:bg-accent-blue/30 transition-all flex items-center justify-center gap-1"
                      >
                        <Check size={10} /> ACCEPT
                      </button>
                      <button 
                        onClick={() => handleAction(inv.id, 'rejected')}
                        className="px-2 py-1 bg-white/5 text-white/40 text-[9px] font-black rounded-lg hover:bg-white/10 transition-all"
                      >
                        <X size={10} />
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default Invitations;
