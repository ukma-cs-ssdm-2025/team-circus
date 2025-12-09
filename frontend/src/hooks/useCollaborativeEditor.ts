import { useCallback, useEffect, useRef, useState } from "react";
import type { Awareness } from "y-protocols/awareness";
import type * as Y from "yjs";
import {
  type CollaborativeUser,
  YjsCollaborativeEditor,
} from "../services/YjsWebsocketProvider";

export interface UseCollaborativeEditorOptions {
  documentId: string;
  user: CollaborativeUser;
  baseUrl?: string;
}

export interface RemoteUser {
  id: string;
  name?: string;
  color?: string;
  cursorPosition?: number;
}

export function useCollaborativeEditor(options: UseCollaborativeEditorOptions) {
  const [content, setContentState] = useState("");
  const [isConnected, setIsConnected] = useState(false);
  const [remoteUsers, setRemoteUsers] = useState<RemoteUser[]>([]);
  const [yDoc, setYDoc] = useState<Y.Doc | null>(null);

  const providerRef = useRef<YjsCollaborativeEditor | null>(null);

  const syncRemoteUsers = useCallback(
    (awareness: Awareness) => {
      const states = Array.from(awareness.getStates().values());
      const others = states
        .map((state: any) => state?.user)
        .filter((user) => user && user.id !== options.user.id)
        .map((user) => ({
          id: user.id,
          name: user.name,
          color: user.color,
          cursorPosition: user.cursor ?? user.cursorPosition,
        }));
      setRemoteUsers(others);
    },
    [options.user.id],
  );

  useEffect(() => {
    if (!options.documentId || !options.user?.id) {
      return;
    }

    const collaborativeProvider = new YjsCollaborativeEditor(
      options.documentId,
      options.user,
      options.baseUrl,
    );
    providerRef.current = collaborativeProvider;
    const text = collaborativeProvider.getText();
    const awareness = collaborativeProvider.getProvider().awareness;

    setYDoc(collaborativeProvider.getDoc());
    setContentState(text.toString());
    syncRemoteUsers(awareness);

    const handleTextChange = () => {
      setContentState(text.toString());
    };
    text.observe(handleTextChange);

    const handleAwarenessChange = () => syncRemoteUsers(awareness);
    awareness.on("change", handleAwarenessChange);

    const handleStatusChange = (event: any) => {
      setIsConnected(event?.status === "connected");
    };
    collaborativeProvider.getProvider().on("status", handleStatusChange);

    return () => {
      text.unobserve(handleTextChange);
      awareness.off("change", handleAwarenessChange);
      collaborativeProvider.getProvider().off("status", handleStatusChange);
      collaborativeProvider.destroy();
      providerRef.current = null;
      setYDoc(null);
      setIsConnected(false);
      setRemoteUsers([]);
    };
  }, [
    options.baseUrl,
    options.documentId,
    options.user.id,
    options.user.name,
    syncRemoteUsers,
  ]);

  const setContent = useCallback((value: string) => {
    const target = providerRef.current?.getText();
    if (!target) {
      setContentState(value);
      return;
    }
    target.delete(0, target.length);
    target.insert(0, value);
  }, []);

  const insertText = useCallback((index: number, value: string) => {
    const target = providerRef.current?.getText();
    if (!target) {
      return;
    }
    target.insert(index, value);
  }, []);

  const deleteText = useCallback((index: number, length: number) => {
    const target = providerRef.current?.getText();
    if (!target) {
      return;
    }
    target.delete(index, length);
  }, []);

  const updateCursorPosition = useCallback((position: number) => {
    const awareness = providerRef.current?.getProvider().awareness;
    if (!awareness) {
      return;
    }
    const currentState = awareness.getLocalState() || {};
    awareness.setLocalState({
      ...currentState,
      cursor: position,
    });
  }, []);

  const awareness = providerRef.current?.getProvider().awareness ?? null;

  return {
    content,
    setContent,
    insertText,
    deleteText,
    updateCursorPosition,
    isConnected,
    remoteUsers,
    yDoc,
    awareness,
  };
}
