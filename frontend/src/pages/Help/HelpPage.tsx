// src/pages/Help/HelpPage.tsx - Redirect to SourcePicker
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/config/routes';

export function HelpPage() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate(ROUTES.HELP_UPLOAD, { replace: true });
  }, [navigate]);

  return null;
}
