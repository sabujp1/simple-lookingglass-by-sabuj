'use client';

import { useState, useEffect, useCallback } from 'react';

interface WebSocketMessage {
  type: string;
  data?: any;
  error?: string;
}

interface UseWebSocketOptions {
  onMessage?: (message: WebSocketMessage) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  onError?: (error: Event) => void;
}

export function useWebSocket(queryId: string | null, options?: UseWebSocketOptions) {
  const [connected, setConnected] = useState(false);
  const [messages, setMessages] = useState<WebSocketMessage[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);

  const connect = useCallback(() => {
    if (!queryId) return;

    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080';
    const socket = new WebSocket(`${wsUrl}/ws?query_id=${queryId}`);

    socket.onopen = () => {
      setConnected(true);
      options?.onConnect?.();
    };

    socket.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        setMessages((prev) => [...prev, message]);
        options?.onMessage?.(message);
      } catch (e) {
        console.error('Failed to parse WebSocket message:', e);
      }
    };

    socket.onclose = () => {
      setConnected(false);
      options?.onDisconnect?.();
    };

    socket.onerror = (error) => {
      options?.onError?.(error);
    };

    setWs(socket);

    return () => {
      socket.close();
    };
  }, [queryId, options]);

  const disconnect = useCallback(() => {
    ws?.close();
    setWs(null);
    setConnected(false);
  }, [ws]);

  const send = useCallback((message: WebSocketMessage) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message));
    }
  }, [ws]);

  useEffect(() => {
    if (queryId) {
      connect();
    }
    return () => {
      disconnect();
    };
  }, [queryId, connect, disconnect]);

  return {
    connected,
    messages,
    send,
    disconnect,
    reconnect: connect,
  };
}