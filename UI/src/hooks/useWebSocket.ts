import { useWebSocketContext, ConnectionStatus } from '@/contexts/WebSocketContext';
import { Transaction } from '@/lib/api';

interface UseWebSocketReturn {
  transactions: Transaction[];
  status: ConnectionStatus;
  totalCount: number;
  unreadCount: number;
  clearTransactions: () => void;
  reconnect: () => void;
  markAsRead: () => void;
}

export const useWebSocket = (): UseWebSocketReturn => {
  return useWebSocketContext();
};

export type { ConnectionStatus };
