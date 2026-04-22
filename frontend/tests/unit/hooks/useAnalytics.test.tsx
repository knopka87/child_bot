// tests/unit/hooks/useAnalytics.test.tsx
import { describe, it, expect } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useAnalytics } from '@/hooks/useAnalytics';

describe('useAnalytics', () => {
  it('throws error when used outside AnalyticsProvider', () => {
    expect(() => renderHook(() => useAnalytics())).toThrow(
      'useAnalytics must be used within AnalyticsProvider'
    );
  });
});
