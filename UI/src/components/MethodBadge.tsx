import { cn } from '@/lib/utils';

interface MethodBadgeProps {
  method: string;
  className?: string;
}

export const MethodBadge = ({ method, className }: MethodBadgeProps) => {
  const methodLower = method.toLowerCase();
  
  const methodClass = {
    get: 'method-get',
    post: 'method-post',
    put: 'method-put',
    delete: 'method-delete',
    patch: 'method-patch',
    head: 'method-head',
    options: 'method-options',
  }[methodLower] || 'method-head';

  return (
    <span className={cn('method-badge', methodClass, className)}>
      {method.toUpperCase()}
    </span>
  );
};
