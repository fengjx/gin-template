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
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Copy,
  FileText,
  HardDriveUpload,
  RefreshCw,
  Search,
  Trash2,
  UploadCloud,
} from 'lucide-react';
import { startTransition, useRef, useState } from 'react';
import { toast } from 'sonner';
import { api, isApiError } from '../api/client';
import type { FileItem } from '../api/types';
import { formatDateTime, formatFileSize } from '../utils/format';

const columns: DataColumn<FileItem>[] = [
  {
    key: 'original_name',
    header: '原文件名',
    cell: (record) => <div className="font-medium text-foreground">{record.original_name}</div>,
    minWidth: '220px',
  },
  {
    key: 'content_type',
    header: '类型',
    cell: (record) => <Badge variant="info">{record.content_type}</Badge>,
    minWidth: '180px',
  },
  {
    key: 'size',
    header: '大小',
    cell: (record) => formatFileSize(record.size),
    minWidth: '120px',
  },
  {
    key: 'path',
    header: '存储路径',
    cell: (record) => (
      <div className="flex items-center gap-2">
        <span className="max-w-[360px] truncate text-muted-foreground">{record.path}</span>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={async () => {
            try {
              await navigator.clipboard.writeText(record.path);
              toast.success('路径已复制');
            } catch (error) {
              toast.error((error as Error).message);
            }
          }}
        >
          <Copy className="h-4 w-4" />
          复制
        </Button>
      </div>
    ),
    minWidth: '340px',
  },
  {
    key: 'ctime',
    header: '上传时间',
    cell: (record) => formatDateTime(record.ctime),
    minWidth: '180px',
  },
];

export function FilesPage() {
  const [inputKeyword, setInputKeyword] = useState('');
  const [keyword, setKeyword] = useState('');
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const queryClient = useQueryClient();
  const files = useQuery({
    queryKey: ['files', keyword],
    queryFn: () => (keyword ? api.searchFiles(keyword) : api.listFiles()),
  });
  const deleteMutation = useMutation({
    mutationFn: api.deleteFile,
    onSuccess: async () => {
      toast.success('文件已删除');
      await queryClient.invalidateQueries({ queryKey: ['files'] });
    },
    onError: (error) => {
      if (!isApiError(error)) {
        toast.error((error as Error).message);
      }
    },
  });

  async function handleFileChange(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }
    try {
      setUploading(true);
      await api.uploadFile(file);
      toast.success('上传成功');
      await queryClient.invalidateQueries({ queryKey: ['files'] });
      event.target.value = '';
    } catch (error) {
      if (!isApiError(error)) {
        toast.error((error as Error).message);
      }
    } finally {
      setUploading(false);
    }
  }

  const totalFiles = files.data?.total ?? files.data?.items.length ?? 0;
  const totalStorage = files.data?.items.reduce((sum, item) => sum + item.size, 0) ?? 0;

  return (
    <div className="space-y-6">
      <PageHeader
        badge="Assets"
        title="文件管理"
        description="管理上传资产、文件分布和存储占用，支持按文件名或路径搜索，并可直接删除无用文件。"
        actions={
          <>
            <Button
              variant="outline"
              onClick={() => files.refetch()}
              disabled={files.isFetching || uploading}
            >
              <RefreshCw className={cn('h-4 w-4', files.isFetching && 'animate-spin')} />
              刷新列表
            </Button>
            <Button onClick={() => fileInputRef.current?.click()} disabled={uploading}>
              <UploadCloud className="h-4 w-4" />
              {uploading ? '上传中...' : '上传文件'}
            </Button>
            <input ref={fileInputRef} type="file" className="hidden" onChange={handleFileChange} />
          </>
        }
      />

      <section className="grid gap-4 md:grid-cols-3">
        <StatCard
          label="当前文件数"
          value={totalFiles}
          hint="当前筛选结果中的文件数量"
          icon={FileText}
          tone="primary"
        />
        <StatCard
          label="占用存储"
          value={formatFileSize(totalStorage)}
          hint="实时汇总当前列表容量"
          icon={HardDriveUpload}
        />
        <StatCard
          label="检索条件"
          value={keyword || '全部文件'}
          hint="支持搜索文件名与存储路径"
          icon={Search}
          tone="accent"
        />
      </section>

      {files.isError ? (
        <Alert variant="destructive">
          <AlertIcon className="h-4 w-4" />
          <AlertTitle>文件列表加载失败</AlertTitle>
          <AlertDescription>{(files.error as Error).message}</AlertDescription>
        </Alert>
      ) : null}

      <DataTable<FileItem>
        data={files.data?.items ?? []}
        columns={[
          ...columns,
          {
            key: 'actions',
            header: '操作',
            cell: (record) => (
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="ghost" size="sm" className="text-rose-600 hover:text-rose-600">
                    <Trash2 className="h-4 w-4" />
                    删除
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>确认删除这个文件吗？</AlertDialogTitle>
                    <AlertDialogDescription>
                      删除后无法恢复，文件路径为 {record.path}。
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>取消</AlertDialogCancel>
                    <AlertDialogAction
                      onClick={async () => {
                        await deleteMutation.mutateAsync(record.id);
                      }}
                    >
                      删除文件
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            ),
            minWidth: '130px',
          },
        ]}
        rowKey={(row) => row.id}
        loading={files.isFetching || deleteMutation.isPending || uploading}
        emptyTitle="没有找到匹配的文件"
        emptyDescription="可以尝试更换搜索关键词，或直接上传新的文件资产。"
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
                placeholder="搜索文件名或存储路径"
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
            <div className="text-sm text-muted-foreground">共 {totalFiles} 个文件</div>
          </form>
        }
      />
    </div>
  );
}
