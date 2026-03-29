import { PageHeader } from '@/components/shared/page-header';
import { StatCard } from '@/components/shared/stat-card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { cn } from '@/lib/utils';
import { useQuery } from '@tanstack/react-query';
import {
  Activity,
  ArrowUpRight,
  FileText,
  RefreshCw,
  ShieldCheck,
  Users,
  Workflow,
} from 'lucide-react';
import { Link } from 'react-router-dom';
import { api } from '../api/client';
import { formatFileSize, formatTextFallback } from '../utils/format';

function getServiceStatusVariant(status?: string) {
  if (status === 'ok') {
    return 'success' as const;
  }
  if (!status) {
    return 'outline' as const;
  }
  return 'warning' as const;
}

export function DashboardPage() {
  const status = useQuery({ queryKey: ['status'], queryFn: api.systemStatus });
  const about = useQuery({ queryKey: ['about'], queryFn: api.about });
  const notice = useQuery({ queryKey: ['notice'], queryFn: api.notice });
  const users = useQuery({ queryKey: ['dashboard-users'], queryFn: api.listUsers });
  const files = useQuery({ queryKey: ['dashboard-files'], queryFn: api.listFiles });

  const serviceStatus = status.data?.status;
  const userTotal = users.data?.total ?? users.data?.items.length ?? 0;
  const fileTotal = files.data?.total ?? files.data?.items.length ?? 0;
  const storageTotal = files.data?.items.reduce((sum, item) => sum + item.size, 0) ?? 0;

  function reloadAll() {
    void Promise.all([
      status.refetch(),
      about.refetch(),
      notice.refetch(),
      users.refetch(),
      files.refetch(),
    ]);
  }

  return (
    <div className="space-y-6">
      <PageHeader
        badge="Workspace"
        title="系统概览"
        description="把服务状态、资源规模、公告内容和联调入口集中到一个后台工作台里，"
        actions={
          <>
            <Button
              variant="outline"
              onClick={reloadAll}
              disabled={status.isFetching || about.isFetching || notice.isFetching}
            >
              <RefreshCw
                className={cn(
                  'h-4 w-4',
                  (status.isFetching || about.isFetching || notice.isFetching) && 'animate-spin',
                )}
              />
              刷新数据
            </Button>
            <Button asChild>
              <a href="/docs" target="_blank" rel="noreferrer">
                <Workflow className="h-4 w-4" />
                打开 OpenAPI
              </a>
            </Button>
          </>
        }
      />

      <section className="grid gap-4 xl:grid-cols-[minmax(0,1.4fr)_repeat(3,minmax(0,1fr))]">
        <Card className="overflow-hidden border-0 bg-hero-grid text-white shadow-glow xl:col-span-1">
          <CardContent className="space-y-8 p-8">
            <div className="flex flex-wrap items-center gap-3">
              <Badge
                variant={getServiceStatusVariant(serviceStatus)}
                className="border-white/10 bg-white/10 text-white"
              >
                服务状态 {formatTextFallback(serviceStatus)}
              </Badge>
              <Badge variant="outline" className="border-white/10 bg-white/5 text-slate-200">
                Trace Ready
              </Badge>
            </div>
            <div className="flex flex-wrap gap-3">
              <Button variant="secondary" onClick={reloadAll}>
                <RefreshCw className="h-4 w-4" />
                立即刷新
              </Button>
              <Button
                asChild
                variant="outline"
                className="border-white/15 bg-white/5 text-white hover:bg-white/10 hover:text-white"
              >
                <Link to="/users">
                  查看用户趋势
                  <ArrowUpRight className="h-4 w-4" />
                </Link>
              </Button>
            </div>
          </CardContent>
        </Card>

        <StatCard
          label="用户总数"
          value={userTotal}
          hint="系统内已注册账号"
          icon={Users}
          tone="primary"
        />
        <StatCard label="文件总数" value={fileTotal} hint="当前已托管文件数" icon={FileText} />
        <StatCard
          label="占用存储"
          value={formatFileSize(storageTotal)}
          hint="按当前结果实时计算"
          icon={Activity}
          tone="accent"
        />
      </section>

      <section className="grid gap-4 lg:grid-cols-[0.95fr_1.05fr]">
        <Card className="overflow-hidden">
          <CardHeader>
            <CardTitle>系统公告</CardTitle>
          </CardHeader>
          <CardContent>
            {notice.isLoading ? (
              <div className="space-y-3">
                <Skeleton className="h-5 w-40" />
                <Skeleton className="h-5 w-full" />
                <Skeleton className="h-5 w-5/6" />
              </div>
            ) : (
              <p className="whitespace-pre-wrap text-sm leading-7 text-muted-foreground">
                {notice.data?.value?.trim() || '暂无公告'}
              </p>
            )}
          </CardContent>
        </Card>

        <Card className="overflow-hidden">
          <CardHeader>
            <CardTitle>系统介绍</CardTitle>
          </CardHeader>
          <CardContent>
            {about.isLoading ? (
              <div className="space-y-3">
                <Skeleton className="h-5 w-32" />
                <Skeleton className="h-5 w-full" />
                <Skeleton className="h-5 w-full" />
                <Skeleton className="h-5 w-3/4" />
              </div>
            ) : (
              <p className="whitespace-pre-wrap text-sm leading-7 text-muted-foreground">
                {about.data?.value?.trim() || '暂无介绍'}
              </p>
            )}
          </CardContent>
        </Card>
      </section>
    </div>
  );
}
