import { afterEach, describe, expect, it, vi } from 'vitest';

const { toastError } = vi.hoisted(() => ({
  toastError: vi.fn(),
}));

vi.mock('sonner', () => ({
  toast: {
    error: toastError,
  },
}));

import {
  api,
  clearAccessToken,
  isApiError,
  refreshSession,
  resetApiClientForTest,
  setAccessToken,
  setApiAuthFailureHandler,
} from './client';

describe('api client', () => {
  afterEach(() => {
    vi.restoreAllMocks();
    resetApiClientForTest();
    toastError.mockReset();
  });

  it('refreshSession 成功时写入 access token', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(
      new Response(
        JSON.stringify({
          status: 0,
          msg: 'ok',
          data: {
            access_token: 'token-1',
            expires_at: 1767225600,
            user: { uid: 1 },
          },
        }),
        { status: 200, headers: { 'Content-Type': 'application/json' } },
      ),
    );

    const result = await refreshSession();

    expect(result?.access_token).toBe('token-1');
  });

  it('401 时会刷新会话并重试请求', async () => {
    setAccessToken('expired-token');
    const fetchMock = vi.spyOn(globalThis, 'fetch');
    fetchMock
      .mockResolvedValueOnce(
        new Response('{}', { status: 401, headers: { 'Content-Type': 'application/json' } }),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            status: 0,
            msg: 'ok',
            data: {
              access_token: 'fresh-token',
              expires_at: 1767225600,
              user: { uid: 1 },
            },
          }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        ),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({ status: 0, msg: 'ok', data: { uid: 1, username: 'admin' } }),
          {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          },
        ),
      );

    const result = await api.me();

    expect(result).toEqual({ uid: 1, username: 'admin' });
    expect(fetchMock).toHaveBeenCalledTimes(3);
    expect(fetchMock.mock.calls[2]?.[1]).toMatchObject({
      credentials: 'include',
    });
    const headers = new Headers(fetchMock.mock.calls[2]?.[1]?.headers);
    expect(headers.get('Authorization')).toBe('Bearer fresh-token');
  });

  it('会根据错误码统一弹出 tips', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(
      new Response(
        JSON.stringify({
          status: 100104,
          msg: '没有权限',
          details: '需要管理员权限',
          data: null,
        }),
        { status: 403, headers: { 'Content-Type': 'application/json' } },
      ),
    );

    await expect(api.listUsers()).rejects.toSatisfy((error) => isApiError(error));

    expect(toastError).toHaveBeenCalledWith('没有权限执行该操作', {
      description: '需要管理员权限',
    });
  });

  it('鉴权失效时会清理登录态并通知上层', async () => {
    setAccessToken('expired-token');
    const onAuthFailure = vi.fn();
    setApiAuthFailureHandler(onAuthFailure);

    vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            status: 100103,
            msg: '请先登录',
            details: '请先登录',
            data: null,
          }),
          { status: 401, headers: { 'Content-Type': 'application/json' } },
        ),
      )
      .mockResolvedValueOnce(new Response(null, { status: 204 }));

    await expect(api.me()).rejects.toSatisfy((error) => isApiError(error));

    expect(onAuthFailure).toHaveBeenCalledTimes(1);
    expect(toastError).toHaveBeenCalledWith('登录状态已失效', {
      description: '请重新登录后继续操作。',
    });
    clearAccessToken();
  });

  it('相同错误会在短时间内去重提示', async () => {
    vi.useFakeTimers();
    const fetchMock = vi.spyOn(globalThis, 'fetch');
    const problem = () =>
      new Response(
        JSON.stringify({
          status: 100007,
          msg: '服务内部错误',
          details: '服务内部错误',
          data: null,
        }),
        {
          status: 500,
          headers: { 'Content-Type': 'application/json' },
        },
      );

    fetchMock.mockResolvedValueOnce(problem()).mockResolvedValueOnce(problem());

    await expect(api.systemStatus()).rejects.toSatisfy((error) => isApiError(error));
    await expect(api.about()).rejects.toSatisfy((error) => isApiError(error));

    expect(toastError).toHaveBeenCalledTimes(1);
    vi.useRealTimers();
  });
});
