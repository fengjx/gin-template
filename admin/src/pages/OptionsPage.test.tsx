import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { OptionsPage } from './OptionsPage';

const { createOptionMock, listOptionsMock, toastErrorMock, toastSuccessMock, updateOptionMock } =
  vi.hoisted(() => ({
    listOptionsMock: vi.fn(),
    createOptionMock: vi.fn(),
    updateOptionMock: vi.fn(),
    toastSuccessMock: vi.fn(),
    toastErrorMock: vi.fn(),
  }));

vi.mock('../api/client', () => ({
  api: {
    listOptions: () => listOptionsMock(),
    createOption: (payload: unknown) => createOptionMock(payload),
    updateOption: (payload: unknown) => updateOptionMock(payload),
  },
  isApiError: () => false,
}));

vi.mock('sonner', () => ({
  toast: {
    success: toastSuccessMock,
    error: toastErrorMock,
  },
}));

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <OptionsPage />
    </QueryClientProvider>,
  );
}

describe('OptionsPage', () => {
  afterEach(() => {
    listOptionsMock.mockReset();
    createOptionMock.mockReset();
    updateOptionMock.mockReset();
    toastSuccessMock.mockReset();
    toastErrorMock.mockReset();
  });

  it('可以新增 string 类型配置', async () => {
    listOptionsMock.mockResolvedValue([
      {
        id: 1,
        option_key: 'about',
        option_value: 'Gin + React 同构脚手架',
        description: '关于信息',
        is_public: true,
        type: 'string',
        status: 'online',
        ctime: 1,
        utime: 1,
      },
    ]);
    createOptionMock.mockResolvedValue({
      id: 2,
      option_key: 'site_name',
      option_value: 'Gin Template',
      description: '站点名称',
      is_public: true,
      type: 'string',
      status: 'online',
      ctime: 1,
      utime: 1,
    });

    const user = userEvent.setup();
    renderPage();

    await user.click(await screen.findByRole('button', { name: '新增配置' }));
    await user.type(screen.getByLabelText('Key'), 'site_name');
    await user.type(screen.getByLabelText('值'), 'Gin Template');
    await user.type(screen.getByLabelText('描述'), '站点名称');
    await user.click(screen.getByRole('switch'));
    await user.click(screen.getByRole('button', { name: '创建配置' }));

    await waitFor(() =>
      expect(createOptionMock).toHaveBeenCalledWith({
        key: 'site_name',
        value: 'Gin Template',
        description: '站点名称',
        is_public: true,
        type: 'string',
        status: 'online',
      }),
    );
    expect(toastSuccessMock).toHaveBeenCalledWith('配置已创建');
  });

  it('json 类型非法时阻止提交并展示错误', async () => {
    listOptionsMock.mockResolvedValue([]);
    const user = userEvent.setup();

    renderPage();

    await user.click(await screen.findByRole('button', { name: '新增配置' }));
    await user.type(screen.getByLabelText('Key'), 'site_profile');
    await user.selectOptions(screen.getByLabelText('类型'), 'json');
    fireEvent.change(screen.getByLabelText('值'), { target: { value: '{invalid' } });

    expect(await screen.findByText('JSON 格式无效')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '创建配置' })).toBeDisabled();
    expect(createOptionMock).not.toHaveBeenCalled();
  });
});
