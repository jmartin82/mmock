import { cn } from '@/lib/utils';
import { Wifi, WifiOff } from 'lucide-react';
import type { ConnectionStatus as ConnectionStatusType } from '@/hooks/useWebSocket';

interface ConnectionStatusProps {
  status: ConnectionStatusType;
  className?: string;
}

export const ConnectionStatus = ({ status, className }: ConnectionStatusProps) => {
  const statusConfig = {
    connected: {
      dotClass: 'connection-connected',
      label: 'Connected',
      Icon: Wifi,
    },
    disconnected: {
      dotClass: 'connection-disconnected',
      label: 'Disconnected',
      Icon: WifiOff,
    },
    connecting: {
      dotClass: 'connection-connecting',
      label: 'Connecting...',
      Icon: Wifi,
    },
  };

  const config = statusConfig[status];

  return (
    <div className={cn('flex items-center gap-2', className)}>
      <span className={cn('connection-dot', config.dotClass)} />
      <config.Icon className="h-4 w-4 text-muted-foreground" />
      <span className="text-sm text-muted-foreground">{config.label}</span>
    </div>
  );
};
