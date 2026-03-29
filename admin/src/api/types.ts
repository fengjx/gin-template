import type { components } from './generated';

export type ApiEnvelope<T> = {
  status: number;
  msg: string;
  details?: string;
  data: T;
};

export type AuthResponse = components['schemas']['AuthResponse'];
export type FileItem = components['schemas']['File'];
export type FileListResponse = components['schemas']['FileListResponse'];
export type MessageResponse = components['schemas']['MessageResponse'];
export type OptionItem = components['schemas']['Option'];
export type Problem = ApiEnvelope<null>;
export type User = components['schemas']['User'];
export type UserCreateRequest = components['schemas']['UserCreateRequest'];
export type UserListResponse = components['schemas']['UserListResponse'];
export type UserUpdateRequest = components['schemas']['UserUpdateRequest'];
