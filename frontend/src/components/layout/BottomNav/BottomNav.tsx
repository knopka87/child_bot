// src/components/layout/BottomNav/BottomNav.tsx
import { Home, Trophy, Users, User } from 'lucide-react';
import { useNavigate, useLocation } from 'react-router-dom';
import { ROUTES } from '@/config/routes';
import { useHaptics } from '@/lib/platform/haptics';

const tabs = [
  { path: ROUTES.HOME, label: 'Главная', icon: Home },
  { path: ROUTES.ACHIEVEMENTS, label: 'Достижения', icon: Trophy },
  { path: ROUTES.FRIENDS, label: 'Друзья', icon: Users },
  { path: ROUTES.PROFILE, label: 'Профиль', icon: User },
];

export function BottomNav() {
  const navigate = useNavigate();
  const location = useLocation();
  const { onButtonClick } = useHaptics();

  const handleNavigate = (route: string) => {
    onButtonClick();
    navigate(route);
  };

  return (
    <nav className="bottom-nav-bar fixed bottom-0 left-0 right-0 bg-white border-t border-[#DFE6E9] flex z-50 safe-area-bottom">
      {tabs.map((tab) => {
        const isActive =
          location.pathname === tab.path ||
          (tab.path !== ROUTES.HOME && location.pathname.startsWith(tab.path));
        const Icon = tab.icon;
        return (
          <button
            key={tab.path}
            onClick={() => handleNavigate(tab.path)}
            className={`flex-1 flex flex-col items-center py-2 pt-3 gap-0.5 transition-colors ${
              isActive ? 'text-[#6C5CE7]' : 'text-[#636e72]'
            }`}
          >
            <Icon size={22} strokeWidth={2} />
            <span className="text-[11px] font-medium">{tab.label}</span>
          </button>
        );
      })}
    </nav>
  );
}
