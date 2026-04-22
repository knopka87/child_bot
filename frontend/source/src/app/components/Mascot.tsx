import { motion } from "motion/react";
import mascotImg from "figma:asset/83306635cb88e5262396a24e1a26e96dbc87efa5.png";

interface MascotProps {
  size?: "sm" | "md" | "lg" | "xl";
  message?: string;
  className?: string;
}

const sizeMap = { sm: 72, md: 100, lg: 130, xl: 200 };

export function Mascot({ size = "md", message, className = "" }: MascotProps) {
  const px = sizeMap[size];

  return (
    <div className={`flex flex-col items-center gap-0 ${className}`}>
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
      >
        <img
          src={mascotImg}
          alt="Маскот"
          width={px}
          height={px}
          className="object-contain drop-shadow-lg"
          style={{ width: px, height: px }}
        />
      </motion.div>
    </div>
  );
}