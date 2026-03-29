import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { LoginPage } from './LoginPage';

const loginMock = vi.fn();
const navigateMock = vi.fn();

vi.mock('../features/auth/AuthContext', () => ({
  useAuth: () => ({
    login: loginMock,
  }),
}));

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom');
  return {
    ...actual,
    useNavigate: () => navigateMock,
  };
});

describe('LoginPage', () => {
  afterEach(() => {
    loginMock.mockReset();
    navigateMock.mockReset();
  });

  it('提交登录表单后调用登录逻辑并跳转', async () => {
    loginMock.mockResolvedValue(undefined);
    const user = userEvent.setup();

    render(
      <MemoryRouter initialEntries={['/login?redirect=%2Fusers']}>
        <LoginPage />
      </MemoryRouter>,
    );

    await user.type(screen.getByLabelText('账号'), 'admin');
    await user.type(screen.getByLabelText('密码'), 'secret');
    await user.click(screen.getByRole('button', { name: '登录进入工作台' }));

    expect(loginMock).toHaveBeenCalledWith('admin', 'secret');
    expect(navigateMock).toHaveBeenCalledWith('/users');
  });

  it('登录失败时展示错误信息', async () => {
    loginMock.mockRejectedValue(new Error('账号或密码错误'));
    const user = userEvent.setup();

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>,
    );

    await user.type(screen.getByLabelText('账号'), 'admin');
    await user.type(screen.getByLabelText('密码'), 'wrong');
    await user.click(screen.getByRole('button', { name: '登录进入工作台' }));

    expect(await screen.findByText('账号或密码错误')).toBeInTheDocument();
  });
});
