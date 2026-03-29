import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter, Navigate, Route, Routes, useLocation } from 'react-router-dom';
import { Toaster } from 'sonner';
import { AppShell } from '../components/AppShell';
import { LoadingScreen } from '../components/shared/loading-screen';
import { AuthProvider, useAuth } from '../features/auth/AuthContext';
import { DashboardPage } from '../pages/DashboardPage';
import { DocsPage } from '../pages/DocsPage';
import { FilesPage } from '../pages/FilesPage';
import { LoginPage } from '../pages/LoginPage';
import { OAuthCallbackPage } from '../pages/OAuthCallbackPage';
import { OptionsPage } from '../pages/OptionsPage';
import { PprofPage } from '../pages/PprofPage';
import { ProfilePage } from '../pages/ProfilePage';
import { UsersPage } from '../pages/UsersPage';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      staleTime: 10_000,
    },
  },
});

export function ProtectedLayout() {
  const { user, ready } = useAuth();
  const location = useLocation();
  if (!ready) {
    return <LoadingScreen message="正在恢复登录态..." fullscreen />;
  }
  if (!user) {
    const redirect = `${location.pathname}${location.search}${location.hash}`;
    return <Navigate to={`/login?redirect=${encodeURIComponent(redirect)}`} replace />;
  }
  return <AppShell />;
}

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/oauth/callback" element={<OAuthCallbackPage />} />
            <Route element={<ProtectedLayout />}>
              <Route path="/" element={<DashboardPage />} />
              <Route path="/users" element={<UsersPage />} />
              <Route path="/files" element={<FilesPage />} />
              <Route path="/options" element={<OptionsPage />} />
              <Route path="/profile" element={<ProfilePage />} />
              <Route path="/docs" element={<DocsPage />} />
              <Route path="/pprof" element={<PprofPage />} />
            </Route>
          </Routes>
        </BrowserRouter>
        <Toaster richColors position="top-right" closeButton />
      </AuthProvider>
    </QueryClientProvider>
  );
}
