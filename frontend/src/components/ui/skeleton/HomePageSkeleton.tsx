// src/components/ui/skeleton/HomePageSkeleton.tsx
import { Skeleton } from './Skeleton';

export function HomePageSkeleton() {
  return (
    <div className="flex flex-col min-h-screen bg-[#E8E4FF]">
      {/* Header Skeleton */}
      <div className="px-4 pt-4 pb-3 bg-white/80 backdrop-blur-sm">
        <div className="flex items-center justify-between mb-3">
          {/* Avatar */}
          <Skeleton variant="circular" width={40} height={40} />

          {/* Stats */}
          <div className="flex gap-3">
            <Skeleton className="h-8 w-20 rounded-full" />
            <Skeleton className="h-8 w-20 rounded-full" />
          </div>
        </div>

        {/* Level Progress Bar */}
        <Skeleton className="h-8 w-full rounded-full" />
      </div>

      <div className="flex-1 flex flex-col justify-between pb-20 px-4">
        {/* Mascot Battle Section */}
        <div className="mt-8 space-y-4">
          {/* Villain */}
          <div className="flex justify-center">
            <Skeleton className="h-32 w-32 rounded-2xl" />
          </div>

          {/* Mascot */}
          <div className="flex justify-center">
            <Skeleton className="h-28 w-28 rounded-2xl" />
          </div>

          {/* Speech bubble */}
          <div className="flex justify-center">
            <Skeleton className="h-16 w-3/4 rounded-2xl" />
          </div>
        </div>

        {/* Action Buttons */}
        <div className="space-y-3 mt-8">
          <Skeleton className="h-14 w-full rounded-2xl" />
          <Skeleton className="h-14 w-full rounded-2xl" />
        </div>
      </div>

      {/* Bottom Nav */}
      <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200">
        <div className="flex justify-around items-center h-16 px-4">
          {[1, 2, 3, 4].map((i) => (
            <Skeleton key={i} variant="circular" width={32} height={32} />
          ))}
        </div>
      </div>
    </div>
  );
}
