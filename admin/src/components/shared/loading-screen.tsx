import { cn } from '@/lib/utils';
import { LoaderCircle } from 'lucide-react';

type LoadingScreenProps = {
  message?: string;
  fullscreen?: boolean;
  className?: string;
};

export function LoadingScreen({
  message = '正在加载内容...',
  fullscreen = false,
  className,
}: LoadingScreenProps) {
  return (
    <div
      className={cn(
        'flex items-center justify-center',
        fullscreen ? 'min-h-screen' : 'min-h-[280px]',
        className,
      )}
    >
      <div className="glass-panel flex items-center gap-3 rounded-full px-5 py-3 text-sm text-muted-foreground">
        <LoaderCircle className="h-4 w-4 animate-spin text-primary" />
        <span>{message}</span>
      </div>
    </div>
  );
}
