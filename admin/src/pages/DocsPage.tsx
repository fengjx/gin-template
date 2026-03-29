import { PageHeader } from '@/components/shared/page-header';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { ExternalLink, Workflow } from 'lucide-react';

export function DocsPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        badge="OpenAPI"
        title="接口文档"
        description="通过 OpenAPI 文档内嵌在当前后台中，方便前后端联调、接口排查和快速跳转。"
        actions={
          <Button asChild>
            <a href="/docs" target="_blank" rel="noreferrer">
              <ExternalLink className="h-4 w-4" />
              新窗口打开
            </a>
          </Button>
        }
      />

      <Card className="overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center justify-between border-b border-border/60 px-6 py-4">
            <div>
              <div className="text-sm font-semibold text-foreground">内嵌文档视图</div>
              <div className="text-sm text-muted-foreground">
                当前地址 `/docs`，可直接用于联调和错误排查。
              </div>
            </div>
            <div className="hidden items-center gap-2 rounded-full border border-primary/15 bg-primary/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-primary sm:flex">
              <Workflow className="h-3.5 w-3.5" />
              API READY
            </div>
          </div>
          <iframe title="docs" src="/docs" className="min-h-[78vh] w-full border-0 bg-white" />
        </CardContent>
      </Card>
    </div>
  );
}
