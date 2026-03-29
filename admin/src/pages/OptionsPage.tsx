import { type DataColumn, DataTable } from '@/components/shared/data-table';
import { PageHeader } from '@/components/shared/page-header';
import { StatCard } from '@/components/shared/stat-card';
import { Alert, AlertDescription, AlertIcon, AlertTitle } from '@/components/ui/alert';
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
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Edit3, LockKeyhole, RefreshCw, Settings2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { api, isApiError } from '../api/client';
import type { OptionItem } from '../api/types';
import { formatDateTime, formatTextFallback } from '../utils/format';

const columns: DataColumn<OptionItem>[] = [
  {
    key: 'option_key',
    header: 'Key',
    cell: (record) => <div className="font-medium text-foreground">{record.option_key}</div>,
    minWidth: '240px',
  },
  {
    key: 'option_value',
    header: '值',
    cell: (record) => (
      <p className="max-w-[420px] truncate text-muted-foreground">
        {formatTextFallback(record.option_value)}
      </p>
    ),
    minWidth: '320px',
  },
  {
    key: 'description',
    header: '描述',
    cell: (record) => formatTextFallback(record.description),
    minWidth: '240px',
  },
  {
    key: 'is_public',
    header: '公开',
    cell: (record) => (
      <Badge variant={record.is_public ? 'success' : 'outline'}>
        {record.is_public ? '是' : '否'}
      </Badge>
    ),
    minWidth: '120px',
  },
  {
    key: 'utime',
    header: '更新时间',
    cell: (record) => formatDateTime(record.utime),
    minWidth: '180px',
  },
];

type OptionFormValues = {
  key: string;
  value: string;
  description?: string;
  is_public?: boolean;
};

export function OptionsPage() {
  const [editingOption, setEditingOption] = useState<OptionItem | null>(null);
  const [formValues, setFormValues] = useState<OptionFormValues>({
    key: '',
    value: '',
    description: '',
    is_public: false,
  });
  const queryClient = useQueryClient();
  const options = useQuery({ queryKey: ['options'], queryFn: api.listOptions });
  const mutation = useMutation({
    mutationFn: api.updateOption,
    onSuccess: async () => {
      toast.success('配置已更新');
      await queryClient.invalidateQueries({ queryKey: ['options'] });
    },
    onError: (error) => {
      if (!isApiError(error)) {
        toast.error((error as Error).message);
      }
    },
  });

  const optionTotal = options.data?.length ?? 0;
  const publicOptionTotal = options.data?.filter((item) => item.is_public).length ?? 0;

  function openEditor(option: OptionItem) {
    setEditingOption(option);
    setFormValues({
      key: option.option_key,
      value: option.option_value,
      description: option.description,
      is_public: option.is_public,
    });
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    await mutation.mutateAsync(formValues);
    setEditingOption(null);
  }

  return (
    <div className="space-y-6">
      <PageHeader
        badge="Config"
        title="系统配置"
        description="查看和维护系统级配置项，支持公开标记、说明信息与配置值编辑。"
        actions={
          <Button
            variant="outline"
            onClick={() => options.refetch()}
            disabled={options.isFetching || mutation.isPending}
          >
            <RefreshCw className={cn('h-4 w-4', options.isFetching && 'animate-spin')} />
            刷新配置
          </Button>
        }
      />

      <section className="grid gap-4 md:grid-cols-3">
        <StatCard
          label="配置项总数"
          value={optionTotal}
          hint="系统配置仓库中的全部键值"
          icon={Settings2}
          tone="primary"
        />
        <StatCard
          label="公开配置项"
          value={publicOptionTotal}
          hint="可对外读取或展示的配置"
          icon={LockKeyhole}
        />
        <StatCard
          label="当前状态"
          value={editingOption ? '编辑中' : '待操作'}
          hint={editingOption?.option_key || '选择任意配置项开始编辑'}
          icon={Edit3}
          tone="accent"
        />
      </section>

      {options.isError ? (
        <Alert variant="destructive">
          <AlertIcon className="h-4 w-4" />
          <AlertTitle>配置列表加载失败</AlertTitle>
          <AlertDescription>{(options.error as Error).message}</AlertDescription>
        </Alert>
      ) : null}

      <DataTable<OptionItem>
        data={options.data ?? []}
        columns={[
          ...columns,
          {
            key: 'actions',
            header: '操作',
            cell: (record) => (
              <Button variant="ghost" size="sm" onClick={() => openEditor(record)}>
                <Edit3 className="h-4 w-4" />
                编辑
              </Button>
            ),
            minWidth: '120px',
          },
        ]}
        rowKey={(row) => row.id}
        loading={options.isFetching || mutation.isPending}
        emptyTitle="暂无配置项"
        emptyDescription="当前环境还没有可编辑的系统配置。"
      />

      <Sheet
        open={!!editingOption}
        onOpenChange={(open) => {
          if (!open) {
            setEditingOption(null);
          }
        }}
      >
        <SheetContent side="right" className="w-full max-w-2xl">
          <SheetHeader>
            <SheetTitle>
              {editingOption ? `编辑配置 · ${editingOption.option_key}` : '编辑配置'}
            </SheetTitle>
            <SheetDescription>
              修改配置值、说明信息和公开状态，保存后会立即刷新当前列表。
            </SheetDescription>
          </SheetHeader>

          <form id="option-form" className="space-y-5 py-8" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label htmlFor="option-key">Key</Label>
              <Input id="option-key" value={formValues.key} disabled />
            </div>

            <div className="space-y-2">
              <Label htmlFor="option-value">值</Label>
              <Textarea
                id="option-value"
                value={formValues.value}
                onChange={(event) =>
                  setFormValues((value) => ({ ...value, value: event.target.value }))
                }
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="option-description">描述</Label>
              <Input
                id="option-description"
                value={formValues.description ?? ''}
                onChange={(event) =>
                  setFormValues((value) => ({ ...value, description: event.target.value }))
                }
                placeholder="补充配置说明，方便后续维护"
              />
            </div>

            <div className="flex items-center justify-between rounded-[1.4rem] border border-border/70 bg-muted/25 px-4 py-4">
              <div>
                <div className="font-medium text-foreground">公开可读</div>
                <div className="text-sm text-muted-foreground">
                  开启后，该配置项可以按公开配置进行读取。
                </div>
              </div>
              <Switch
                checked={!!formValues.is_public}
                onCheckedChange={(checked) =>
                  setFormValues((value) => ({ ...value, is_public: checked }))
                }
              />
            </div>
          </form>

          <SheetFooter>
            <Button variant="outline" onClick={() => setEditingOption(null)}>
              取消
            </Button>
            <Button form="option-form" type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? '保存中...' : '保存变更'}
            </Button>
          </SheetFooter>
        </SheetContent>
      </Sheet>
    </div>
  );
}
