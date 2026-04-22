// src/components/ui/skeleton/ListPageSkeleton.tsx
import { Skeleton } from './Skeleton';

interface ListPageSkeletonProps {
  itemCount?: number;
  showHeader?: boolean;
  showBottomNav?: boolean;
}

export function ListPageSkeleton({
  itemCount = 5,
  showHeader = true,
  showBottomNav = true
}: ListPageSkeletonProps) {
  return (
    <div className="flex flex-col min-h-screen bg-[#F5F3FF]">
      {/* Header */}
      {showHeader && (
        <div className="flex items-center gap-3 px-4 py-4 bg-white/90">
          <Skeleton variant="circular" width={32} height={32} />
          <Skeleton className="h-6 w-32" />
        </div>
      )}

      {/* List Items */}
      <div className="flex-1 px-4 py-6 space-y-4">
        {Array.from({ length: itemCount }).map((_, i) => (
          <div key={i} className="bg-white rounded-2xl p-4">
            <div className="flex items-center gap-4">
              <Skeleton variant="circular" width={48} height={48} />
              <div className="flex-1 space-y-2">
                <Skeleton className="h-5 w-3/4" />
                <Skeleton className="h-4 w-1/2" />
              </div>
              <Skeleton variant="circular" width={24} height={24} />
            </div>
          </div>
        ))}
      </div>

      {/* Bottom Nav */}
      {showBottomNav && (
        <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200">
          <div className="flex justify-around items-center h-16 px-4">
            {[1, 2, 3, 4].map((i) => (
              <Skeleton key={i} variant="circular" width={32} height={32} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
