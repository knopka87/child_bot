// src/components/ui/skeleton/ProfilePageSkeleton.tsx
import { Skeleton } from './Skeleton';

export function ProfilePageSkeleton() {
  return (
    <div className="flex flex-col min-h-screen bg-[#F5F3FF]">
      {/* Header */}
      <div className="flex items-center gap-3 px-4 py-4 bg-white/90">
        <Skeleton variant="circular" width={32} height={32} />
        <Skeleton className="h-6 w-24" />
      </div>

      <div className="flex-1 px-4 py-6 space-y-6">
        {/* Avatar Section */}
        <div className="flex flex-col items-center space-y-3">
          <Skeleton variant="circular" width={96} height={96} />
          <Skeleton className="h-8 w-32" />
          <Skeleton className="h-5 w-24" />
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-2 gap-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="bg-white rounded-2xl p-4">
              <Skeleton className="h-10 w-12 mb-2" />
              <Skeleton className="h-4 w-full" />
            </div>
          ))}
        </div>

        {/* Actions */}
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-white rounded-2xl p-4 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Skeleton variant="circular" width={40} height={40} />
                <div>
                  <Skeleton className="h-5 w-24 mb-2" />
                  <Skeleton className="h-4 w-32" />
                </div>
              </div>
              <Skeleton variant="circular" width={24} height={24} />
            </div>
          ))}
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
