import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { ProtectedLayout } from './App';

const useAuthMock = vi.fn();

vi.mock('../features/auth/AuthContext', () => ({
  useAuth: () => useAuthMock(),
  AuthProvider: ({ children }: { children: React.ReactNode }) => children,
}));

vi.mock('../components/AppShell', () => ({
  AppShell: () => <div>app shell</div>,
}));

describe('ProtectedLayout', () => {
  afterEach(() => {
    useAuthMock.mockReset();
  });

  it('在未恢复登录态时展示加载文案', () => {
    useAuthMock.mockReturnValue({ user: null, ready: false });

    render(
      <MemoryRouter initialEntries={['/users']}>
        <Routes>
          <Route path="/users" element={<ProtectedLayout />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByText('正在恢复登录态...')).toBeInTheDocument();
  });

  it('未登录时跳转到登录页', () => {
    useAuthMock.mockReturnValue({ user: null, ready: true });

    render(
      <MemoryRouter initialEntries={['/users?tab=all']}>
        <Routes>
          <Route path="/login" element={<div>login page</div>} />
          <Route path="/users" element={<ProtectedLayout />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByText('login page')).toBeInTheDocument();
  });

  it('已登录时渲染工作台', () => {
    useAuthMock.mockReturnValue({ user: { uid: 1 }, ready: true });

    render(
      <MemoryRouter initialEntries={['/users']}>
        <Routes>
          <Route path="/users" element={<ProtectedLayout />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByText('app shell')).toBeInTheDocument();
  });
});
