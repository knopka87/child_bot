// src/components/ui/skeleton/AchievementsPageSkeleton.tsx
import { Skeleton } from './Skeleton';

export function AchievementsPageSkeleton() {
  return (
    <div className="flex flex-col min-h-screen bg-[#F5F3FF]">
      {/* Header */}
      <div className="flex items-center gap-3 px-4 py-4 bg-white/90">
        <Skeleton variant="circular" width={32} height={32} />
        <Skeleton className="h-6 w-32" />
      </div>

      {/* Stats */}
      <div className="px-4 py-6 bg-white/50">
        <div className="flex justify-around">
          <div className="text-center">
            <Skeleton className="h-8 w-16 mx-auto mb-2" />
            <Skeleton className="h-4 w-24 mx-auto" />
          </div>
          <div className="text-center">
            <Skeleton className="h-8 w-16 mx-auto mb-2" />
            <Skeleton className="h-4 w-24 mx-auto" />
          </div>
        </div>
      </div>

      {/* Achievements Grid */}
      <div className="flex-1 px-4 py-6 space-y-8">
        {[1, 2, 3].map((shelf) => (
          <div key={shelf}>
            {/* Shelf */}
            <div className="relative">
              <div className="grid grid-cols-4 gap-3 mb-2">
                {[1, 2, 3, 4].map((item) => (
                  <div key={item} className="flex flex-col items-center">
                    <Skeleton variant="circular" width={64} height={64} className="mb-2" />
                    <Skeleton className="h-3 w-12" />
                  </div>
                ))}
              </div>
              {/* Shelf line */}
              <Skeleton className="h-2 w-full rounded-full" />
            </div>
          </div>
        ))}
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
