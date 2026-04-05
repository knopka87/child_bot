// src/pages/Check/CheckPage.tsx - Redirect to ScenarioSelection
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/config/routes';

export function CheckPage() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate(ROUTES.CHECK_SCENARIO, { replace: true });
  }, [navigate]);

  return null;
}
