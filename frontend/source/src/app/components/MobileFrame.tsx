import { Outlet } from "react-router";

export function MobileFrame() {
  return (
    <div className="flex justify-center min-h-screen bg-[#E8E4FF]">
      <div className="w-full max-w-[390px] min-h-screen bg-background flex flex-col shadow-xl">
        <Outlet />
      </div>
    </div>
  );
}
