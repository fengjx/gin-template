import { type DataColumn, DataTable } from '@/components/shared/data-table';
import { PageHeader } from '@/components/shared/page-header';
import { StatCard } from '@/components/shared/stat-card';
import { Alert, AlertDescription, AlertIcon, AlertTitle } from '@/components/ui/alert';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { Switch } from '@/components/ui/switch';
import { cn } from '@/lib/utils';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Edit3,
  Plus,
  RefreshCw,
  Search,
  ShieldCheck,
  Trash2,
  UserCheck,
  Users as UsersIcon,
} from 'lucide-react';
import { startTransition, useMemo, useState } from 'react';
import { toast } from 'sonner';
import { api, isApiError } from '../api/client';
import type { User, UserCreateRequest, UserUpdateRequest } from '../api/types';
import { useAuth } from '../features/auth/AuthContext';
import { formatDateTime, formatTextFallback } from '../utils/format';

type UserFormValues = {
  username: string;
  email: string;
  password: string;
  display_name: string;
  role: string;
  status: string;
  email_verified: boolean;
};

const defaultFormValues: UserFormValues = {
  username: '',
  email: '',
  password: '',
  display_name: '',
  role: 'user',
  status: 'active',
  email_verified: false,
};

const roleOptions = [
  { value: 'user', label: '普通用户' },
  { value: 'admin', label: '管理员' },
  { value: 'root', label: 'Root' },
];

const statusOptions = [
  { value: 'active', label: '正常' },
  { value: 'locked', label: '禁用' },
];

export function UsersPage() {
  const { user: currentUser } = useAuth();
  const [inputKeyword, setInputKeyword] = useState('');
  const [keyword, setKeyword] = useState('');
  const [sheetOpen, setSheetOpen] = useState(false);
  const [sheetMode, setSheetMode] = useState<'create' | 'edit'>('create');
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [deletingUser, setDeletingUser] = useState<User | null>(null);
  const [formValues, setFormValues] = useState<UserFormValues>(defaultFormValues);
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: ['users', keyword],
    queryFn: () => (keyword ? api.searchUsers(keyword) : api.listUsers()),
  });

  const saveMutation = useMutation({
    mutationFn: async (payload: {
      mode: 'create' | 'edit';
      uid?: number;
      values: UserFormValues;
    }) => {
      if (payload.mode === 'create') {
        const body: UserCreateRequest = {
          username: payload.values.username.trim(),
          email: payload.values.email.trim(),
          password: payload.values.password,
          display_name: payload.values.display_name.trim() || undefined,
          role: payload.values.role,
          status: payload.values.status,
          email_verified: payload.values.email_verified,
        };
        return api.createUser(body);
      }

      const body: UserUpdateRequest = {
        email: payload.values.email.trim(),
        display_name: payload.values.display_name.trim(),
        role: payload.values.role,
        status: payload.values.status,
        email_verified: payload.values.email_verified,
      };
      if (payload.values.password.trim()) {
        body.password = payload.values.password;
      }
      if (!payload.uid) {
        throw new Error('编辑用户时缺少 uid');
      }
      return api.updateUser(payload.uid, body);
    },
    onSuccess: async (_, payload) => {
      toast.success(payload.mode === 'create' ? '用户已创建' : '用户已更新');
      await queryClient.invalidateQueries({ queryKey: ['users'] });
      closeSheet();
    },
    onError: (error) => {
      if (!isApiError(error)) {
        toast.error((error as Error).message);
      }
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (uid: number) => api.deleteUser(uid),
    onSuccess: async () => {
      toast.success('用户已删除');
      await queryClient.invalidateQueries({ queryKey: ['users'] });
      setDeletingUser(null);
    },
    onError: (error) => {
      if (!isApiError(error)) {
        toast.error((error as Error).message);
      }
    },
  });

  const totalUsers = query.data?.total ?? query.data?.items.length ?? 0;
  const adminUsers = query.data?.items.filter((item) => item.role === 'admin').length ?? 0;
  const canManageRoot = currentUser?.role === 'root';

  const availableRoles = useMemo(
    () => roleOptions.filter((option) => canManageRoot || option.value !== 'root'),
    [canManageRoot],
  );

  const columns: DataColumn<User>[] = [
    {
      key: 'username',
      header: '用户名',
      cell: (record) => <div className="font-medium text-foreground">{record.username}</div>,
      minWidth: '160px',
    },
    {
      key: 'email',
      header: '邮箱',
      cell: (record) => <span className="text-muted-foreground">{record.email}</span>,
      minWidth: '220px',
    },
    {
      key: 'display_name',
      header: '显示名',
      cell: (record) => formatTextFallback(record.display_name),
      minWidth: '160px',
    },
    {
      key: 'role',
      header: '角色',
      cell: (record) => (
        <Badge
          variant={
            record.role === 'root' ? 'destructive' : record.role === 'admin' ? 'warning' : 'info'
          }
        >
          {record.role}
        </Badge>
      ),
      minWidth: '120px',
    },
    {
      key: 'status',
      header: '状态',
      cell: (record) => (
        <Badge variant={record.status === 'active' ? 'success' : 'outline'}>{record.status}</Badge>
      ),
      minWidth: '120px',
    },
    {
      key: 'email_verified',
      header: '邮箱验证',
      cell: (record) => (
        <Badge variant={record.email_verified ? 'success' : 'outline'}>
          {record.email_verified ? '已验证' : '未验证'}
        </Badge>
      ),
      minWidth: '140px',
    },
    {
      key: 'ctime',
      header: '创建时间',
      cell: (record) => formatDateTime(record.ctime),
      minWidth: '180px',
    },
    {
      key: 'actions',
      header: '操作',
      cell: (record) => (
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => openEditSheet(record)}
            disabled={!canEditUser(record)}
          >
            <Edit3 className="h-4 w-4" />
            编辑
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="text-destructive hover:text-destructive"
            onClick={() => setDeletingUser(record)}
            disabled={!canDeleteUser(record)}
          >
            <Trash2 className="h-4 w-4" />
            删除
          </Button>
        </div>
      ),
      minWidth: '160px',
    },
  ];

  function canEditUser(target: User) {
    if (!currentUser) {
      return false;
    }
    return target.role !== 'root' || currentUser.role === 'root';
  }

  function canDeleteUser(target: User) {
    if (!currentUser) {
      return false;
    }
    if (target.uid === currentUser.uid) {
      return false;
    }
    return target.role !== 'root' || currentUser.role === 'root';
  }

  function closeSheet() {
    setSheetOpen(false);
    setEditingUser(null);
    setFormValues(defaultFormValues);
  }

  function openCreateSheet() {
    setSheetOpen(true);
    setSheetMode('create');
    setEditingUser(null);
    setFormValues(defaultFormValues);
  }

  function openEditSheet(target: User) {
    setSheetOpen(true);
    setSheetMode('edit');
    setEditingUser(target);
    setFormValues({
      username: target.username,
      email: target.email,
      password: '',
      display_name: target.display_name ?? '',
      role: canManageRoot || target.role !== 'root' ? target.role : 'admin',
      status: target.status,
      email_verified: target.email_verified,
    });
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (sheetMode === 'create') {
      await saveMutation.mutateAsync({ mode: 'create', values: formValues });
      return;
    }
    if (!editingUser) {
      return;
    }
    await saveMutation.mutateAsync({ mode: 'edit', uid: editingUser.uid, values: formValues });
  }

  return (
    <div className="space-y-6">
      <PageHeader
        badge="Identity"
        title="用户管理"
        description="统一查看账号状态、身份角色、邮箱验证与注册时间，支持按用户名或邮箱快速检索。"
        actions={
          <div className="flex flex-wrap items-center gap-3">
            <Button
              variant="outline"
              onClick={() => query.refetch()}
              disabled={query.isFetching || saveMutation.isPending || deleteMutation.isPending}
            >
              <RefreshCw className={cn('h-4 w-4', query.isFetching && 'animate-spin')} />
              刷新列表
            </Button>
            <Button onClick={openCreateSheet}>
              <Plus className="h-4 w-4" />
              新增用户
            </Button>
          </div>
        }
      />

      <section className="grid gap-4 md:grid-cols-3">
        <StatCard
          label="当前结果数"
          value={totalUsers}
          hint="按照当前筛选条件统计"
          icon={UsersIcon}
          tone="primary"
        />
        <StatCard
          label="管理员数量"
          value={adminUsers}
          hint="角色为 admin 的账号数"
          icon={ShieldCheck}
        />
        <StatCard
          label="检索条件"
          value={keyword || '全部用户'}
          hint="支持搜索用户名或邮箱"
          icon={UserCheck}
          tone="accent"
        />
      </section>

      {query.isError ? (
        <Alert variant="destructive">
          <AlertIcon className="h-4 w-4" />
          <AlertTitle>用户列表加载失败</AlertTitle>
          <AlertDescription>{(query.error as Error).message}</AlertDescription>
        </Alert>
      ) : null}

      <DataTable<User>
        data={query.data?.items ?? []}
        columns={columns}
        rowKey={(row) => String(row.uid)}
        loading={query.isFetching || saveMutation.isPending || deleteMutation.isPending}
        emptyTitle="没有找到匹配的用户"
        emptyDescription="可以尝试更换关键词，或清空检索条件后重新查看全部用户。"
        toolbar={
          <form
            className="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between"
            onSubmit={(event) => {
              event.preventDefault();
              startTransition(() => setKeyword(inputKeyword.trim()));
            }}
          >
            <div className="flex flex-1 flex-col gap-3 sm:flex-row">
              <Input
                value={inputKeyword}
                onChange={(event) => {
                  const value = event.target.value;
                  setInputKeyword(value);
                  if (!value) {
                    startTransition(() => setKeyword(''));
                  }
                }}
                placeholder="搜索用户名或邮箱"
                className="sm:max-w-md"
              />
              <Button type="submit">
                <Search className="h-4 w-4" />
                搜索
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setInputKeyword('');
                  startTransition(() => setKeyword(''));
                }}
              >
                清空筛选
              </Button>
            </div>
            <div className="text-sm text-muted-foreground">共 {totalUsers} 名用户</div>
          </form>
        }
      />

      <Sheet
        open={sheetOpen}
        onOpenChange={(open) => {
          if (!open) {
            closeSheet();
          }
          setSheetOpen(open);
        }}
      >
        <SheetContent side="right" className="w-full max-w-2xl overflow-y-auto">
          <SheetHeader>
            <SheetTitle>
              {sheetMode === 'create' ? '新增用户' : `编辑用户 · ${editingUser?.username ?? ''}`}
            </SheetTitle>
            <SheetDescription>
              {sheetMode === 'create'
                ? '创建新的系统账号，并设置角色、状态和邮箱验证状态。'
                : '维护用户的邮箱、显示名、角色、状态与密码信息。'}
            </SheetDescription>
          </SheetHeader>

          <form id="user-form" className="space-y-5 py-8" onSubmit={handleSubmit}>
            <div className="grid gap-5 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="user-username">用户名</Label>
                <Input
                  id="user-username"
                  value={formValues.username}
                  onChange={(event) =>
                    setFormValues((value) => ({ ...value, username: event.target.value }))
                  }
                  disabled={sheetMode === 'edit'}
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="user-email">邮箱</Label>
                <Input
                  id="user-email"
                  type="email"
                  value={formValues.email}
                  onChange={(event) =>
                    setFormValues((value) => ({ ...value, email: event.target.value }))
                  }
                  required
                />
              </div>
            </div>

            <div className="grid gap-5 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="user-display-name">显示名</Label>
                <Input
                  id="user-display-name"
                  value={formValues.display_name}
                  onChange={(event) =>
                    setFormValues((value) => ({ ...value, display_name: event.target.value }))
                  }
                  placeholder="选填，用于页面展示"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="user-password">
                  {sheetMode === 'create' ? '登录密码' : '重置密码'}
                </Label>
                <Input
                  id="user-password"
                  type="password"
                  value={formValues.password}
                  onChange={(event) =>
                    setFormValues((value) => ({ ...value, password: event.target.value }))
                  }
                  placeholder={sheetMode === 'create' ? '请输入初始密码' : '留空表示不修改'}
                  required={sheetMode === 'create'}
                />
              </div>
            </div>

            <div className="grid gap-5 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="user-role">角色</Label>
                <select
                  id="user-role"
                  value={formValues.role}
                  onChange={(event) =>
                    setFormValues((value) => ({ ...value, role: event.target.value }))
                  }
                  className="h-11 w-full rounded-full border border-input bg-background px-4 text-sm text-foreground outline-none ring-offset-background transition-colors focus:ring-2 focus:ring-ring focus:ring-offset-2"
                >
                  {availableRoles.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="user-status">状态</Label>
                <select
                  id="user-status"
                  value={formValues.status}
                  onChange={(event) =>
                    setFormValues((value) => ({ ...value, status: event.target.value }))
                  }
                  className="h-11 w-full rounded-full border border-input bg-background px-4 text-sm text-foreground outline-none ring-offset-background transition-colors focus:ring-2 focus:ring-ring focus:ring-offset-2"
                >
                  {statusOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="flex items-center justify-between rounded-[1.4rem] border border-border/70 bg-muted/25 px-4 py-4">
              <div>
                <div className="font-medium text-foreground">邮箱已验证</div>
                <div className="text-sm text-muted-foreground">
                  开启后，系统将该账号视为已完成邮箱验证。
                </div>
              </div>
              <Switch
                checked={formValues.email_verified}
                onCheckedChange={(checked) =>
                  setFormValues((value) => ({ ...value, email_verified: checked }))
                }
              />
            </div>
          </form>

          <SheetFooter>
            <Button variant="outline" onClick={closeSheet}>
              取消
            </Button>
            <Button form="user-form" type="submit" disabled={saveMutation.isPending}>
              {saveMutation.isPending
                ? '保存中...'
                : sheetMode === 'create'
                  ? '创建用户'
                  : '保存变更'}
            </Button>
          </SheetFooter>
        </SheetContent>
      </Sheet>

      <AlertDialog open={!!deletingUser} onOpenChange={(open) => !open && setDeletingUser(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>确认删除用户</AlertDialogTitle>
            <AlertDialogDescription>
              {deletingUser
                ? `即将删除用户 ${deletingUser.username}，该操作不可撤销。`
                : '该操作不可撤销。'}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteMutation.isPending}>取消</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              disabled={deleteMutation.isPending}
              onClick={async () => {
                if (!deletingUser) {
                  return;
                }
                await deleteMutation.mutateAsync(deletingUser.uid);
              }}
            >
              {deleteMutation.isPending ? '删除中...' : '确认删除'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
