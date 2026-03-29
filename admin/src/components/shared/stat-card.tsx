import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import type { LucideIcon } from 'lucide-react';

type StatCardProps = {
  label: string;
  value: string | number;
  hint?: string;
  icon: LucideIcon;
  tone?: 'default' | 'primary' | 'accent';
};

const toneMap: Record<NonNullable<StatCardProps['tone']>, string> = {
  default: 'from-white via-white to-slate-50',
  primary: 'from-teal-50 via-white to-cyan-50',
  accent: 'from-slate-900 via-slate-900 to-teal-900 text-white',
};

export function StatCard({ label, value, hint, icon: Icon, tone = 'default' }: StatCardProps) {
  const accent = tone === 'accent';

  return (
    <Card
      className={cn(
        'overflow-hidden bg-gradient-to-br',
        toneMap[tone],
        accent && 'border-slate-800 text-white',
      )}
    >
      <CardContent className="flex items-start justify-between p-6">
        <div className="space-y-3">
          <p
            className={cn(
              'text-xs font-semibold uppercase tracking-[0.24em]',
              accent ? 'text-slate-300' : 'text-muted-foreground',
            )}
          >
            {label}
          </p>
          <div
            className={cn(
              'text-3xl font-semibold tracking-tight',
              accent ? 'text-white' : 'text-foreground',
            )}
          >
            {value}
          </div>
          {hint ? (
            <p className={cn('text-sm', accent ? 'text-slate-300' : 'text-muted-foreground')}>
              {hint}
            </p>
          ) : null}
        </div>
        <div
          className={cn(
            'flex h-12 w-12 items-center justify-center rounded-2xl border',
            accent
              ? 'border-white/10 bg-white/10 text-white'
              : 'border-primary/10 bg-primary/10 text-primary',
          )}
        >
          <Icon className="h-5 w-5" />
        </div>
      </CardContent>
    </Card>
  );
}
