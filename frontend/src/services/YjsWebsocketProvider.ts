import * as Y from "yjs";
import { WebsocketProvider } from "y-websocket";
import { ENV } from "../config/env";

export type CollaborativeUser = {
  id: string;
  name?: string;
  color?: string;
};

export class YjsCollaborativeEditor {
  private doc: Y.Doc;
  private text: Y.Text;
  private provider: WebsocketProvider;

  constructor(documentId: string, user: CollaborativeUser, baseUrl?: string) {
    this.doc = new Y.Doc();
    this.text = this.doc.getText("content");
    const wsUrl = this.buildWebsocketUrl(documentId, baseUrl);

    this.provider = new WebsocketProvider(wsUrl, documentId, this.doc, {
      connect: true,
    });

    this.provider.awareness.setLocalStateField("user", {
      id: user.id,
      name: user.name || user.id,
      color: user.color || randomColor(user.id),
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

  private buildWebsocketUrl(documentId: string, baseUrl?: string) {
    const endpoint = baseUrl || ENV.API_BASE_URL;
    const url = new URL(endpoint);
    url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
    // y-websocket appends the room to the URL; keep base path without the id
    url.pathname = "/ws/documents";
    return url.toString();
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
