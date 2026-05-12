import React, { useEffect, useState, useCallback } from 'react';
import YouTube from 'react-youtube';
import type { YouTubeProps, YouTubePlayer } from 'react-youtube';
import { fetchStationState } from '../../api/client';
import type { StationState } from '../../api/client';
import { Play, Pause, Volume2, Headphones } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

interface PlayerProps {
  stationId: string;
}

const Player: React.FC<PlayerProps> = ({ stationId }) => {
  const [state, setState] = useState<StationState | null>(null);
  const [player, setPlayer] = useState<YouTubePlayer | null>(null);
  const [volume, setVolume] = useState(50);
  const [isPlaying, setIsPlaying] = useState(false);

  const refreshState = useCallback(() => {
    fetchStationState(stationId)
      .then(newState => {
        setState(newState);
      })
      .catch(err => {
        console.error('Failed to sync state', err);
      });
  }, [stationId]);

  useEffect(() => {
    refreshState();
    const interval = setInterval(refreshState, 15000);
    return () => clearInterval(interval);
  }, [refreshState]);

  const onReady: YouTubeProps['onReady'] = (event) => {
    setPlayer(event.target);
    event.target.setVolume(volume);
    if (state) {
      event.target.seekTo(state.currentSeconds, true);
      event.target.playVideo();
    }
  };

  const onStateChange: YouTubeProps['onStateChange'] = (event) => {
    if (event.data === 0) refreshState();
    setIsPlaying(event.data === 1);
  };

  const togglePlay = () => {
    if (isPlaying) player?.pauseVideo();
    else player?.playVideo();
  };

  return (
    <div className="flex flex-col items-center justify-center w-full max-w-xl relative px-4 py-4">
      {/* Optimized Background Glow */}
      <div className="absolute w-[200px] h-[200px] md:w-[350px] md:h-[350px] bg-accent-blue/10 rounded-full blur-[60px] pointer-events-none" />

      <div className="relative z-10 w-full flex flex-col items-center space-y-6">
        <AnimatePresence mode="wait">
          <motion.div 
            key={state?.currentTrack.id || 'loading'}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="text-center space-y-2"
          >
            <div className="flex items-center justify-center gap-2 opacity-30">
              <Headphones size={10} className="text-accent-cyan" />
              <span className="text-[8px] font-black uppercase tracking-[0.4em] text-accent-cyan">Digital Stream</span>
            </div>
            <h2 className="text-xl md:text-3xl lg:text-4xl font-black text-white leading-tight tracking-tighter filter drop-shadow-xl line-clamp-2 px-4">
              {state?.currentTrack.title || 'Connecting...'}
            </h2>
          </motion.div>
        </AnimatePresence>

        <div className="relative group">
          <motion.button 
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={togglePlay}
            className="relative w-28 h-28 md:w-36 md:h-36 bg-white text-space-black rounded-full flex items-center justify-center transition-all duration-300 z-10 shadow-xl"
          >
            {isPlaying ? <Pause size={40} fill="currentColor" /> : <Play size={40} className="translate-x-1" fill="currentColor" />}
          </motion.button>
        </div>

        <div className="w-full max-w-[280px]">
          <div className="flex items-center gap-3 glass rounded-xl p-3 border-white/5">
            <Volume2 size={14} className="text-accent-blue opacity-50" />
            <input 
              type="range" 
              min="0" 
              max="100" 
              value={volume}
              onChange={(e) => {
                const val = parseInt(e.target.value);
                setVolume(val);
                player?.setVolume(val);
              }}
              className="flex-1 h-0.5 bg-white/10 rounded-full appearance-none accent-accent-blue cursor-pointer"
            />
          </div>
        </div>
      </div>

      <div className="hidden">
        {state && (
          <YouTube 
            videoId={state.currentTrack.id} 
            opts={{ 
              height: '0', 
              width: '0', 
              playerVars: { 
                autoplay: 1, 
                controls: 0, 
                modestbranding: 1,
                origin: window.location.origin,
                enablejsapi: 1
              }
            }} 
            onReady={onReady} 
            onStateChange={onStateChange}
          />
        )}
      </div>
    </div>
  );
};

export default Player;
