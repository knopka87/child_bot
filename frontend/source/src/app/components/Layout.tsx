import { Outlet, useLocation, useNavigate } from "react-router";
import { Home, Trophy, Users, User } from "lucide-react";

const tabs = [
  { path: "/", label: "Главная", icon: Home },
  { path: "/achievements", label: "Достижения", icon: Trophy },
  { path: "/friends", label: "Друзья", icon: Users },
  { path: "/profile", label: "Профиль", icon: User },
];

export function Layout() {
  const location = useLocation();
  const navigate = useNavigate();

  const isTabVisible = tabs.some(
    (t) => location.pathname === t.path || location.pathname === t.path + "/"
  );
  const showTabs =
    isTabVisible ||
    location.pathname.startsWith("/profile/") ||
    location.pathname.startsWith("/achievements");

  return (
    <div className="flex justify-center min-h-screen bg-[#E8E4FF]">
      <div className="w-full max-w-[390px] min-h-screen bg-background flex flex-col relative shadow-xl">
        <div className="flex-1 overflow-y-auto pb-20">
          <Outlet />
        </div>
        {showTabs && (
          <nav className="fixed bottom-0 left-1/2 -translate-x-1/2 w-full max-w-[390px] bg-white border-t border-border flex z-50 safe-area-bottom">
            {tabs.map((tab) => {
              const isActive =
                location.pathname === tab.path ||
                (tab.path !== "/" && location.pathname.startsWith(tab.path));
              const Icon = tab.icon;
              return (
                <button
                  key={tab.path}
                  onClick={() => navigate(tab.path)}
                  className={`flex-1 flex flex-col items-center py-2 pt-3 gap-0.5 transition-colors ${
                    isActive ? "text-primary" : "text-muted-foreground"
                  }`}
                >
                  <Icon size={22} />
                  <span className="text-[11px]">{tab.label}</span>
                </button>
              );
            })}
          </nav>
        )}
      </div>
    </div>
  );
}