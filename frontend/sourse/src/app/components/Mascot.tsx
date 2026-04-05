import { motion } from "motion/react";

interface MascotProps {
  size?: "sm" | "md" | "lg";
  message?: string;
  className?: string;
}

export function Mascot({ size = "md", message, className = "" }: MascotProps) {
  const sizeMap = { sm: 60, md: 96, lg: 120 };
  const px = sizeMap[size];

  return (
    <div className={`flex flex-col items-center gap-1 ${className}`}>
      {message && (
        <motion.div
          initial={{ opacity: 0, y: 5 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white rounded-2xl px-3 py-1.5 shadow-md max-w-[180px] relative"
        >
          <p className="text-[11px] text-center text-[#2D3436]">{message}</p>
          <div className="absolute -bottom-1.5 left-1/2 -translate-x-1/2 w-3 h-3 bg-white rotate-45 shadow-sm" />
        </motion.div>
      )}
      <motion.div
        animate={{ y: [0, -6, 0] }}
        transition={{ repeat: Infinity, duration: 2, ease: "easeInOut" }}
        style={{ width: px, height: px }}
        className="relative"
      >
        <svg viewBox="0 0 120 120" width={px} height={px}>
          {/* Body */}
          <ellipse cx="60" cy="70" rx="38" ry="40" fill="#6C5CE7" />
          {/* Belly */}
          <ellipse cx="60" cy="78" rx="25" ry="26" fill="#A29BFE" />
          {/* Left eye */}
          <circle cx="45" cy="55" r="12" fill="white" />
          <circle cx="47" cy="54" r="6" fill="#2D3436" />
          <circle cx="49" cy="52" r="2" fill="white" />
          {/* Right eye */}
          <circle cx="75" cy="55" r="12" fill="white" />
          <circle cx="73" cy="54" r="6" fill="#2D3436" />
          <circle cx="75" cy="52" r="2" fill="white" />
          {/* Beak / mouth */}
          <path d="M55 67 Q60 73 65 67" stroke="#FD79A8" strokeWidth="2.5" fill="none" strokeLinecap="round" />
          {/* Ears / horns */}
          <path d="M30 45 Q25 20 45 35" fill="#6C5CE7" />
          <path d="M90 45 Q95 20 75 35" fill="#6C5CE7" />
          {/* Cheeks */}
          <circle cx="35" cy="68" r="6" fill="#FD79A8" opacity="0.4" />
          <circle cx="85" cy="68" r="6" fill="#FD79A8" opacity="0.4" />
          {/* Feet */}
          <ellipse cx="45" cy="108" rx="10" ry="5" fill="#FDCB6E" />
          <ellipse cx="75" cy="108" rx="10" ry="5" fill="#FDCB6E" />
          {/* Star on belly */}
          <polygon
            points="60,60 62,66 68,66 63,70 65,76 60,72 55,76 57,70 52,66 58,66"
            fill="#FDCB6E"
          />
        </svg>
      </motion.div>
    </div>
  );
}
