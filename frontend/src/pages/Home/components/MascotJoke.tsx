// src/pages/Home/components/MascotJoke.tsx
import { motion } from 'framer-motion';
import { Mascot } from '@/components/ui/Mascot';

interface MascotJokeProps {
  joke: string;
}

const jokes = [
  'Почему учебник грустит? Потому что у него слишком много проблем!',
  'Что сказал ноль восьмёрке? Красивый пояс!',
  'Какой предмет самый вкусный? ИЗО-бражение торта!',
  'Почему карандаш пошёл в школу? Чтобы стать острее!',
];

export function MascotJoke({ joke }: MascotJokeProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.3 }}
      className="bg-[#F0F4FF] rounded-2xl p-3 flex gap-3 items-center mb-4"
    >
      <Mascot size="sm" className="flex-shrink-0" />
      <p className="text-[12px] text-[#2D3436]">{joke}</p>
    </motion.div>
  );
}

export { jokes };
