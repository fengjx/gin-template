import { cn } from '@/lib/utils';
import { type VariantProps, cva } from 'class-variance-authority';
import { AlertCircle } from 'lucide-react';
import * as React from 'react';

const alertVariants = cva(
  'relative w-full rounded-[1.4rem] border px-4 py-4 text-sm shadow-sm [&>svg+div]:translate-y-[-2px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg~*]:pl-8',
  {
    variants: {
      variant: {
        default: 'border-border/70 bg-card text-card-foreground',
        destructive: 'border-rose-200 bg-rose-50 text-rose-700 [&>svg]:text-rose-600',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  },
);

const Alert = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement> & VariantProps<typeof alertVariants>
>(({ className, variant, ...props }, ref) => (
  <div ref={ref} role="alert" className={cn(alertVariants({ variant }), className)} {...props} />
));
Alert.displayName = 'Alert';

const AlertTitle = React.forwardRef<HTMLHeadingElement, React.HTMLAttributes<HTMLHeadingElement>>(
  ({ className, ...props }, ref) => (
    <h5
      ref={ref}
      className={cn('mb-1 font-semibold leading-none tracking-tight', className)}
      {...props}
    />
  ),
);
AlertTitle.displayName = 'AlertTitle';

const AlertDescription = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
  <div ref={ref} className={cn('text-sm [&_p]:leading-6', className)} {...props} />
));
AlertDescription.displayName = 'AlertDescription';

function AlertIcon(props: React.ComponentProps<typeof AlertCircle>) {
  return <AlertCircle {...props} />;
}

export { Alert, AlertDescription, AlertIcon, AlertTitle };
