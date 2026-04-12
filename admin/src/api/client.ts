import { toast } from 'sonner';
import type {
  ApiEnvelope,
  AuthResponse,
  FileItem,
  FileListResponse,
  MessageResponse,
  OptionItem,
  OptionWriteRequest,
  Problem,
  User,
  UserCreateRequest,
  UserListResponse,
  UserUpdateRequest,
} from './types';

type ApiRequestInit = RequestInit & {
  skipErrorTip?: boolean;
};

type AuthFailureHandler = () => void;

const ERROR_TIP_DEDUPE_MS = 3_000;

export class ApiError extends Error {
  status: number;
  problem?: Problem;

  constructor(status: number, message: string, problem?: Problem) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.problem = problem;
  }
}

let accessToken = '';
let authFailureHandler: AuthFailureHandler | null = null;
const recentErrorTips = new Map<string, number>();

export function setAccessToken(token: string) {
  accessToken = token;
}

export function clearAccessToken() {
  accessToken = '';
}

export function setApiAuthFailureHandler(handler: AuthFailureHandler | null) {
  authFailureHandler = handler;
}

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError;
}

export function resetApiClientForTest() {
  accessToken = '';
  authFailureHandler = null;
  recentErrorTips.clear();
}

function tipKey(error: ApiError, title: string, description?: string) {
  return [error.status, error.problem?.msg, title, description].filter(Boolean).join(':');
}

function shouldShowTip(key: string) {
  const now = Date.now();
  for (const [cachedKey, cachedAt] of recentErrorTips.entries()) {
    if (now - cachedAt > ERROR_TIP_DEDUPE_MS) {
      recentErrorTips.delete(cachedKey);
    }
  }
  const previousAt = recentErrorTips.get(key);
  if (previousAt && now - previousAt < ERROR_TIP_DEDUPE_MS) {
    return false;
  }
  recentErrorTips.set(key, now);
  return true;
}

function resolveErrorTip(error: ApiError) {
  const msg = error.problem?.msg?.trim();
  const details = error.problem?.details?.trim() || error.message;

  switch (error.status) {
    case 401:
      return {
        title: '登录状态已失效',
        description: '请重新登录后继续操作。',
      };
    case 403:
      return {
        title: '没有权限执行该操作',
        description: details || '请联系管理员开通对应权限。',
      };
    case 429:
      return {
        title: '操作过于频繁',
        description: details || '请稍后再试。',
      };
    case 500:
      return {
        title: '服务开小差了',
        description: '请稍后重试；如果持续出现，请联系管理员排查。',
      };
    case 400:
      return {
        title: msg || '请求参数无效',
        description: details || '请检查输入内容后重试。',
      };
    default:
      return {
        title: msg || details || '请求失败',
        description: details && details !== msg ? details : undefined,
      };
  }
}

function isApiEnvelope(value: unknown): value is ApiEnvelope<unknown> {
  if (!value || typeof value !== 'object') {
    return false;
  }
  return 'status' in value && 'msg' in value && 'data' in value;
}

function handleApiError(error: ApiError, options: { skipErrorTip?: boolean } = {}) {
  if (error.status === 401) {
    clearAccessToken();
    authFailureHandler?.();
  }

  if (options.skipErrorTip) {
    return;
  }

  const tip = resolveErrorTip(error);
  const key = tipKey(error, tip.title, tip.description);
  if (!shouldShowTip(key)) {
    return;
  }

  toast.error(tip.title, {
    description: tip.description,
  });
}

async function readResponseBody(
  response: Response,
): Promise<ApiEnvelope<unknown> | string | undefined> {
  if (response.status === 204) {
    return undefined;
  }
  const contentType = response.headers.get('Content-Type') || '';
  if (contentType.includes('application/json')) {
    return (await response.json()) as ApiEnvelope<unknown>;
  }
  const text = await response.text();
  return text || undefined;
}

async function parseResponse<T>(
  response: Response,
  options: { skipErrorTip?: boolean } = {},
): Promise<T> {
  const data = await readResponseBody(response);
  if (!response.ok) {
    const problem = isApiEnvelope(data) ? (data as Problem) : undefined;
    const error = new ApiError(
      response.status,
      problem?.details ||
        problem?.msg ||
        (typeof data === 'string' ? data : '') ||
        response.statusText,
      problem,
    );
    handleApiError(error, options);
    throw error;
  }
  if (isApiEnvelope(data)) {
    return data.data as T;
  }
  return data as T;
}

async function performFetch(input: RequestInfo, init: ApiRequestInit = {}) {
  const { skipErrorTip: _skipErrorTip, ...requestInit } = init;
  const headers = new Headers(init.headers);
  if (accessToken) {
    headers.set('Authorization', `Bearer ${accessToken}`);
  }
  headers.set('X-Trace-Id', crypto.randomUUID());
  return fetch(input, {
    ...requestInit,
    credentials: 'include',
    headers,
  });
}

export async function refreshSession(): Promise<AuthResponse | null> {
  const response = await fetch('/api/v1/auth/refresh', {
    method: 'POST',
    credentials: 'include',
    headers: { 'X-Trace-Id': crypto.randomUUID() },
  });
  if (response.status === 204) {
    clearAccessToken();
    return null;
  }
  if (!response.ok) {
    clearAccessToken();
    return null;
  }
  const payload = (await response.json()) as ApiEnvelope<AuthResponse>;
  setAccessToken(payload.data.access_token);
  return payload.data;
}

async function request<T>(input: RequestInfo, init: ApiRequestInit = {}, retry = true): Promise<T> {
  const response = await performFetch(input, init);
  if (response.status === 401 && retry) {
    const refreshed = await refreshSession();
    if (refreshed) {
      return request<T>(input, init, false);
    }
  }
  return parseResponse<T>(response, init);
}

export const api = {
  login: (payload: { identifier: string; password: string }) =>
    request<AuthResponse>('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }),
  register: (payload: {
    username: string;
    email: string;
    password: string;
    display_name?: string;
  }) =>
    request<AuthResponse>('/api/v1/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }),
  logout: () =>
    request<MessageResponse>('/api/v1/auth/logout', {
      method: 'POST',
    }),
  me: () => request<User>('/api/v1/users/me'),
  updateMe: (displayName: string) =>
    request<User>('/api/v1/users/me', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ display_name: displayName }),
    }),
  listUsers: () => request<UserListResponse>('/api/v1/users'),
  searchUsers: (q: string) =>
    request<UserListResponse>(`/api/v1/users/search?q=${encodeURIComponent(q)}`),
  createUser: (payload: UserCreateRequest) =>
    request<User>('/api/v1/users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }),
  updateUser: (uid: number, payload: UserUpdateRequest) =>
    request<User>(`/api/v1/users/${uid}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }),
  deleteUser: (uid: number) =>
    request<MessageResponse>(`/api/v1/users/${uid}`, {
      method: 'DELETE',
    }),
  listFiles: () => request<FileListResponse>('/api/v1/files'),
  searchFiles: (q: string) =>
    request<FileListResponse>(`/api/v1/files/search?q=${encodeURIComponent(q)}`),
  uploadFile: async (file: File) => {
    const body = new FormData();
    body.append('file', file);
    return request<FileItem>('/api/v1/files/upload', { method: 'POST', body });
  },
  deleteFile: (id: number) =>
    request<MessageResponse>(`/api/v1/files/${id}`, {
      method: 'DELETE',
    }),
  listOptions: () => request<OptionItem[]>('/api/v1/options'),
  createOption: (payload: OptionWriteRequest) =>
    request<OptionItem>('/api/v1/options', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }),
  updateOption: (payload: OptionWriteRequest) =>
    request<OptionItem>('/api/v1/options', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }),
  systemStatus: () => request<{ status: string }>('/api/v1/system/status'),
  about: () => request<{ value: string }>('/api/v1/system/about'),
  notice: () => request<{ value: string }>('/api/v1/system/notice'),
  pprofURL: () => request<{ url: string }>('/api/v1/system/pprof-url'),
};
