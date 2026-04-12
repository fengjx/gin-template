import { JSONEditor, getJSONError } from '@/components/shared/JsonEditor';
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
import { Edit3, LockKeyhole, Plus, RefreshCw, Settings2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { api, isApiError } from '../api/client';
import type { OptionItem, OptionWriteRequest } from '../api/types';
import { formatDateTime, formatTextFallback } from '../utils/format';

const typeOptions: Array<{ value: OptionItem['type']; label: string }> = [
  { value: 'string', label: 'string' },
  { value: 'json', label: 'json' },
];

const statusOptions: Array<{ value: OptionItem['status']; label: string }> = [
  { value: 'online', label: '上线' },
  { value: 'offline', label: '下线' },
];

const typeLabels: Record<OptionItem['type'], string> = {
  string: 'string',
  json: 'json',
};

const statusLabels: Record<OptionItem['status'], string> = {
  online: '上线',
  offline: '下线',
};

function normalizeOption(item: OptionItem): OptionItem {
  return {
    ...item,
    type: item.type || 'string',
    status: item.status || 'online',
  };
}

const columns: DataColumn<OptionItem>[] = [
  {
    key: 'option_key',
    header: 'Key',
    cell: (record) => <div className="font-medium text-foreground">{record.option_key}</div>,
    minWidth: '220px',
  },
  {
    key: 'option_value',
    header: '值',
    cell: (record) => (
      <p className="max-w-[360px] truncate text-muted-foreground">
        {formatTextFallback(record.option_value)}
      </p>
    ),
    minWidth: '280px',
  },
  {
    key: 'type',
    header: '类型',
    cell: (record) => <Badge variant="info">{typeLabels[record.type]}</Badge>,
    minWidth: '120px',
  },
  {
    key: 'status',
    header: '状态',
    cell: (record) => (
      <Badge variant={record.status === 'online' ? 'success' : 'outline'}>
        {statusLabels[record.status]}
      </Badge>
    ),
    minWidth: '120px',
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

type OptionFormValues = OptionWriteRequest;

const defaultFormValues: OptionFormValues = {
  key: '',
  value: '',
  description: '',
  is_public: false,
  type: 'string',
  status: 'online',
};

export function OptionsPage() {
  const [sheetOpen, setSheetOpen] = useState(false);
  const [sheetMode, setSheetMode] = useState<'create' | 'edit'>('create');
  const [editingOption, setEditingOption] = useState<OptionItem | null>(null);
  const [formValues, setFormValues] = useState<OptionFormValues>(defaultFormValues);
  const queryClient = useQueryClient();
  const options = useQuery({ queryKey: ['options'], queryFn: api.listOptions });
  const saveMutation = useMutation({
    mutationFn: async (payload: { mode: 'create' | 'edit'; values: OptionFormValues }) => {
      if (payload.mode === 'create') {
        return api.createOption({
          ...payload.values,
          key: payload.values.key.trim(),
          description: payload.values.description?.trim() || '',
        });
      }

      return api.updateOption({
        ...payload.values,
        key: payload.values.key.trim(),
        description: payload.values.description?.trim() || '',
      });
    },
    onSuccess: async (_, payload) => {
      toast.success(payload.mode === 'create' ? '配置已创建' : '配置已更新');
      await queryClient.invalidateQueries({ queryKey: ['options'] });
      closeSheet();
    },
    onError: (error) => {
      if (!isApiError(error)) {
        toast.error((error as Error).message);
      }
    },
  });

  const normalizedOptions = (options.data ?? []).map(normalizeOption);
  const optionTotal = normalizedOptions.length;
  const onlineOptionTotal = normalizedOptions.filter((item) => item.status === 'online').length;
  const jsonError = formValues.type === 'json' ? getJSONError(formValues.value) : '';
  const canSubmit =
    !saveMutation.isPending &&
    !!formValues.key.trim() &&
    !!formValues.value &&
    (!jsonError || formValues.type !== 'json');

  function closeSheet() {
    setSheetOpen(false);
    setEditingOption(null);
    setSheetMode('create');
    setFormValues(defaultFormValues);
  }

  function openCreateSheet() {
    setSheetMode('create');
    setEditingOption(null);
    setFormValues(defaultFormValues);
    setSheetOpen(true);
  }

  function openEditSheet(option: OptionItem) {
    const normalized = normalizeOption(option);
    setSheetMode('edit');
    setEditingOption(normalized);
    setFormValues({
      key: normalized.option_key,
      value: normalized.option_value,
      description: normalized.description,
      is_public: normalized.is_public,
      type: normalized.type,
      status: normalized.status,
    });
    setSheetOpen(true);
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!canSubmit) {
      return;
    }
    await saveMutation.mutateAsync({ mode: sheetMode, values: formValues });
  }

  return (
    <div className="space-y-6">
      <PageHeader
        badge="Config"
        title="系统配置"
        description="查看和维护系统级配置项，支持新增配置、JSON 编辑和上下线管理。"
        actions={
          <div className="flex flex-wrap items-center gap-3">
            <Button onClick={openCreateSheet}>
              <Plus className="h-4 w-4" />
              新增配置
            </Button>
            <Button
              variant="outline"
              onClick={() => options.refetch()}
              disabled={options.isFetching || saveMutation.isPending}
            >
              <RefreshCw className={cn('h-4 w-4', options.isFetching && 'animate-spin')} />
              刷新配置
            </Button>
          </div>
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
          label="在线配置项"
          value={onlineOptionTotal}
          hint="当前可被业务读取的系统配置"
          icon={LockKeyhole}
        />
        <StatCard
          label="当前状态"
          value={sheetOpen ? (sheetMode === 'create' ? '新建中' : '编辑中') : '待操作'}
          hint={editingOption?.option_key || '可新增或编辑任意配置项'}
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
        data={normalizedOptions}
        columns={[
          ...columns,
          {
            key: 'actions',
            header: '操作',
            cell: (record) => (
              <Button variant="ghost" size="sm" onClick={() => openEditSheet(record)}>
                <Edit3 className="h-4 w-4" />
                编辑
              </Button>
            ),
            minWidth: '120px',
          },
        ]}
        rowKey={(row) => String(row.id)}
        loading={options.isFetching || saveMutation.isPending}
        emptyTitle="暂无配置项"
        emptyDescription="当前环境还没有可编辑的系统配置。"
      />

      <Sheet
        open={sheetOpen}
        onOpenChange={(open) => {
          if (!open) {
            closeSheet();
          }
        }}
      >
        <SheetContent
          side="right"
          className="flex h-full w-full max-w-2xl flex-col overflow-hidden"
        >
          <SheetHeader>
            <SheetTitle>
              {sheetMode === 'create'
                ? '新增配置'
                : editingOption
                  ? `编辑配置 · ${editingOption.option_key}`
                  : '编辑配置'}
            </SheetTitle>
            <SheetDescription>
              维护配置值、类型、公开状态和上下线状态；保存后会立即刷新当前列表。
            </SheetDescription>
          </SheetHeader>

          <div className="min-h-0 flex-1 overflow-y-auto px-1">
            <form id="option-form" className="space-y-5 py-8" onSubmit={handleSubmit}>
              <div className="space-y-2">
                <Label htmlFor="option-key">Key</Label>
                <Input
                  id="option-key"
                  value={formValues.key}
                  disabled={sheetMode === 'edit'}
                  onChange={(event) =>
                    setFormValues((value) => ({ ...value, key: event.target.value }))
                  }
                  placeholder="例如：site_profile"
                  required
                />
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <label className="space-y-2">
                  <span className="text-sm font-medium text-foreground">类型</span>
                  <select
                    aria-label="类型"
                    value={formValues.type}
                    onChange={(event) =>
                      setFormValues((value) => ({
                        ...value,
                        type: event.target.value as OptionItem['type'],
                      }))
                    }
                    className="h-11 w-full rounded-[1.35rem] border border-input bg-background/80 px-4 text-sm text-foreground outline-none ring-offset-background transition-colors focus:ring-2 focus:ring-ring focus:ring-offset-2"
                  >
                    {typeOptions.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </select>
                </label>

                <label className="space-y-2">
                  <span className="text-sm font-medium text-foreground">状态</span>
                  <select
                    aria-label="状态"
                    value={formValues.status}
                    onChange={(event) =>
                      setFormValues((value) => ({
                        ...value,
                        status: event.target.value as OptionItem['status'],
                      }))
                    }
                    className="h-11 w-full rounded-[1.35rem] border border-input bg-background/80 px-4 text-sm text-foreground outline-none ring-offset-background transition-colors focus:ring-2 focus:ring-ring focus:ring-offset-2"
                  >
                    {statusOptions.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </select>
                </label>
              </div>

              <div className="space-y-2">
                <Label htmlFor="option-value">值</Label>
                {formValues.type === 'json' ? (
                  <JSONEditor
                    id="option-value"
                    value={formValues.value}
                    onChange={(value) => setFormValues((current) => ({ ...current, value }))}
                    disabled={saveMutation.isPending}
                  />
                ) : (
                  <Textarea
                    id="option-value"
                    value={formValues.value}
                    onChange={(event) =>
                      setFormValues((value) => ({ ...value, value: event.target.value }))
                    }
                    required
                  />
                )}
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

              <div className="grid gap-4 pb-4 md:grid-cols-2">
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

                <div className="rounded-[1.4rem] border border-border/70 bg-muted/25 px-4 py-4">
                  <div className="font-medium text-foreground">当前运行状态</div>
                  <div className="mt-1 text-sm text-muted-foreground">
                    下线后，该配置项不会再被对外业务读口返回。
                  </div>
                  <div className="mt-3">
                    <Badge variant={formValues.status === 'online' ? 'success' : 'outline'}>
                      {statusLabels[formValues.status]}
                    </Badge>
                  </div>
                </div>
              </div>
            </form>
          </div>

          <SheetFooter className="border-t border-border/60 bg-background/95 pt-4">
            <Button variant="outline" onClick={closeSheet}>
              取消
            </Button>
            <Button form="option-form" type="submit" disabled={!canSubmit}>
              {saveMutation.isPending
                ? sheetMode === 'create'
                  ? '创建中...'
                  : '保存中...'
                : sheetMode === 'create'
                  ? '创建配置'
                  : '保存变更'}
            </Button>
          </SheetFooter>
        </SheetContent>
      </Sheet>
    </div>
  );
}
