import { PageHeader } from '@/components/shared/page-header';
import { Alert, AlertDescription, AlertIcon, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { cn } from '@/lib/utils';
import { useQuery } from '@tanstack/react-query';
import { Activity, ExternalLink, RefreshCw } from 'lucide-react';
import { api } from '../api/client';

export function PprofPage() {
  const pprof = useQuery({ queryKey: ['pprof-url'], queryFn: api.pprofURL });
  const pprofURL = pprof.data?.url;

  return (
    <div className="space-y-6">
      <PageHeader
        badge="Profiler"
        title="pprof 监控"
        description="通过系统接口动态获取监控地址，在后台内嵌查看分析页，同时支持单独新窗口打开。"
        actions={
          <>
            <Button variant="outline" onClick={() => pprof.refetch()} disabled={pprof.isFetching}>
              <RefreshCw className={cn('h-4 w-4', pprof.isFetching && 'animate-spin')} />
              刷新地址
            </Button>
            <Button asChild disabled={!pprofURL}>
              <a href={pprofURL || '#'} target="_blank" rel="noreferrer">
                <ExternalLink className="h-4 w-4" />
                新窗口打开
              </a>
            </Button>
          </>
        }
      />

      <Card className="overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center justify-between gap-4 border-b border-border/60 px-6 py-4">
            <div>
              <div className="text-sm font-semibold text-foreground">内嵌监控视图</div>
              <div className="text-sm text-muted-foreground">
                {pprofURL ? `当前地址 ${pprofURL}` : '正在加载监控地址'}
              </div>
            </div>
            <div className="hidden items-center gap-2 rounded-full border border-primary/15 bg-primary/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-primary sm:flex">
              <Activity className="h-3.5 w-3.5" />
              PPROF READY
            </div>
          </div>

          {pprof.isLoading ? (
            <div className="space-y-3 p-6">
              <Skeleton className="h-5 w-40" />
              <Skeleton className="h-5 w-full" />
              <Skeleton className="h-[78vh] w-full" />
            </div>
          ) : pprof.isError ? (
            <div className="p-6">
              <Alert variant="destructive">
                <AlertIcon className="h-4 w-4" />
                <AlertTitle>监控地址加载失败</AlertTitle>
                <AlertDescription>{(pprof.error as Error).message}</AlertDescription>
              </Alert>
            </div>
          ) : (
            <iframe
              title="pprof"
              src={pprofURL}
              className="min-h-[78vh] w-full border-0 bg-white"
            />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
