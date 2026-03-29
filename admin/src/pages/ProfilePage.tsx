import { PageHeader } from '@/components/shared/page-header';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { CalendarClock, Mail, Save, Shield, Sparkles, UserCircle2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { api, isApiError } from '../api/client';
import { useAuth } from '../features/auth/AuthContext';
import { formatDateTime } from '../utils/format';

export function ProfilePage() {
  const { user, reload } = useAuth();
  const [displayName, setDisplayName] = useState(user?.display_name ?? '');
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    setDisplayName(user?.display_name ?? '');
  }, [user?.display_name]);

  if (!user) {
    return null;
  }

  const profileName = user.display_name || user.username;

  return (
    <div className="space-y-6">
      <PageHeader
        badge="Profile"
        title="个人中心"
        description="维护当前账号的展示信息、角色状态和基础资料，保存后会立即刷新当前登录态。"
      />

      <section className="grid gap-4 xl:grid-cols-[0.95fr_1.05fr]">
        <Card className="overflow-hidden bg-gradient-to-br from-white via-white to-teal-50">
          <CardContent className="space-y-6 p-6">
            <div className="flex flex-col gap-5 sm:flex-row sm:items-center">
              <Avatar className="h-20 w-20 border border-primary/10 shadow-soft">
                <AvatarFallback className="text-xl">
                  {profileName.charAt(0).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <div className="space-y-3">
                <div>
                  <h2 className="text-2xl font-semibold text-foreground">{profileName}</h2>
                  <p className="mt-1 flex items-center gap-2 text-sm text-muted-foreground">
                    <Mail className="h-4 w-4" />
                    {user.email}
                  </p>
                </div>
                <div className="flex flex-wrap gap-2">
                  <Badge variant={user.role === 'admin' ? 'warning' : 'info'}>{user.role}</Badge>
                  <Badge variant={user.status === 'active' ? 'success' : 'outline'}>
                    {user.status}
                  </Badge>
                  <Badge variant={user.email_verified ? 'success' : 'outline'}>
                    {user.email_verified ? '邮箱已验证' : '邮箱未验证'}
                  </Badge>
                </div>
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              {[
                { label: '用户名', value: user.username, icon: UserCircle2 },
                { label: '创建时间', value: formatDateTime(user.ctime), icon: CalendarClock },
                { label: '更新时间', value: formatDateTime(user.utime), icon: Sparkles },
                { label: '权限角色', value: user.role, icon: Shield },
              ].map((item) => {
                const Icon = item.icon;
                return (
                  <div
                    key={item.label}
                    className="rounded-[1.4rem] border border-border/70 bg-background/70 p-4"
                  >
                    <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-2xl bg-primary/10 text-primary">
                      <Icon className="h-4 w-4" />
                    </div>
                    <div className="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
                      {item.label}
                    </div>
                    <div className="mt-2 text-base font-medium text-foreground">{item.value}</div>
                  </div>
                );
              })}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>更新资料</CardTitle>
            <CardDescription>
              显示名会出现在后台头部、个人资料卡片以及部分管理交互中。
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form
              className="space-y-5"
              onSubmit={async (event) => {
                event.preventDefault();
                try {
                  setSaving(true);
                  await api.updateMe(displayName);
                  await reload();
                  toast.success('个人信息已更新');
                } catch (error) {
                  if (!isApiError(error)) {
                    toast.error((error as Error).message);
                  }
                } finally {
                  setSaving(false);
                }
              }}
            >
              <div className="space-y-2">
                <Label htmlFor="display-name">显示名</Label>
                <Input
                  id="display-name"
                  value={displayName}
                  onChange={(event) => setDisplayName(event.target.value)}
                  placeholder="请输入对外展示的名称"
                  required
                />
              </div>

              <div className="rounded-[1.4rem] border border-border/70 bg-muted/20 p-4 text-sm leading-6 text-muted-foreground">
                更新完成后，导航头部和个人资料展示会自动同步，无需重新登录。
              </div>

              <Button type="submit" disabled={saving}>
                <Save className="h-4 w-4" />
                {saving ? '保存中...' : '保存变更'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </section>
    </div>
  );
}
