import React, { useEffect, useState, useRef } from 'react';
import { Send, MessageSquare } from 'lucide-react';

interface Message {
  type: string;
  stationId: string;
  payload: {
    user: string;
    text: string;
  };
}

interface ChatProps {
  stationId: string;
  username: string;
}

const Chat: React.FC<ChatProps> = ({ stationId, username }) => {
  const [messages, setMessages] = useState<Message['payload'][]>([]);
  const [inputText, setInputText] = useState('');
  const ws = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    let socket: WebSocket | null = null;
    let timeoutId: number | null = null;
    let isMounted = true;

    const connect = () => {
      if (!isMounted) return;
      
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host;
      socket = new WebSocket(`${protocol}//${host}/ws?station=${stationId}`);
      ws.current = socket;

      socket.onmessage = (event) => {
        try {
          const msg: Message = JSON.parse(event.data);
          if (msg.type === 'chat') {
            setMessages((prev) => [...prev, msg.payload]);
          }
        } catch (error) {
          console.error('Failed to parse WS message', error);
        }
      };

      socket.onclose = () => {
        if (isMounted) {
          timeoutId = window.setTimeout(connect, 3000);
        }
      };

      socket.onerror = (err) => {
        console.error('WS Error:', err);
      };
    };

    connect();

    return () => {
      isMounted = false;
      if (timeoutId) clearTimeout(timeoutId);
      if (socket) {
        socket.onclose = null;
        socket.onmessage = null;
        socket.onerror = null;
        const s = socket; // Local ref
        if (s.readyState === WebSocket.OPEN) {
          s.close();
        } else if (s.readyState === WebSocket.CONNECTING) {
          s.onopen = () => {
            try { s.close(); } catch { /* Ignore close error */ }
          };
        }
      }
    };
  }, [stationId]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const sendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputText.trim()) return;

    const msg = {
      type: 'chat',
      payload: {
        user: username || 'Anonymous',
        text: inputText,
      },
    };

    ws.current?.send(JSON.stringify(msg));
    setInputText('');
  };

  return (
    <div className="flex flex-col h-full glass rounded-[32px] overflow-hidden border-white/5">
      <div className="px-6 py-4 border-b border-white/5 bg-white/[0.02] flex items-center justify-between">
        <div className="flex items-center gap-2">
          <MessageSquare size={14} className="text-accent-blue" />
          <h3 className="text-[10px] font-black uppercase tracking-[0.2em] text-white">Live Space</h3>
        </div>
        <div className="flex items-center gap-2 bg-black/20 px-3 py-1 rounded-full border border-white/5">
          <div className="w-1 h-1 bg-green-500 rounded-full animate-pulse" />
          <span className="text-[8px] font-black uppercase tracking-widest text-green-500">Active</span>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-6 py-4 space-y-4 custom-scrollbar">
        {messages.map((m, i) => (
          <div key={i} className="flex flex-col gap-1">
            <span className="text-[9px] font-black text-accent-blue uppercase tracking-widest ml-1">{m.user}</span>
            <div className="bg-white/[0.03] px-4 py-2.5 rounded-2xl rounded-tl-none border border-white/5 max-w-[90%]">
              <p className="text-xs text-slate-300 leading-relaxed font-medium">{m.text}</p>
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <form onSubmit={sendMessage} className="p-4 bg-white/[0.02] border-t border-white/5">
        <div className="relative">
          <input 
            type="text" 
            value={inputText} 
            onChange={(e) => setInputText(e.target.value)}
            placeholder="Say something..."
            className="w-full bg-black/40 border border-white/10 rounded-2xl px-5 py-3 pr-12 text-xs focus:outline-none focus:border-accent-blue/30 transition-all placeholder:text-slate-700"
          />
          <button 
            type="submit"
            className="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 flex items-center justify-center text-accent-blue hover:text-white transition-colors"
          >
            <Send size={16} />
          </button>
        </div>
      </form>
    </div>
  );
};

export default Chat;
