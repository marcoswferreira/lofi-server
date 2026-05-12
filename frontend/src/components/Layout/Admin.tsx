import React, { useState, useEffect, useCallback } from 'react';
import { fetchStations, createStation, updateStation, deleteStation, shareStation } from '../../api/client';
import type { Station, Track } from '../../api/client';
import { Plus, Trash2, Edit2, Save, X, Music, Radio, Share2 } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const Admin: React.FC = () => {
  const [stations, setStations] = useState<Station[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState<Partial<Station>>({
    name: '',
    description: '',
    playlist: []
  });
  const [shareEmail, setShareEmail] = useState<{[key: string]: string}>({});

  const loadStations = useCallback(() => {
    fetchStations()
      .then(data => {
        setStations(data);
      })
      .catch(err => {
        console.error(err);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    loadStations();
  }, [loadStations]);

  const handleSave = async () => {
    try {
      if (editingId) {
        await updateStation(editingId, formData);
      } else {
        await createStation(formData);
      }
      setEditingId(null);
      setFormData({ name: '', description: '', playlist: [] });
      loadStations();
    } catch {
      alert('Error saving station');
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure?')) return;
    try {
      await deleteStation(id);
      loadStations();
    } catch {
      alert('Error deleting station');
    }
  };

  const handleShare = async (id: string) => {
    const email = shareEmail[id];
    if (!email) return;
    try {
      await shareStation(id, email);
      alert('Invitation sent!');
      setShareEmail({ ...shareEmail, [id]: '' });
    } catch {
      alert('Error sharing: User not found or already shared.');
    }
  };

  const startEdit = (station: Station) => {
    setEditingId(station.id);
    setFormData(station);
  };

  const addTrack = () => {
    const playlist = [...(formData.playlist || []), { id: '', title: '', duration: 300000000000 }];
    setFormData({ ...formData, playlist });
  };

  const updateTrack = (index: number, field: keyof Track, value: string | number) => {
    const playlist = [...(formData.playlist || [])];
    playlist[index] = { ...playlist[index], [field]: value };
    setFormData({ ...formData, playlist });
  };

  const removeTrack = (index: number) => {
    const playlist = [...(formData.playlist || [])];
    playlist.splice(index, 1);
    setFormData({ ...formData, playlist });
  };

  if (loading) return <div className="text-white p-8">Loading dashboard...</div>;

  return (
    <div className="p-4 md:p-8 max-w-6xl mx-auto space-y-8">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-4xl font-black text-white tracking-tighter">ADMIN PORTAL</h1>
          <p className="text-accent-blue/60 font-mono text-sm tracking-widest uppercase">Station Management Control</p>
        </div>
        {!editingId && (
          <motion.button 
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={() => { setEditingId(null); setFormData({ name: '', description: '', playlist: [] }); }}
            className="flex items-center gap-2 px-6 py-3 bg-white text-space-black rounded-full font-bold shadow-lg"
          >
            <Plus size={18} /> CREATE STATION
          </motion.button>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* List of Stations */}
        <div className="lg:col-span-1 space-y-4">
          <h2 className="text-xs font-black text-white/30 uppercase tracking-[0.3em] px-2">Active Stations</h2>
          {stations.map(station => (
            <motion.div 
              key={station.id}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              className={`glass p-4 rounded-2xl border ${editingId === station.id ? 'border-accent-blue bg-accent-blue/10' : 'border-white/5'} transition-all`}
            >
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="text-white font-bold">{station.name}</h3>
                  <p className="text-white/40 text-xs line-clamp-1">{station.description}</p>
                  <p className="text-accent-cyan text-[10px] mt-2 font-mono">{station.playlist.length} TRACKS</p>
                </div>
                <div className="flex gap-2">
                  <button onClick={() => startEdit(station)} className="p-2 hover:bg-white/10 rounded-full text-white/60 transition-colors">
                    <Edit2 size={16} />
                  </button>
                  <button onClick={() => handleDelete(station.id)} className="p-2 hover:bg-red-500/20 rounded-full text-red-500/60 transition-colors">
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>

              {station.ownerId && (
                <div className="mt-3 pt-3 border-t border-white/5 flex gap-2">
                  <input 
                    value={shareEmail[station.id] || ''}
                    onChange={e => setShareEmail({ ...shareEmail, [station.id]: e.target.value })}
                    placeholder="Share with email..."
                    className="flex-1 bg-black/20 border border-white/5 rounded-lg px-3 py-1 text-[10px] text-white outline-none focus:border-accent-blue"
                  />
                  <button 
                    onClick={() => handleShare(station.id)}
                    className="p-1.5 bg-accent-blue/20 text-accent-blue rounded-lg hover:bg-accent-blue/30 transition-all"
                  >
                    <Share2 size={12} />
                  </button>
                </div>
              )}
            </motion.div>
          ))}
        </div>

        {/* Editor Form */}
        <div className="lg:col-span-2">
          <AnimatePresence mode="wait">
            <motion.div 
              key={editingId || 'new'}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="glass p-6 md:p-8 rounded-[2rem] border border-white/10 space-y-6"
            >
              <div className="flex items-center gap-3 border-b border-white/5 pb-4">
                <Radio className="text-accent-blue" />
                <h2 className="text-xl font-black text-white">{editingId ? 'EDIT STATION' : 'NEW STATION'}</h2>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <label className="text-[10px] font-black text-white/30 uppercase tracking-widest ml-2">Station Name</label>
                  <input 
                    value={formData.name}
                    onChange={e => setFormData({...formData, name: e.target.value})}
                    placeholder="e.g. Midnight Beats"
                    className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white focus:border-accent-blue outline-none transition-all"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-[10px] font-black text-white/30 uppercase tracking-widest ml-2">Description</label>
                  <input 
                    value={formData.description}
                    onChange={e => setFormData({...formData, description: e.target.value})}
                    placeholder="Describe the vibe..."
                    className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white focus:border-accent-blue outline-none transition-all"
                  />
                </div>
              </div>

              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <h3 className="text-xs font-black text-white/30 uppercase tracking-widest ml-2">Playlist Configuration</h3>
                  <button 
                    onClick={addTrack}
                    className="flex items-center gap-1 text-[10px] font-black bg-accent-blue/20 text-accent-blue px-3 py-1 rounded-full hover:bg-accent-blue/30 transition-all"
                  >
                    <Plus size={12} /> ADD TRACK
                  </button>
                </div>

                <div className="space-y-3 max-h-[400px] overflow-y-auto pr-2 custom-scrollbar">
                  {formData.playlist?.map((track, idx) => (
                    <div key={idx} className="flex flex-col md:flex-row gap-3 bg-white/5 p-4 rounded-2xl border border-white/5">
                      <div className="flex-1 space-y-2">
                        <label className="text-[8px] font-black text-white/20 uppercase">YouTube ID</label>
                        <input 
                          value={track.id}
                          onChange={e => updateTrack(idx, 'id', e.target.value)}
                          placeholder="e.g. jfKfPfyJRdk"
                          className="w-full bg-black/20 border border-white/5 rounded-lg px-3 py-2 text-white text-xs font-mono"
                        />
                      </div>
                      <div className="flex-[2] space-y-2">
                        <label className="text-[8px] font-black text-white/20 uppercase">Track Title</label>
                        <input 
                          value={track.title}
                          onChange={e => updateTrack(idx, 'title', e.target.value)}
                          placeholder="Song name - Artist"
                          className="w-full bg-black/20 border border-white/5 rounded-lg px-3 py-2 text-white text-xs"
                        />
                      </div>
                      <div className="flex-shrink-0 flex items-end">
                        <button 
                          onClick={() => removeTrack(idx)}
                          className="p-3 text-red-500/40 hover:text-red-500 hover:bg-red-500/10 rounded-xl transition-all"
                        >
                          <Trash2 size={16} />
                        </button>
                      </div>
                    </div>
                  ))}
                  {(!formData.playlist || formData.playlist.length === 0) && (
                    <div className="text-center py-12 bg-white/5 rounded-3xl border border-dashed border-white/10">
                      <Music className="mx-auto text-white/10 mb-2" size={32} />
                      <p className="text-white/20 text-xs uppercase font-black tracking-widest">No tracks added yet</p>
                    </div>
                  )}
                </div>
              </div>

              <div className="flex gap-4 pt-4">
                <button 
                  onClick={handleSave}
                  className="flex-1 flex items-center justify-center gap-2 py-4 bg-accent-blue text-white rounded-2xl font-black shadow-xl shadow-accent-blue/20 hover:brightness-110 transition-all"
                >
                  <Save size={18} /> SAVE STATION
                </button>
                {editingId && (
                  <button 
                    onClick={() => { setEditingId(null); setFormData({ name: '', description: '', playlist: [] }); }}
                    className="px-6 py-4 bg-white/5 text-white/60 rounded-2xl font-bold hover:bg-white/10 transition-all"
                  >
                    <X size={18} />
                  </button>
                )}
              </div>
            </motion.div>
          </AnimatePresence>
        </div>
      </div>
    </div>
  );
};

export default Admin;
