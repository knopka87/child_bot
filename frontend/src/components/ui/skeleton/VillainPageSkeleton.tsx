// src/components/ui/skeleton/VillainPageSkeleton.tsx
import { Skeleton } from './Skeleton';

export function VillainPageSkeleton() {
  return (
    <div className="flex flex-col min-h-screen bg-gradient-to-b from-[#1a1a2e] to-[#16213e]">
      {/* Header */}
      <div className="flex items-center gap-3 px-4 py-4">
        <Skeleton variant="circular" width={32} height={32} className="bg-white/20" />
      </div>

      <div className="flex-1 flex flex-col items-center justify-center px-4 space-y-6">
        {/* Villain Image */}
        <Skeleton className="h-48 w-48 rounded-3xl bg-white/20" />

        {/* Villain Name */}
        <Skeleton className="h-8 w-40 bg-white/20" />

        {/* Health Bar */}
        <div className="w-full max-w-xs">
          <div className="flex justify-between mb-2">
            <Skeleton className="h-5 w-16 bg-white/20" />
            <Skeleton className="h-5 w-20 bg-white/20" />
          </div>
          <Skeleton className="h-6 w-full rounded-full bg-white/20" />
        </div>

        {/* Description */}
        <div className="w-full max-w-sm space-y-2">
          <Skeleton className="h-4 w-full bg-white/20" />
          <Skeleton className="h-4 w-5/6 bg-white/20" />
          <Skeleton className="h-4 w-4/6 bg-white/20" />
        </div>

        {/* Battle Button */}
        <Skeleton className="h-14 w-full max-w-xs rounded-2xl bg-white/20 mt-8" />
      </div>
    </div>
  );
}
