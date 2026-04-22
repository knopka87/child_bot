// tests/unit/components/Button.test.tsx
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ConfigProvider } from '@vkontakte/vkui';
import { Button } from '@/components/ui/Button/Button';

const renderWithVKUI = (component: React.ReactElement) => {
  return render(
    <ConfigProvider platform="android">{component}</ConfigProvider>
  );
};

describe('Button', () => {
  it('renders with text', () => {
    renderWithVKUI(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick when clicked', () => {
    const onClick = vi.fn();
    renderWithVKUI(<Button onClick={onClick}>Click me</Button>);

    fireEvent.click(screen.getByText('Click me'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('is disabled when disabled prop is true', () => {
    renderWithVKUI(<Button disabled>Click me</Button>);
    const button = screen.getByRole('button');

    expect(button).toBeDisabled();
  });
});
