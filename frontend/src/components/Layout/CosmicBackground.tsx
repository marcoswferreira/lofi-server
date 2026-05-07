import { motion } from 'framer-motion';

const CosmicBackground = () => {
  // Ultra-lightweight noise texture (Data URI)
  const noiseTexture = "data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E";

  return (
    <div className="fixed inset-0 -z-50 bg-[#050508] overflow-hidden">
      {/* Animated Nebulae (CSS Only) */}
      <motion.div 
        animate={{ 
          scale: [1, 1.1, 1],
          opacity: [0.2, 0.3, 0.2],
          x: [0, 50, 0],
          y: [0, 30, 0]
        }}
        transition={{ duration: 20, repeat: Infinity, ease: "easeInOut" }}
        className="absolute -top-[20%] -left-[10%] w-[80%] h-[80%] bg-[radial-gradient(circle,rgba(122,162,247,0.15)_0%,transparent_70%)] blur-[100px]"
      />
      <motion.div 
        animate={{ 
          scale: [1, 1.2, 1],
          opacity: [0.15, 0.25, 0.15],
          x: [0, -40, 0],
          y: [0, -60, 0]
        }}
        transition={{ duration: 25, repeat: Infinity, ease: "easeInOut", delay: 2 }}
        className="absolute -bottom-[10%] -right-[5%] w-[70%] h-[70%] bg-[radial-gradient(circle,rgba(187,154,247,0.1)_0%,transparent_70%)] blur-[100px]"
      />

      {/* Hero Planet (CSS) */}
      <div 
        className="absolute top-1/4 right-[10%] w-64 h-64 md:w-96 md:h-96 rounded-full"
        style={{
          background: 'radial-gradient(circle at 30% 30%, #1a1b26 0%, #050508 100%)',
          boxShadow: 'inset -20px -20px 50px rgba(0,0,0,0.8), 0 0 40px rgba(122,162,247,0.1)'
        }}
      >
         <div className="absolute inset-0 rounded-full bg-accent-blue/5 blur-3xl scale-110" />
      </div>

      {/* Noise Texture Overlay */}
      <div 
        className="absolute inset-0 opacity-[0.03] mix-blend-overlay pointer-events-none"
        style={{ backgroundImage: `url("${noiseTexture}")` }}
      />
    </div>
  );
};

export default CosmicBackground;
