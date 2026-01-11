import { cn } from '@/lib/utils';

interface StatusBadgeProps {
  statusCode: number;
  className?: string;
}

export const StatusBadge = ({ statusCode, className }: StatusBadgeProps) => {
  const getStatusClass = (code: number): string => {
    if (code >= 200 && code < 300) return 'status-2xx';
    if (code >= 300 && code < 400) return 'status-3xx';
    if (code >= 400 && code < 500) return 'status-4xx';
    if (code >= 500) return 'status-5xx';
    return 'text-muted-foreground';
  };

  return (
    <span className={cn('font-mono text-sm font-medium', getStatusClass(statusCode), className)}>
      {statusCode}
    </span>
  );
};
