import { cn } from '@/lib/utils';
import type { ReactNode } from 'react';

type PageHeaderProps = {
  badge?: string;
  title: string;
  description: string;
  actions?: ReactNode;
  className?: string;
};

export function PageHeader({ badge, title, description, actions, className }: PageHeaderProps) {
  return (
    <div
      className={cn('flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between', className)}
    >
      <div className="space-y-3">
        {badge ? (
          <div className="inline-flex rounded-full border border-primary/20 bg-primary/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-primary">
            {badge}
          </div>
        ) : null}
        <div className="space-y-2">
          <h1 className="font-display text-3xl font-semibold tracking-tight text-foreground sm:text-4xl">
            {title}
          </h1>
          <p className="max-w-3xl text-sm leading-7 text-muted-foreground sm:text-base">
            {description}
          </p>
        </div>
      </div>
      {actions ? <div className="flex flex-wrap items-center gap-3">{actions}</div> : null}
    </div>
  );
}
