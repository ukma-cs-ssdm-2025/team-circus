import * as Y from "yjs";
import { WebsocketProvider } from "y-websocket";
import { ENV } from "../config/env";

export type CollaborativeUser = {
  id: string;
  name?: string;
  color?: string;
  cursor?: number;
  cursorPosition?: number;
  role?: string;
};

type ProviderOptions = {
  baseUrl?: string;
  path?: string;
  queryParams?: Record<string, string>;
  roomParams?: Record<string, string>;
};

export class YjsCollaborativeEditor {
  private doc: Y.Doc;
  private text: Y.Text;
  private provider: WebsocketProvider;

  constructor(
    documentId: string,
    user: CollaborativeUser,
    options?: ProviderOptions,
  ) {
    this.doc = new Y.Doc();
    this.text = this.doc.getText("content");
    const { roomParams, ...restOptions } = options || {};
    const wsUrl = this.buildWebsocketUrl(documentId, restOptions);
    const roomName = this.buildRoomName(documentId, roomParams);

    this.provider = new WebsocketProvider(wsUrl, roomName, this.doc, {
      connect: true,
    });

    this.provider.awareness.setLocalStateField("user", {
      id: user.id,
      name: user.name || user.id,
      color: user.color || randomColor(user.id),
      role: user.role,
    });
  }

  getText() {
    return this.text;
  }

  getProvider() {
    return this.provider;
  }

  getDoc() {
    return this.doc;
  }

  destroy() {
    this.provider.destroy();
    this.doc.destroy();
  }

  private buildWebsocketUrl(
    _documentId: string,
    options: ProviderOptions = {},
  ) {
    const { baseUrl, path, queryParams } = options;
    const endpoint = baseUrl || ENV.API_BASE_URL;
    const url = new URL(endpoint);
    url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
    // y-websocket appends the room to the URL; keep base path without the id
    url.pathname = path || "/ws/documents";
    if (queryParams && Object.keys(queryParams).length > 0) {
      const search = new URLSearchParams(queryParams);
      url.search = search.toString();
    }
    return url.toString();
  }

  private buildRoomName(documentId: string, roomParams?: Record<string, string>) {
    if (!roomParams || Object.keys(roomParams).length === 0) {
      return documentId;
    }
    const search = new URLSearchParams(roomParams);
    return `${documentId}?${search.toString()}`;
  }
}

export function randomColor(seed: string = "") {
  let hash = 0;
  for (let i = 0; i < seed.length; i += 1) {
    // simple string hash to keep colors stable per user
    hash = seed.charCodeAt(i) + ((hash << 5) - hash);
    hash &= hash;
  }

  const hue = Math.abs(hash) % 360;
  return `hsl(${hue}, 70%, 60%)`;
}
