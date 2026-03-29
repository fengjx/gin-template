import { createContext, useContext, useEffect, useState } from 'react';
import {
  api,
  clearAccessToken,
  refreshSession,
  setAccessToken,
  setApiAuthFailureHandler,
} from '../../api/client';
import type { AuthResponse, User } from '../../api/types';

type AuthContextValue = {
  user: User | null;
  ready: boolean;
  login: (identifier: string, password: string) => Promise<void>;
  register: (
    username: string,
    email: string,
    password: string,
    displayName?: string,
  ) => Promise<void>;
  logout: () => Promise<void>;
  reload: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

function applyAuth(data: AuthResponse | null, setUser: (value: User | null) => void) {
  if (!data) {
    clearAccessToken();
    setUser(null);
    return;
  }
  setAccessToken(data.access_token);
  setUser(data.user);
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    refreshSession()
      .then((data) => applyAuth(data, setUser))
      .finally(() => setReady(true));
  }, []);

  useEffect(() => {
    setApiAuthFailureHandler(() => {
      setUser(null);
    });
    return () => {
      setApiAuthFailureHandler(null);
    };
  }, []);

  async function login(identifier: string, password: string) {
    const data = await api.login({ identifier, password });
    applyAuth(data, setUser);
  }

  async function register(username: string, email: string, password: string, displayName?: string) {
    const data = await api.register({ username, email, password, display_name: displayName });
    applyAuth(data, setUser);
  }

  async function logout() {
    await api.logout();
    applyAuth(null, setUser);
  }

  async function reload() {
    const me = await api.me();
    setUser(me);
  }

  return (
    <AuthContext.Provider value={{ user, ready, login, register, logout, reload }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('AuthContext 未初始化');
  }
  return context;
}
