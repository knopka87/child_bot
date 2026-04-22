import { motion } from "motion/react";
import villainImg from "figma:asset/5a3db19afccf8ffff4ba4ca72544adb630bcb325.png";

interface VillainProps {
  size?: "sm" | "md" | "lg" | "xl" | "2xl";
  className?: string;
  defeated?: boolean;
}

const sizeMap = { sm: 72, md: 100, lg: 140, xl: 200, "2xl": 280 };

export function Villain({ size = "md", className = "", defeated = false }: VillainProps) {
  const px = sizeMap[size];

  return (
    <motion.div
      animate={defeated ? {} : { y: [0, -5, 0] }}
      transition={{ repeat: Infinity, duration: 2.5, ease: "easeInOut" }}
      className={className}
    >
      <img
        src={villainImg}
        alt="Злодей"
        width={px}
        height={px}
        className={`object-contain drop-shadow-lg ${defeated ? "grayscale opacity-50" : ""}`}
        style={{ width: px, height: px }}
      />
    </motion.div>
  );
}