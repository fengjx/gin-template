import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { cn } from '@/lib/utils';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { type ReactNode, useEffect, useMemo, useState } from 'react';

export type DataColumn<T> = {
  key: string;
  header: ReactNode;
  cell: (row: T) => ReactNode;
  className?: string;
  headerClassName?: string;
  minWidth?: string;
};

type DataTableProps<T> = {
  data: T[];
  columns: DataColumn<T>[];
  rowKey: (row: T) => string;
  toolbar?: ReactNode;
  loading?: boolean;
  emptyTitle?: string;
  emptyDescription?: string;
  defaultPageSize?: number;
  pageSizeOptions?: number[];
  className?: string;
};

export function DataTable<T>({
  data,
  columns,
  rowKey,
  toolbar,
  loading = false,
  emptyTitle = '暂无数据',
  emptyDescription = '稍后再试，或调整当前筛选条件。',
  defaultPageSize = 10,
  pageSizeOptions = [10, 20, 50],
  className,
}: DataTableProps<T>) {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(defaultPageSize);

  const totalPages = Math.max(1, Math.ceil(data.length / pageSize));
  const paginatedData = useMemo(() => {
    const start = (page - 1) * pageSize;
    return data.slice(start, start + pageSize);
  }, [data, page, pageSize]);

  useEffect(() => {
    if (page > totalPages) {
      setPage(totalPages);
    }
  }, [page, totalPages]);

  return (
    <Card className={cn('overflow-hidden', className)}>
      {toolbar ? <div className="border-b border-border/60 px-6 py-4">{toolbar}</div> : null}
      <CardContent className="space-y-5 p-0">
        {loading ? (
          <div className="space-y-4 p-6">
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-20 w-full" />
            <Skeleton className="h-20 w-full" />
            <Skeleton className="h-20 w-full" />
          </div>
        ) : data.length === 0 ? (
          <div className="flex min-h-[280px] flex-col items-center justify-center gap-3 px-6 text-center">
            <div className="rounded-full border border-dashed border-border bg-muted/30 px-4 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-muted-foreground">
              Empty
            </div>
            <div className="space-y-2">
              <h3 className="text-lg font-semibold text-foreground">{emptyTitle}</h3>
              <p className="max-w-md text-sm leading-6 text-muted-foreground">{emptyDescription}</p>
            </div>
          </div>
        ) : (
          <>
            <Table className="min-w-[920px]">
              <TableHeader>
                <TableRow className="hover:bg-transparent">
                  {columns.map((column) => (
                    <TableHead
                      key={column.key}
                      className={column.headerClassName}
                      style={column.minWidth ? { minWidth: column.minWidth } : undefined}
                    >
                      {column.header}
                    </TableHead>
                  ))}
                </TableRow>
              </TableHeader>
              <TableBody>
                {paginatedData.map((row) => (
                  <TableRow key={rowKey(row)}>
                    {columns.map((column) => (
                      <TableCell key={column.key} className={column.className}>
                        {column.cell(row)}
                      </TableCell>
                    ))}
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            <div className="flex flex-col gap-4 border-t border-border/60 px-6 py-4 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
              <div>
                共 {data.length} 条记录，当前第 {page} / {totalPages} 页
              </div>

              <div className="flex flex-wrap items-center gap-3">
                <label className="flex items-center gap-2">
                  <span>每页</span>
                  <select
                    value={pageSize}
                    onChange={(event) => {
                      setPageSize(Number(event.target.value));
                      setPage(1);
                    }}
                    className="h-10 rounded-full border border-input bg-background px-3 text-foreground outline-none ring-offset-background transition-colors focus:ring-2 focus:ring-ring focus:ring-offset-2"
                  >
                    {pageSizeOptions.map((option) => (
                      <option key={option} value={option}>
                        {option}
                      </option>
                    ))}
                  </select>
                </label>

                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((value) => Math.max(1, value - 1))}
                    disabled={page <= 1}
                  >
                    <ChevronLeft className="h-4 w-4" />
                    上一页
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((value) => Math.min(totalPages, value + 1))}
                    disabled={page >= totalPages}
                  >
                    下一页
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
