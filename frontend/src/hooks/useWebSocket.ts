import { useEffect, useRef, useState, useCallback } from 'react';
import type { WSMessage, Pool, Opportunity } from '../types';

type WebSocketStatus = 'connecting' | 'connected' | 'disconnected' | 'error';

interface UseWebSocketOptions {
  onPoolUpdate?: (pool: Pool) => void;
  onOpportunityAlert?: (opportunity: Opportunity) => void;
  autoReconnect?: boolean;
  reconnectInterval?: number;
}

export function useWebSocket(
  endpoint: 'pools' | 'opportunities',
  options: UseWebSocketOptions = {}
) {
  const {
    onPoolUpdate,
    onOpportunityAlert,
    autoReconnect = true,
    reconnectInterval = 5000,
  } = options;

  const [status, setStatus] = useState<WebSocketStatus>('disconnected');
  const [lastMessage, setLastMessage] = useState<WSMessage | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    const wsUrl = `${protocol}//${host}/ws/${endpoint}`;

    setStatus('connecting');

    try {
      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        setStatus('connected');
      };

      ws.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data);
          setLastMessage(message);

          if (message.type === 'pool_update' && onPoolUpdate && message.data) {
            onPoolUpdate(message.data as Pool);
          }

          if (message.type === 'opportunity_alert' && onOpportunityAlert && message.data) {
            onOpportunityAlert(message.data as Opportunity);
          }
        } catch {
          // Silently ignore parse errors in production
        }
      };

      ws.onerror = () => {
        setStatus('error');
      };

      ws.onclose = () => {
        setStatus('disconnected');

        if (autoReconnect) {
          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, reconnectInterval);
        }
      };

      wsRef.current = ws;
    } catch {
      setStatus('error');
    }
  }, [endpoint, onPoolUpdate, onOpportunityAlert, autoReconnect, reconnectInterval]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }

    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
  }, []);

  const send = useCallback((message: unknown) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    }
  }, []);

  useEffect(() => {
    connect();

    return () => {
      disconnect();
    };
  }, [connect, disconnect]);

  return {
    status,
    lastMessage,
    send,
    connect,
    disconnect,
    isConnected: status === 'connected',
  };
}

// Hook for pool updates
export function usePoolUpdates(onUpdate?: (pool: Pool) => void) {
  return useWebSocket('pools', { onPoolUpdate: onUpdate });
}

// Hook for opportunity alerts
export function useOpportunityAlerts(onAlert?: (opportunity: Opportunity) => void) {
  return useWebSocket('opportunities', { onOpportunityAlert: onAlert });
}
