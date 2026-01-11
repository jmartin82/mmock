import { createContext, useContext, useState, useEffect, useCallback, useRef, ReactNode } from 'react';
import { getBaseUrl, Transaction } from '@/lib/api';

export type ConnectionStatus = 'connected' | 'disconnected' | 'connecting';

interface WebSocketContextType {
  transactions: Transaction[];
  status: ConnectionStatus;
  totalCount: number;
  unreadCount: number;
  clearTransactions: () => void;
  reconnect: () => void;
  markAsRead: () => void;
}

const WebSocketContext = createContext<WebSocketContextType | null>(null);

export const WebSocketProvider = ({ children }: { children: ReactNode }) => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [status, setStatus] = useState<ConnectionStatus>('disconnected');
  const [unreadCount, setUnreadCount] = useState(0);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return;

    setStatus('connecting');

    const baseUrl = getBaseUrl();
    const wsUrl = baseUrl
      .replace(/^https:\/\//, 'wss://')
      .replace(/^http:\/\//, 'ws://') + '/echo';

    try {
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        setStatus('connected');
        if (reconnectTimeoutRef.current) {
          clearTimeout(reconnectTimeoutRef.current);
          reconnectTimeoutRef.current = null;
        }
      };

      ws.onmessage = (event) => {
        try {
          const transaction: Transaction = JSON.parse(event.data);
          setTransactions((prev) => [transaction, ...prev]);
          
          // Increment unread count if tab is hidden
          if (document.hidden) {
            setUnreadCount((prev) => prev + 1);
          }
        } catch (e) {
          console.error('Failed to parse WebSocket message:', e);
        }
      };

      ws.onclose = () => {
        setStatus('disconnected');
        wsRef.current = null;
        // Auto-reconnect after 3 seconds
        reconnectTimeoutRef.current = setTimeout(() => {
          connect();
        }, 3000);
      };

      ws.onerror = () => {
        setStatus('disconnected');
      };
    } catch (e) {
      setStatus('disconnected');
      console.error('WebSocket connection error:', e);
    }
  }, []);

  const clearTransactions = useCallback(() => {
    setTransactions([]);
    setUnreadCount(0);
  }, []);

  const reconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.close();
    }
    connect();
  }, [connect]);

  const markAsRead = useCallback(() => {
    setUnreadCount(0);
  }, []);

  // Reset unread count when tab becomes visible
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (!document.hidden) {
        setUnreadCount(0);
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
  }, []);

  // Update document title based on unread count
  useEffect(() => {
    if (unreadCount > 0) {
      document.title = `(${unreadCount}) New Requests - MMock Console`;
    } else {
      document.title = 'MMock Console';
    }
  }, [unreadCount]);

  // Connect on mount
  useEffect(() => {
    connect();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
    };
  }, [connect]);

  return (
    <WebSocketContext.Provider
      value={{
        transactions,
        status,
        totalCount: transactions.length,
        unreadCount,
        clearTransactions,
        reconnect,
        markAsRead,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocketContext = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useWebSocketContext must be used within a WebSocketProvider');
  }
  return context;
};
